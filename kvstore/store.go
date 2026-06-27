package kvstore

import (
	"sync"
)

type KVStore struct {
	mu    sync.RWMutex
	store map[string]string
}

func NewKVStore() *KVStore {
	return &KVStore{
		store: make(map[string]string),
	}
}

func (k *KVStore) Get(key string) (string, bool) {
	k.mu.RLock()
	defer k.mu.RUnlock()
	value, ok := k.store[key]
	return value, ok

}

func (k *KVStore) Set(key, value string) {
	k.mu.Lock()
	defer k.mu.Unlock()
	k.store[key] = value
}

func (k *KVStore) Delete(key string) {
	k.mu.Lock()
	defer k.mu.Unlock()
	delete(k.store, key)
}
