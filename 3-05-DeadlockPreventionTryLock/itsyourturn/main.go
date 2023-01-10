package main

import (
	"context"
	"flag"
	"log"
	"strings"
	"sync"
	"time"
)

// A contrived file type
type File struct {
	path string
	data [10]byte
	mu   sync.Mutex
}

// copyFile uses a backoff strategy to avoid a deadlock when acquiring both locks
func copyFile(task string, source, target *File, withBackoff, withOrderedLock bool) (backoffs int) {
	// Try a backoff-retry loop to acquire the lock
	var lock1 bool
	backoffs = 0

	first := source
	second := target

	if withOrderedLock {
		if strings.Compare(source.path, target.path) > 0 {
			first = target
			second = source
		}
	}

	if !withBackoff {
		first.mu.Lock()
		// This time, I add no artificial delay here, as I want to
		// collect some stats, rather than provoking a deadlock
		time.Sleep(time.Nanosecond)
		second.mu.Lock()
	} else {
		for {
			if first.mu.TryLock() {
				lock1 = true

				time.Sleep(time.Nanosecond)
				if second.mu.TryLock() {
					break
				}
			}
			// At this point, at least one of the locks is still not acquired.
			// Unlock the first mutex if necessary, back off for a while, then try again
			if lock1 {
				first.mu.Unlock()
				lock1 = false
			}
			backoffs++
			time.Sleep(time.Millisecond)
		}
	}

	copy(target.data[:], source.data[:])
	time.Sleep(time.Microsecond)

	first.mu.Unlock()
	second.mu.Unlock()

	return backoffs
}

func main() {
	orig := &File{path: "original"}
	bck := &File{path: "backup"}
	c, _ := context.WithTimeout(context.Background(), time.Second)
	var backups, restores, backoffB, backoffR int

	var withBackoff bool
	flag.BoolVar(&withBackoff, "backoff", false, "use backoff strategy")
	var withOrderedLock bool
	flag.BoolVar(&withOrderedLock, "ordered", false, "use ordered lock strategy")
	flag.Parse()

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				backoffB += copyFile("backup", orig, bck, withBackoff, withOrderedLock)
				backups++
			}
		}
	}(c)

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				backoffR += copyFile("restore", bck, orig, withBackoff, withOrderedLock)
				restores++
			}
		}
	}(c)

	<-c.Done()
	log.Printf("%8d  backups with %8d backoff rounds\n%8d restores with %8d backoff rounds", backups, backoffB, restores, backoffR)
}

func init() {
	log.SetFlags(0)
}
