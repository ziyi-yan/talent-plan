.PHONY: all

all: test bench

test:
	go test

bench:
	go test -bench Benchmark -run xx -count 5 -benchmem

profile:
	go test -bench BenchmarkMergeSort -run xx -count 5 -benchmem -blockprofile block.out -cpuprofile cpu.out -memprofile mem.out -outputdir _pprof/${OUTPUT_DIR}
