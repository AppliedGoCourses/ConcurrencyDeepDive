package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// A contrived file type
type File struct {
	path string
	data [10]byte
	mu   sync.Mutex
}

func copyFileTryLock(task string, source, target *File) (backoffs int) {
	// get exclusive access to both files
	var lock1 bool
	for {
		if source.mu.TryLock() {
			lock1 = true

			if target.mu.TryLock() {
				break
			}
		}
		// At this point, at least one of the locks is still not acquired.
		// Unlock the source mutex if necessary, back off for a while, then try again
		if lock1 {
			source.mu.Unlock()
			lock1 = false
		}
		backoffs++
		time.Sleep(time.Millisecond)
	}

	// simulate copying data between files
	copy(target.data[:], source.data[:])
	time.Sleep(time.Millisecond)

	// release the file locks again
	target.mu.Unlock()
	source.mu.Unlock()

	return backoffs
}

func cp(ctx context.Context,
	task string,
	source, target *File,
	result chan<- string) {

	var count, backoffs int
	for {
		select {
		case <-ctx.Done():
			result <- fmt.Sprintf("%8d  %ss with %8d back-off rounds", count, task, backoffs)
			return
		default:
			backoffs += copyFileTryLock(task, source, target)
			count++
		}
	}
}

func main() {
	orig := &File{path: "original"}
	bck := &File{path: "backup"}
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	result := make(chan string)

	go cp(ctx, " backup", orig, bck, result)
	go cp(ctx, "restore", bck, orig, result)

	<-ctx.Done()
	fmt.Println(<-result)
	fmt.Println(<-result)
}

func init() {
	log.SetFlags(0)
}
