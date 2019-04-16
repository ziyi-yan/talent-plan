package main

import (
	"pingcap/talentplan/tidb/mergesort/kway"
)

// MergeSort performs the merge sort algorithm.
// Please supplement this function to accomplish the home work.
func MergeSort(src []int64) {
	kway.Sort(src)
}
