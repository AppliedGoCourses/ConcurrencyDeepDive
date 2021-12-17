package main

import (
	"fmt"
	"math/rand"
	"time"
)

// A security scanner - serial version

// Each work item is a whole file to scan.
type goFile struct {
	name    string
	content string
}

// All the "real" security scanning work is done here.
// Sometimes, the scan finds something, sometimes it doesn't.
func mockScan() string {
	// Simulate a scan result
	if rand.Intn(100) > 90 {
		return "ALERT - vulnerability found"
	}
	return "OK - Oll Korrect"
	// other sources claim that "OK" comes from the
	// Scots phrase "och aye" (oh yes), or from
	// Greek "όλα καλά" (óla kalá), meaning "all good".
}

func scanSQLInjection(data goFile) string {
	// scan the data and print the result
	return fmt.Sprintf("SQL injection scan: %s scanned, result: %s", data.name, mockScan())
}

func scanTimingExploits(data goFile) string {
	return fmt.Sprintf("Timing exploits scan: %s scanned, result: %s", data.name, mockScan())
}

func scanAuth(data goFile) string {
	return fmt.Sprintf("Authentication scan: %s scanned, result: %s", data.name, mockScan())
}

func main() {

	si := []goFile{
		// featuring the worst package names ever
		{name: "utils.go", content: "package utils\n\nfunc Util() {}"},
		{name: "helper.go", content: "package helper\n\nfunc Help() {}"},
		{name: "misc.go", content: "package misc\n\nfunc Misc() {}"},
		{name: "various.go", content: "package various\n\nfunc Various() {}"},
	}

	for _, d := range si {
		fmt.Println(scanSQLInjection(d))
		fmt.Println(scanTimingExploits(d))
		fmt.Println(scanAuth(d))
	}

	fmt.Println("main: done")

}

func init() {
	rand.Seed(time.Now().UnixNano())
}
