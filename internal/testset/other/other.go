// Package other contributes to package testset.
package other

type Nested struct {
	Another
	Foo bool
}

type Another struct {
	*Nested
	Bar int
}
