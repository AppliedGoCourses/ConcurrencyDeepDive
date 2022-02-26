package main

import (
	"fmt"
	"runtime"
	"time"
)

// Fix: use a buffered channel to collect results

func webSearch(url, query string, c chan<- string) {
	// Do a fake search and pass the "result" to the channel
	fmt.Printf("Goroutine: trying to write result from %s to the channel.\n", url)
	c <- fmt.Sprintf("Queried server '%s' for '%s'", url, query) // instead of a real search
	fmt.Printf("Goroutine: wrote result from %s to the channel.\n", url)
}

func main() {
	fmt.Println("\n*** Fix: buffered channel ***")

	// create a BUFFERED result channel
	// large enough to hold all results
	res := make(chan string, 3)

	// Run the same query concurrently against three search servers
	servers := []string{"server1", "server2", "server3"}
	query := "appliedgo"

	for _, server := range servers {
		go webSearch(server, query, res)
	}

	fmt.Println(runtime.NumGoroutine(), "goroutines running (main plus three spawned goroutines)")

	// now collect the fastest result, and ONLY that.
	fmt.Println("Search result:", <-res)

	// give the goroutines a chance to finish
	time.Sleep(1 * time.Millisecond)

	// now let's see what remains
	fmt.Println(runtime.NumGoroutine(), "goroutine running (main)")
}
