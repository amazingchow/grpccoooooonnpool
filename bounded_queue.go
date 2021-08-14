package gpool

import (
	"sync"

	"github.com/gammazero/deque"
)

// BoundedQueue is a threadsafe bounded queue.
type BoundedQueue struct {
	cond *sync.Cond

	q   *deque.Deque
	len uint32
	cap uint32
}

// NewBoundedQueue returns a new BoundedQueue instance.
func NewBoundedQueue(cap uint32) *BoundedQueue {
	if cap == 0 {
		cap = 800
	}

	q := &BoundedQueue{
		cond: sync.NewCond(&sync.Mutex{}),
		q:    deque.New(int(cap)),
		len:  0,
		cap:  cap,
	}

	return q
}

// Push adds conn into the queue.
func (q *BoundedQueue) Push(x int) {
	q.cond.L.Lock()
	for uint32(q.q.Len()) == q.Cap() {
		// P1: queue is full now, wait for consumers to pop conn.
		q.cond.Wait()
	}
	defer q.cond.L.Unlock()

	q.q.PushBack(x)
	// P2 -> P3: tell consumers that there is conn enqueued.
	q.cond.Broadcast()
}

// Pop gets conn from the queue.
func (q *BoundedQueue) Pop() int {
	q.cond.L.Lock()
	for q.q.Len() == 0 {
		// P3: queue is empty now, wait for producers to push conn.
		// TODO: implement q.cond.WaitWithTimeout(timeout time.Duration)
		q.cond.Wait()
	}
	defer q.cond.L.Unlock()

	x := q.q.PopFront().(int)

	// P4 -> P1: tell producers that there is conn dequeued.
	q.cond.Broadcast()
	return x
}

// Len gets the length of the queue.
func (q *BoundedQueue) Len() uint32 {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	return uint32(q.q.Len())
}

// Len gets the capacity of the queue.
func (q *BoundedQueue) Cap() uint32 {
	return uint32(q.q.Cap())
}
