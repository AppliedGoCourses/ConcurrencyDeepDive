package main

import (
	"fmt"
	"math/rand"
	"time"
)

/* True Fan-out

Goal: send the same data to multiple goroutines.

Use case: a source code security scanner runs different checks on the same source code.

For this scenario, we need a distributor that distributes the same data to multiple channels.

*/

type scanItem struct {
	name string
	data string
}

// All the "real" security scanning work is done here.
// Sometimes, the scan finds something, sometimes it doesn't.
func mockScan() string {
	// Simulate intense security scanning
	<-time.After(100 * time.Millisecond)
	// Simulate a scan result
	if rand.Intn(100) > 90 {
		return "ALERT - vulnerability found"
	}
	return "OK - Oll Korrect" // other sources claim that "OK" comes from the Scots phrase "och aye" (oh yes), or from Greek "όλα καλά" (óla kalá), meaning "all good".
}

// Sample worker: SQL injections
func scanSQLInjection(data <-chan scanItem) {
	for d := range data {
		// scan the data and print the result
		fmt.Printf("SQL injection scan: %s scanned, result: %s\n", d.name, mockScan())
	}
}

// Sample worker find buffer overflows
func scanBufferOverflow(data <-chan scanItem) {
	for d := range data {
		fmt.Printf("Buffer overflow scan: %s scanned, result: %s\n", d.name, mockScan())
	}
}

// sample worker: find flaws in authentication code
func scanAuth(data <-chan scanItem) {
	for d := range data {
		fmt.Printf("Timing attack scan: %s scanned, result: %s\n", d.name, mockScan())
	}
}

type FanOut[T any] struct {
	chans []chan T
}

// Register receives a function that receives a channel.
// It creates a new channel, adds the channel to the list of channels,
// and starts the function as a goroutine, passing the created channel to it.
func (f *FanOut[T]) Register(fn func(<-chan T)) {
	ch := make(chan T, 1)
	f.chans = append(f.chans, ch)
	go fn(ch)
}

func NewFanOut[T any](data <-chan T) *FanOut[T] {
	f := &FanOut[T]{}

	go func() {
		// Close all outgoing channels when this function exits
		defer func() {
			fmt.Println("FanOut: closing all channels")
			for _, ch := range f.chans {
				close(ch)
			}
		}()

		// Read from the incoming channel and send to all outgoing channels.
		// If the incoming channel is closed, close the outgoing channels .
		for d := range data {
			for _, ch := range f.chans {
				ch <- d
			}
		}
		fmt.Println("FanOut: done")
	}()

	return f
}

// func broadcast[T any](ch <-chan T, n int) []<-chan T {
// 	chans := make([]chan T, 0, n)
// 	for i := 0; i < n; i++ {
// 		chans = append(chans, make(chan T))
// 	}

// 	 distribute := func(ch <-chan int, chans []chan<- int) {
// 		// Close every channel when the execution ends.
// 		defer func(chans []chan<- int) {
// 			for _, c := range chans {
// 				close(c)
// 			}
// 		}(cs)

// 		for {
// 			for _, c := range chans {
// 				select {
// 				case val, ok := <-ch:
// 					if !ok {
// 						return
// 					}

// 					c <- val
// 				}
// 			}
// 		}
// 	}

// 	go distribute(ch, chans)

// 	return chans
// }

func main() {

	ch := make(chan scanItem)
	fan := NewFanOut[scanItem](ch)
	fan.Register(scanSQLInjection)
	fan.Register(scanBufferOverflow)
	fan.Register(scanAuth)

	si := []scanItem{
		{name: "main.go", data: "package main\n\nfunc main() {\n\tprintln(\"Hello, world!\")\n}"},
		{name: "utils.go", data: "package utils\n\nfunc Util() {\n\tprintln(\"Hello, world!\")\n}"},
		{name: "helper.go", data: "package helper\n\nfunc Help() {\n\tprintln(\"Hello, world!\")\n}"},
		{name: "misc.go", data: "package misc\n\nfunc Misc() {\n\tprintln(\"Hello, world!\")\n}"},
		{name: "various.go", data: "package various\n\nfunc Various() {\n\tprintln(\"Hello, world!\")\n}"},
	}

	for _, d := range si {
		ch <- d
	}
	fmt.Println("main: done")

}

func init() {
	rand.Seed(time.Now().UnixNano())
}
