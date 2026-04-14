package daemon

import (
	"context"
	"testing"
	"time"
)

func BenchmarkRunLoopTicks(b *testing.B) {
	count := 0
	fn := func(ctx context.Context) {
		count++
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		count = 0
		ft := newFakeTicker()
		ctx, cancel := context.WithCancel(context.Background())

		done := make(chan struct{})
		go func() {
			runLoop(ctx, ft, fn)
			close(done)
		}()

		ft.ch <- time.Now()
		ft.ch <- time.Now()
		ft.ch <- time.Now()
		cancel()
		<-done
	}
}

func BenchmarkRunLoopCancelImmediate(b *testing.B) {
	fn := func(ctx context.Context) {}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ft := newFakeTicker()
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		done := make(chan struct{})
		go func() {
			runLoop(ctx, ft, fn)
			close(done)
		}()
		<-done
	}
}
