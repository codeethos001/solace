package cache

import (
	"sync"
	"time"
)

type CacheEntry struct {
	Value      interface{}
	ExpiresAt  time.Time
	LastUpdate time.Time
}

type Cache struct {
	mu       sync.RWMutex
	data     map[string]*CacheEntry
	ttl      time.Duration
	maxSize  int
}

func New(ttl time.Duration) *Cache {
	c := &Cache{
		data:    make(map[string]*CacheEntry),
		ttl:     ttl,
		maxSize: 1000,
	}

	go c.cleanup()

	return c
}

func NewWithSize(ttl time.Duration, maxSize int) *Cache {
	c := &Cache{
		data:    make(map[string]*CacheEntry),
		ttl:     ttl,
		maxSize: maxSize,
	}

	go c.cleanup()

	return c
}

func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Simple eviction: remove oldest entry if cache is full
	if len(c.data) >= c.maxSize {
		var oldestKey string
		var oldestTime time.Time

		for k, entry := range c.data {
			if oldestTime.IsZero() || entry.LastUpdate.Before(oldestTime) {
				oldestTime = entry.LastUpdate
				oldestKey = k
			}
		}

		if oldestKey != "" {
			delete(c.data, oldestKey)
		}
	}

	now := time.Now()
	c.data[key] = &CacheEntry{
		Value:      value,
		ExpiresAt:  now.Add(c.ttl),
		LastUpdate: now,
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.data[key]
	if !exists {
		return nil, false
	}

	if time.Now().After(entry.ExpiresAt) {
		return nil, false // cleanup handles
	}

	return entry.Value, true
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}

func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = make(map[string]*CacheEntry)
}

func (c *Cache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.data)
}

func (c *Cache) Has(key string) bool {
	_, ok := c.Get(key)
	return ok
}

func (c *Cache) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()

		now := time.Now()
		for key, entry := range c.data {
			if now.After(entry.ExpiresAt) {
				delete(c.data, key)
			}
		}

		c.mu.Unlock()
	}
}

func (c *Cache) SetTTL(ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ttl = ttl
}

func (c *Cache) SetMaxSize(maxSize int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.maxSize = maxSize
}

type CacheStats struct {
	Size    int
	MaxSize int
	TTL     time.Duration
}

func (c *Cache) Stats() CacheStats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return CacheStats{
		Size:    len(c.data),
		MaxSize: c.maxSize,
		TTL:     c.ttl,
	}
}
