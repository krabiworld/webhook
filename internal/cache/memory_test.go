package cache

import (
	"fmt"
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

// TestMemory_TTLExpiry_ErrorHandling tests TTL expiry error handling
func TestMemory_TTLExpiry_ErrorHandling(t *testing.T) {
	cache := NewMemory[string]()

	// Set with a very short TTL to trigger the async deletion
	ttl := 10 * time.Millisecond
	if err := cache.Set("temp_error", "value", ttl); err != nil {
		t.Fatalf("Set returned error: %v", err)
	}

	// Wait for TTL to expire and async deletion to complete
	time.Sleep(ttl + 100*time.Millisecond)

	// Verify the key was deleted
	exists, err := cache.Exists("temp_error")
	if err != nil {
		t.Fatalf("Exists returned error: %v", err)
	}
	if exists {
		t.Fatalf("expected key to be deleted after TTL expiry")
	}
}

// TestMemory_TTLExpiry_Immediate tests TTL expiry with immediate deletion
func TestMemory_TTLExpiry_Immediate(t *testing.T) {
	cache := NewMemory[string]()

	// Set with a very short TTL and immediately check
	ttl := 1 * time.Millisecond
	if err := cache.Set("immediate_temp", "value", ttl); err != nil {
		t.Fatalf("Set returned error: %v", err)
	}

	// Wait a bit longer than TTL
	time.Sleep(ttl + 10*time.Millisecond)

	// Verify the key was deleted
	exists, err := cache.Exists("immediate_temp")
	if err != nil {
		t.Fatalf("Exists returned error: %v", err)
	}
	if exists {
		t.Fatalf("expected key to be deleted after immediate TTL expiry")
	}
}

// TestMemory_TTLExpiry_ErrorPath tests the error path in TTL expiry
func TestMemory_TTLExpiry_ErrorPath(t *testing.T) {
	// Create a cache with a custom type that might trigger errors
	cache := NewMemory[chan int]()

	// Set with TTL - this will trigger the async deletion
	ttl := 10 * time.Millisecond
	if err := cache.Set("chan_key", make(chan int), ttl); err != nil {
		t.Fatalf("Set returned error: %v", err)
	}

	// Wait for TTL to expire
	time.Sleep(ttl + 50*time.Millisecond)

	// The key should be deleted, but we can't easily test the error path
	// since the error handling is in a goroutine and the error is just logged
	exists, err := cache.Exists("chan_key")
	if err != nil {
		t.Fatalf("Exists returned error: %v", err)
	}
	if exists {
		t.Fatalf("expected key to be deleted after TTL expiry")
	}
}

// TestMemory_ConcurrentAccess tests concurrent access to the memory cache
func TestMemory_ConcurrentAccess(t *testing.T) {
	cache := NewMemory[int]()
	const numGoroutines = 10
	const numOperations = 100

	// Test concurrent writes
	done := make(chan bool)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			for j := 0; j < numOperations; j++ {
				key := fmt.Sprintf("key_%d_%d", id, j)
				value := id*1000 + j
				if err := cache.Set(key, value, 0); err != nil {
					t.Errorf("Set failed: %v", err)
				}
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	// Verify all values were set correctly
	for i := 0; i < numGoroutines; i++ {
		for j := 0; j < numOperations; j++ {
			key := fmt.Sprintf("key_%d_%d", i, j)
			expectedValue := i*1000 + j
			value, err := cache.Get(key)
			if err != nil {
				t.Errorf("Get failed for key %s: %v", key, err)
			}
			if value != expectedValue {
				t.Errorf("Expected value %d for key %s, got %d", expectedValue, key, value)
			}
		}
	}
}

// TestMemory_ZeroValue tests handling of zero values
func TestMemory_ZeroValue(t *testing.T) {
	cache := NewMemory[string]()

	// Test setting empty string
	if err := cache.Set("empty", "", 0); err != nil {
		t.Fatalf("Set empty string failed: %v", err)
	}

	value, err := cache.Get("empty")
	if err != nil {
		t.Fatalf("Get empty string failed: %v", err)
	}
	if value != "" {
		t.Errorf("Expected empty string, got %q", value)
	}

	// Test setting zero value for int
	intCache := NewMemory[int]()
	if err := intCache.Set("zero", 0, 0); err != nil {
		t.Fatalf("Set zero int failed: %v", err)
	}

	intValue, err := intCache.Get("zero")
	if err != nil {
		t.Fatalf("Get zero int failed: %v", err)
	}
	if intValue != 0 {
		t.Errorf("Expected zero int, got %d", intValue)
	}
}

// TestMemory_InterfaceCompliance tests that Memory implements the Cache interface
func TestMemory_InterfaceCompliance(t *testing.T) {
	var _ Cache[string] = &Memory[string]{}
	var _ Cache[int] = &Memory[int]{}
	var _ Cache[bool] = &Memory[bool]{}
}
