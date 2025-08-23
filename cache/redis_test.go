package cache

import (
	"fmt"
	"testing"
	"time"
	"webhook/codec"

	"github.com/alicebob/miniredis/v2"
)

func TestRedisStringValues(t *testing.T) {
	s := miniredis.RunT(t)

	redisCache, err := NewRedis[string](fmt.Sprintf("redis://%s", s.Addr()), codec.StringCodec{})
	if err != nil {
		t.Skipf("Redis not available: %v", err)
	}

	testKey := "test_string_key"
	testValue := "hello world"

	err = redisCache.Set(testKey, testValue, 1*time.Minute)
	if err != nil {
		t.Fatalf("Failed to set string value: %v", err)
	}

	value, err := redisCache.Get(testKey)
	if err != nil {
		t.Fatalf("Failed to get string value: %v", err)
	}

	if value != testValue {
		t.Errorf("Expected %s, got %s", testValue, value)
	}

	exists, err := redisCache.Exists(testKey)
	if err != nil {
		t.Fatalf("Failed to check existence: %v", err)
	}
	if !exists {
		t.Error("Key should exist")
	}

	err = redisCache.Delete(testKey)
	if err != nil {
		t.Fatalf("Failed to delete key: %v", err)
	}

	exists, err = redisCache.Exists(testKey)
	if err != nil {
		t.Fatalf("Failed to check existence after deletion: %v", err)
	}
	if exists {
		t.Error("Key should not exist after deletion")
	}
}
