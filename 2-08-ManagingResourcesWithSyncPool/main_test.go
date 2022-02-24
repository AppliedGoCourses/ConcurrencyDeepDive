package main

import (
	"os"
	"testing"
)

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

func TestMain(m *testing.M) {
	stdout, _ = os.Open(os.DevNull)
	os.Exit(m.Run())
}
