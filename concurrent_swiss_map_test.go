package csmap_test

import (
	"strconv"
	"sync"
	"testing"

	"github.com/mhmtszr/concurrent-swiss-map"
)

func TestHas(t *testing.T) {
	myMap := csmap.Create[int, string]()
	myMap.Store(1, "test")
	if !myMap.Has(1) {
		t.Fatal("1 should exists")
	}
}

func TestLoad(t *testing.T) {
	myMap := csmap.Create[int, string]()
	myMap.Store(1, "test")
	if *myMap.Load(1) != "test" {
		t.Fatal("1 should test")
	}
	if myMap.Load(2) != nil {
		t.Fatal("2 should not exist")
	}
}

func TestDelete(t *testing.T) {
	myMap := csmap.Create[int, string]()
	myMap.Store(1, "test")
	ok1 := myMap.Delete(20)
	ok2 := myMap.Delete(1)
	if myMap.Has(1) {
		t.Fatal("1 should be deleted")
	}
	if ok1 {
		t.Fatal("ok1 should be false")
	}
	if !ok2 {
		t.Fatal("ok2 should be true")
	}
}

func TestBasicConcurrentWriteDeleteCount(t *testing.T) {
	myMap := csmap.Create[int, string](
		csmap.WithShardCount(5),
	)

	var wg sync.WaitGroup
	wg.Add(1000000)
	for i := 0; i < 1000000; i++ {
		i := i
		go func() {
			defer wg.Done()
			myMap.Store(i, strconv.Itoa(i))
		}()
	}
	wg.Wait()
	wg.Add(1000000)
	for i := 0; i < 1000000; i++ {
		i := i
		go func() {
			defer wg.Done()
			if !myMap.Has(i) {
				t.Error(strconv.Itoa(i) + " should exist")
				return
			}
		}()
	}

	wg.Wait()
	wg.Add(1000000)

	for i := 0; i < 1000000; i++ {
		i := i
		go func() {
			defer wg.Done()
			myMap.Delete(i)
		}()
	}

	wg.Wait()
	wg.Add(1000000)

	for i := 0; i < 1000000; i++ {
		i := i
		go func() {
			defer wg.Done()
			if myMap.Has(i) {
				t.Error(strconv.Itoa(i) + " should not exist")
				return
			}
		}()
	}

	wg.Wait()
}
