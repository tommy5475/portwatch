package daemon

import (
	"testing"
)

func BenchmarkEpochAdvance(b *testing.B) {
	e := newEpoch()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.advance()
	}
}

func BenchmarkEpochGeneration(b *testing.B) {
	e := newEpoch()
	e.advance()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.generation()
	}
}

func BenchmarkEpochConcurrentAdvanceGeneration(b *testing.B) {
	e := newEpoch()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%2 == 0 {
				e.advance()
			} else {
				_ = e.generation()
			}
			i++
		}
	})
}
