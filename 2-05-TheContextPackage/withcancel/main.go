package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/sync/errgroup"
)

func queryShard(ctx context.Context, query string, shard int, shardRes chan<- string) {
	start := time.Now()

	queryTime := time.After(time.Duration(5+rand.Intn(2)) * time.Second)

	for {
		select {
		// If the context is canceled, the channel returned by ctx.Done() is closed
		// and starts emitting the zero value.
		case <-ctx.Done():
			fmt.Printf("queryShard: %s on shard %d after %s\n", ctx.Err(), shard, time.Since(start))
			// Do any cleanup here. Close files and connections etc.
			return
		case <-time.After(time.Duration(rand.Intn(1500)+500) * time.Millisecond):
			shardRes <- fmt.Sprintf("queryShard: found an occurrence of %s in shard %d at index %d", query, shard, 100000+rand.Intn(899999))
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

	// Create a background context.
	bgctx := context.Background()

	// Derive a child context that has the cancellation feature.
	// The second return value is a function that can be called
	// to cancel the context.
	ctx, cancel := context.WithCancel(bgctx)

	queries := []string{"pid=5543", "HTTP_418", "CON_RST"}
	logs := make(chan string)
	start := time.Now()

	for _, query := range queries {
		q := query
		for shard := 0; shard < numShards; shard++ {
			sh := shard
			g.Go(func() error {
				queryShard(ctx, q, sh, logs)
				return nil
			})
		}
	}

	limit := 1

	g.Go(func() error {
		for i := 1; i <= limit; i++ {
			fmt.Printf("Result %d: %s\n", i, <-logs)
		}
		fmt.Printf("Receiving goroutine: all results received. Query completed after %s\n", time.Since(start))
		// At this point, calling cancel() will cause the context to be canceled.
		// All goroutines started by g will receive the cancellation signal through
		// the channel that ctx.Done() returns. See queryShard() for the other end
		// of the cancel mechanism.
		cancel()
		return nil
	})

	g.Wait()
	fmt.Printf("All goroutines finished after %s\n", time.Since(start))
}
