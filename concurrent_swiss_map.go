package csmap

import (
	"sync"

	"github.com/dolthub/swiss"
)

type CsMap[K comparable, V any] struct {
	hasher     func(key K) uint64
	shards     []*shard[K, V]
	shardCount uint64
}

type shard[K comparable, V any] struct {
	items *swiss.Map[K, V]
	sync.RWMutex
}

func Create[K comparable, V any](options ...func(options *CsMap[K, V])) *CsMap[K, V] {
	defaultHasher := NewDefaultHasher[K]()

	m := CsMap[K, V]{
		hasher:     defaultHasher.h.Hash,
		shardCount: 32,
	}
	for _, option := range options {
		option(&m)
	}

	m.shards = make([]*shard[K, V], m.shardCount)

	for i := 0; i < int(m.shardCount); i++ {
		m.shards[i] = &shard[K, V]{items: swiss.NewMap[K, V](0)}
	}
	return &m
}

func WithShardCount[K comparable, V any](count uint64) func(csMap *CsMap[K, V]) {
	return func(csMap *CsMap[K, V]) {
		csMap.shardCount = count
	}
}

func WithCustomHasher[K comparable, V any](h func(key K) uint64) func(csMap *CsMap[K, V]) {
	return func(csMap *CsMap[K, V]) {
		csMap.hasher = h
	}
}

func (m *CsMap[K, V]) getShard(key K) *shard[K, V] {
	u := m.hasher(key)
	return m.shards[u%m.shardCount]
}

func (m *CsMap[K, V]) Store(key K, value V) {
	shard := m.getShard(key)
	shard.Lock()
	defer shard.Unlock()
	shard.items.Put(key, value)
}

func (m *CsMap[K, V]) Delete(key K) bool {
	shard := m.getShard(key)
	shard.Lock()
	defer shard.Unlock()
	return shard.items.Delete(key)
}

func (m *CsMap[K, V]) Load(key K) *V {
	shard := m.getShard(key)
	shard.RLock()
	defer shard.RUnlock()
	value, ok := shard.items.Get(key)
	if !ok {
		return nil
	}
	return &value
}

func (m *CsMap[K, V]) Has(key K) bool {
	shard := m.getShard(key)
	shard.RLock()
	defer shard.RUnlock()
	return shard.items.Has(key)
}

func (m *CsMap[K, V]) Count() int {
	count := 0
	for i := 0; i < len(m.shards); i++ {
		shard := m.shards[i]
		shard.RLock()
		count += shard.items.Count()
		shard.RUnlock()
	}
	return count
}

func (m *CsMap[K, V]) SetIfAbsent(key K, value V) {
	shard := m.getShard(key)
	shard.Lock()
	defer shard.Unlock()
	_, ok := shard.items.Get(key)
	if !ok {
		shard.items.Put(key, value)
	}
}

func (m *CsMap[K, V]) IsEmpty() bool {
	return m.Count() == 0
}

// Range If the callback function returns true iteration will stop.
// TODO: currently it only stops the current shard
func (m *CsMap[K, V]) Range(f func(key K, value V) (stop bool)) {
	for i := 0; i < len(m.shards); i++ {
		shard := m.shards[i]
		shard.RLock()
		shard.items.Iter(f)
		shard.RUnlock()
	}
}
