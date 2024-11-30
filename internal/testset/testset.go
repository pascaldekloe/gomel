// Package testset provides data structures for testing.
package testset

type Bytes struct {
	A [0]byte
	B [1]byte
	C [2]byte
}

type GenericInts[T int32|int64] struct {
	A T
	B T
}
