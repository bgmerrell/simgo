package simgo

import "container/heap"

// eventQueue implements the heap.Interface interface by way of the following
// methods: Len(), Less(), Swap(), Push(), Pop().
var _ heap.Interface = (*eventQueue)(nil)

// eventQueueItem is an event which has been schedule and thus has some
// additional relevant metadata.
type eventQueueItem struct {
	// The event that is scheduled
	*Event
	// The time the event is scheduled for
	time uint64
	// The priority of the scheduled event
	priority int
	// The incremental event ID
	eid EventID
	// The index in the event queue (used for performance gains)
	idx int
}

// NewEventQueueItem returns a new eventQueueItem given an event, time,
// priority code, and EID.
func NewEventQueueItem(v *Event, time uint64, priority int, eid EventID) *eventQueueItem {
	return &eventQueueItem{
		Event:    v,
		time:     time,
		priority: priority,
		eid:      eid,
	}
}

// A eventQueue implements heap.Interface and holds ScheduledEvent objects.
// Use heap.Push() and heap.Pop() to manipulate.
type eventQueue []*eventQueueItem

// Len returns the length of the event queue.
func (eq eventQueue) Len() int { return len(eq) }

// Less returns whether the event at index i is "less than" (i.e., has a
// lower priority) than the event at index j.
func (eq eventQueue) Less(i, j int) bool {
	if eq[i].time < eq[j].time {
		return true
	}
	if eq[i].time == eq[j].time {
		if eq[i].priority < eq[j].priority {
			return true
		}
		if eq[i].priority == eq[j].priority {
			if eq[i].eid < eq[j].eid {
				return true
			}
		}
	}
	return false
}

// Swap swaps the events at the given indexes.
func (eq eventQueue) Swap(i, j int) {
	eq[i], eq[j] = eq[j], eq[i]
	eq[i].idx = i
	eq[j].idx = j
}

// Push pushes the given value, x, to the event queue.
func (eq *eventQueue) Push(x interface{}) {
	n := len(*eq)
	eventQueueItem := x.(*eventQueueItem)
	eventQueueItem.idx = n
	*eq = append(*eq, eventQueueItem)
}

// Pop pops a value from the event queue and returns it.
func (eq *eventQueue) Pop() interface{} {
	if len(*eq) <= 0 {
		return nil
	}
	old := *eq
	n := len(old)
	eventQueueItem := old[n-1]
	eventQueueItem.idx = -1 // for safety
	*eq = old[0 : n-1]
	return eventQueueItem
}
