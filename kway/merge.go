package kway

import (
	"runtime"
	"sort"
	"sync"
)

// Sort sorts the array with multiple goroutines using K-way Merge Sort algorithm.
func Sort(data []int64) {
	var (
		n    = runtime.NumCPU()
		size = len(data) / n
	)

	var wg sync.WaitGroup
	wg.Add(n)

	segs := split(data, n)

	curr := data
	for _, seg := range segs {
		go func() {
			sort.Sort(Int64Slice(seg))
			wg.Done()
		}()
		curr = curr[size:]
	}

	wg.Wait()
	sorted := Merge(segs...)

	copy(data, sorted)
}

func split(data []int64, n int) (segs [][]int64) {
	curr := data
	size := len(data) / n
	for i := 0; i < n; i++ {
		segs = append(segs, curr[:size])
		curr = curr[size:]
	}
	return
}

// Merge performs a K-way Merge using Tournament Tree
// algorithm (https://en.wikipedia.org/wiki/K-way_merge_algorithm#Tournament_Tree).
func Merge(arrs ...[]int64) []int64 {
	panic("implement me")
}

// Int64Slice attaches the methods of Interface to []int64, sorting in increasing order.
type Int64Slice []int64

func (p Int64Slice) Len() int           { return len(p) }
func (p Int64Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Int64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
