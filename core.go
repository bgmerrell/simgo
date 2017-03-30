package simulago

// An environment is the execution environment for an event-based simulation.
// The passing of time is simulated by stepping from event to event.
type environment struct {
	// Current time count
	now uint64
	// Event ID counter
	eid uint64
	// The list of all currently scheduled events
	queue eventQueue
	// TODO: self._active_proc = None
}

func (env *environment) Step() {
}
