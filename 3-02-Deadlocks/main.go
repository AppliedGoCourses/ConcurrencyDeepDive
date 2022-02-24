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

func copyFile(task string, source, target *File) {
	// get exclusive access to both files
	log.Printf("%s: lock source %s\n", task, source.path)
	source.mu.Lock()
	// simulate time for opening the source file
	// to highly increase the chances for a deadlock
	<-time.After(time.Millisecond * 100)
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
	<-time.After(time.Millisecond * 100)
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
	go copyFile("backup", orig, bck)
	copyFile("restore", bck, orig)
	// go copyFileOrderedLock("backup", orig, bck)
	// copyFileOrderedLock("restore", bck, orig)
}

func init() {
	log.SetFlags(0)
}
