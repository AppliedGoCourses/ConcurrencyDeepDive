package main

import (
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/sync/errgroup"
)

/*
This example simulates a distributed log database. Every node contains a
number of shards, and data is distributed across these shards based on the
time stamp of a log message.

When a query arrives that does not specify a time span, the query is
executed concurrently across all nodes and shards.
*/

// queryShard queries a single shard on the current node.
func queryShard(query string, shard int, shardRes chan<- string) {
	start := time.Now()

	// Simulate a query that takes between 5 and 7 seconds to search through the entire shard.
	// This timer needs to be set outside the loop, otherwise
	// the loop would create one new timer for each iteration
	// and thus never be able to end.
	queryTime := time.After(time.Duration(5+rand.Intn(2)) * time.Second)

	for {
		select {
		// Simulate a shard that takes between 500ms and 2 seconds to find the
		// next occurrence of the query data in the log database.
		//
		case <-time.After(time.Duration(rand.Intn(1500)+500) * time.Millisecond):
			shardRes <- fmt.Sprintf("queryShard: found an occurrence of %s in shard %d at index %d", query, shard, 100000+rand.Intn(899999))
		// We simulate this through time.After.
		case <-queryTime:
			fmt.Printf("queryShard: finished query '%s' on shard %d after %s\n", query, shard, time.Since(start))
			return
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	var g errgroup.Group
	const numShards = 5

	// Some random queries for log entries.
	queries := []string{"pid=5543", "HTTP_418", "CON_RST"}

	// The result channel returns the retrieved log entries.
	logs := make(chan string, 1000)

	// Save the current time to measure the duration of the query.
	start := time.Now()

	// Start a goroutine for each query.
	for _, query := range queries {
		q := query // capture loop variable for the closure
		for shard := 0; shard < numShards; shard++ {
			sh := shard
			g.Go(func() error {
				queryShard(q, sh, logs)
				return nil
			})
		}
	}

	// Set a limit for the number of results.
	limit := 1

	// Print the results.
	g.Go(func() error {
		for i := 1; i <= limit; i++ {
			fmt.Printf("Result %d: %s\n", i, <-logs)
		}
		// The desired number of results is reached, cancel the context.
		fmt.Printf("Receiving goroutine: all results received. Query completed after %s\n", time.Since(start))
		return nil
	})

	g.Wait() // error check omitted because no goroutine returns an error
	fmt.Printf("All goroutines finished after %s\n", time.Since(start))
}
