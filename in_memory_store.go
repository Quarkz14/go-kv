package main

import (
	"errors"
	"sync"
)

type InMemoryKeyValueStore struct {
	store map[string]any
	mu    sync.Mutex
}

var KeyNotInMemory = errors.New("key doesn't exist")

func NewKeyValueStore() *InMemoryKeyValueStore {
	return &InMemoryKeyValueStore{
		map[string]any{},
		sync.Mutex{},
	}
}

func (i *InMemoryKeyValueStore) GetAllKeys() []string {
	i.mu.Lock()
	defer i.mu.Unlock()
	var keys []string
	for key := range i.store {
		keys = append(keys, key)
	}
	return keys
}

func (i *InMemoryKeyValueStore) GetValue(key string) (any, error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	value := i.store[key]
	if value == nil {
		return nil, KeyNotInMemory
	}
	return value, nil
}

func (i *InMemoryKeyValueStore) PutOrCreateValue(key string, value any) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.store[key] = value
}

func (i *InMemoryKeyValueStore) DeleteValue(key string) {
  i.mu.Lock()
  defer i.mu.Unlock()
  delete(i.store, key)
}
