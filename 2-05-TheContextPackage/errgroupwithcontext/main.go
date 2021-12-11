package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AppliedGoCourses/ConcurrencyDeepDive/mockdb"
	"golang.org/x/sync/errgroup"
)

// checkDBstatus is intended to be run as a goroutine inside an ErrGroup,
// so it needs a result channel but no error channel. The error is
// returned by the function itself and handled by the ErrGroup.
func checkDBstatus(ctx context.Context, conn string, res chan<- string) error {

	db, err := mockdb.Open(conn)
	if err != nil {
		return fmt.Errorf("checkDBstatus: cannot open DB: %w", err)
	}
	defer db.Close()

	select {
	status, err := db.Status()
	if err != nil {
		return fmt.Errorf("checkDBstatus: cannot check status: %w", err)
	}
	res <- status
	return nil
}

func main() {

	// Create a parent context.
	pctx := context.Background()

	// Create a context with a timeout of 500 milliseconds.
	ctx, cancel := context.WithTimeout(pctx, 500*time.Millisecond)

	// We have no need for canceling running goroutines.
	// We use it here to cancel all goroutines in case the user hits Ctrl-C.
	go exitOnSignal(cancel)

	// Create an errgroup with the given context.
	// The second return value is the derived context, which
	// we do not need here.
	g, _ := errgroup.WithContext(ctx)

	conns := []string{"db1", "db2", "db3", "db4", "db5", "db6"}

	res := make(chan string, len(conns))

	for _, conn := range conns {
		c := conn // to allow the closer to grab the CURRENT value of conn
		g.Go(func() error {
			return checkDBstatus(ctx, c, res)
		})
	}

	done := make(chan struct{})
	go func() {
		for {
			r, ok := <-res
			if !ok {
				// res is closed and drained, work is done.
				// tell the main goroutine that we're done.
				close(done)
				return
			}
			fmt.Println(r)
		}
	}()

	err := g.Wait()
	if err != nil {
		log.Printf("Monitor error: %s\n", err)
	}
	close(res)
	<-done
	fmt.Println("\nDone.")
}

// exitOnSignal listens for signals and calls cancel on the context.
func exitOnSignal(cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	sig := <-c
	cancel()
	log.Fatalf("received signal %s, exiting", sig)
}
