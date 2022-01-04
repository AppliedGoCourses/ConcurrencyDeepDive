package main

import "testing"

func Benchmark_batchQueryWithoutPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		batchQueryWithoutPool()
	}
}

func Benchmark_batchQueryWithAutoPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		batchQueryWithAutoPool()
	}
}
