package daemon

import "testing"

func BenchmarkQueuePush(b *testing.B) {
	q := newQueue[int](b.N + 1)
	b.ResetTimer()
	for i := range b.N {
		q.Push(i)
	}
}

func BenchmarkQueuePushPop(b *testing.B) {
	q := newQueue[int](64)
	b.ResetTimer()
	for i := range b.N {
		q.Push(i)
		q.Pop() //nolint:errcheck
	}
}

func BenchmarkQueueDrain(b *testing.B) {
	const size = 64
	q := newQueue[int](size)
	for i := range size {
		q.Push(i)
	}
	b.ResetTimer()
	for range b.N {
		for i := range size {
			q.Push(i)
		}
		q.Drain()
	}
}
