package fillslice

import "testing"

func Benchmark_MakeFillSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = MakeFillSlice()
	}
}

func Benchmark_AppendFillSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = AppendFillSlice()
	}
}

func Benchmark_ConcurrentFillSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ConcurrentFillSlice()
	}
}
