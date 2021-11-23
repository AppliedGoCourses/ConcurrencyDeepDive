package main

import (
	"fmt"
	"time"
)

func nilChannel(c chan int) {

	// This is our sender goroutine.
	go func(ch chan int) {
		fmt.Println("Trying to send something to a nil channel...")
		ch <- 42
		fmt.Println("Done sending to nil channel. (You should not see this message.)")
	}(c)

	fmt.Println("Trying to receive something from a nil channel...")
	n := <-c
	fmt.Printf("Received %d from nil channel. (You should not see this message.)\n", n)
}

func sendToClosedChannel(c chan<- int) {
	defer func() {
		r := recover()
		if r != nil {
			fmt.Printf("Panic: %v\n", r)
		}
	}()

	c <- 1
	c <- 2
	c <- 3
	close(c)
	// wait a bit to let the receiver read some zero values
	time.Sleep(100 * time.Microsecond) // adjust this if your computer is faster or slower than mine
	// Now try to send something through the closed channel.
	c <- 4
}

func receiveFromClosedChannel(c chan int) {
	for i := 1; i < 5; i++ {
		n := <-c
		fmt.Println("Read from closed channel:", n)
	}
}

func receiveFromClosedChannelCommaOk(c chan int) {
	for i := 1; i < 5; i++ {
		n, ok := <-c
		if !ok {
			fmt.Println("Comma,ok: the channel is closed")
			return
		}
		fmt.Println("Comma,ok: read", n)
	}
}

func main() {

	fmt.Println("A send to a nil channel blocks forever. ")
	fmt.Println("A receive from a nil channel blocks forever. ")

	var ch chan int // the zero value is nil
	go nilChannel(ch)

	fmt.Println("A send to a closed channel panics.")
	fmt.Println("A receive from a closed channel returns the zero value immediately.")

	c := make(chan int, 10)
	close(c)

	// Sorry for the unwieldy function names...
	go sendToClosedChannel(c)
	go receiveFromClosedChannel(c)
	receiveFromClosedChannelCommaOk(c)
	time.Sleep(time.Millisecond)

}
