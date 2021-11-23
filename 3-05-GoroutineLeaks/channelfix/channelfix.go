package main

import (
	"fmt"
	"runtime"
)

// Fix: use a channel to tell the goroutines to stop.

// webSearch gets a second channel of type "empty struct".
// We do not need to send anything through that channel.
// Rather, when the channel is closed, it starts emitting zero values
// of type struct{}.
// The select case that waits for this channel then unblocks and stops the goroutine.
func webSearch(url, query string, c chan<- string, done chan struct{}) {
	select {
	case c <- fmt.Sprintf("Queried server '%s' for '%s'", url, query):
		// The fake result has been delivered. Exit the goroutine.
		return
	case <-done:
		// The cancel channel started emitting zero values. Exit the goroutine.
		return
	}

}

func spawnThreeSearches(query string, res chan string, done chan struct{}) {
	servers := []string{"server1", "server2", "server3"}
	for _, server := range servers {
		go channelFixWebSearch(server, query, res, done)
	}
}

func run() {
	fmt.Println("\n*** Fix: buffered channel ***")

	// create a BUFFERED result channel
	// large enough to hold all results
	res := make(chan string)
	done := make(chan struct{})

	// Run the same query concurrently against three search servers
	spawnThreeSearches("appliedgo", res, done)

	fmt.Println(runtime.NumGoroutine(), "goroutines running (main plus three spawned goroutines)")

	// now collect the fastest result, and ONLY that.
	fmt.Println("Result:", <-res)
	close(done)
}
