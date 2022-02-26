package main

import (
	"fmt"
	"runtime"
	"time"
)

// Fix: use a channel to tell the goroutines to stop.

// webSearch gets a second channel of type "empty struct".
// We do not need to send anything through that channel.
// Rather, when the channel is closed, it starts emitting zero values
// of type struct{}.
// The select case that waits for this channel then unblocks and stops the goroutine.
func webSearch(url, query string, c chan<- string, done chan struct{}) {
	fmt.Printf("Goroutine for %s: before select\n", url)
	select {
	case c <- fmt.Sprintf("Queried server '%s' for '%s'", url, query):
		// The fake result has been delivered. Exit the goroutine.
		fmt.Printf("Goroutine for %s: wrote the result. Exiting.\n", url)
		return
	case <-done:
		// The cancel channel started emitting zero values. Exit the goroutine.
		fmt.Printf("Goroutine for %s: read something from the 'done' channel. Exiting.\n", url)
		return
	}
}

func main() {
	fmt.Println("\n*** Fix: done channel ***")

	res := make(chan string)

	// Create a channel for sending a stop signal to the goroutinesgorot
	done := make(chan struct{})

	// Run the same query concurrently against three search servers
	servers := []string{"server1", "server2", "server3"}
	query := "appliedgo"

	for _, server := range servers {
		go webSearch(server, query, res, done)
	}

	fmt.Println(runtime.NumGoroutine(), "goroutines running (main plus three spawned goroutines)")

	// now collect the fastest result, and ONLY that.
	fmt.Println("Search result:", <-res)
	close(done)

	// give the goroutines a chance to finish
	time.Sleep(1 * time.Millisecond)

	// now let's see what remains
	fmt.Println(runtime.NumGoroutine(), "goroutine running (main)")

}
