package simulago

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
	now uint64
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
	// self._now, _, _, event = heappop(self._queue)
	fmt.Println("Popping event...")
	item := env.queue.Pop().(*eventQueueItem)
	env.now = item.time
	fmt.Println("now:", env.now)

	// Process the event callbacks
	fmt.Printf("Processing %d callback(s)...\n", len(item.callbacks))
	for _, callback := range item.callbacks {
		callback()
	}
}

func (env *Environment) Schedule(v *Event, priority int, delay uint64) {
	fmt.Println("Pushing event...")
	env.queue.Push(NewEventQueueItem(v, env.now+delay, priority, env.eid.Next()))
}
