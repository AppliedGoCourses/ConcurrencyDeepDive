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

func copyFile(task string, source, target *File) {
	// get exclusive access to both files
	log.Printf("%s: lock source %s\n", task, source.path)
	source.mu.Lock()
	// simulate time for opening the source file
	// to highly increase the chances for a deadlock
	time.Sleep(time.Millisecond)
	log.Printf("%s: lock target %s\n", task, target.path)
	target.mu.Lock()

	// simulate copying data between files
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
		copyFile("backup", orig, bck)
		done <- struct{}{}
	}()
	copyFile("restore", bck, orig)

	<-done
}

func init() {
	log.SetFlags(0)
}
