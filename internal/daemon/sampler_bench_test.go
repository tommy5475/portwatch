package daemon

import "testing"

func BenchmarkSamplerRecord(b *testing.B) {
	s := newSampler()
	for i := 0; i < b.N; i++ {
		s.record(float64(i))
	}
}

func BenchmarkSamplerStats(b *testing.B) {
	s := newSampler()
	for i := 0; i < 1000; i++ {
		s.record(float64(i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.stats()
	}
}

func BenchmarkSamplerMean(b *testing.B) {
	s := newSampler()
	for i := 0; i < 1000; i++ {
		s.record(float64(i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.mean()
	}
}
