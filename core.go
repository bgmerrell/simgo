package simgo

import "fmt"

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
}

func NewEnvironment() *Environment {
	return new(Environment)
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
