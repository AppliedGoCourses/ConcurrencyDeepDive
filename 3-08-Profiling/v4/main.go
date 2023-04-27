package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"runtime/pprof"
)

func isPrime(n int) bool {
	if n <= 1 {
		return false
	}
	if n == 2 || n == 3 {
		return true
	}
	root := int(math.Sqrt(float64(n)))
	for i := 2; i <= root; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func primeCheckerGoroutinePerNumber(start, end int, results chan<- int) {
	for i := start; i <= end; i++ {
		go func(n int) {
			if isPrime(n) {
				results <- n
			}
		}(i)
	}
}

func primeCheckerWorkerPool(start, end, workers int, results chan<- int) {
	jobs := make(chan int,workers)

	for w := 0; w < workers; w++ {
		go func() {
			for n := range jobs {
				if isPrime(n) {
					results <- n
				}
			}
		}()
	}

	for i := start; i <= end; i++ {
		jobs <- i
	}

	close(jobs)
}

func main() {
	profile := flag.Bool("profile", false, "Enable CPU profiling")
	method := flag.String("method", "goroutine", "Choose implementation: goroutine or workerpool")
	start := flag.Int("start", 1, "Start of the range")
	end := flag.Int("end", 10000, "End of the range")
	workers := flag.Int("workers", 4, "Number of workers for worker pool implementation")
	flag.Parse()

	if *profile {
		f, err := os.Create("cpuprofile")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	results := make(chan int)
	done := make(chan bool)
	go func() {
		for range results {
		}
		done <- true
}()

	switch *method {
	case "goroutine":
		primeCheckerGoroutinePerNumber(*start, *end, results)
	case "workerpool":
		primeCheckerWorkerPool(*start, *end, *workers, results)
	default:
		fmt.Fprintln(os.Stderr, "Invalid method. Choose 'goroutine' or 'workerpool'.")
		return
	}

	for i := *start; i <= *end; i++ {
		<-results
	}

	close(results)
	<-done
}

