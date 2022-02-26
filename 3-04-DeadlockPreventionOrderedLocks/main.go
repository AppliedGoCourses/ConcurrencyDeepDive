package main

import (
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

func copyFileOrderedLock(task string, source, target *File) {
	// determine the order of locking by ordering the file paths lexically
	first := source
	second := target
	if strings.Compare(source.path, target.path) > 0 {
		first = target
		second = source
	}

	log.Printf("%s: lock first %s\n", task, first.path)
	first.mu.Lock()
	time.Sleep(time.Millisecond)
	log.Printf("%s: lock second %s\n", task, second.path)
	second.mu.Lock()

	copy(target.data[:], source.data[:])

	log.Printf("%s: unlock second %s\n", task, second.path)
	second.mu.Unlock()
	log.Printf("%s: unlock first %s\n", task, first.path)
	first.mu.Unlock()
}

func main() {
	orig := &File{path: "original"}
	bck := &File{path: "backup"}
	done := make(chan struct{})

	go func() {
		copyFileOrderedLock("backup", orig, bck)
		done <- struct{}{}
	}()
	copyFileOrderedLock("restore", bck, orig)

	<-done
}

func init() {
	log.SetFlags(0)
}
