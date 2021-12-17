package main

/* Trivial "fan-out" and fan-in.

This is standard go channel usage.

Fan-out:

Pass a single channel to multiple goroutines and
start sending data to the channel.

Resulting behavior:
- All goroutines can pick some work from the channel
- Every goroutine works on unique work items; there is no work duplication.

Fan-in:

Pass a single result channel to multiple goroutines and
have them send their results to the channel.

The receiving side can then read the channel in a loop.
*/

import "fmt"

// Data to work with.
type data int

// worker receives work from wch and sends the result to res.
func worker(wch <-chan data, res chan<- data) {
	for {
		// Receive a work item.
		w, ok := <-wch
		// If the channel is closed, we're done.
		if !ok {
			return
		}
		// Process the work item.
		w *= 2
		// Send the result back.
		res <- w
	}
}

func main() {
	// Data we'll process.
	work := []data{1, 2, 3, 4, 5}

	// Number of workers.
	const numWorkers = 3

	// Create input and result channels
	wch := make(chan data, len(work))
	res := make(chan data, len(work))

	// Start the workers.
	for i := 0; i < numWorkers; i++ {
		go worker(wch, res)
	}

	// Distribute the work.
	for _, w := range work {
		wch <- w
	}
	close(wch)

	// Receive the results.
	// The number of results is known, hence there is no need for a wait group.
	for range work {
		w := <-res
		fmt.Println(w)
	}
}
