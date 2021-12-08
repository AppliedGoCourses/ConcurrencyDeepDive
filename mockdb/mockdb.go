// Package mockdb implements a fully non-functional mock database client.
package mockdb

import (
	"fmt"
	"math/rand"
	"time"
)

// MockDB represents a mock database.
type MockDB struct {
	name string
}

// Close closes the connection to the MockDB server m.
func (m *MockDB) Close() {}

// Status returns the current status of the MockDB server m.
// The check may fail with a certain probability, in which case Status
// returns an error.
func (m *MockDB) Status() (string, error) {

	if rand.Float64() > 0.8 {
		return "", fmt.Errorf("Error checking status of %s", m.name)
	}

	states := []string{"starting", "running", "sleeping", "blocked", "stopping", "stopped"}
	return fmt.Sprintf("Server %s: %s", m.name, states[rand.Intn(len(states))]), nil
}

// Open opens a connection to a MockDB server, defined by connection string "conn".
// If no connection can be established, Open returns an error.
// The call may get delayed by up to one second, to support
// demonstrating timeout behavior.
func Open(conn string) (*MockDB, error) {

	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

	if rand.Float64() > 0.8 {
		return nil, fmt.Errorf("Error connecting to %s", conn)
	}
	return &MockDB{name: conn}, nil
}

// The use of init() is generally frowned upon. If an executable contains
// multiple init functions (e.g. by importing multiple libraries that contain
// init functions), the sequence of invoking them is undefined,
// hence they can be a source of subtle, hard-to-replicate bugs.
// Here, init() is used for seeding the random number generator.
// This is one of the few legitimate uses of init().
func init() {
	rand.Seed(time.Now().UnixNano())
}
