package csmap_test

// import (
//	"fmt"
//	"runtime"
//	"strconv"
//	"sync"
//	"testing"
//
//	"github.com/mhmtszr/concurrent-swiss-map"
//)
//

// var table = []struct {
//	total    int
//	deletion int
// }{
//	{
//		total:    100,
//		deletion: 100,
//	},
//	{
//		total:    5000000,
//		deletion: 5000000,
//	},
//}

// func PrintMemUsage() {
//	var m runtime.MemStats
//	runtime.ReadMemStats(&m)
//	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
//	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
//	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
//	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
//	fmt.Printf("\tNumGC = %v\n", m.NumGC)
//}
//
// func bToMb(b uint64) uint64 {
//	return b / 1024 / 1024
//}

// func BenchmarkConcurrentSwissMapGoMaxProcs1(b *testing.B) {
//	runtime.GOMAXPROCS(1)
//	debug.SetGCPercent(-1)
//	debug.SetMemoryLimit(math.MaxInt64)
//	for _, v := range table {
//		b.Run(fmt.Sprintf("total: %d deletion: %d", v.total, v.deletion), func(b *testing.B) {
//			for i := 0; i < b.N; i++ {
//				m1 := csmap.Create[int, string]()
//				var wg sync.WaitGroup
//				wg.Add(3)
//				go func() {
//					defer wg.Done()
//					var wg2 sync.WaitGroup
//					wg2.Add(v.total)
//					for i := 0; i < v.total; i++ {
//						i := i
//						go func() {
//							defer wg2.Done()
//							m1.Store(i, strconv.Itoa(i))
//						}()
//					}
//					wg2.Wait()
//				}()
//
//				go func() {
//					defer wg.Done()
//					var wg2 sync.WaitGroup
//					wg2.Add(v.total)
//					for i := 0; i < v.total; i++ {
//						i := i
//						go func() {
//							defer wg2.Done()
//							m1.Store(i, strconv.Itoa(i))
//						}()
//					}
//					wg2.Wait()
//				}()
//
//				go func() {
//					defer wg.Done()
//					var wg2 sync.WaitGroup
//					wg2.Add(v.total)
//					for i := 0; i < v.total; i++ {
//						i := i
//						go func() {
//							defer wg2.Done()
//							m1.Store(10, strconv.Itoa(i))
//							m1.Delete(10)
//						}()
//					}
//					wg2.Wait()
//				}()
//				wg.Wait()
//
//				wg.Add(v.deletion + v.total)
//				for i := 0; i < v.deletion; i++ {
//					i := i
//					go func() {
//						defer wg.Done()
//						m1.Delete(i)
//					}()
//				}
//
//				for i := 0; i < v.total; i++ {
//					i := i
//					go func() {
//						defer wg.Done()
//						m1.Load(i)
//					}()
//				}
//				wg.Wait()
//			}
//		})
//	}
//	PrintMemUsage()
//}

// func BenchmarkSyncMapGoMaxProcs1(b *testing.B) {
//	runtime.GOMAXPROCS(1)
//	debug.SetGCPercent(-1)
//	debug.SetMemoryLimit(math.MaxInt64)
//	for _, v := range table {
//		b.Run(fmt.Sprintf("total: %d deletion: %d", v.total, v.deletion), func(b *testing.B) {
//			for i := 0; i < b.N; i++ {
//				var m1 sync.Map
//				var wg sync.WaitGroup
//				wg.Add(3)
//				go func() {
//					defer wg.Done()
//					var wg2 sync.WaitGroup
//					wg2.Add(v.total)
//					for i := 0; i < v.total; i++ {
//						i := i
//						go func() {
//							defer wg2.Done()
//							m1.Store(i, strconv.Itoa(i))
//						}()
//					}
//					wg2.Wait()
//				}()
//
//				go func() {
//					defer wg.Done()
//					var wg2 sync.WaitGroup
//					wg2.Add(v.total)
//					for i := 0; i < v.total; i++ {
//						i := i
//						go func() {
//							defer wg2.Done()
//							m1.Store(i, strconv.Itoa(i))
//						}()
//					}
//					wg2.Wait()
//				}()
//
//				go func() {
//					defer wg.Done()
//					var wg2 sync.WaitGroup
//					wg2.Add(v.total)
//					for i := 0; i < v.total; i++ {
//						i := i
//						go func() {
//							defer wg2.Done()
//							m1.Store(10, strconv.Itoa(i))
//							m1.Delete(10)
//						}()
//					}
//					wg2.Wait()
//				}()
//				wg.Wait()
//
//				wg.Add(v.deletion + v.total)
//				for i := 0; i < v.deletion; i++ {
//					i := i
//					go func() {
//						defer wg.Done()
//						m1.Delete(i)
//					}()
//				}
//
//				for i := 0; i < v.total; i++ {
//					i := i
//					go func() {
//						defer wg.Done()
//						m1.Load(i)
//					}()
//				}
//				wg.Wait()
//			}
//		})
//	}
//	PrintMemUsage()
//}

