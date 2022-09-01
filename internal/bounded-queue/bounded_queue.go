package boundedq

import (
	"sync"
	"time"

	"github.com/gammazero/deque"
)

// BoundedQueue is a thread-safe bounded queue.
type BoundedQueue struct {
	cond *sync.Cond

	q   *deque.Deque
	len uint32
	cap uint32
}

// NewBoundedQueue returns a new BoundedQueue instance.
func NewBoundedQueue(cap uint32) *BoundedQueue {
	cap = FindNearestPowerOf2(cap)
	if cap == 0 {
		cap = 64
	}
	return &BoundedQueue{
		cond: sync.NewCond(&sync.Mutex{}),
		q:    deque.New(int(cap)),
		len:  0,
		cap:  cap,
	}
}

func FindNearestPowerOf2(n uint32) uint32 {
	if (n & (n - 1)) == 0 {
		return n
	}
	var k uint32 = 1
	for (k << 1) < n {
		k <<= 1
	}
	return k << 1
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

// Pop gets and removes conn from the queue.
// If waitTime <= 0, which means block waiting for Pop.
func (q *BoundedQueue) Pop(wait bool, waitTime int64 /* in millisecs */) int {
	q.cond.L.Lock()
	for q.q.Len() == 0 {
		if !wait {
			q.cond.L.Unlock()
			return -1
		}
		// P3: queue is empty now, wait for producers to push conn.
		if waitTime > 0 {
			waitNotify := make(chan struct{}, 1)
			go func() {
				q.cond.Wait()
				waitNotify <- struct{}{}
			}()

			t := time.NewTicker(time.Duration(waitTime) * time.Millisecond)
			defer t.Stop()
			select {
			case <-t.C:
				q.cond.L.Unlock()
				return -1
			case <-waitNotify:
			}
		} else {
			q.cond.Wait()
		}
	}
	defer q.cond.L.Unlock()

	x := q.q.PopFront().(int)
	// P4 -> P1: tell producers that there is conn dequeued.
	q.cond.Broadcast()
	return x
}

// Len returns the length of the queue.
func (q *BoundedQueue) Len() uint32 {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	return uint32(q.q.Len())
}

// Len returns the capacity of the queue.
func (q *BoundedQueue) Cap() uint32 {
	return uint32(q.q.Cap())
}
