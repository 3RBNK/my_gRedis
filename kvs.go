package main

import (
	"sync"
	"time"
)

type CacheItem struct {
	value interface{}
	exp   *time.Time
}

type KVS struct {
	mu    sync.RWMutex
	items map[string]CacheItem
}

func NewKVS(cleanupTime time.Duration) *KVS {
	rc := &KVS{
		items: make(map[string]CacheItem),
	}

	go rc.StartGC(cleanupTime)

	return rc
}

func (rc *KVS) StartGC(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		rc.mu.Lock()
		now := time.Now()
		for k, v := range rc.items {
			if v.exp != nil && now.After(*v.exp) {
				delete(rc.items, k)
			}
		}
		rc.mu.Unlock()
	}
}

func (rc *KVS) Delete(key string) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	delete(rc.items, key)
}

func (rc *KVS) Get(key string) (interface{}, bool) {
	rc.mu.RLock()
	defer rc.mu.RUnlock()

	item, exists := rc.items[key]
	if !exists {
		return nil, false
	}

	if item.exp != nil && time.Now().After(*item.exp) {
		return nil, false
	}

	return item.value, true
}

func (rc *KVS) Put(key string, value interface{}, ttl time.Duration) {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	var exp *time.Time
	if ttl > 0 {
		t := time.Now().Add(ttl)
		exp = &t
	}

	rc.items[key] = CacheItem{
		value: value,
		exp:   exp,
	}
}
