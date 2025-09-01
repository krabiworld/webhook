package cache

import (
	"testing"
	"time"
)

func TestMemory_SetGetExistsDeleteFlush(t *testing.T) {
	cache := NewMemory[string]()

	// Initially, key should not exist
	exists, err := cache.Exists("foo")
	if err != nil {
		t.Fatalf("Exists returned error: %v", err)
	}
	if exists {
		t.Fatalf("expected key to not exist initially")
	}

	// Set value without TTL
	if err := cache.Set("foo", "bar", 0); err != nil {
		t.Fatalf("Set returned error: %v", err)
	}

	// Exists should be true
	exists, err = cache.Exists("foo")
	if err != nil {
		t.Fatalf("Exists returned error: %v", err)
	}
	if !exists {
		t.Fatalf("expected key to exist after set")
	}

	// Get should return value
	got, err := cache.Get("foo")
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if got != "bar" {
		t.Fatalf("expected value 'bar', got %q", got)
	}

	// Delete should remove key
	if err := cache.Delete("foo"); err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}
	exists, err = cache.Exists("foo")
	if err != nil {
		t.Fatalf("Exists returned error: %v", err)
	}
	if exists {
		t.Fatalf("expected key to not exist after delete")
	}

	// Set multiple values and then flush
	if err := cache.Set("a", "1", 0); err != nil {
		t.Fatalf("Set returned error: %v", err)
	}
	if err := cache.Set("b", "2", 0); err != nil {
		t.Fatalf("Set returned error: %v", err)
	}
	if err := cache.Flush(); err != nil {
		t.Fatalf("Flush returned error: %v", err)
	}
	for _, k := range []string{"a", "b"} {
		exists, err := cache.Exists(k)
		if err != nil {
			t.Fatalf("Exists(%s) returned error: %v", k, err)
		}
		if exists {
			t.Fatalf("expected key %s to not exist after flush", k)
		}
	}
}

func TestMemory_TTLExpiry(t *testing.T) {
	cache := NewMemory[string]()

	// Set with a short TTL
	ttl := 50 * time.Millisecond
	if err := cache.Set("temp", "value", ttl); err != nil {
		t.Fatalf("Set returned error: %v", err)
	}

	// Immediately should exist
	exists, err := cache.Exists("temp")
	if err != nil {
		t.Fatalf("Exists returned error: %v", err)
	}
	if !exists {
		t.Fatalf("expected key to exist immediately after set with TTL")
	}

	// Wait past TTL and verify deletion occurred
	time.Sleep(ttl + 50*time.Millisecond)

	exists, err = cache.Exists("temp")
	if err != nil {
		t.Fatalf("Exists returned error: %v", err)
	}
	if exists {
		t.Fatalf("expected key to be deleted after TTL expiry")
	}
}
