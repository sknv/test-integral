package integral

import (
	"context"
	"sync"
	"time"
)

const defaultTtl = 30 * time.Second

type (
	KVStore interface {
		Get(key interface{}) (interface{}, bool)
		Put(key, value interface{})
		Remove(key interface{})
	}

	KVMemoryStoreValue struct {
		inner interface{}
		read  chan struct{}
	}

	KVMemoryStore struct {
		ttl   time.Duration
		store *sync.Map
	}
)

func NewKVStore() KVStore {
	return NewKVStoreWithTtl(defaultTtl)
}

func NewKVStoreWithTtl(ttl time.Duration) KVStore {
	return &KVMemoryStore{
		ttl:   ttl,
		store: new(sync.Map),
	}
}

func (s *KVMemoryStore) Get(key interface{}) (interface{}, bool) {
	storeValue, ok := s.store.Load(key)
	if !ok {
		return nil, false
	}

	s.Remove(key)
	value := storeValue.(*KVMemoryStoreValue)
	value.read <- struct{}{}
	return value.inner, true
}

func (s *KVMemoryStore) Put(key, value interface{}) {
	storeValue := &KVMemoryStoreValue{
		inner: value,
		read:  make(chan struct{}),
	}
	s.store.Store(key, storeValue)

	go s.removeValueAfterTimeout(key, storeValue.read)
}

func (s *KVMemoryStore) Remove(key interface{}) {
	s.store.Delete(key)
}

func (s *KVMemoryStore) removeValueAfterTimeout(key interface{}, valueRead <-chan struct{}) {
	ctx, cancel := context.WithTimeout(context.Background(), s.ttl)
	defer cancel()

	select {
	case <-valueRead:
		return
	case <-ctx.Done():
		s.Remove(key)
	}
}
