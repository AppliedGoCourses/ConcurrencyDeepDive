package fillslice

import (
	"context"
	"runtime"
	"runtime/pprof"
	"strconv"
)

func MakeFillSlice() []int {
	s := make([]int, 1000000)
	for i := range s {
		s[i] = i
	}
	return s
}

func AppendFillSlice() []int {
	var s []int
	for i := 0; i < 1000000; i++ {
		s = append(s, i)
	}
	return s
}

func ConcurrentFillSlice() []int {
	cpu := runtime.NumCPU()
	size := 1000000
	s := make([]int, cpu*size)
	// concurrently iterate over disjoint parts of the slice
	for i := 0; i < cpu; i++ {
		i := i //
		pprof.Do(context.Background(), pprof.Labels("part", strconv.Itoa(i)), func(context.Context) {
			go func(part int) {
				for j := 0; j < size; j++ {
					s[part*size+j] = part*size + j
				}
			}(i)
		})
	}
	return s
}
