package cache

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
	"webhook/internal/codec"

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

// TestRedis_Get_RedisError tests Get when Redis returns an error other than redis.Nil
func TestRedis_Get_RedisError(t *testing.T) {
	s := miniredis.RunT(t)

	redisCache, err := NewRedis[string](fmt.Sprintf("redis://%s", s.Addr()), codec.StringCodec{})
	if err != nil {
		t.Skipf("Redis not available: %v", err)
	}

	// Close the miniredis instance to force an error
	s.Close()

	_, err = redisCache.Get("test_key")
	if err == nil {
		t.Fatal("expected error when Redis is closed, got nil")
	}
	// The exact error might vary, so just check it contains the right parts
	if !strings.Contains(err.Error(), "dial tcp") {
		t.Errorf("unexpected error message: %v", err)
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

// TestRedis_NewRedis_InvalidURL tests NewRedis with invalid URL
func TestRedis_NewRedis_InvalidURL(t *testing.T) {
	_, err := NewRedis[string]("invalid://url", codec.StringCodec{})
	if err == nil {
		t.Fatal("expected error for invalid URL, got nil")
	}
	// The exact error message might vary, so just check it contains the right parts
	if !strings.Contains(err.Error(), "failed to parse redis url") {
		t.Errorf("unexpected error message: %v", err)
	}
}

// TestRedis_NewRedis_ConnectionFailure tests NewRedis with unreachable Redis
func TestRedis_NewRedis_ConnectionFailure(t *testing.T) {
	_, err := NewRedis[string]("redis://localhost:9999", codec.StringCodec{})
	if err == nil {
		t.Fatal("expected error for unreachable Redis, got nil")
	}
	if err.Error() != "failed to connect to redis: dial tcp 127.0.0.1:9999: connect: connection refused" {
		// The exact error message might vary by OS, so just check it contains the right parts
		if !errors.Is(err, errors.New("failed to connect to redis")) &&
			!errors.Is(err, errors.New("connection refused")) {
			t.Errorf("unexpected error message: %v", err)
		}
	}
}

// TestRedis_Delete_NonExistentKey tests Delete with non-existent key
func TestRedis_Delete_NonExistentKey(t *testing.T) {
	s := miniredis.RunT(t)

	redisCache, err := NewRedis[string](fmt.Sprintf("redis://%s", s.Addr()), codec.StringCodec{})
	if err != nil {
		t.Skipf("Redis not available: %v", err)
	}

	// Delete non-existent key should not error
	if err := redisCache.Delete("non_existent_key"); err != nil {
		t.Fatalf("delete non-existent key failed: %v", err)
	}
}

// TestRedis_Delete_ExistingKey tests Delete with existing key
func TestRedis_Delete_ExistingKey(t *testing.T) {
	s := miniredis.RunT(t)

	redisCache, err := NewRedis[string](fmt.Sprintf("redis://%s", s.Addr()), codec.StringCodec{})
	if err != nil {
		t.Skipf("Redis not available: %v", err)
	}

	// Set a key first
	if err := redisCache.Set("existing_key", "value", time.Minute); err != nil {
		t.Fatalf("set failed: %v", err)
	}

	// Delete existing key should not error
	if err := redisCache.Delete("existing_key"); err != nil {
		t.Fatalf("delete existing key failed: %v", err)
	}

	// Verify key was deleted
	exists, err := redisCache.Exists("existing_key")
	if err != nil {
		t.Fatalf("exists check failed: %v", err)
	}
	if exists {
		t.Error("key should not exist after deletion")
	}
}

// TestRedis_Delete_MultipleKeys tests Delete with multiple keys
func TestRedis_Delete_MultipleKeys(t *testing.T) {
	s := miniredis.RunT(t)

	redisCache, err := NewRedis[string](fmt.Sprintf("redis://%s", s.Addr()), codec.StringCodec{})
	if err != nil {
		t.Skipf("Redis not available: %v", err)
	}

	// Set multiple keys
	keys := []string{"key1", "key2", "key3"}
	for _, key := range keys {
		if err := redisCache.Set(key, "value", time.Minute); err != nil {
			t.Fatalf("set %s failed: %v", key, err)
		}
	}

	// Delete all keys
	for _, key := range keys {
		if err := redisCache.Delete(key); err != nil {
			t.Fatalf("delete %s failed: %v", key, err)
		}
	}

	// Verify all keys were deleted
	for _, key := range keys {
		exists, err := redisCache.Exists(key)
		if err != nil {
			t.Fatalf("exists check for %s failed: %v", key, err)
		}
		if exists {
			t.Errorf("key %s should not exist after deletion", key)
		}
	}
}

// TestRedis_Delete_ResultZero tests Delete when Redis returns result 0
func TestRedis_Delete_ResultZero(t *testing.T) {
	s := miniredis.RunT(t)

	redisCache, err := NewRedis[string](fmt.Sprintf("redis://%s", s.Addr()), codec.StringCodec{})
	if err != nil {
		t.Skipf("Redis not available: %v", err)
	}

	// Try to delete a key that doesn't exist
	// This should trigger the "Key not found for deletion" log message
	if err := redisCache.Delete("non_existent_key_for_logging"); err != nil {
		t.Fatalf("delete non-existent key failed: %v", err)
	}

	// The key should not exist
	exists, err := redisCache.Exists("non_existent_key_for_logging")
	if err != nil {
		t.Fatalf("exists check failed: %v", err)
	}
	if exists {
		t.Error("key should not exist")
	}
}

// TestRedis_Exists_Error tests Exists when Redis returns an error
func TestRedis_Exists_Error(t *testing.T) {
	s := miniredis.RunT(t)

	redisCache, err := NewRedis[string](fmt.Sprintf("redis://%s", s.Addr()), codec.StringCodec{})
	if err != nil {
		t.Skipf("Redis not available: %v", err)
	}

	// Close the miniredis instance to force an error
	s.Close()

	_, err = redisCache.Exists("test_key")
	if err == nil {
		t.Fatal("expected error when Redis is closed, got nil")
	}
	// The exact error might vary, so just check it contains the right parts
	if !strings.Contains(err.Error(), "failed to check existence of key test_key") {
		t.Errorf("unexpected error message: %v", err)
	}
}

// TestRedis_Flush_Error tests Flush when Redis returns an error
func TestRedis_Flush_Error(t *testing.T) {
	s := miniredis.RunT(t)

	redisCache, err := NewRedis[string](fmt.Sprintf("redis://%s", s.Addr()), codec.StringCodec{})
	if err != nil {
		t.Skipf("Redis not available: %v", err)
	}

	// Close the miniredis instance to force an error
	s.Close()

	err = redisCache.Flush()
	if err == nil {
		t.Fatal("expected error when Redis is closed, got nil")
	}
	// The exact error might vary, so just check it contains the right parts
	if !strings.Contains(err.Error(), "failed to flush Redis") {
		t.Errorf("unexpected error message: %v", err)
	}
}

// TestRedis_InterfaceCompliance tests that Redis implements the Cache interface
func TestRedis_InterfaceCompliance(t *testing.T) {
	// We can't create a Redis instance without a real connection for this test,
	// but we can verify the type assertion compiles
	var _ Cache[string] = (*Redis[string])(nil)
	var _ Cache[int] = (*Redis[int])(nil)
	var _ Cache[bool] = (*Redis[bool])(nil)
}
