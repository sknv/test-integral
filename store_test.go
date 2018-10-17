package integral

import (
	"testing"
	"time"
)

func TestPutAndGet(t *testing.T) {
	store := NewKVStore()
	key, val := "foo", "bar"
	store.Put(key, val)

	value, ok := store.Get(key)
	if !ok {
		t.Error("Value must present")
	}

	if val != value {
		t.Errorf("Expected %s, got %s", val, value)
	}
}

func TestRemove(t *testing.T) {
	store := NewKVStore()
	key, val := "foo", "bar"
	store.Put(key, val)

	store.Remove(key)
	_, ok := store.Get(key)
	if ok {
		t.Error("Value must not present")
	}
}

func TestValueShouldBeRemovedAfterGet(t *testing.T) {
	store := NewKVStore()
	key, val := "foo", "bar"
	store.Put(key, val)

	_, ok := store.Get(key)
	if !ok {
		t.Error("Value must present")
	}

	_, ok = store.Get(key)
	if ok {
		t.Error("Value must not present")
	}
}

func TestValueShouldPresentBeforeTimeout(t *testing.T) {
	store := NewKVStoreWithTtl(4 * time.Second)
	key, val := "foo", "bar"
	store.Put(key, val)

	time.Sleep(2 * time.Second)
	_, ok := store.Get(key)
	if !ok {
		t.Error("Value must present before the timeout exceeds")
	}
}

func TestValueShouldBeRemovedAfterTimeout(t *testing.T) {
	store := NewKVStoreWithTtl(2 * time.Second)
	key, val := "foo", "bar"
	store.Put(key, val)

	time.Sleep(4 * time.Second)
	_, ok := store.Get(key)
	if ok {
		t.Error("Value must be removed after the timeout exceeds")
	}
}

func TestPutAndGetFromMultipleGoroutines(t *testing.T) {
	store := NewKVStore()
	for i := 0; i < 100; i++ {
		go func(i int) {
			store.Put(i, i+1)
			time.Sleep(time.Second)

			value, ok := store.Get(i)
			if !ok {
				t.Error("Value must present")
			}

			intVal := value.(int)
			if i+1 != value {
				t.Errorf("Expected %d, got %d", i+1, intVal)
			}
		}(i)
	}
	time.Sleep(5 * time.Second)
}
