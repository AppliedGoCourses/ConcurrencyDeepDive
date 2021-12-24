package main

import (
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/sync/errgroup"
)

type goFile struct {
	name string
	data string
}

func mockScan() string {
	<-time.After(time.Duration(rand.Intn(100)) * time.Millisecond)
	if rand.Intn(100) > 90 {
		return "ALERT - vulnerability found"
	}
	return "OK"
}

func scanSQLInjection(data <-chan goFile, res chan<- string) error {
	for d := range data {
		res <- fmt.Sprintf("SQL injection scan: %s scanned, result: %s", d.name, mockScan())
	}
	close(res)
	return nil
}

func scanTimingExploits(data <-chan goFile, res chan<- string) error {
	for d := range data {
		res <- fmt.Sprintf("Timing exploits scan: %s scanned, result: %s", d.name, mockScan())
	}
	close(res)
	return nil
}

func scanAuth(data <-chan goFile, res chan<- string) error {
	for d := range data {
		res <- fmt.Sprintf("Authentication scan: %s scanned, result: %s", d.name, mockScan())
	}
	close(res)
	return nil
}

// fanOut takes a channel of type T, the number of channels to create,
// and the size of each channel. It returns a slice of channels of type T,
// and starts a goroutine that copies each item from the input channel to
// each of the output channels.
// When the input channel is closed, the output channels are closed as well,
// and the goroutine ends.
func fanOut[T any](ch chan T, n, cap int) []chan T {
	chans := make([]chan T, 0, n)
	for i := 0; i < n; i++ {
		chans = append(chans, make(chan T, cap))
	}

	go func() {
		// read the input channel
		for item := range ch {
			// and copy each item to each of the output channels
			for _, c := range chans {
				select {
				case c <- item:
					// continue the loop
				case <-time.After(90 * time.Millisecond):
					fmt.Println("Timeout")
					// continue the loop, current item is lost
				}
			}
		}
		// close all the output channels to signal that we're done
		for _, c := range chans {
			close(c)
		}
	}()

	return chans
}

// fanIn takes a list of channels and spawns a goroutine for each of them
// to listen for input and pass it on to the output channel.
// When all channels are closed, the output channel is closed as well
func fanIn[T any](chans ...chan T) chan T {
	res := make(chan T)
	var g errgroup.Group

	for _, c := range chans {
		// for each channel, start a goroutine that copies data from the channel
		// to the output channel
		c := c
		g.Go(func() error {
			for s := range c {
				res <- s
			}
			return nil
		})
	}

	// start a goroutine that closes the output channel when all the goroutines are done,
	// to signal that work is done
	go func() {
		g.Wait()
		close(res)
	}()

	return res
}

func main() {

	si := []goFile{
		{name: "main.go", data: "package main\n\nfunc main() {\n\tprintln(\"Hello, world!\")\n}"},
		{name: "utils.go", data: "package utils\n\nfunc Util() {\n\tprintln(\"Hello, world!\")\n}"},
		{name: "helper.go", data: "package helper\n\nfunc Help() {\n\tprintln(\"Hello, world!\")\n}"},
		{name: "misc.go", data: "package misc\n\nfunc Misc() {\n\tprintln(\"Hello, world!\")\n}"},
		{name: "various.go", data: "package various\n\nfunc Various() {\n\tprintln(\"Hello, world!\")\n}"},
	}

	input := make(chan goFile, len(si))
	res1 := make(chan string, len(si))
	res2 := make(chan string, len(si))
	res3 := make(chan string, len(si))

	// Create a fanout, get the list of channels back
	chans := fanOut(input, 3, len(si))

	// Start the goroutines and connect each one to one of the fanout channels.
	var g errgroup.Group
	g.Go(func() error {
		return scanSQLInjection(chans[0], res1)
	})
	g.Go(func() error {
		return scanTimingExploits(chans[1], res2)
	})
	g.Go(func() error {
		return scanAuth(chans[2], res3)
	})

	// send the data to the fanout
	g.Go(func() error {
		for _, d := range si {
			input <- d
		}
		close(input)
		return nil
	})

	// collect the results
	g.Go(func() error {
		// collect the results via fan-in
		res := fanIn(res1, res2, res3)
		for r := range res {
			fmt.Println(r)
		}
		return nil
	})

	err := g.Wait()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("main: done")

}

func init() {
	rand.Seed(time.Now().UnixNano())
}
