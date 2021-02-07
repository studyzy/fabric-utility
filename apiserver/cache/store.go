package cache

import (
	"sync"

	"github.com/golang/groupcache/lru"
)

// Store represts a cache instance.
type Store interface {
	Get(key interface{}) (interface{}, bool)
	Set(key, value interface{})
	Size() int
	SetNotExist(key interface{}, newValGenerator func() interface{}) (interface{}, bool)
}

type lruStore struct {
	mu    sync.Mutex
	cache *lru.Cache
}

// NewLRU creates a LRU memory cache, which is safe for concurrent access.
func NewLRU(size int) Store {
	return &lruStore{
		cache: lru.New(size),
	}
}

// NewLRUWithOnEvictedCallback creates a LRU memory cache with a OnEvicted callback
func NewLRUWithOnEvictedCallback(size int, onEvicted func(key, value interface{})) Store {
	s := &lruStore{
		cache: lru.New(size),
	}

	s.cache.OnEvicted = func(key lru.Key, value interface{}) { onEvicted(key, value) }

	return s
}

func (l *lruStore) Get(key interface{}) (value interface{}, ok bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.cache.Get(key)
}

func (l *lruStore) Set(key, value interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.cache.Add(key, value)
}

func (l *lruStore) Size() int {
	return l.cache.Len()
}

func (l *lruStore) SetNotExist(key interface{}, newValGenerator func() interface{}) (interface{}, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	cached, ok := l.cache.Get(key)
	if ok {
		return cached, false
	}

	newVal := newValGenerator()
	if newVal == nil {
		return nil, false
	}
	l.cache.Add(key, newVal)
	return newVal, true
}
