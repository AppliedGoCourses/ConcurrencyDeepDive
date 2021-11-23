package main

import (
	"fmt"
	"time"
)

func spawnInALoop(urls []string) {
	for _, url := range urls {
		// Our goroutine is a closure, and closures can
		// access variables in the outer scope, so we can
		// grab the URL here to give each goroutine an ID, right?
		// (Hint: wrong.)
		go func() {
			fmt.Println("Fetching", url)
		}() // <- Don't forget to call () the closure
	}
}

func spawnInALoopFixed(urls []string) {
	for _, url := range urls {
		// Always pass any start value properly as a
		// function parameter. This way, nothing can go wrong.
		go func(url string) {
			fmt.Println("Fetching", url)
		}(url) // We pass the URL here as an argument
	}
	time.Sleep(100 * time.Millisecond)
}

func main() {
	fmt.Println("*** Spawn in a loop ***")

	urls := []string{"https://appliedgo.com", "https://appliedgo.net", "https://golang.org", "https://go.dev"}
	spawnInALoop(urls)

	time.Sleep(100 * time.Millisecond)
	fmt.Println("\n*** Fixed ***")
	spawnInALoopFixed(urls)
}
