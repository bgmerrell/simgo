package simgo

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

func (eq eventQueue) Len() int { return len(eq) }

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

func (eq eventQueue) Swap(i, j int) {
	eq[i], eq[j] = eq[j], eq[i]
	eq[i].idx = i
	eq[j].idx = j
}

func (eq *eventQueue) Push(x interface{}) {
	n := len(*eq)
	eventQueueItem := x.(*eventQueueItem)
	eventQueueItem.idx = n
	*eq = append(*eq, eventQueueItem)
}

func (eq *eventQueue) Pop() interface{} {
	old := *eq
	n := len(old)
	eventQueueItem := old[n-1]
	eventQueueItem.idx = -1 // for safety
	*eq = old[0 : n-1]
	return eventQueueItem
}