// func BenchmarkRWMutexMapGoMaxProcs1(b *testing.B) {
//	runtime.GOMAXPROCS(1)
//	debug.SetGCPercent(-1)
//	debug.SetMemoryLimit(math.MaxInt64)
//	for _, v := range table {
//		b.Run(fmt.Sprintf("total: %d deletion: %d", v.total, v.deletion), func(b *testing.B) {
//			for i := 0; i < b.N; i++ {
//				m1 := CreateTestRWMutexMap()
//				var wg sync.WaitGroup
//				wg.Add(3)
//				go func() {
//					defer wg.Done()
//					var wg2 sync.WaitGroup
//					wg2.Add(v.total)
//					for i := 0; i < v.total; i++ {
//						i := i
//						go func() {
//							defer wg2.Done()
//							m1.Store(i, strconv.Itoa(i))
//						}()
//					}
//					wg2.Wait()
//				}()
//
//				go func() {
//					defer wg.Done()
//					var wg2 sync.WaitGroup
//					wg2.Add(v.total)
//					for i := 0; i < v.total; i++ {
//						i := i
//						go func() {
//							defer wg2.Done()
//							m1.Store(i, strconv.Itoa(i))
//						}()
//					}
//					wg2.Wait()
//				}()
//
//				go func() {
//					defer wg.Done()
//					var wg2 sync.WaitGroup
//					wg2.Add(v.total)
//					for i := 0; i < v.total; i++ {
//						i := i
//						go func() {
//							defer wg2.Done()
//							m1.Store(10, strconv.Itoa(i))
//							m1.Delete(10)
//						}()
//					}
//					wg2.Wait()
//				}()
//				wg.Wait()
//
//				wg.Add(v.deletion + v.total)
//				for i := 0; i < v.deletion; i++ {
//					i := i
//					go func() {
//						defer wg.Done()
//						m1.Delete(i)
//					}()
//				}
//
//				for i := 0; i < v.total; i++ {
//					i := i
//					go func() {
//						defer wg.Done()
//						m1.Load(i)
//					}()
//				}
//				wg.Wait()
//			}
//		})
//	}
//	PrintMemUsage()
//}

// func BenchmarkConcurrentSwissMapGoMaxProcsCore(b *testing.B) {
//	debug.SetGCPercent(-1)
//	debug.SetMemoryLimit(math.MaxInt64)
//	for _, v := range table {
//		b.Run(fmt.Sprintf("total: %d deletion: %d", v.total, v.deletion), func(b *testing.B) {
//			for i := 0; i < b.N; i++ {
//				m1 := csmap.Create[int, string]()
//				var wg sync.WaitGroup
//				wg.Add(3)
//				go func() {
//					defer wg.Done()
//					var wg2 sync.WaitGroup
//					wg2.Add(v.total)
//					for i := 0; i < v.total; i++ {
//						i := i
//						go func() {
//							defer wg2.Done()
//							m1.Store(i, strconv.Itoa(i))
//						}()
//					}
//					wg2.Wait()
//				}()
//
//				go func() {
//					defer wg.Done()
//					var wg2 sync.WaitGroup
//					wg2.Add(v.total)
//					for i := 0; i < v.total; i++ {
//						i := i
//						go func() {
//							defer wg2.Done()
//							m1.Store(i, strconv.Itoa(i))
//						}()
//					}
//					wg2.Wait()
//				}()
//
//				go func() {
//					defer wg.Done()
//					var wg2 sync.WaitGroup
//					wg2.Add(v.total)
//					for i := 0; i < v.total; i++ {
//						i := i
//						go func() {
//							defer wg2.Done()
//							m1.Store(10, strconv.Itoa(i))
//							m1.Delete(10)
//						}()
//					}
//					wg2.Wait()
//				}()
//				wg.Wait()
//
//				wg.Add(v.deletion + v.total)
//				for i := 0; i < v.deletion; i++ {
//					i := i
//					go func() {
//						defer wg.Done()
//						m1.Delete(i)
//					}()
//				}
//
//				for i := 0; i < v.total; i++ {
//					i := i
//					go func() {
//						defer wg.Done()
//						m1.Load(i)
//					}()
//				}
//				wg.Wait()
//			}
//		})
//	}
//	PrintMemUsage()
//}

