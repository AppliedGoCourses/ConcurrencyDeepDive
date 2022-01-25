package main

// Run 'go get -u -v -t github.com/gosuri/uilive' to install
// realtime terminal output

import (
	"fmt"
	"runtime"
	"time"

	"github.com/gosuri/uilive"
)

func worker() {
	for {
		// always busy! (Until activating the time.Sleep
		// call below, a blocking operation)
		time.Sleep(1 * time.Millisecond)
		// unix.Getpid() // a system call
	}
}

// showStats writes out the number of goroutines, threads and CPU's every second
func showStats() {
	term := uilive.New()
	term.RefreshInterval = 1 * time.Second
	term.Start()
	n := 0
	for {
		fmt.Fprintln(term, n, "s")
		g := runtime.NumGoroutine()
		fmt.Fprintln(term.Newline(), g, "goroutines")
		t := systemThreadCount()
		fmt.Fprintln(term.Newline(), t, "threads")
		c := runtime.NumCPU()
		fmt.Fprintln(term.Newline(), c, "CPUs")
		term.Flush()
		n++
		time.Sleep(1 * time.Second)
	}
}

func main() {

	// Display the current number of threads and goroutines every second.
	go showStats()

	// Spawn large numbers of workers
	for i := 0; i < 100000; i++ {
		go worker()
	}

	time.Sleep(10 * time.Second) // Use Ctrl-C to exit earlier
}
