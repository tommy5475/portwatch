package daemon

import (
	"testing"
	"time"
)

func BenchmarkScheduleNextSuccess(b *testing.B) {
	s, _ := newSchedule(5*time.Second, 60*time.Second, 2.0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.next(false)
	}
}

func BenchmarkScheduleNextFailure(b *testing.B) {
	s, _ := newSchedule(5*time.Second, 60*time.Second, 2.0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.next(true)
		if s.interval() >= 60*time.Second {
			s.reset()
		}
	}
}

func BenchmarkScheduleReset(b *testing.B) {
	s, _ := newSchedule(5*time.Second, 60*time.Second, 2.0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.reset()
	}
}
