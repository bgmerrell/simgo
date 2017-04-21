package simgo

import (
	"github.com/juju/errgo"
)

const (
	// Built-in event priorities
	PriorityUrgent = iota
	PriorityNormal
)

// conditionEvaluateFn is the evaluate function signature used for conditions
type conditionEvaluateFn func(events []*Event, count int) bool

// TODO: Move EventValue to its own package or refactor.  Currently, it's too
// easy to use the struct fields directly when the methods should really be
// the only things used.

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

// Succeeds sets the event's value, marks it as successful and schedules it for
// processing by the environment. Returns the event instance along with any
// errors.
func (e *Event) Succeed(val interface{}) (*Event, error) {
	if !e.Value.isPending {
		return e, errgo.Newf("%s has already been triggered", e)
	}
	e.Value.Set(val)
	e.env.Schedule(e, PriorityNormal, 0)
	return e, nil
}

// Fail sets the provided EventValue as the events value, marks the event as
// failed, and schedules it for processing by the environment.  The event
// instance is returned along with any errors.
func (e *Event) Fail(eventValue *EventValue) (*Event, error) {
	if !e.Value.isPending {
		return nil, errgo.Newf("%s has already been triggered", e)
	}

	errVal, err := eventValue.Get()
	if err != nil {
		return nil, err
	}
	if _, ok := errVal.(error); ok {
		e.Value.Set(errVal)
	} else {
		return nil, errgo.Newf("%#v is not an error", errVal)
	}
	e.env.Schedule(e, PriorityNormal, 0)
	return e, nil
}

// isOK returns whether the event is OK, which means that the EventValue is
// not an error.
func (e *Event) isOK() bool {
	_, isErr := e.Value.val.(error)
	return !isErr
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

// Schedule schedules the timeout event for the provided environment.
func (to *Timeout) Schedule(env *Environment) {
	env.Schedule(to.Event, PriorityNormal, to.delay)
}

// A Process processes an event yielding process function.
//
// A user implements a process function which is a coroutine function that can
// suspend its execution by yielding an event (using ProcComm.Yield()).
// Process will take care of resuming the process function with the value of
// that event once it has happened.
type Process struct {
	*Event
	env *Environment
	pc  *ProcComm
}

// NewProcess returns a new Process given an Environment and a ProcComm
// (which is used to communicate between the process function coroutine and the
// Process).
func NewProcess(env *Environment, pc *ProcComm) *Process {
	return &Process{
		NewEvent(env),
		env,
		pc,
	}
}

// Init initializes the process.  The process's Event is automatically
// triggered and scheduled.
func (p *Process) Init() {
	initEvent := NewEvent(p.env)
	initEvent.callbacks = append(initEvent.callbacks, p.resume)
	initEvent.Value.Set(nil)
	p.env.Schedule(initEvent, PriorityUrgent, 0)
}

// resume takes care of resuming the process function with the value of the
// provided Event.
func (p *Process) resume(event *Event) {
	p.env.ActiveProcess = p
	defer func() {
		p.env.ActiveProcess = nil
	}()
	for {
		// event value is already triggered, no need to check err
		eventVal, _ := event.Value.Get()
		if nextEvent, ok := p.pc.Resume(eventVal); !ok {
			// Set the process value to nil if it hasn't been set
			// already.
			if p.Event.Value.isPending {
				p.Event.Value.Set(nil)
			}
			p.env.Schedule(p.Event, PriorityNormal, 0)
			break
		} else {
			event = nextEvent
		}

		if event.callbacks != nil {
			// The event has not yet been triggered. Register
			// callback to resume the process if that happens.
			event.callbacks = append(event.callbacks, p.resume)
			return
		}
	}
}

// ReturnValue returns the value returned by the process function.
func (p *Process) ReturnValue() interface{} {
	return p.pc.returnValue
}

// ProcWrapper is function that turns a user process function into a coroutine
// that can suspend its execution by yielding an event (using
// ProcComm.Yield()).
//
// See the examples directory for example usage.
//
func ProcWrapper(env *Environment, procFn func(*Environment, *ProcComm) interface{}) *ProcComm {
	pc := NewProcComm()
	go func() {
		// An initial yield imitates coroutine behavior of not
		// executing the coroutine body upon creation.
		pc.Yield(nil)
		pc.Finish(procFn(env, pc))
	}()
	return pc
}

// A Condition embeds an event that gets triggered once the "evaluate"
// condition function returns true on the list of events.
//
// The value of the condition event is an instance of ConditionValue which
// allows convenient access to the input events and their values. The
// ConditionValue will only contain entries for those events that occurred
// before the condition is processed.
//
// If one of the events fails, the condition also fails and... TODO: does what?
//
// The evaluate function receives the list of target events and the number
// of processed events in this list: evaluate(events, processedCount). If it
// returns true, the condition is triggered.
//
// Condition events can be nested.
type Condition struct {
	*Event
	evaluateFn conditionEvaluateFn
	events     []*Event
	count      int
}

func NewCondition(env *Environment, evaluateFn conditionEvaluateFn, events []*Event) *Condition {
	c := &Condition{
		NewEvent(env),
		evaluateFn,
		make([]*Event, 0),
		0,
	}

	if len(c.events) == 0 {
		// Immediately succeed if no events are provided
		c.Event.Succeed(nil)
		return c
	}

	// TODO: Simpy has a check that all events are in the same environment.
	// Do we need one?

	// Check if the condition is met for each processed event.  Attach
	// check() as a callback otherwise.
	for _, event := range c.events {
		event.callbacks = append(event.callbacks, c.check)
	}
	c.Event.callbacks = append(c.Event.callbacks, c.buildValue)

	return c
}

// check checks if the condition was already met and schedules the event if
// it was.
func (c *Condition) check(event *Event) {
	if !c.Event.Value.isPending {
		return
	}

	c.count++

	if !c.Event.isOK() {
		// TODO: Use "Defused?"
		// c.Event.Defused = true
		c.Event.Fail(c.Event.Value)
	} else if c.evaluateFn(c.events, c.count) {
		c.Event.Succeed(nil)
	}
}

func (c *Condition) buildValue(event *Event) {
	c.removeCheckCallbacks()
	if c.Event.isOK() {
		// TODO: Set and populate the event value to a ConditionValue
		// simpy:
		//   self._value = ConditionValue()
		//   self._populate_value(self._value)
	}
}

// removeCheckCallbacks removes check() callbacks from events recursively.
//
// Once the condition has triggered, the condition's events no longer need
// to have check() callbacks. Removing the check() callbacks is important
// to break circular references between the condition and untriggered events.
func (c *Condition) removeCheckCallbacks() {
	// TODO: Implement me
	/*
	   for event in self._events:
	       if event.callbacks and self._check in event.callbacks:
	           event.callbacks.remove(self._check)
	       if isinstance(event, Condition):
	           event._remove_check_callbacks()

	*/
}
