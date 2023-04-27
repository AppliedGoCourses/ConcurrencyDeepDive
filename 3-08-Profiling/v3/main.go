package main

import (
	"fmt"
	"math/rand"
	_ "net/http/pprof"
	"os"
	"runtime/pprof"
	"sync"
)

var mu sync.Mutex
var sum int

func calculateChunkSumMutex(data []int, wg *sync.WaitGroup) {
	defer wg.Done()
	localSum := 0
	for _, num := range data {
		localSum += num
	}
	mu.Lock()
	sum += localSum
	mu.Unlock()
}


func calculateSumMutex(data []int) {

	// open mutex profile file and start profiling
	fm, err := os.Create("cpu_mutex.prof")
	if err != nil {
		panic(err)
	}
	defer fm.Close()
	pprof.StartCPUProfile(fm)

	// calculate sum of data concurrently in chunks of chunkSize
	var wg sync.WaitGroup
	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}
		wg.Add(1)
		go calculateChunkSumMutex(data[i:end], &wg)
	}
	wg.Wait()

	// stop the mutex profile
	pprof.StopCPUProfile()
	return sum
}

func calculateChunkSumChannel(data []int, result chan<- int) {
	localSum := 0
	for _, num := range data {
		localSum += num
	}
	result <- localSum
}

func calculateSumChannel() int {
	// open channel profile file and start profiling
	fc, err := os.Create("cpu_channel.prof")
	if err != nil {
		panic(err)
	}
	defer fc.Close()
	pprof.StartCPUProfile(fc)

	// calculate sum of data concurrently in chunks of chunkSize
	result := make(chan int, len(data)/chunkSize)
	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}
		go calculateChunkSumChannel(data[i:end], result)
	}

	// sum up all partial sums from the channel
	sum := 0
	for i := 0; i < len(data)/chunkSize; i++ {
		sum += <-result
	}

	pprof.StopCPUProfile()
	return sum
}

func main() {
	// create slice with random numbers
	data := make([]int, 1000000)
	for i := range data {
		data[i] = rand.Intn(100)
	}

	fmt.Println(profileMutex(data))
	fmt.Println(profileChannel(data))
}
