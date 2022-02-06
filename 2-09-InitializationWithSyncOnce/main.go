package main

import (
	"fmt"
	"log"
	"os"
	"sync"
)

// FileBuffer keeps a file's contents in memory.
// The file is loaded only once.
type FileBuffer struct {
	data []byte
	once sync.Once
}

// BadGetFile loads a file if not already cached,
// and returns the file contents as a byte slice.
// When called concurrently, it may load the same file multiple times.
func (f *FileBuffer) BadGetFile() []byte {
	if f == nil {
		log.Fatalln("receiver must not be nil")
	}
	if f.data != nil {
		return f.data
	}
	// A very simple way to demonstrate that BadGetFile()
	// is not safe for concurrent use.
	fmt.Println("BadGetFile: loading 'data'")

	var err error
	f.data, err = os.ReadFile("data")
	// quick  & dirty error handling for brevity
	if err != nil {
		log.Fatalln(err)
		return nil
	}
	return f.data
}

// GoodGetFile uses a sync.Once object to ensure that the file is loaded
// only once.
func (f *FileBuffer) GoodGetFile() []byte {
	if f == nil {
		log.Fatalln("receiver must not be nil")
	}
	f.once.Do(func() {
		fmt.Println("GoodGetFile: loading 'data'")

		var err error
		f.data, err = os.ReadFile("data")
		if err != nil {
			log.Fatalln(err)
		}
	})
	return f.data
}

func main() {
	b := FileBuffer{}

	var wg sync.WaitGroup
	for i := 1; i < 1000; i++ {
		wg.Add(1)
		go func() {
			// Test this call with both BadGetFile() and GoodGetFile().
			data := b.BadGetFile()
			if data == nil {
				log.Println("error: file not loaded")
			}
		}()
		wg.Done()
	}
	wg.Wait()
}