// func BenchmarkSyncMapGoMaxProcsCore(b *testing.B) {
//	debug.SetGCPercent(-1)
//	debug.SetMemoryLimit(math.MaxInt64)
//	for _, v := range table {
//		b.Run(fmt.Sprintf("total: %d deletion: %d", v.total, v.deletion), func(b *testing.B) {
//			for i := 0; i < b.N; i++ {
//				var m1 sync.Map
//				var wg sync.WaitGroup
//				wg.Add(3)
//				go func() {
//					defer wg.Done()
//					var wg2 sync.WaitGroup
//					wg2.Add(v.total)
//					for i := 0; i < v.total; i++ {
//						i := i
//						go func() {
//							defer wg2.Done()
//							m1.Store(i, strconv.Itoa(i))
//						}()
//					}
//					wg2.Wait()
//				}()
//
//				go func() {
//					defer wg.Done()
//					var wg2 sync.WaitGroup
//					wg2.Add(v.total)
//					for i := 0; i < v.total; i++ {
//						i := i
//						go func() {
//							defer wg2.Done()
//							m1.Store(i, strconv.Itoa(i))
//						}()
//					}
//					wg2.Wait()
//				}()
//
//				go func() {
//					defer wg.Done()
//					var wg2 sync.WaitGroup
//					wg2.Add(v.total)
//					for i := 0; i < v.total; i++ {
//						i := i
//						go func() {
//							defer wg2.Done()
//							m1.Store(10, strconv.Itoa(i))
//							m1.Delete(10)
//						}()
//					}
//					wg2.Wait()
//				}()
//				wg.Wait()
//
//				wg.Add(v.deletion + v.total)
//				for i := 0; i < v.deletion; i++ {
//					i := i
//					go func() {
//						defer wg.Done()
//						m1.Delete(i)
//					}()
//				}
//
//				for i := 0; i < v.total; i++ {
//					i := i
//					go func() {
//						defer wg.Done()
//						m1.Load(i)
//					}()
//				}
//				wg.Wait()
//			}
//		})
//	}
//	PrintMemUsage()
//}

// func BenchmarkRWMutexMapGoMaxProcsCore(b *testing.B) {
//	debug.SetGCPercent(-1)
//	debug.SetMemoryLimit(math.MaxInt64)
//	for _, v := range table {
//		b.Run(fmt.Sprintf("total: %d deletion: %d", v.total, v.deletion), func(b *testing.B) {
//			for i := 0; i < b.N; i++ {
//				m1 := CreateTestRWMutexMap()
//				var wg sync.WaitGroup
//				wg.Add(3)
//				go func() {
//					defer wg.Done()
//					var wg2 sync.WaitGroup
//					wg2.Add(v.total)
//					for i := 0; i < v.total; i++ {
//						i := i
//						go func() {
//							defer wg2.Done()
//							m1.Store(i, strconv.Itoa(i))
//						}()
//					}
//					wg2.Wait()
//				}()
//
//				go func() {
//					defer wg.Done()
//					var wg2 sync.WaitGroup
//					wg2.Add(v.total)
//					for i := 0; i < v.total; i++ {
//						i := i
//						go func() {
//							defer wg2.Done()
//							m1.Store(i, strconv.Itoa(i))
//						}()
//					}
//					wg2.Wait()
//				}()
//
//				go func() {
//					defer wg.Done()
//					var wg2 sync.WaitGroup
//					wg2.Add(v.total)
//					for i := 0; i < v.total; i++ {
//						i := i
//						go func() {
//							defer wg2.Done()
//							m1.Store(10, strconv.Itoa(i))
//							m1.Delete(10)
//						}()
//					}
//					wg2.Wait()
//				}()
//				wg.Wait()
//
//				wg.Add(v.deletion + v.total)
//				for i := 0; i < v.deletion; i++ {
//					i := i
//					go func() {
//						defer wg.Done()
//						m1.Delete(i)
//					}()
//				}
//
//				for i := 0; i < v.total; i++ {
//					i := i
//					go func() {
//						defer wg.Done()
//						m1.Load(i)
//					}()
//				}
//				wg.Wait()
//			}
//		})
//	}
//	PrintMemUsage()
//}

// type TestRWMutexMap struct {
//	m map[int]string
//	sync.RWMutex
//}
//
// func CreateTestRWMutexMap() *TestRWMutexMap {
//	return &TestRWMutexMap{
//		m: make(map[int]string),
//	}
//}
//
// func (m *TestRWMutexMap) Store(key int, value string) {
//	m.Lock()
//	defer m.Unlock()
//	m.m[key] = value
//}
//
// func (m *TestRWMutexMap) Delete(key int) {
//	m.Lock()
//	defer m.Unlock()
//	delete(m.m, key)
//}
//
// func (m *TestRWMutexMap) Load(key int) *string {
//	m.RLock()
//	defer m.RUnlock()
//	s, ok := m.m[key]
//	if !ok {
//		return nil
//	}
//	return &s
//}
