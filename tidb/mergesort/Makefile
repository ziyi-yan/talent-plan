.PHONY: all

all: test bench

test:
	go test

bench:
	go test -bench Benchmark -run xx -count 5 -benchmem

profile:
	go test -bench BenchmarkMergeSort -run xx -benchmem -blockprofile block.out -cpuprofile cpu.out -memprofile mem.out -trace=trace.out -outputdir _pprof/${OUTPUT_DIR}

benchcmp:
	benchcmp normal.bench merge.bench
