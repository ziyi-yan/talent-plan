package kway

// Tournament maintains a tournament tree with loser nodes
type Tournament struct {
	tree []Node
	size int
	arrs [][]int64
}

// Node stores losers of each match in the tournament tree
type Node struct {
	element  int64
	arrIndex int
}

// NewTournament create a new tournament tree
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

// Replace puts next element of the arrs[arrIndex] array to its corresponding leaf node
func (tm *Tournament) Replace(arrIndex int, element int64) {
	tm.setElement(arrIndex, element)
	tm.promote(tm.toTreeIndex(arrIndex))
}

// Winner returns the current winner of the tournament
func (tm *Tournament) Winner() (int64, int) {
	winner := tm.tree[0]
	return winner.element, winner.arrIndex
}

func (tm *Tournament) setElement(arrIndex int, element int64) {
	tm.tree[tm.toTreeIndex(arrIndex)] = Node{element, arrIndex}
}

func (tm *Tournament) toTreeIndex(idx int) int {
	return tm.size - idx
}

// initialBattle fills up the leaf nodes and runs the first tournament building the loser tree
func (tm *Tournament) initialBattle(arrs [][]int64) {
	for i := 0; i < len(arrs); i++ {
		tm.setElement(i, arrs[i][0])
	}
	root := 1
	winner := tm.build(root)
	tm.tree[0] = winner
}

// build builds the loser tree rooted at index node
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

// promote updates the loser tree from a leaf node index to the root
func (tm *Tournament) promote(index int) {
	curr := index / 2
	promotingNode := tm.tree[index]
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
