package main

import (
	"fmt"
	"log"
	"os"
	"runtime/pprof"
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

func main() {
	cpuf, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal("could not create CPU file:", err)
	}
	defer cpuf.Close()

	err = pprof.StartCPUProfile(cpuf)
	if err != nil {
		log.Fatal("could not start CPU profile:", err)
	}

	m := MakeFillSlice()
	s := AppendFillSlice()

	profs := pprof.Profiles()
	for _, prof := range profs {
		fmt.Printf("%s: %d\n", prof.Name(), prof.Count())
	}

	pprof.StopCPUProfile()

	allocf, err := os.Create("allocs.prof")
	if err != nil {
		log.Fatal("could not create allocs file:", err)
	}
	defer allocf.Close()
	err = pprof.Lookup("allocs").WriteTo(allocf, 0)
	if err != nil {
		log.Fatal("could not write allocs profile:", err)
	}

	heapf, err := os.Create("heap.prof")
	if err != nil {
		log.Fatal("could not create heap file:", err)
	}
	defer heapf.Close()
	err = pprof.Lookup("heap").WriteTo(heapf, 0)
	if err != nil {
		log.Fatal("could not write heap profile:", err)
	}

	fmt.Println(m[0], s[0]) // to use the variables and avoid optimization
}
