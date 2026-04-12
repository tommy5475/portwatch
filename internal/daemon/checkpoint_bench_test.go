package daemon

import "testing"

func BenchmarkCheckpointMark(b *testing.B) {
	cp := newCheckpoint()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cp.Mark("scan")
	}
}

func BenchmarkCheckpointCount(b *testing.B) {
	cp := newCheckpoint()
	cp.Mark("scan")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cp.Count("scan")
	}
}

func BenchmarkCheckpointMarkMultiKey(b *testing.B) {
	cp := newCheckpoint()
	keys := []string{"scan.start", "scan.complete", "alert.sent", "state.saved"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cp.Mark(keys[i%len(keys)])
	}
}
