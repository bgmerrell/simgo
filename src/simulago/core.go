package simulago

type Environment struct {
	// Current time count
	Now uint64
	// Event ID counter
	EID uint64
	// The list of all currently scheduled events
	Queue []Event
	// TODO: self._active_proc = None
}
