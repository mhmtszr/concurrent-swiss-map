// Copyright 2022 Dolthub, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// This file incorporates work covered by the following copyright and
// permission notice:
//
// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.18 || go1.19
// +build go1.18 go1.19

package maphash

import (
	"math/rand"
	"unsafe"
)

type hashfn func(unsafe.Pointer, uintptr) uintptr

func getRuntimeHasher[K comparable]() (h hashfn) {
	a := any(make(map[K]struct{}))
	i := (*mapiface)(unsafe.Pointer(&a))
	h = i.typ.hasher
	return
}

var hashSeed = rand.Int()

func newHashSeed() uintptr {
	return uintptr(hashSeed)
}

// noescape hides a pointer from escape analysis. It is the identity function
// but escape analysis doesn't think the output depends on the input.
// noescape is inlined and currently compiles down to zero instructions.
// USE CAREFULLY!
// This was copied from the runtime (via pkg "strings"); see issues 23382 and 7921.
//
//go:nosplit
//go:nocheckptr
func noescape(p unsafe.Pointer) unsafe.Pointer {
	x := uintptr(p)
	return unsafe.Pointer(x ^ 0)
}

type mapiface struct {
	typ *maptype
	val *hmap
}

// go/src/runtime/type.go
type maptype struct {
	key        *_type
	elem       *_type
	bucket     *_type
	hasher     func(unsafe.Pointer, uintptr) uintptr
	typ        _type
	flags      uint32
	bucketsize uint16
	keysize    uint8
	elemsize   uint8
}

// go/src/runtime/map.go
type hmap struct {
	buckets    unsafe.Pointer
	oldbuckets unsafe.Pointer
	extra      unsafe.Pointer
	count      int
	nevacuate  uintptr
	hash0      uint32
	noverflow  uint16
	flags      uint8
	B          uint8
}

// go/src/runtime/type.go
type tflag uint8
type nameOff int32
type typeOff int32

// go/src/runtime/type.go
type _type struct {
	equal      func(unsafe.Pointer, unsafe.Pointer) bool
	gcdata     *byte
	size       uintptr
	ptrdata    uintptr
	hash       uint32
	str        nameOff
	ptrToThis  typeOff
	tflag      tflag
	align      uint8
	fieldAlign uint8
	kind       uint8
}
