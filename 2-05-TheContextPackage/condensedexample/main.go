package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func send(ctx context.Context, res chan<- int) {
	for {
		select {
		case <-ctx.Done():
			return
		case res <- 1:
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func receive(ctx context.Context, res <-chan int) {
	for {
		select {
		case <-ctx.Done():
			return
		case result := <-res:
			fmt.Print(result, " ")
		}
	}
}

func delayedCancel(cancel context.CancelFunc) {
	rand.Seed(time.Now().UnixNano())
	<-time.After(time.Duration(rand.Intn(1000)) * time.Millisecond)
	cancel()
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	res := make(chan int)
	go send(ctx, res)
	go receive(ctx, res)
	go delayedCancel(cancel)
	<-ctx.Done()
	fmt.Println(ctx.Err())
}
