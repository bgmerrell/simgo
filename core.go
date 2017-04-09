package simgo

import (
	"fmt"

	"github.com/juju/errgo"
)

type EventID uint64

func (eid *EventID) Next() EventID {
	(*eid)++
	return *eid
}

// An environment is the execution environment for an event-based simulation.
// The passing of time is simulated by stepping from event to event.
type Environment struct {
	// Current time count
	Now uint64
	// Event ID counter
	eid EventID
	// The list of all currently scheduled events
	queue eventQueue
	// TODO: self._active_proc = None
	// shouldStop tells the Environment when it's time to stop running
	shouldStop bool
}

func NewEnvironment() *Environment {
	return new(Environment)
}

/*
def run(self, until=None):
        """Executes :meth:`step()` until the given criterion *until* is met.

        - If it is ``None`` (which is the default), this method will return
          when there are no further events to be processed.

        - If it is an :class:`~simpy.events.Event`, the method will continue
          stepping until this event has been triggered and will return its
          value.  Raises a :exc:`RuntimeError` if there are no further events
          to be processed and the *until* event was not triggered.

        - If it is a number, the method will continue stepping
          until the environment's time reaches *until*.

        """
*/

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
				at = uint64(intAt)
			} else {
				at, _ = until.(uint64)
			}
			if at <= env.Now {
				return nil, errgo.Newf(`"until" value (%d) should be greater than the current simulation time (%d)`, at, env.Now)
			}
			untilEvent = NewEvent(env)
			untilEvent.Value.Set(nil)
			fmt.Printf("Scheduling to end at: %d\n", at)
			env.Schedule(untilEvent, PriorityUrgent, at-env.Now)

		}
		untilEvent.callbacks = append(untilEvent.callbacks, env.stopSimulation)
	}
	for !env.shouldStop {
		env.Step()
	}
	return untilEvent.Value.Get()
}

func (env *Environment) Step() {
	fmt.Println("Step() called...")
	item := env.queue.Pop().(*eventQueueItem)
	env.Now = item.time
	fmt.Println("Now:", env.Now)

	// Process the event callbacks
	fmt.Printf("Processing %d callback(s)...\n", len(item.callbacks))
	callbacks := make([]func(...interface{}), len(item.callbacks))
	copy(callbacks, item.callbacks)
	item.callbacks = nil
	for _, callback := range callbacks {
		callback()
	}
}

func (env *Environment) Schedule(v *Event, priority int, delay uint64) {
	fmt.Println("Pushing event...")
	env.queue.Push(NewEventQueueItem(v, env.Now+delay, priority, env.eid.Next()))
}

func (env *Environment) stopSimulation(...interface{}) {
	env.shouldStop = true
}
