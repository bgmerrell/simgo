// Package pcomm is for process communication with process functions.
package pcomm

import "fmt"

type commState int

const (
	stateSuspended commState = iota
	stateRunning
)

type PCommunicator struct {
	// communication channel
	yieldCh  chan interface{}
	resumeCh chan interface{}
	state    commState
}

func New() *PCommunicator {
	return &PCommunicator{
		yieldCh:  make(chan interface{}),
		resumeCh: make(chan interface{}),
		state:    stateSuspended,
	}
}

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

func (pc *PCommunicator) Resume(x interface{}) (interface{}, bool) {
	fmt.Println("Resuming...")
	pc.resumeCh <- x
	yieldedVal, ok := <-pc.yieldCh
	return yieldedVal, ok
}

func (pc *PCommunicator) Finish() {
	close(pc.yieldCh)
}
