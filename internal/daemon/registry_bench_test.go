package daemon

import (
	"fmt"
	"testing"
)

func BenchmarkRegistryRegister(b *testing.B) {
	r := newRegistry()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Register(fmt.Sprintf("entry-%d", i), i, "bench")
	}
}

func BenchmarkRegistryGet(b *testing.B) {
	r := newRegistry()
	for i := 0; i < 1000; i++ {
		r.Register(fmt.Sprintf("entry-%d", i), i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = r.Get(fmt.Sprintf("entry-%d", i%1000))
	}
}

func BenchmarkRegistryFilterByTag(b *testing.B) {
	r := newRegistry()
	for i := 0; i < 500; i++ {
		r.Register(fmt.Sprintf("tagged-%d", i), i, "active")
	}
	for i := 0; i < 500; i++ {
		r.Register(fmt.Sprintf("plain-%d", i), i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = r.FilterByTag("active")
	}
}
