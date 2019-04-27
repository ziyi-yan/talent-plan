package main

import (
	"bytes"
	"container/heap"
	"fmt"
	"strconv"
	"strings"
)

// URLTop10 .
func URLTop10(nWorkers int) RoundsArgs {
	// YOUR CODE HERE :)
	// And don't forget to document your idea.
	var args RoundsArgs
	// round 1: do url count
	args = append(args, RoundArgs{
		MapFunc:    URLCountMap,
		ReduceFunc: URLCountReduce,
		NReduce:    nWorkers,
	})
	// round 2: sort and get the 10 most frequent URLs
	args = append(args, RoundArgs{
		MapFunc:    URLTop10Map,
		ReduceFunc: URLTop10Reduce,
		NReduce:    1,
	})
	return args
}

func URLCountMap(_ string, contents string) []KeyValue {
	lines := strings.Split(string(contents), "\n")
	m := make(map[string]int)
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if len(l) == 0 {
			continue
		}
		m[l] = m[l] + 1
	}
	kvs := make([]KeyValue, 0, len(lines))
	for url, count := range m {
		kvs = append(kvs, KeyValue{Key: url, Value: strconv.Itoa(count)})
	}
	return kvs
}

func URLCountReduce(key string, values []string) string {
	var count int
	for _, v := range values {
		c, err := strconv.Atoi(v)
		if err != nil {
			panic(fmt.Sprintf("cannot convert string %q to int", v))
		}
		count += c
	}
	return fmt.Sprintf("%s %s\n", key, strconv.Itoa(count))
}

const n = 10

func URLTop10Map(_ string, contents string) []KeyValue {
	lines := strings.Split(contents, "\n")
	var values []string
	for _, l := range lines {
		trim := strings.TrimSpace(l)
		if len(trim) == 0 {
			continue
		}
		values = append(values, trim)
	}

	urls := Top10URL(values)

	kvs := make([]KeyValue, 0, len(urls))
	for _, url := range urls {
		kvs = append(kvs, KeyValue{"", fmt.Sprintf("%s %d", url.URL, url.Count)})
	}
	return kvs
}

func URLTop10Reduce(_ string, values []string) string {
	urls := Top10URL(values)

	var sorted []URLCount
	for len(urls) > 0 {
		url := heap.Pop(&urls).(URLCount)
		sorted = append(sorted, url)
	}

	buf := new(bytes.Buffer)
	// write more frequent urls to output first for testing
	for i := len(sorted) - 1; i >= 0; i-- {
		url := sorted[i]
		fmt.Fprintf(buf, "%s: %d\n", url.URL, url.Count)
	}
	return buf.String()
}

func Top10URL(values []string) URLCountMinHeap {
	if len(values) <= n {
		urls := make([]URLCount, 0, len(values))
		for _, v := range values {
			url, count := splitLine(v)
			urls = append(urls, URLCount{url, count})
		}
		h := URLCountMinHeap(urls)
		heap.Init(&h)
		return h
	}

	urls := make([]URLCount, 0, n)
	for i := 0; i < n; i++ {
		url, count := splitLine(values[i])
		urls = append(urls, URLCount{url, count})
	}

	h := URLCountMinHeap(urls)
	heap.Init(&h)
	for i := n; i < len(values); i++ {
		url, count := splitLine(values[i])
		newVal := URLCount{url, count}
		if URLCountLess(newVal, h[0]) { // or `count <= h[0].Count` ?
			continue
		}
		h[0] = URLCount{url, count} // replace the minimum value with new value
		heap.Fix(&h, 0)
	}
	return h
}

type URLCount struct {
	URL   string
	Count int
}

type URLCountMinHeap []URLCount

func (h URLCountMinHeap) Len() int {
	return len(h)
}

func (h URLCountMinHeap) Less(i int, j int) bool {
	return URLCountLess(h[i], h[j])
}

func URLCountLess(lhs, rhs URLCount) bool {
	if lhs.Count == rhs.Count {
		return lhs.URL >= rhs.URL // for testing, per utils.go:36
	}

	return lhs.Count <= rhs.Count
}

func (h URLCountMinHeap) Swap(i int, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *URLCountMinHeap) Push(x interface{}) {
	*h = append(*h, x.(URLCount))
}

func (h *URLCountMinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func splitLine(line string) (string, int) {
	tmp := strings.Split(line, " ")
	url, cnt := tmp[0], tmp[1]
	count, err := strconv.Atoi(cnt)
	if err != nil {
		panic(fmt.Sprintf("cannot convert cnt %q to int", cnt))
	}
	return url, count
}
