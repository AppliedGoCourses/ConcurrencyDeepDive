package main

import (
	"fmt"
	"runtime"
	"time"
)

// This function shall run concurrently, with different URLs.
func webSearch(url, query string, c chan<- string) {
	// Do a fake search and pass the "result" to the channel
	fmt.Printf("Goroutine: trying to write result from %s to the channel.\n", url)
	c <- fmt.Sprintf("Queried server '%s' for '%s'", url, query) // instead of a real search
	fmt.Printf("Goroutine: wrote result from %s to the channel.\n", url)
}

func main() {
	fmt.Println("\n*** Goroutine leak ***")

	// create a result channel
	res := make(chan string)

	// Run the same query concurrently against three search servers
	servers := []string{"server1", "server2", "server3"}
	query := "appliedgo"

	for _, server := range servers {
		go webSearch(server, query, res)
	}

	fmt.Println(runtime.NumGoroutine(), "goroutines running (main plus three spawned goroutines)")

	// now collect the fastest result, and ONLY that.
	fmt.Println("Read the first search result:", <-res)

	// give the goroutines a chance to finish if they can
	time.Sleep(100 * time.Millisecond)
	// I am intentionally NOT using a WaitGroup here, because that would result in a deadlock.
	// Rather, I want to show the goroutine count without distractions.

	// now let's see what remains
	fmt.Println(runtime.NumGoroutine(), "goroutines running (main plus two blocked goroutines)")
}
