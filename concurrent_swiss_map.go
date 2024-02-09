package csmap

import (
	"sync"

	"github.com/mhmtszr/concurrent-swiss-map/maphash"

	"github.com/mhmtszr/concurrent-swiss-map/swiss"
)

type CsMap[K comparable, V any] struct {
	hasher     func(key K) uint64
	shards     []shard[K, V]
	shardCount uint64
	size       uint64
}

type shard[K comparable, V any] struct {
	items *swiss.Map[K, V]
	*sync.RWMutex
}

func Create[K comparable, V any](options ...func(options *CsMap[K, V])) *CsMap[K, V] {
	m := CsMap[K, V]{
		hasher:     maphash.NewHasher[K]().Hash,
		shardCount: 32,
	}
	for _, option := range options {
		option(&m)
	}

	m.shards = make([]shard[K, V], m.shardCount)

	for i := 0; i < int(m.shardCount); i++ {
		m.shards[i] = shard[K, V]{items: swiss.NewMap[K, V](uint32((m.size / m.shardCount) + 1)), RWMutex: &sync.RWMutex{}}
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

func WithSize[K comparable, V any](size uint64) func(csMap *CsMap[K, V]) {
	return func(csMap *CsMap[K, V]) {
		csMap.size = size
	}
}

func (m *CsMap[K, V]) getShard(key K) HashShardPair[K, V] {
	u := m.hasher(key)
	return HashShardPair[K, V]{
		hash:  u,
		shard: m.shards[u%m.shardCount],
	}
}

func (m *CsMap[K, V]) Store(key K, value V) {
	hashShardPair := m.getShard(key)
	shard := hashShardPair.shard
	shard.Lock()
	shard.items.PutWithHash(key, value, hashShardPair.hash)
	shard.Unlock()
}

func (m *CsMap[K, V]) Delete(key K) bool {
	hashShardPair := m.getShard(key)
	shard := hashShardPair.shard
	shard.Lock()
	defer shard.Unlock()
	return shard.items.DeleteWithHash(key, hashShardPair.hash)
}

func (m *CsMap[K, V]) DeleteIf(key K, condition func(value V) bool) bool {
	hashShardPair := m.getShard(key)
	shard := hashShardPair.shard
	shard.Lock()
	defer shard.Unlock()
	value, ok := shard.items.GetWithHash(key, hashShardPair.hash)
	if ok && condition(value) {
		return shard.items.DeleteWithHash(key, hashShardPair.hash)
	}
	return false
}

func (m *CsMap[K, V]) Load(key K) (V, bool) {
	hashShardPair := m.getShard(key)
	shard := hashShardPair.shard
	shard.RLock()
	defer shard.RUnlock()
	return shard.items.GetWithHash(key, hashShardPair.hash)
}

func (m *CsMap[K, V]) Has(key K) bool {
	hashShardPair := m.getShard(key)
	shard := hashShardPair.shard
	shard.RLock()
	defer shard.RUnlock()
	return shard.items.HasWithHash(key, hashShardPair.hash)
}

func (m *CsMap[K, V]) Clear() {
	for i := range m.shards {
		shard := m.shards[i]

		shard.Lock()
		shard.items.Clear()
		shard.Unlock()
	}
}

func (m *CsMap[K, V]) Count() int {
	count := 0
	for i := range m.shards {
		shard := m.shards[i]
		shard.RLock()
		count += shard.items.Count()
		shard.RUnlock()
	}
	return count
}

func (m *CsMap[K, V]) SetIfAbsent(key K, value V) {
	hashShardPair := m.getShard(key)
	shard := hashShardPair.shard
	shard.Lock()
	_, ok := shard.items.GetWithHash(key, hashShardPair.hash)
	if !ok {
		shard.items.PutWithHash(key, value, hashShardPair.hash)
	}
	shard.Unlock()
}

func (m *CsMap[K, V]) SetIf(key K, conditionFn func(previousVale V, previousFound bool) (value V, set bool)) {
	hashShardPair := m.getShard(key)
	shard := hashShardPair.shard
	shard.Lock()
	value, found := shard.items.GetWithHash(key, hashShardPair.hash)
	value, ok := conditionFn(value, found)
	if ok {
		shard.items.PutWithHash(key, value, hashShardPair.hash)
	}
	shard.Unlock()
}

func (m *CsMap[K, V]) SetIfPresent(key K, value V) {
	hashShardPair := m.getShard(key)
	shard := hashShardPair.shard
	shard.Lock()
	_, ok := shard.items.GetWithHash(key, hashShardPair.hash)
	if ok {
		shard.items.PutWithHash(key, value, hashShardPair.hash)
	}
	shard.Unlock()
}

func (m *CsMap[K, V]) IsEmpty() bool {
	return m.Count() == 0
}

// Range If the callback function returns true iteration will stop.
func (m *CsMap[K, V]) Range(f func(key K, value V) (stop bool)) {
	for i := range m.shards {
		shard := m.shards[i]
		shard.RLock()
		stop := shard.items.Iter(f)
		if stop {
			shard.RUnlock()
			return
		}
		shard.RUnlock()
	}
}

type HashShardPair[K comparable, V any] struct {
	shard shard[K, V]
	hash  uint64
}
