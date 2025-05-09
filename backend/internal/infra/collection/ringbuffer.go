package collection

import "sync"

type RingBuffer[T any] struct {
	data  []T
	start int
	size  int
	cap   int
	mu    sync.RWMutex
}

func NewRingBuffer[T any](capacity int) *RingBuffer[T] {
	return &RingBuffer[T]{
		data: make([]T, capacity),
		cap:  capacity,
	}
}

func (rb *RingBuffer[T]) Push(item T) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	if rb.size < rb.cap {
		rb.data[(rb.start+rb.size)%rb.cap] = item
		rb.size++
	} else {
		rb.data[rb.start] = item
		rb.start = (rb.start + 1) % rb.cap
	}
}

func (rb *RingBuffer[T]) GetAll() []T {
	rb.mu.RLock()
	defer rb.mu.RUnlock()

	out := make([]T, rb.size)
	for i := range rb.size {
		out[i] = rb.data[(rb.start+i)%rb.cap]
	}
	return out
}

func (rb *RingBuffer[T]) Len() int {
	rb.mu.RLock()
	defer rb.mu.RUnlock()
	return rb.size
}

func (rb *RingBuffer[T]) Clear() int {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	clear(rb.data)

	tmpSize := rb.size
	rb.start = 0
	rb.size = 0

	return tmpSize
}
