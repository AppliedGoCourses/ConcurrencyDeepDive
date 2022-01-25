package main

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// Deliberately ignoring all errors here.
// In unlikely case the thread count is 0, add error handling
// where you see blank identifiers and inspect the error
func systemThreadCount() int {
	f, _ := os.Open("/proc/self/stat")
	b, _ := ioutil.ReadAll(f)
	stats := strings.Split(string(b), " ")
	threads, _ := strconv.Atoi(stats[19])
	// number of threads is in the 20th field of /stat output.
	// see manpage proc(5) -> entry /proc/[pid]/stat
	return threads
}
