package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func worker(n int, wg *sync.WaitGroup) {

	defer wg.Done()

	for i := 0; i < rand.Intn(10)+10; i++ {
		fmt.Print(n)
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go worker(i, &wg)
	}

	wg.Wait()
	fmt.Println("\nDone.")
}
