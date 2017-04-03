package simulago

import (
	"fmt"

	"github.com/bgmerrell/simulago/pcomm"
)

const (
	PriorityNormal int = iota
	PriorityUrgent
)

type (
	eventCmd     int
	PriorityCode int
)

// PendingValue is a unique object for a pending Event value
type PendingValue struct{}

// An Event is an event that may happen at some point in time.
//
//    An event
//
//    - may happen (i.e., triggered is False),
//    - is going to happen (i.e., triggered True) or
//    - has happened (i.e., processed True).
//
// Every event is bound to an environment (env) and is initially not triggered.
// Events are scheduled for processing by the environment after they are
// triggered by either `succeed`, `fail` or `trigger`. These methods also set
// the `ok` flag and the value of the event.
//
// An event has a list of `callbacks`. Once an event gets processed, all
// callbacks will be called with the event as the single argument. Callbacks
// can check if the event was successful by examining `ok` and do further
// processing with the value it has produced.
//
// TODO: Talk about how events are finalized/defused (?) after being processed.
type Event struct {
	// The environment the event lives in
	env *Environment
	// List of functions that are called when the event is processed.
	callbacks []func(...interface{})
}

// Timeout embeds an event and adds a delay
type Timeout struct {
	*Event
	delay uint64
}

// NewTimeout returns a new Timeout object given an environment, delay and an
// Event value.
func NewTimeout(env *Environment, delay uint64) Timeout {
	return Timeout{
		&Event{env, make([]func(...interface{}), 0)},
		delay,
	}
}

func (to *Timeout) Schedule(env *Environment) {
	env.Schedule(to.Event, PriorityNormal, to.delay)
}

type Process struct {
	*Event
	env *Environment
	pc  *pcomm.PCommunicator
}

func NewProcess(env *Environment, pc *pcomm.PCommunicator) *Process {
	return &Process{
		&Event{env, make([]func(...interface{}), 0)},
		env,
		pc,
	}
}

func (p *Process) Init() {
	fmt.Println("Adding callback...")
	p.Event.callbacks = append(p.Event.callbacks, p.resume)
	p.env.Schedule(p.Event, PriorityUrgent, 0)
}

func (p *Process) resume(...interface{}) {
	for event := range p.pc.Ch {
		fmt.Println("Ranging through channel")
		if event == nil {
			// The other end was just waiting, so continue on
			continue
		} else {
			p.Event = event.(*Event)
		}

		if p.Event.callbacks != nil {
			// The event has not yet been triggered. Register
			// callback to resume the process if that happens.
			fmt.Println("Adding callback from resume()...")
			p.Event.callbacks = append(p.Event.callbacks, p.resume)
			break
		}
	}
}

func ProcWrapper(env *Environment, procFn func(*Environment, *pcomm.PCommunicator)) *pcomm.PCommunicator {
	pc := pcomm.New()
	go func() {
		pc.Wait()
		procFn(env, pc)
		pc.Close()
	}()
	return pc
}
