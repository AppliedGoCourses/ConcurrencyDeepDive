package main

import (
	"fmt"
	"runtime"
)

// Later, we will redefine this func for a fix, hence we assign it to a variable.
func webSearch(url, query string, c chan<- string) {
	// Do a fake search and pass the "result" to the channel
	c <- fmt.Sprintf("Queried server '%s' for '%s'", url, query) // instead of a real search
}

func spawnThreeSearches(query string, res chan string) {
	servers := []string{"server1", "server2", "server3"}
	for _, server := range servers {
		go webSearch(server, query, res)
	}
}

func run() {
	fmt.Println("\n*** Fix: buffered channel ***")

	// create a BUFFERED result channel
	// large enough to hold all results
	res := make(chan string, 3)

	// Run the same query concurrently against three search servers
	spawnThreeSearches("appliedgo", res)

	fmt.Println(runtime.NumGoroutine(), "goroutines running (main plus three spawned goroutines)")

	// now collect the fastest result, and ONLY that.
	fmt.Println("Result:", <-res)
}
