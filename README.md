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
    csmap.WithShardCount[int, string](32), // default 32,
    csmap.WithCustomHasher[int, string](func(key int) uint64 {
        return 0
    }), // default maphash,
    csmap.WithSize[int, string](1000), // default 0
)
myMap.Store(10, "test")
myMap.Load(10)
myMap.Delete(10)
myMap.Has(10)
myMap.IsEmpty()
myMap.SetIfAbsent(10, "test")
myMap.Range(func(key int, value string) (stop bool) {})
myMap.Count()

```

## Basic Architecture
![img.png](img.png)

## Benchmark Test
Benchmark was made on:
- Apple M1 Max
- 32 GB memory

Benchmark test results can be obtained by running [this file](concurrent_swiss_map_benchmark_test.go) on local computers.

![benchmark.png](benchmark.png)

[doc-img]: https://godoc.org/github.com/mhmtszr/concurrent-swiss-map?status.svg
[doc]: https://godoc.org/github.com/mhmtszr/concurrent-swiss-map
[ci-img]: https://github.com/mhmtszr/concurrent-swiss-map/actions/workflows/build-test.yml/badge.svg
[ci]: https://github.com/mhmtszr/concurrent-swiss-map/actions/workflows/build-test.yml
[cov-img]: https://codecov.io/gh/mhmtszr/concurrent-swiss-map/branch/master/graph/badge.svg
[cov]: https://codecov.io/gh/mhmtszr/concurrent-swiss-map
[go-report-img]: https://goreportcard.com/badge/github.com/mhmtszr/concurrent-swiss-map
[go-report]: https://goreportcard.com/report/github.com/mhmtszr/concurrent-swiss-map