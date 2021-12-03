package main

import (
	"fmt"
	"log"

	"github.com/AppliedGoCourses/ConcurrencyDeepDive/mockdb"
	"golang.org/x/sync/errgroup"
)

// checkDBstatus is intended to be run as a goroutine inside an ErrGroup,
// so it needs a result channel but no error channel. The error is
// returned by the function itself and handled by the ErrGroup.
func checkDBstatus(conn string, res chan<- string) error {

	db, err := mockdb.Open(conn)
	if err != nil {
		return fmt.Errorf("checkDBstatus: cannot open DB: %w", err)
	}
	defer db.Close()

	status, err := db.Status()
	if err != nil {
		return fmt.Errorf("checkDBstatus: cannot check status: %w", err)
	}
	res <- status
	return nil
}

func main() {

	var g errgroup.Group

	conns := []string{"db1", "db2", "db3", "db4", "db5", "db6"}

	res := make(chan string)

	for _, conn := range conns {
		c := conn // to allow the closer to grab the CURRENT value of conn
		g.Go(func() error {
			return checkDBstatus(c, res)
		})
	}

	done := make(chan struct{})
	go func() {
		for {
			r, ok := <-res
			if !ok {
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
