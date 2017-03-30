package simulago

type eventCmd int

// PendingValue is a unique object for a pending Event value
type PendingValue struct{}

// An event is an event that may happen at some point in time.
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
type event struct {
	// The environment the event lives in
	env *Environment
	// List of functions that are called when the event is processed.
	callbacks []func(...interface{})
}

// Timeout embeds an event and adds a delay
type Timeout struct {
	*event
	delay uint64
}

// NewTimeout returns a new Timeout object given an environment, delay and an
// event value.
func NewTimeout(env *Environment, delay uint64, value interface{}) Timeout {
	return Timeout{
		&event{env, nil},
		delay,
	}
}

type Process struct {
	*event
	env *Environment
	fn  func(*Environment)
}

func NewProcess(env *Environment, fn func(*Environment)) *Process {
	return &Process{
		&event{env, nil},
		env,
		fn,
	}
}
