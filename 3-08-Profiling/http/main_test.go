package main

import "testing"

func Benchmark_fillSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fillSlice()
	}
}

func Benchmark_fillSliceAppend(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fillSliceAppend()
	}
}
