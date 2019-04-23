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

	if len(data) < 10*n {
		sort.Sort(Int64Slice(data))
		return
	}

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
			idxArr[winnerIdx]++
		} else {
			tournamentTree.Replace(winnerIdx, math.MaxInt64)
		}
	}

	return resultArr
}

// Int64Slice attaches the methods of Interface to []int64, sorting in increasing order.
type Int64Slice []int64

func (p Int64Slice) Len() int           { return len(p) }
func (p Int64Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Int64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
