package daemon

import (
	"sync"
	"testing"
	"time"
)

func TestTokenClaimSucceeds(t *testing.T) {
	tok := newToken("tok-1", time.Second, nil)
	if !tok.Claim() {
		t.Fatal("expected first claim to succeed")
	}
}

func TestTokenClaimIsIdempotent(t *testing.T) {
	tok := newToken("tok-2", time.Second, nil)
	tok.Claim()
	if tok.Claim() {
		t.Fatal("expected second claim to fail")
	}
}

func TestTokenIsClaimed(t *testing.T) {
	tok := newToken("tok-3", time.Second, nil)
	if tok.IsClaimed() {
		t.Fatal("should not be claimed initially")
	}
	tok.Claim()
	if !tok.IsClaimed() {
		t.Fatal("should be claimed after Claim()")
	}
}

func TestTokenTimesOut(t *testing.T) {
	expired := make(chan string, 1)
	tok := newToken("tok-4", 20*time.Millisecond, func(id string) {
		expired <- id
	})

	select {
	case id := <-expired:
		if id != "tok-4" {
			t.Fatalf("unexpected id %q", id)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("onExpire callback was not called")
	}

	if !tok.IsTimedOut() {
		t.Fatal("expected token to be timed out")
	}
	if tok.Claim() {
		t.Fatal("claim should fail after timeout")
	}
}

func TestTokenClaimPreventsTimeout(t *testing.T) {
	called := false
	tok := newToken("tok-5", 50*time.Millisecond, func(_ string) { called = true })
	tok.Claim()
	time.Sleep(100 * time.Millisecond)
	if called {
		t.Fatal("onExpire should not fire after successful claim")
	}
}

func TestTokenAgeIsPositive(t *testing.T) {
	tok := newToken("tok-6", time.Second, nil)
	time.Sleep(5 * time.Millisecond)
	if tok.Age() <= 0 {
		t.Fatal("age should be positive after sleep")
	}
}

func TestTokenDefaultTTLOnInvalidArgs(t *testing.T) {
	tok := newToken("tok-7", -1, nil)
	if tok.ttl != 30*time.Second {
		t.Fatalf("expected default TTL 30s, got %v", tok.ttl)
	}
}

func TestTokenConcurrentClaim(t *testing.T) {
	tok := newToken("tok-8", time.Second, nil)
	var wg sync.WaitGroup
	results := make([]bool, 10)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			results[i] = tok.Claim()
		}(i)
	}
	wg.Wait()
	successCount := 0
	for _, ok := range results {
		if ok {
			successCount++
		}
	}
	if successCount != 1 {
		t.Fatalf("expected exactly 1 successful claim, got %d", successCount)
	}
}
