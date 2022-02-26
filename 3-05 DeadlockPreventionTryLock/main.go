package main

import (
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

func copyFileTryLock(task string, source, target *File) {
	// get exclusive access to both files
	var lock1 bool
	for {
		if source.mu.TryLock() {
			log.Printf("%s: source %s is locked\n", task, source.path)
			lock1 = true

			// simulate time for opening the source file
			// to highly increase the chances for a deadlock
			time.Sleep(time.Millisecond)

			if target.mu.TryLock() {
				log.Printf("%s: target %s is locked\n", task, target.path)
				break
			}
		}
		// At this point, at least one of the locks is still not acquired.
		// Unlock the source mutex if necessary, back off for a while, then try again
		if lock1 {
			log.Printf("%s: unlock source %s and back off\n", task, source.path)
			source.mu.Unlock()
			lock1 = false
		}
		time.Sleep(time.Millisecond)
	}

	// simulate copying data between files
	log.Printf("%s: copy data\n", task)
	copy(target.data[:], source.data[:])

	// release the file locks again
	log.Printf("%s: unlock source %s\n", task, source.path)
	target.mu.Unlock()
	log.Printf("%s: unlock target %s\n", task, target.path)
	source.mu.Unlock()
}

func main() {
	orig := &File{path: "original"}
	bck := &File{path: "backup"}
	done := make(chan struct{})

	go func() {
		copyFileTryLock("backup", orig, bck)
		done <- struct{}{}
	}()
	copyFileTryLock("restore", bck, orig)

	<-done
}

func init() {
	log.SetFlags(0)
}
