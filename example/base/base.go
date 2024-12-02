package main

import (
	"hash/fnv"

	csmap "github.com/mhmtszr/concurrent-swiss-map"
)

func main() {
	myMap := csmap.New[string, int](
		// set the number of map shards. the default value is 32.
		csmap.WithShardCount[string, int](32),

		// if don't set custom hasher, use the built-in maphash.
		csmap.WithCustomHasher[string, int](func(key string) uint64 {
			hash := fnv.New64a()
			hash.Write([]byte(key))
			return hash.Sum64()
		}),

		// set the total capacity, every shard map has total capacity/shard count capacity. the default value is 0.
		csmap.WithSize[string, int](1000),
	)

	key := "swiss-map"
	myMap.Store(key, 10)

	val, ok := myMap.Load(key)
	println("load val:", val, "exists:", ok)

	deleted := myMap.Delete(key)
	println("deleted:", deleted)

	ok = myMap.Has(key)
	println("has:", ok)

	empty := myMap.IsEmpty()
	println("empty:", empty)

	myMap.SetIfAbsent(key, 11)

	myMap.Range(func(key string, value int) (stop bool) {
		println("range:", key, value)
		return true
	})

	count := myMap.Count()
	println("count:", count)

	// Output:
	// load val: 10 exists: true
	// deleted: true
	// has: false
	// empty: true
	// range: swiss-map 11
	// count: 1
}
