package daemon

import (
	"sync"
	"testing"
	"time"
)

func TestCacheSetAndGet(t *testing.T) {
	c := newCache(time.Second)
	c.Set("k", 42)
	v, ok := c.Get("k")
	if !ok {
		t.Fatal("expected key to be present")
	}
	if v.(int) != 42 {
		t.Fatalf("expected 42, got %v", v)
	}
}

func TestCacheMiss(t *testing.T) {
	c := newCache(time.Second)
	_, ok := c.Get("missing")
	if ok {
		t.Fatal("expected miss for unknown key")
	}
}

func TestCacheExpiry(t *testing.T) {
	c := newCache(10 * time.Millisecond)
	c.Set("k", "val")
	time.Sleep(20 * time.Millisecond)
	_, ok := c.Get("k")
	if ok {
		t.Fatal("expected entry to have expired")
	}
}

func TestCacheDelete(t *testing.T) {
	c := newCache(time.Second)
	c.Set("k", true)
	c.Delete("k")
	_, ok := c.Get("k")
	if ok {
		t.Fatal("expected key to be deleted")
	}
}

func TestCacheLen(t *testing.T) {
	c := newCache(time.Second)
	c.Set("a", 1)
	c.Set("b", 2)
	c.Set("c", 3)
	if n := c.Len(); n != 3 {
		t.Fatalf("expected 3, got %d", n)
	}
}

func TestCacheLenEvictsExpired(t *testing.T) {
	c := newCache(10 * time.Millisecond)
	c.Set("x", 1)
	time.Sleep(20 * time.Millisecond)
	if n := c.Len(); n != 0 {
		t.Fatalf("expected 0 after expiry, got %d", n)
	}
}

func TestCacheDefaultTTL(t *testing.T) {
	c := newCache(0) // should default to 60s
	if c.ttl != 60*time.Second {
		t.Fatalf("expected default TTL 60s, got %v", c.ttl)
	}
}

func TestCacheConcurrentAccess(t *testing.T) {
	c := newCache(time.Second)
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			key := string(rune('a' + n%26))
			c.Set(key, n)
			c.Get(key)
		}(i)
	}
	wg.Wait()
}
