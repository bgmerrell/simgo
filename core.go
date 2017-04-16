package simgo

import (
	"container/heap"

	"github.com/juju/errgo"
)

// An EventID is a unique numerical ID for each event.
type EventID uint64

// Next returns the next sequential event ID.
func (eid *EventID) Next() EventID {
	(*eid)++
	return *eid
}

// An environment is the execution environment for an event-based simulation.
// The passing of time is simulated by stepping from event to event.
type Environment struct {
	// Current time count
	Now uint64
	// ActiveProcess is the currently active process
	ActiveProcess *Process
	// Event ID counter
	eid EventID
	// The list of all currently scheduled events
	queue eventQueue
	// shouldStop tells the Environment when it's time to stop running
	shouldStop bool
}

// NewEnvironment returns an Environment with default values.
func NewEnvironment() *Environment {
	return new(Environment)
}

// Run executes Step() until a condition below is met for the provided "until"
// value:
//
//     - If it is "None``, this method will return when there are no further
//       events to be processed.
//
//     - If it is an Event, the method will continue stepping until this event
//       has been triggered and will return its value.  Returns an error if
//       there are no further events to be processed and the "until" event was
//       not triggered.
//
//     - If it is a number, the method will continue stepping until the
//       environment's time reaches *until*.
func (env *Environment) Run(until interface{}) (interface{}, error) {
	var (
		untilEvent *Event
		at         uint64
	)
	if until != nil {
		switch until := until.(type) {
		default:
			return nil, errgo.New(`"until" value must be a number or an *Event`)
		case *Event:
			untilEvent = until
			if untilEvent.callbacks == nil {
				// "until" event has already been processed.
				return untilEvent.Value.Get()
			}
		case int, uint64:
			if intAt, ok := until.(int); ok {
				if intAt < 0 {
					return nil, errgo.Newf(`"until" value (%d) must be positive`, intAt)
				}
				at = uint64(intAt)
			} else {
				at, _ = until.(uint64)
			}
			if at <= env.Now {
				return nil, errgo.Newf(`"until" value (%d) must be greater than the current simulation time (%d)`, at, env.Now)
			}
			untilEvent = NewEvent(env)
			untilEvent.Value.Set(nil)
			env.Schedule(untilEvent, PriorityUrgent, at-env.Now)

		}
		untilEvent.callbacks = append(untilEvent.callbacks, env.stopSimulation)
	}
	for !env.shouldStop {
		env.Step()
	}
	if untilEvent != nil {
		return untilEvent.Value.Get()
	}
	return nil, nil
}

// Step processes the next event.
func (env *Environment) Step() {
	var eqItem *eventQueueItem
	switch item := heap.Pop(&env.queue).(type) {
	case nil:
		// We're out of event queue items
		env.shouldStop = true
		return
	case *eventQueueItem:
		eqItem = item
	default:
		// Should never happen
		panic("Unknown type from event queue")
	}
	env.Now = eqItem.time

	// Process the event callbacks
	callbacks := make([]func(*Event), len(eqItem.callbacks))
	copy(callbacks, eqItem.callbacks)
	eqItem.callbacks = nil
	for _, callback := range callbacks {
		callback(eqItem.Event)
	}
}

// Schedule adds the provided Event to the event priority queue.  A priority
// and delay for the event is also provided.
func (env *Environment) Schedule(v *Event, priority int, delay uint64) {
	heap.Push(&env.queue, NewEventQueueItem(v, env.Now+delay, priority, env.eid.Next()))
}

// stopSimulation is a special callback that tells the Environment that it's
// time to stop the simulation.
func (env *Environment) stopSimulation(_ *Event) {
	env.shouldStop = true
}
