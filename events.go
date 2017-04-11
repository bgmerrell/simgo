package simgo

import (
	"fmt"

	"github.com/bgmerrell/simgo/pcomm"
	"github.com/juju/errgo"
)

const (
	PriorityNormal int = iota
	PriorityUrgent
)

type (
	eventCmd     int
	PriorityCode int
)

// EventValue holds the value state for an Event.  If the value is pending it
// means that the event
type EventValue struct {
	val       interface{}
	isPending bool
}

// NewEventValue returns a pending EventValue.
func NewEventValue() *EventValue {
	return &EventValue{nil, true}
}

// Set sets the underlying value (and sets isPending accordingly).
func (ev *EventValue) Set(value interface{}) {
	ev.val = value
	ev.isPending = false
}

// Get returns the underlying event value along with an error if the value is
// still pending.
func (ev *EventValue) Get() (interface{}, error) {
	var err error
	if ev.isPending {
		err = errgo.New("event value is still pending")
	}
	return ev.val, err
}

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
	callbacks []func(*Event)
	// Value holds the event's value
	Value *EventValue
}

// NewEvent returns a new Event object with default values.
func NewEvent(env *Environment) *Event {
	return &Event{
		env,
		make([]func(*Event), 0),
		NewEventValue(),
	}
}

// Timeout embeds an event and adds a delay
type Timeout struct {
	*Event
	delay uint64
}

// NewTimeout returns a new Timeout object given an environment, delay and an
// Event value.  The event is automatically triggered when this function is
// called.
func NewTimeout(env *Environment, delay uint64, value interface{}) Timeout {
	return Timeout{
		&Event{
			env,
			make([]func(*Event), 0),
			&EventValue{value, false},
		},
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
		NewEvent(env),
		env,
		pc,
	}
}

// Init initializes the process.  The event is automatically triggered and
// scheduled.
func (p *Process) Init() {
	fmt.Println("Adding callback...")
	p.Event.callbacks = append(p.Event.callbacks, p.resume)
	p.Event.Value.Set(nil)
	p.env.Schedule(p.Event, PriorityUrgent, 0)
}

func (p *Process) resume(event *Event) {
	fmt.Println("in resume")
	for {
		// event value is already triggered, no need to check err
		eventVal, _ := event.Value.Get()
		if nextEvent, ok := p.pc.Resume(eventVal); !ok {
			fmt.Println("proc finished...")
			break
		} else if nextEvent == nil {
			// The other end was just waiting, so continue on
			fmt.Println("Continuing...")
		} else {
			p.Event = nextEvent.(*Event)
		}

		if p.Event.callbacks != nil {
			// The event has not yet been triggered. Register
			// callback to resume the process if that happens.
			fmt.Println("Adding callback from resume()...")
			p.Event.callbacks = append(p.Event.callbacks, p.resume)
			return
		}
	}
}

func ProcWrapper(env *Environment, procFn func(*Environment, *pcomm.PCommunicator)) *pcomm.PCommunicator {
	pc := pcomm.New()
	go func() {
		pc.Yield(nil)
		procFn(env, pc)
		pc.Finish()
	}()
	return pc
}
