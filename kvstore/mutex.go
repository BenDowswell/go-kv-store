//go:build !singlechan && !dualchan

 

package kvstore

 

import (

    "sync"

)

 

type mutexStore struct {

    mu    sync.RWMutex

    store map[string]string

}

 

func NewKVStore() Store {

    return &mutexStore{

        store: make(map[string]string),

    }

}

 

func (k *mutexStore) Get(key string) (string, bool) {

    k.mu.RLock()

    defer k.mu.RUnlock()

    value, ok := k.store[key]

    return value, ok

 

}

 

func (k *mutexStore) Set(key, value string) {

    k.mu.Lock()

    defer k.mu.Unlock()

    k.store[key] = value

}

 

func (k *mutexStore) Delete(key string) {

    k.mu.Lock()

    defer k.mu.Unlock()

    delete(k.store, key)

}