package daemon

import "testing"

func BenchmarkHooksFire(b *testing.B) {
	h := newHooks()
	h.Register(HookAfterScan, func(_ hookEvent, _ any) {})
	h.Register(HookAfterScan, func(_ hookEvent, _ any) {})
	h.Register(HookAfterScan, func(_ hookEvent, _ any) {})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Fire(HookAfterScan, nil)
	}
}

func BenchmarkHooksRegister(b *testing.B) {
	h := newHooks()
	fn := func(_ hookEvent, _ any) {}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Register(HookBeforeScan, fn)
	}
}

func BenchmarkHooksFireNoSubscribers(b *testing.B) {
	h := newHooks()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Fire(HookOnChange, nil)
	}
}
