package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
)

type application struct {
	port      string
	pprofPort string
}

func fillSlice() []int {
	s := make([]int, 1000000)
	for i := range s {
		s[i] = i
	}
	return s
}

func fillSliceAppend() []int {
	var s []int
	for i := 0; i < 1000000; i++ {
		s = append(s, i)
	}
	return s
}

func fillHandler(w http.ResponseWriter, r *http.Request) {
	s := fillSlice()
	fmt.Fprintf(w, "Random element: %d", s[rand.Intn(len(s))])
}

func appendHandler(w http.ResponseWriter, r *http.Request) {
	s := fillSliceAppend()
	fmt.Fprintf(w, "Random element: %d", s[rand.Intn(len(s))])
}

// runServer sets up routing and starts an HTTP server.
func runServer(app *application) {
	// Do not use the default router. The pprof package automatically adds endpoints to the default router. This will leak the pprof endpoints in the main server whose port might be exposed to the public. Instead, create a separate router and server for pprof.
	mux := http.NewServeMux()
	mux.HandleFunc("/fill", fillHandler)
	mux.HandleFunc("/append", appendHandler)
	go func() {
		println("map server started on port,", app.port)
		log.Fatal(http.ListenAndServe("localhost:"+app.port, mux))
	}()
}

// runPprofServer adds endpoints for pprof and starts a seprate HTTP server listening on port 6060.
func runPprofServer(app *application) {
	mux := http.NewServeMux()
	// Add the pprof endpoint. The handler "pprof.Index" serves the endpoints for /debug/pprof/cmdline, /debug/pprof/profile, /debug/pprof/symbol, and /debug/pprof/trace implicitly.
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	go func() {
		println("pprof server started on port", app.pprofPort)
		log.Fatal(http.ListenAndServe("localhost:"+app.pprofPort, mux))
	}()
}

func main() {
	app := application{
		port:      "8080",
		pprofPort: "6060",
	}

	runServer(&app)
	runPprofServer(&app)

	waitForSigint()
	os.Exit(0)
}

// waitForSigint waits for the SIGINT signal (usually sent by hitting Ctrl-C on the terminal) and then returns.
func waitForSigint() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	log.Println("signal:", <-c)
}
