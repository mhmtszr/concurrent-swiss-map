package csmap_test

import (
	"strconv"
	"sync"
	"testing"

	csmap "github.com/mhmtszr/concurrent-swiss-map"
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
	v, ok := myMap.Load(1)
	v2, ok2 := myMap.Load(2)
	if v != "test" || !ok {
		t.Fatal("1 should test")
	}
	if v2 != "" || ok2 {
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

func TestSetIfAbsent(t *testing.T) {
	myMap := csmap.Create[int, string]()
	myMap.SetIfAbsent(1, "test")
	if !myMap.Has(1) {
		t.Fatal("1 should be exist")
	}
}

func TestSetIfPresent(t *testing.T) {
	myMap := csmap.Create[int, string]()
	myMap.SetIfPresent(1, "test")
	if myMap.Has(1) {
		t.Fatal("1 should be not exist")
	}

	myMap.Store(1, "test")
	myMap.SetIfPresent(1, "new-test")
	val, _ := myMap.Load(1)
	if val != "new-test" {
		t.Fatal("val should be new-test")
	}
}

func TestSetIf(t *testing.T) {
	myMap := csmap.Create[int, string]()
	valueA := "value a"
	myMap.SetIf(1, func(previousVale string, previousFound bool) (value string, set bool) {
		// operate like  a SetIfAbsent...
		if !previousFound {
			return valueA, true
		}
		return "", false
	})
	value, _ := myMap.Load(1)
	if value != valueA {
		t.Fatal("value should value a")
	}

	myMap.SetIf(1, func(previousVale string, previousFound bool) (value string, set bool) {
		// operate like  a SetIfAbsent...
		if !previousFound {
			return "bad", true
		}
		return "", false
	})
	value, _ = myMap.Load(1)
	if value != valueA {
		t.Fatal("value should value a")
	}
}

func TestDeleteIf(t *testing.T) {
	myMap := csmap.Create[int, string]()
	myMap.Store(1, "value b")
	ok1 := myMap.DeleteIf(20, func(value string) bool {
		t.Fatal("condition function should not have been called")
		return false
	})
	if ok1 {
		t.Fatal("ok1 should be false")
	}

	ok2 := myMap.DeleteIf(1, func(value string) bool {
		if value != "value b" {
			t.Fatal("condition function arg should be tests")
		}
		return false // don't delete
	})
	if ok2 {
		t.Fatal("ok1 should be false")
	}

	ok3 := myMap.DeleteIf(1, func(value string) bool {
		if value != "value b" {
			t.Fatal("condition function arg should be tests")
		}
		return true // delete the entry
	})
	if !ok3 {
		t.Fatal("ok2 should be true")
	}
}

func TestCount(t *testing.T) {
	myMap := csmap.Create[int, string]()
	myMap.SetIfAbsent(1, "test")
	myMap.SetIfAbsent(2, "test2")
	if myMap.Count() != 2 {
		t.Fatal("count should be 2")
	}
}

func TestIsEmpty(t *testing.T) {
	myMap := csmap.Create[int, string]()
	if !myMap.IsEmpty() {
		t.Fatal("map should be empty")
	}
}

func TestRangeStop(t *testing.T) {
	myMap := csmap.Create[int, string](
		csmap.WithShardCount[int, string](1),
	)
	myMap.SetIfAbsent(1, "test")
	myMap.SetIfAbsent(2, "test2")
	myMap.SetIfAbsent(3, "test2")
	total := 0
	myMap.Range(func(key int, value string) (stop bool) {
		total++
		return true
	})
	if total != 1 {
		t.Fatal("total should be 1")
	}
}

func TestRange(t *testing.T) {
	myMap := csmap.Create[int, string]()
	myMap.SetIfAbsent(1, "test")
	myMap.SetIfAbsent(2, "test2")
	total := 0
	myMap.Range(func(key int, value string) (stop bool) {
		total++
		return
	})
	if total != 2 {
		t.Fatal("total should be 2")
	}
}

func TestCustomHasherWithRange(t *testing.T) {
	myMap := csmap.Create[int, string](
		csmap.WithCustomHasher[int, string](func(key int) uint64 {
			return 0
		}),
	)
	myMap.SetIfAbsent(1, "test")
	myMap.SetIfAbsent(2, "test2")
	myMap.SetIfAbsent(3, "test2")
	myMap.SetIfAbsent(4, "test2")
	total := 0
	myMap.Range(func(key int, value string) (stop bool) {
		total++
		return true
	})
	if total != 1 {
		t.Fatal("total should be 1, because currently range stops current shard only.")
	}
}

func TestBasicConcurrentWriteDeleteCount(t *testing.T) {
	myMap := csmap.Create[int, string](
		csmap.WithShardCount[int, string](32),
		csmap.WithSize[int, string](1000),
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

func TestClear(t *testing.T) {
	myMap := csmap.Create[int, string]()
	loop := 10000
	for i := 0; i < loop; i++ {
		myMap.Store(i, "test")
	}

	myMap.Clear()

	if !myMap.IsEmpty() {
		t.Fatal("count should be true")
	}

	// store again
	for i := 0; i < loop; i++ {
		myMap.Store(i, "test")
	}

	// get again
	for i := 0; i < loop; i++ {
		val, ok := myMap.Load(i)
		if ok != true {
			t.Fatal("ok should be true")
		}

		if val != "test" {
			t.Fatal("val should be test")
		}
	}

	// check again
	count := myMap.Count()
	if count != loop {
		t.Fatal("count should be 1000")
	}
}
