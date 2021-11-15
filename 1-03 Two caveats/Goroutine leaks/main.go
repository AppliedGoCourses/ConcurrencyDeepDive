package main

import (
	"fmt"
	"runtime"
)

type Worker struct {
	Ch chan string
}

func (w *Worker) work() {
	for {
		w.Ch <- "worker here!"
	}
}

func NewWorker() *Worker {
	w := &Worker{Ch: make(chan string, 100)}
	go w.work()
	return w
}

func goroutineleak() {
	w := NewWorker()
	fmt.Println(<-w.Ch)
	// After this point, w goes out of scope
	// but the goroutine continues to exist
	// (and so does w). Goroutine leak!
}

func main() {
	fmt.Println("\n*** Goroutine leak ***")

	fmt.Println(runtime.NumGoroutine())
	goroutineleak()
	fmt.Println(runtime.NumGoroutine())
}
