// Package testset provides data structures for testing.
package testset

type Duo struct {
	A uint64
	B [8]byte
}

type GenericDuo[T int8|int16|int32|int64] struct {
	A T
	B T
}
