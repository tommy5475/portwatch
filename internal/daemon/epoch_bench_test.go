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
	for i := 0; i < 100; i++ {
		e.advance()
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.generation()
	}
}

func BenchmarkEpochConcurrentAdvanceGeneration(b *testing.B) {
	e := newEpoch()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			e.advance()
			_ = e.generation()
		}
	})
}
