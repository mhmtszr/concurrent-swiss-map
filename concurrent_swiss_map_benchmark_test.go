package csmap_test

//
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
//}{
//	{
//		total:    100,
//		deletion: 100,
//	},
//	{
//		total:    300000,
//		deletion: 10000,
//	},
//	{
//		total:    500000,
//		deletion: 500000,
//	},
//	{
//		total:    5000000,
//		deletion: 5000000,
//	},
//}
//
//func BenchmarkConcurrentSwissMapGoMaxProcs1(b *testing.B) {
//	runtime.GOMAXPROCS(1)
//	for _, v := range table {
//		b.Run(fmt.Sprintf("total: %d deletion: %d", v.total, v.deletion), func(b *testing.B) {
//			for i := 0; i < b.N; i++ {
//				m1 := csmap.Create[int, string]()
//
//				go func() {
//					for i := 0; i < v.total; i++ {
//						i := i
//						go func() {
//							m1.Store(i, strconv.Itoa(i))
//						}()
//					}
//				}()
//
//				go func() {
//					for i := 0; i < v.total; i++ {
//						i := i
//						go func() {
//							m1.Store(i, strconv.Itoa(i))
//						}()
//					}
//				}()
//
//				go func() {
//					for i := 0; i < v.total; i++ {
//						i := i
//						go func() {
//							m1.Store(10, strconv.Itoa(i))
//							m1.Delete(10)
//						}()
//					}
//				}()
//
//				for i := 0; i < v.deletion; i++ {
//					i := i
//					go func() {
//						m1.Delete(i)
//					}()
//				}
//
//				for i := 0; i < v.total; i++ {
//					i := i
//					go func() {
//						m1.Load(i)
//					}()
//				}
//			}
//		})
//	}
//}
//
//func BenchmarkSyncMapGoMaxProcs1(b *testing.B) {
//	runtime.GOMAXPROCS(1)
//	for _, v := range table {
//		b.Run(fmt.Sprintf("total: %d deletion: %d", v.total, v.deletion), func(b *testing.B) {
//			for i := 0; i < b.N; i++ {
//				var m1 sync.Map
//
//				go func() {
//					for i := 0; i < v.total; i++ {
//						i := i
//						go func() {
//							m1.Store(i, strconv.Itoa(i))
//						}()
//					}
//				}()
//
//				go func() {
//					for i := 0; i < v.total; i++ {
//						i := i
//						go func() {
//							m1.Store(i, strconv.Itoa(i))
//						}()
//					}
//				}()
//
//				go func() {
//					for i := 0; i < v.total; i++ {
//						i := i
//						go func() {
//							m1.Store(10, strconv.Itoa(i))
//							m1.Delete(10)
//						}()
//					}
//				}()
//
//				for i := 0; i < v.deletion; i++ {
//					i := i
//					go func() {
//						m1.Delete(i)
//					}()
//				}
//
//				for i := 0; i < v.total; i++ {
//					i := i
//					go func() {
//						m1.Load(i)
//					}()
//				}
//			}
//		})
//	}
//}
//
//func BenchmarkRWMutexMapGoMaxProcs1(b *testing.B) {
//	runtime.GOMAXPROCS(1)
//	for _, v := range table {
//		b.Run(fmt.Sprintf("total: %d deletion: %d", v.total, v.deletion), func(b *testing.B) {
//			for i := 0; i < b.N; i++ {
//				m1 := CreateTestRWMutexMap()
//
//				go func() {
//					for i := 0; i < v.total; i++ {
//						i := i
//						go func() {
//							m1.Store(i, strconv.Itoa(i))
//						}()
//					}
//				}()
//
//				go func() {
//					for i := 0; i < v.total; i++ {
//						i := i
//						go func() {
//							m1.Store(i, strconv.Itoa(i))
//						}()
//					}
//				}()
//
//				go func() {
//					for i := 0; i < v.total; i++ {
//						i := i
//						go func() {
//							m1.Store(10, strconv.Itoa(i))
//							m1.Delete(10)
//						}()
//					}
//				}()
//
//				for i := 0; i < v.deletion; i++ {
//					i := i
//					go func() {
//						m1.Delete(i)
//					}()
//				}
//
//				for i := 0; i < v.total; i++ {
//					i := i
//					go func() {
//						m1.Load(i)
//					}()
//				}
//			}
//		})
//	}
//}
//
//func BenchmarkConcurrentSwissMapGoMaxProcsCore(b *testing.B) {
//	for _, v := range table {
//		b.Run(fmt.Sprintf("total: %d deletion: %d", v.total, v.deletion), func(b *testing.B) {
//			for i := 0; i < b.N; i++ {
//				m1 := csmap.Create[int, string]()
//
//				go func() {
//					for i := 0; i < v.total; i++ {
//						i := i
//						go func() {
//							m1.Store(i, strconv.Itoa(i))
//						}()
//					}
//				}()
//
//				go func() {
//					for i := 0; i < v.total; i++ {
//						i := i
//						go func() {
//							m1.Store(i, strconv.Itoa(i))
//						}()
//					}
//				}()
//
//				go func() {
//					for i := 0; i < v.total; i++ {
//						i := i
//						go func() {
//							m1.Store(10, strconv.Itoa(i))
//							m1.Delete(10)
//						}()
//					}
//				}()
//
//				for i := 0; i < v.deletion; i++ {
//					i := i
//					go func() {
//						m1.Delete(i)
//					}()
//				}
//
//				for i := 0; i < v.total; i++ {
//					i := i
//					go func() {
//						m1.Load(i)
//					}()
//				}
//			}
//		})
//	}
//}
//
//func BenchmarkSyncMapGoMaxProcsCore(b *testing.B) {
//	for _, v := range table {
//		b.Run(fmt.Sprintf("total: %d deletion: %d", v.total, v.deletion), func(b *testing.B) {
//			for i := 0; i < b.N; i++ {
//				var m1 sync.Map
//
//				go func() {
//					for i := 0; i < v.total; i++ {
//						i := i
//						go func() {
//							m1.Store(i, strconv.Itoa(i))
//						}()
//					}
//				}()
//
//				go func() {
//					for i := 0; i < v.total; i++ {
//						i := i
//						go func() {
//							m1.Store(i, strconv.Itoa(i))
//						}()
//					}
//				}()
//
//				go func() {
//					for i := 0; i < v.total; i++ {
//						i := i
//						go func() {
//							m1.Store(10, strconv.Itoa(i))
//							m1.Delete(10)
//						}()
//					}
//				}()
//
//				for i := 0; i < v.deletion; i++ {
//					i := i
//					go func() {
//						m1.Delete(i)
//					}()
//				}
//
//				for i := 0; i < v.total; i++ {
//					i := i
//					go func() {
//						m1.Load(i)
//					}()
//				}
//			}
//		})
//	}
//}
//
//func BenchmarkRWMutexMapGoMaxProcsCore(b *testing.B) {
//	for _, v := range table {
//		b.Run(fmt.Sprintf("total: %d deletion: %d", v.total, v.deletion), func(b *testing.B) {
//			for i := 0; i < b.N; i++ {
//				m1 := CreateTestRWMutexMap()
//
//				go func() {
//					for i := 0; i < v.total; i++ {
//						i := i
//						go func() {
//							m1.Store(i, strconv.Itoa(i))
//						}()
//					}
//				}()
//
//				go func() {
//					for i := 0; i < v.total; i++ {
//						i := i
//						go func() {
//							m1.Store(i, strconv.Itoa(i))
//						}()
//					}
//				}()
//
//				go func() {
//					for i := 0; i < v.total; i++ {
//						i := i
//						go func() {
//							m1.Store(10, strconv.Itoa(i))
//							m1.Delete(10)
//						}()
//					}
//				}()
//
//				for i := 0; i < v.deletion; i++ {
//					i := i
//					go func() {
//						m1.Delete(i)
//					}()
//				}
//
//				for i := 0; i < v.total; i++ {
//					i := i
//					go func() {
//						m1.Load(i)
//					}()
//				}
//			}
//		})
//	}
//}
//
//type TestRWMutexMap struct {
//	m map[int]string
//	sync.RWMutex
//}
//
//func CreateTestRWMutexMap() *TestRWMutexMap {
//	return &TestRWMutexMap{
//		m: make(map[int]string),
//	}
//}
//
//func (m *TestRWMutexMap) Store(key int, value string) {
//	m.Lock()
//	defer m.Unlock()
//	m.m[key] = value
//}
//
//func (m *TestRWMutexMap) Delete(key int) {
//	m.Lock()
//	defer m.Unlock()
//	delete(m.m, key)
//}
//
//func (m *TestRWMutexMap) Load(key int) *string {
//	m.RLock()
//	defer m.RUnlock()
//	s, ok := m.m[key]
//	if !ok {
//		return nil
//	}
//	return &s
//}
