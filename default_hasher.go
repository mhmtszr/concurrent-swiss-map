package csmap

import "github.com/mhmtszr/concurrent-swiss-map/maphash"

type DefaultHasher[K comparable] struct {
	h maphash.Hasher[K]
}

func NewDefaultHasher[K comparable]() *DefaultHasher[K] {
	return &DefaultHasher[K]{
		h: maphash.NewHasher[K](),
	}
}
