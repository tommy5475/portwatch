package daemon

import (
	"fmt"
	"testing"
	"time"
)

func BenchmarkRosterJoinLeave(b *testing.B) {
	r := newRoster()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		name := fmt.Sprintf("w-%d", i)
		r.join(name)
		r.leave(name)
	}
}

func BenchmarkRosterCheckin(b *testing.B) {
	r := newRoster()
	r.join("worker")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.checkin("worker")
	}
}

func BenchmarkRosterMarkStale(b *testing.B) {
	r := newRoster()
	for i := 0; i < 100; i++ {
		r.join(fmt.Sprintf("w-%d", i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.markStale(time.Second)
	}
}
