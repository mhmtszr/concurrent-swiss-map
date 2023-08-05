package csmap

import (
	"sync"

	"github.com/dolthub/maphash"
	"github.com/dolthub/swiss"
)

type CsMap[K comparable, V any] struct {
	options *Options
	hasher  maphash.Hasher[K]
	shards  []*shard[K, V]
}

type Options struct {
	shardCount uint64
}

type shard[K comparable, V any] struct {
	items *swiss.Map[K, V]
	sync.RWMutex
}

func Create[K comparable, V any](options ...func(options *Options)) *CsMap[K, V] {
	m := CsMap[K, V]{
		hasher: maphash.NewHasher[K](),
	}
	o := &Options{shardCount: 32}
	for _, option := range options {
		option(o)
	}
	m.options = o

	m.shards = make([]*shard[K, V], m.options.shardCount)

	for i := 0; i < int(m.options.shardCount); i++ {
		m.shards[i] = &shard[K, V]{items: swiss.NewMap[K, V](0)}
	}
	return &m
}

func WithShardCount(count int) func(options *Options) {
	return func(options *Options) {
		options.shardCount = uint64(count)
	}
}

func (m *CsMap[K, V]) getShard(key K) *shard[K, V] {
	u := m.hasher.Hash(key)
	return m.shards[u%m.options.shardCount]
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
