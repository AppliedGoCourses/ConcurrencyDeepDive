package main

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

// Fix: use a context to cancel the goroutines

// Similar to the channel fix but we use the context's Done channel instead.
func webSearch(ctx context.Context, url, query string, c chan<- string) {
	fmt.Printf("Goroutine for %s: before select\n", url)
	select {
	case c <- fmt.Sprintf("Queried server '%s' for '%s'", url, query):
		// The fake result has been delivered. Exit the goroutine.
		fmt.Printf("Goroutine for %s: wrote the result. Exiting.\n", url)
		return
	case <-ctx.Done():
		// The context has been cancelled. Exit the goroutine.
		fmt.Printf("Goroutine for %s: context cancelled. Exiting.\n", url)
		return
	}
}

func concurrentSearch(urls []string, query string) string {

	res := make(chan string)

	// create a context and a cancel func
	ctx, cancel := context.WithCancel(context.Background())

	// cancel any goroutine that still runs when this function exits
	// All goroutines that use the context ctx will then receive
	// a cancel signal through the ctx.Done channel.
	defer cancel()

	for _, url := range urls {
		go webSearch(ctx, url, query, res)
	}

	fmt.Println(runtime.NumGoroutine(), "goroutines running (main plus three spawned goroutines)")

	// now collect the fastest result, and ONLY that.
	result := <-res
	fmt.Println("Search result:", result)
	return result
}

func main() {
	fmt.Println("\n*** Fix: context ***")

	// Run the same query concurrently against three search servers
	servers := []string{"server1", "server2", "server3"}
	query := "appliedgo"

	concurrentSearch(servers, query)

	// give the goroutines a chance to finish
	time.Sleep(1 * time.Millisecond)

	// now let's see what remains
	fmt.Println(runtime.NumGoroutine(), "goroutine running (main)")

}
