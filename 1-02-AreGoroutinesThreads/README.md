# Show goroutines and threads

This code spawns a large number of goroutines and displays the total number of goroutines, threads, and CPU cores (or hardware threads). 

Reading threads is OS dependent and currently only implemented for Linux. 

Hence you either need to `go build` the code, or run `go run main.go threads_linux.go`. 

On MacOS and Windows, the code compiles also (using the respective threads_darwin.go or threads_windows.go file) but the thread count remains at zero.