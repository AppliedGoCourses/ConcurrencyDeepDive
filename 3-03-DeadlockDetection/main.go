package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"sync"
	"time"
)

type File struct {
	path string
	data [10]byte
	mu   sync.Mutex
}

func copyFile(task string, source, target *File) {
	source.mu.Lock()
	<-time.After(time.Millisecond * 100) // provoke a deadlock situation
	target.mu.Lock()

	copy(target.data[:], source.data[:])

	target.mu.Unlock()
	source.mu.Unlock()
}

func main() {
	go func() {
		fmt.Println(http.ListenAndServe("localhost:7070", nil))
	}()
	fmt.Println("Run\ncurl \"http://localhost:7070/debug/pprof/goroutine?debug=2\"\nto get a stack dump")

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
