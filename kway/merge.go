package kway

import (
	"math"
	"runtime"
	"sort"
	"sync"
)

// Sort sorts the array with multiple goroutines using K-way Merge Sort algorithm.
func Sort(data []int64) {
	var (
		n = runtime.NumCPU()
	)

	var wg sync.WaitGroup
	wg.Add(n)

	segs := split(data, n)

	for _, seg := range segs {
		seg := seg
		go func() {
			sort.Sort(Int64Slice(seg))
			wg.Done()
		}()
	}

	wg.Wait()
	sorted := Merge(segs...)

	copy(data, sorted)
}

func split(data []int64, n int) (segs [][]int64) {
	curr := data
	size := len(data) / n
	for i := 0; i < n-1; i++ {
		segs = append(segs, curr[:size])
		curr = curr[size:]
	}
	segs = append(segs, curr) // add remain elements to the last segment
	return
}

// Merge performs a K-way Merge using Tournament Tree
// algorithm (https://en.wikipedia.org/wiki/K-way_merge_algorithm#Tournament_Tree).
func Merge(arrs ...[]int64) []int64 {
	// resultArr, resultIdx
	total := 0
	for _, a := range arrs {
		total += len(a)
	}
	resultArr := make([]int64, total)
	resultIdx := 0

	// arrs, idxArr, tournamentTree
	n := len(arrs)
	idxArr := make([]int, n)
	tournamentTree := NewTournament(arrs)

	// run tournament
	for {
		winner, winnerIdx := tournamentTree.Winner()

		resultArr[resultIdx] = winner
		resultIdx++
		if resultIdx == total {
			break
		}

		if nextElement, arr := idxArr[winnerIdx]+1, arrs[winnerIdx]; nextElement < len(arr) {
			tournamentTree.Replace(winnerIdx, arr[nextElement])
			idxArr[winnerIdx] = nextElement
		} else {
			tournamentTree.Replace(winnerIdx, math.MaxInt64)
		}
	}

	return resultArr
}

type Tournament struct {
	tree []Node
	size int
	arrs [][]int64
}

type Node struct {
	element  int64
	ArrIndex int
}

func NewTournament(arrs [][]int64) *Tournament {
	n := len(arrs)
	size := 1
	nextLine := 2
	for ; nextLine < n; nextLine *= 2 {
		size += nextLine
	}
	remain := n - nextLine/2
	size += 2 * remain

	tm := &Tournament{
		tree: make([]Node, size+1),
		size: size,
		arrs: arrs,
	}
	tm.initialBattle(arrs)
	return tm
}

func (tm *Tournament) Replace(arrIndex int, element int64) {
	tm.setElement(arrIndex, element)
	tm.promote(tm.toTreeIndex(arrIndex), Node{element, arrIndex})
}

func (tm *Tournament) Winner() (int64, int) {
	winner := tm.tree[0]
	return winner.element, winner.ArrIndex
}

func (tm *Tournament) setElement(arrIndex int, element int64) {
	tm.tree[tm.toTreeIndex(arrIndex)] = Node{element, arrIndex}
}

func (tm *Tournament) toTreeIndex(idx int) int {
	return tm.size - idx
}

func (tm *Tournament) initialBattle(arrs [][]int64) {
	for i := 0; i < len(arrs); i++ {
		tm.setElement(i, arrs[i][0])
	}
	root := 1
	winner := tm.build(root)
	tm.tree[0] = winner
}

func (tm *Tournament) build(index int) Node {
	if tm.isLeaf(index) {
		return tm.tree[index]
	}
	leftWinner := tm.build(index * 2)
	rightWinner := tm.build(index*2 + 1)
	winner, loser := play(leftWinner, rightWinner)
	tm.tree[index] = loser
	return winner
}

func (tm *Tournament) isLeaf(index int) bool {
	return len(tm.tree)-len(tm.arrs) <= index && index < len(tm.tree)
}

func (tm *Tournament) promote(index int, node Node) {
	curr := index / 2
	promotingNode := node
	for curr > 0 {
		if promotingNode.element > tm.tree[curr].element {
			loser := promotingNode
			promotingNode = tm.tree[curr]
			tm.tree[curr] = loser
		}
		curr = curr / 2
	}
	tm.tree[0] = promotingNode
}

func play(leftWinner Node, rightWinner Node) (Node, Node) {
	if leftWinner.element < rightWinner.element {
		return leftWinner, rightWinner
	} else {
		return rightWinner, leftWinner
	}
}

// Int64Slice attaches the methods of Interface to []int64, sorting in increasing order.
type Int64Slice []int64

func (p Int64Slice) Len() int           { return len(p) }
func (p Int64Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Int64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
