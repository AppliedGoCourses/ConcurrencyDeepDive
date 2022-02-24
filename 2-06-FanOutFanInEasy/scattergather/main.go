package main

import (
	"fmt"
	"math/rand"
)

// A security scanner - Scatter-Gather version

// Each work item is a whole file to scan.
type goFile struct {
	name    string
	content string
}

// All the "real" security scanning work is done here.
// Sometimes, the scan finds something, sometimes it doesn't.
func mockScan() string {
	if rand.Intn(100) > 90 {
		return "ALERT - vulnerability found"
	}
	return "OK - Oll Korrect"
}

func scanSQLInjection(data goFile, res chan<- string) {
	res <- fmt.Sprintf("SQL injection scan: %s scanned, result: %s", data.name, mockScan())
}

func scanTimingExploits(data goFile, res chan<- string) {
	res <- fmt.Sprintf("Timing exploits scan: %s scanned, result: %s", data.name, mockScan())
}

func scanAuth(data goFile, res chan<- string) {
	res <- fmt.Sprintf("Authentication scan: %s scanned, result: %s", data.name, mockScan())
}

func main() {

	si := []goFile{
		// featuring the worst package names ever
		{name: "utils.go", content: "package utils\n\nfunc Util() {}"},
		{name: "helper.go", content: "package helper\n\nfunc Help() {}"},
		{name: "misc.go", content: "package misc\n\nfunc Misc() {}"},
		{name: "various.go", content: "package various\n\nfunc Various() {}"},
	}

	// The result channel is large enough to hold all results from all workers
	res := make(chan string, len(si)*3)

	// Scatter: send the same data to multiple goroutines
	for _, d := range si {
		d := d
		go scanSQLInjection(d, res)
		go scanTimingExploits(d, res)
		go scanAuth(d, res)
	}

	// Gather: read all results
	for i := 0; i < cap(res); i++ {
		fmt.Println(<-res)
	}
	fmt.Println("main: done")

}
