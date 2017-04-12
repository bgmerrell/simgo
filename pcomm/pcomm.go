// Package pcomm is for process communication with process functions, allowing
// them to behave as coroutines.
package pcomm

import "fmt"

// coroutineState represents the state of the coroutine
type coroutineState int

const (
	stateSuspended coroutineState = iota
	stateRunning
)

// A PCommunicator allows a process to communicate with a process function via
// channels, allowing the process function to behave as a coroutine.
type PCommunicator struct {
	// yieldCh communicates from the coroutine to the process
	yieldCh chan interface{}
	// resumeCh communicates to the courtine from the process
	resumeCh chan interface{}
	// state is the current state of the coroutine (i.e., running vs.
	// suspended).
	state coroutineState
}

// New returns a new PCommunicator with initialized channels and suspended
// state.
func New() *PCommunicator {
	return &PCommunicator{
		yieldCh:  make(chan interface{}),
		resumeCh: make(chan interface{}),
		state:    stateSuspended,
	}
}

// Yield waits for a value from a Resume() call and sets the state to
// stateRunning. Subsequent calls to Yield() communicate the provided value, x,
// such that it can be read from Resume().  Yield() communicates over
// unbuffered channels and may block accordingly.
func (pc *PCommunicator) Yield(x interface{}) interface{} {
	fmt.Println("Yielding...")
	if pc.state == stateRunning {
		pc.state = stateSuspended
		pc.yieldCh <- x
	}
	resumeVal := <-pc.resumeCh
	pc.state = stateRunning
	return resumeVal
}

// Resume communicates the provided value, x, such that it can be read from
// Yield().  Resume then receives a value from Yield().  The value received
// from Yield() is returned along with a bool indicating if the value
// was received from a valid, open channel.  Resume() communicates over
// unbuffered channels and may block accordingly.
func (pc *PCommunicator) Resume(x interface{}) (interface{}, bool) {
	fmt.Println("Resuming...")
	pc.resumeCh <- x
	yieldedVal, ok := <-pc.yieldCh
	return yieldedVal, ok
}

// Finish finalizes the underlying channels such that the process can finish.
func (pc *PCommunicator) Finish() {
	close(pc.yieldCh)
}
