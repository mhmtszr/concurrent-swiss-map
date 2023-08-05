package csmap_test

import (
	csmap "concurrent-swiss-map"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHas(t *testing.T) {
	myMap := csmap.Create[int, string]()
	myMap.Store(1, "test")
	assert.Equal(t, myMap.Has(1), true)
}

func TestLoad(t *testing.T) {
	myMap := csmap.Create[int, string]()
	myMap.Store(1, "test")
	assert.Equal(t, *myMap.Load(1), "test")
	assert.Nil(t, myMap.Load(2))
}

func TestDelete(t *testing.T) {
	myMap := csmap.Create[int, string]()
	myMap.Store(1, "test")
	ok1 := myMap.Delete(20)
	ok2 := myMap.Delete(1)
	assert.Equal(t, myMap.Has(1), false)
	assert.Equal(t, ok1, false)
	assert.Equal(t, ok2, true)
}

func TestBasicConcurrentWriteDeleteCount(t *testing.T) {
	myMap := csmap.Create[int, string](
		csmap.WithShardCount(5),
	)

	for i := 0; i < 1000000; i++ {
		i := i
		go func() {
			myMap.Store(i, strconv.Itoa(i))
		}()
	}

	for i := 0; i < 1000000; i++ {
		i := i
		go func() {
			assert.Equal(t, myMap.Has(i), true)
		}()
	}

	for i := 0; i < 1000000; i++ {
		i := i
		go func() {
			myMap.Delete(i)
		}()
	}

	for i := 0; i < 1000000; i++ {
		i := i
		go func() {
			assert.Equal(t, myMap.Has(i), false)
		}()
	}
}
