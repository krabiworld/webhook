package cache

import (
	"errors"
	"fmt"
	"testing"
	"time"
	"webhook/codec"

	"github.com/alicebob/miniredis/v2"
)

// failingCodec implements codec.Codec[string] and always fails to test error paths
type failingCodec struct{}

func (f failingCodec) Encode(_ string) (string, error) { return "", errors.New("encode failed") }
func (f failingCodec) Decode(_ string) (string, error) { return "", errors.New("decode failed") }

func TestRedis_String_SetGetExistsDelete(t *testing.T) {
	s := miniredis.RunT(t)

	redisCache, err := NewRedis[string](fmt.Sprintf("redis://%s", s.Addr()), codec.StringCodec{})
	if err != nil {
		t.Skipf("Redis not available: %v", err)
	}

	testKey := "test_string_key"
	testValue := "hello world"

	if err := redisCache.Set(testKey, testValue, 1*time.Minute); err != nil {
		t.Fatalf("failed to set string value: %v", err)
	}

	value, err := redisCache.Get(testKey)
	if err != nil {
		t.Fatalf("failed to get string value: %v", err)
	}
	if value != testValue {
		t.Errorf("expected %s, got %s", testValue, value)
	}

	exists, err := redisCache.Exists(testKey)
	if err != nil {
		t.Fatalf("failed to check existence: %v", err)
	}
	if !exists {
		t.Error("key should exist")
	}

	if err := redisCache.Delete(testKey); err != nil {
		t.Fatalf("failed to delete key: %v", err)
	}

	exists, err = redisCache.Exists(testKey)
	if err != nil {
		t.Fatalf("failed to check existence after deletion: %v", err)
	}
	if exists {
		t.Error("key should not exist after deletion")
	}
}

func TestRedis_TTLExpiration(t *testing.T) {
	s := miniredis.RunT(t)

	redisCache, err := NewRedis[string](fmt.Sprintf("redis://%s", s.Addr()), codec.StringCodec{})
	if err != nil {
		t.Skipf("Redis not available: %v", err)
	}

	key := "ttl_key"
	val := "soon gone"
	if err := redisCache.Set(key, val, 50*time.Millisecond); err != nil {
		t.Fatalf("set failed: %v", err)
	}

	// Immediately should exist
	if exists, err := redisCache.Exists(key); err != nil || !exists {
		if err != nil {
			t.Fatalf("exists check failed: %v", err)
		}
		t.Fatal("key should exist before TTL expires")
	}

	// Advance miniredis's virtual clock beyond TTL
	s.FastForward(60 * time.Millisecond)

	// After TTL, Get should return zero value with nil error
	got, err := redisCache.Get(key)
	if err != nil {
		t.Fatalf("get after ttl failed: %v", err)
	}
	if got != "" {
		t.Errorf("expected zero value after expiration, got %q", got)
	}

	exists, err := redisCache.Exists(key)
	if err != nil {
		t.Fatalf("exists after ttl failed: %v", err)
	}
	if exists {
		t.Error("key should not exist after TTL expiration")
	}
}

func TestRedis_GetNonExistentReturnsZeroValue(t *testing.T) {
	s := miniredis.RunT(t)

	redisCache, err := NewRedis[string](fmt.Sprintf("redis://%s", s.Addr()), codec.StringCodec{})
	if err != nil {
		t.Skipf("Redis not available: %v", err)
	}

	val, err := redisCache.Get("does_not_exist")
	if err != nil {
		t.Fatalf("get on non-existent key returned error: %v", err)
	}
	if val != "" {
		t.Errorf("expected zero value for non-existent key, got %q", val)
	}
}

func TestRedis_CodecErrors(t *testing.T) {
	s := miniredis.RunT(t)

	// Use a codec that fails both directions
	redisCache, err := NewRedis[string](fmt.Sprintf("redis://%s", s.Addr()), failingCodec{})
	if err != nil {
		t.Skipf("Redis not available: %v", err)
	}

	if err := redisCache.Set("k", "v", time.Minute); err == nil {
		t.Fatal("expected encode failure on Set, got nil error")
	}

	// Put a raw value into Redis to force decode path
	if err := s.Set("k2", "raw"); err != nil {
		t.Fatal("expected encode failure on Set, got nil error")
	}
	if _, err := redisCache.Get("k2"); err == nil {
		t.Fatal("expected decode failure on Get, got nil error")
	}
}

func TestRedis_Flush(t *testing.T) {
	s := miniredis.RunT(t)

	redisCache, err := NewRedis[string](fmt.Sprintf("redis://%s", s.Addr()), codec.StringCodec{})
	if err != nil {
		t.Skipf("Redis not available: %v", err)
	}

	if err := redisCache.Set("a", "1", time.Minute); err != nil {
		t.Fatalf("set failed: %v", err)
	}
	if err := redisCache.Set("b", "2", time.Minute); err != nil {
		t.Fatalf("set failed: %v", err)
	}

	if err := redisCache.Flush(); err != nil {
		t.Fatalf("flush failed: %v", err)
	}

	existsA, err := redisCache.Exists("a")
	if err != nil {
		t.Fatalf("exists a failed: %v", err)
	}
	if existsA {
		t.Error("expected key 'a' not to exist after flush")
	}

	existsB, err := redisCache.Exists("b")
	if err != nil {
		t.Fatalf("exists b failed: %v", err)
	}
	if existsB {
		t.Error("expected key 'b' not to exist after flush")
	}
}
