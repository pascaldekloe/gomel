// Package testset provides data structures for testing.
package testset

import "github.com/pascaldekloe/gomel/internal/testset/other"

type Bytes struct {
	A [0]byte
	B [1]byte
	C [2]byte
}

type GenericInts[T int32 | int64] struct {
	A T
	B T
}

// Nested has a collision in type name with Sub.
// The main query inherts other packages.
type Nested struct {
	Sub other.Nested
}

// Nested has a collision in type name with Sub.
// The generic queries span multiple packages.
type GenericNested[T other.Nested | Bytes] struct {
	Sub other.Nested
}

// InheritGeneric type T applies to the embedded structure
type InheritGeneric[T int32 | int64] struct {
	GenericInts[T]
}

type BytesAlias Bytes
type FloatAlias float32
