# Concurrent Swiss Map [![GoDoc][doc-img]][doc] [![Build Status][ci-img]][ci] [![Coverage Status][cov-img]][cov] [![Go Report Card][go-report-img]][go-report]

**Concurrent Swiss Map** is an open-source Go library that provides a high-performance, thread-safe generic concurrent hash map implementation designed to handle concurrent access efficiently. It's built with a focus on simplicity, speed, and reliability, making it a solid choice for scenarios where concurrent access to a hash map is crucial.

Uses [dolthub/swiss](https://github.com/dolthub/swiss) map implementation under the hood.

## Installation

Supports 1.18+ Go versions because of Go Generics
```
go get github.com/mhmtszr/concurrent-swiss-map
```

## Usage
New functions will be added soon...
```go

myMap := csmap.Create[int, string](
			csmap.WithShardCount(32) // default 32,
			csmap.WithCustomHasher[int, string](func(key int) uint64 {
                    return 0
            }) // default maphash,
		)
myMap.Store(10, "test")
myMap.Load(10)
myMap.Delete(10)
myMap.Has(10)
myMap.IsEmpty()
myMap.SetIfAbsent(10, "test")
myMap.Range(func(key int, value string) (stop bool) {})

```

## Basic Architecture
![img.png](img.png)

## Benchmark Test
Benchmark was made on:
- Apple M1 Max
- 32 GB memory

Benchmark test results can be obtained by running [this file](concurrent_swiss_map_benchmark_test.go) on local computers.

### GOMAXPROCS = 1

| Map                              | Test Case(Write - Delete) | Execution Time(ns/op) | Allocations  |
|:---------------------------------|:-------------------------:|:---------------------:|:------------:|
| ⚡ Concurrent Swiss Map(32 Shard) |         100 - 100         |      **117334**       |   **733**    
| ⚡ Concurrent Swiss Map(32 Shard) |      300000 - 300000      |     **341128639**     | **2423555**  
| ⚡ Concurrent Swiss Map(32 Shard) |     5000000 - 5000000     |    **6420842000**     | **45035342** 
| Sync Map                         |         100 - 100         |        143586         |     1491     
| Sync Map                         |      300000 - 300000      |      2014041667       |   16468880   
| Sync Map                         |     5000000 - 5000000     |      18961874125      |   96790843   
| RW Mutex Map                     |         100 - 100         |        138229         |     613      
| RW Mutex Map                     |      300000 - 300000      |       559620229       |   2101873    
| RW Mutex Map                     |     5000000 - 5000000     |      14591007000      |   45475146   

### GOMAXPROCS = Core Count

| Map                              | Test Case(Write - Delete) | Execution Time(ns/op) | Allocations  |
|:---------------------------------|:-------------------------:|:---------------------:|:------------:|
| ⚡ Concurrent Swiss Map(32 Shard) |         100 - 100         |      **112274**       |   **733**    
| ⚡ Concurrent Swiss Map(32 Shard) |      300000 - 300000      |     **341237361**     | **2412835**  
| ⚡ Concurrent Swiss Map(32 Shard) |     5000000 - 5000000     |    **6586543375**     | **45002469** 
| Sync Map                         |         100 - 100         |        140694         |     1489     
| Sync Map                         |      300000 - 300000      |       676816958       |   5108574    
| Sync Map                         |     5000000 - 5000000     |      18196143500      |   93150564   
| RW Mutex Map                     |         100 - 100         |        141225         |     613      
| RW Mutex Map                     |      300000 - 300000      |       738111396       |   2398990    
| RW Mutex Map                     |     5000000 - 5000000     |      14398454459      |   45261604   

[doc-img]: https://godoc.org/github.com/mhmtszr/concurrent-swiss-map?status.svg
[doc]: https://godoc.org/github.com/mhmtszr/concurrent-swiss-map
[ci-img]: https://github.com/mhmtszr/concurrent-swiss-map/actions/workflows/build-test.yml/badge.svg
[ci]: https://github.com/mhmtszr/concurrent-swiss-map/actions/workflows/build-test.yml
[cov-img]: https://codecov.io/gh/mhmtszr/concurrent-swiss-map/branch/master/graph/badge.svg
[cov]: https://codecov.io/gh/mhmtszr/concurrent-swiss-map
[go-report-img]: https://goreportcard.com/badge/github.com/mhmtszr/concurrent-swiss-map
[go-report]: https://goreportcard.com/report/github.com/mhmtszr/concurrent-swiss-map