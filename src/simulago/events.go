package simulago

// PendingValue is a unique object for a pending Event value
type PendingValue struct{}

type Event struct {
	// The Environment the event lives in
	Env *Environment
	// List of functions that are called when the event is processed.
	Callbacks []func(...interface{})
	// The Event's response value
	Value interface{}
}
