package simulago

import (
	"container/heap"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func TestEventQueue(t *testing.T) {
	eq := make(eventQueue, 0)
	heap.Init(&eq)
	items := []*eventQueueItem{
		&eventQueueItem{time: 90, priority: 1, eid: 5},  // 0: 3rd
		&eventQueueItem{time: 100, priority: 2, eid: 1}, // 1: 6th
		&eventQueueItem{time: 110, priority: 0, eid: 2}, // 2: 7th
		&eventQueueItem{time: 100, priority: 1, eid: 3}, // 3: 5th
		&eventQueueItem{time: 90, priority: 0, eid: 4},  // 4: 2nd
		&eventQueueItem{time: 90, priority: 0, eid: 0},  // 5: 1st
		&eventQueueItem{time: 100, priority: 0, eid: 6}, // 6: 4th
	}
	var (
		item *eventQueueItem
		i    int
	)
	for i, item = range items {
		heap.Push(&eq, item)
	}
	if eq.Len() != i+1 {
		t.Fatalf("eq.Len() = %d, want: %d", eq.Len(), i)
	}

	expectedIdxOrder := []int{5, 4, 0, 6, 3, 1, 2}
	for i := 0; eq.Len() > 0; i++ {
		item = heap.Pop(&eq).(*eventQueueItem)
		// popped item indexes are always -1
		items[expectedIdxOrder[i]].idx = -1
		if !reflect.DeepEqual(item, items[expectedIdxOrder[i]]) {
			t.Errorf("Test %d: item = %#v, want: %#v\n", i+1, item, items[expectedIdxOrder[i]])
		}
	}
}

func BenchmarkPopPush(b *testing.B) {
	const (
		nItems      = 1000000
		maxTime     = 1000000
		maxPriority = 5
	)
	rand.Seed(time.Now().UnixNano())
	eq := make(eventQueue, 0)
	heap.Init(&eq)
	items := make([]*eventQueueItem, nItems)

	// &eventQueueItem{time: 90, priority: 1, eid: 5},  // 0: 3rd
	for i := 0; i < nItems; i++ {
		items[i] = &eventQueueItem{
			time:     uint64(rand.Intn(maxTime)),
			priority: rand.Intn(maxPriority),
			eid:      uint64(i),
		}
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for _, item := range items {
			heap.Push(&eq, item)
		}
		for eq.Len() > 0 {
			heap.Pop(&eq)
		}
	}
}
