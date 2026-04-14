package daemon

import (
	"testing"
	"time"
)

func BenchmarkBudgetSpend(b *testing.B) {
	budg := newBudget(b.N+1, time.Nanosecond)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		budg.Spend(1)
	}
}

func BenchmarkBudgetRemaining(b *testing.B) {
	budg := newBudget(100, time.Second)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		budg.Remaining()
	}
}

func BenchmarkBudgetReset(b *testing.B) {
	budg := newBudget(100, time.Second)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		budg.Reset()
	}
}
