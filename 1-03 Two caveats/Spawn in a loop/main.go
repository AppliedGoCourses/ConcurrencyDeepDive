package main

import (
	"fmt"
	"runtime"
	"time"
)

func spawnInALoop() {
	for i := 0; i < 10; i++ {
		// Our goroutine is a closure, and closures can
		// access variables in the outer scope, so we can
		// grab i here to give each goroutine an ID, right?
		// (Hint: wrong.)
		go func() {
			fmt.Println("Goroutine", i)
		}() // <- Don't forget to call () the closure
	}
	time.Sleep(100 * time.Millisecond)
}

func spawnInALoopFixed() {
	for i := 0; i < 10; i++ {
		// Always pass any start value properly as a
		// function parameter. This way, nothing can go wrong.
		go func(n int) {
			fmt.Println("Goroutine", n)
		}(i) // We pass i here as an argument
	}
	time.Sleep(100 * time.Millisecond)
}

func main() {
	fmt.Println("*** Spawn in a loop ***")
	spawnInALoop()

	fmt.Println("\n*** Fixed ***")
	spawnInALoopFixed()

}
