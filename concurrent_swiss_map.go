package csmap

import (
	"context"
	"encoding/json"
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

// StoreWithCompute stores the value of the current key using the result of the
// compute function. If the key doesn't exist, it is stored with default value
// of the V's type.
func (m *CsMap[K, V]) StoreWithCompute(key K, compute func(value V) V) {
	hashShardPair := m.getShard(key)
	shard := hashShardPair.shard
	shard.Lock()
	defer shard.Unlock()
	value, _ := shard.items.GetWithHash(key, hashShardPair.hash)
	shard.items.PutWithHash(key, compute(value), hashShardPair.hash)
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

type Tuple[K comparable, V any] struct {
	Key K
	Val V
}

// Range If the callback function returns true iteration will stop.
func (m *CsMap[K, V]) Range(f func(key K, value V) (stop bool)) {
	ch := make(chan Tuple[K, V], m.Count())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	listenCompleted := m.listen(f, ch)
	m.produce(ctx, ch)
	listenCompleted.Wait()
}

func (m *CsMap[K, V]) MarshalJSON() ([]byte, error) {
	tmp := make(map[K]V, m.Count())
	m.Range(func(key K, value V) (stop bool) {
		tmp[key] = value
		return false
	})
	return json.Marshal(tmp)
}

func (m *CsMap[K, V]) UnmarshalJSON(b []byte) error {
	tmp := make(map[K]V, m.Count())

	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}

	for key, val := range tmp {
		m.Store(key, val)
	}
	return nil
}

func (m *CsMap[K, V]) produce(ctx context.Context, ch chan Tuple[K, V]) {
	var wg sync.WaitGroup
	wg.Add(len(m.shards))
	for i := range m.shards {
		go func(i int) {
			defer wg.Done()

			shard := m.shards[i]
			shard.RLock()
			shard.items.Iter(func(k K, v V) (stop bool) {
				select {
				case <-ctx.Done():
					return true
				default:
					ch <- Tuple[K, V]{Key: k, Val: v}
				}
				return false
			})
			shard.RUnlock()
		}(i)
	}
	go func() {
		wg.Wait()
		close(ch)
	}()
}

func (m *CsMap[K, V]) listen(f func(key K, value V) (stop bool), ch chan Tuple[K, V]) *sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for t := range ch {
			if stop := f(t.Key, t.Val); stop {
				return
			}
		}
	}()
	return &wg
}

type HashShardPair[K comparable, V any] struct {
	shard shard[K, V]
	hash  uint64
}
