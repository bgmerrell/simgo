// Process communication
package pcomm

import "fmt"

type PCommunicator struct {
	// communication channel
	Ch chan interface{}
}

func New() *PCommunicator {
	return &PCommunicator{
		Ch: make(chan interface{}),
	}
}

func (pc *PCommunicator) Recv() interface{} {
	return <-pc.Ch
}

func (pc *PCommunicator) Send(x interface{}) {
	fmt.Println("Sending...")
	pc.Ch <- x
}

func (pc *PCommunicator) Close() {
	close(pc.Ch)
}
