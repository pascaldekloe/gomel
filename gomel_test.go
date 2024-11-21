package gomel

import (
	"go/types"
	"log"
	"testing"
)

func TestDuo(t *testing.T) {
	hit, err := Find("github.com/pascaldekloe/gomel/internal/testset.Duo", log.Default())
	if err != nil {
		t.Fatal("lookup error:", err)
	}
	asStruct, ok := hit.(*types.Struct)
	if !ok {
		t.Fatalf("got type %T from Find, want a struct", hit)
	}

	l := LayoutOf(asStruct, &types.StdSizes{WordSize: 8, MaxAlign: 8})
	if len(l.Fields) != 2 {
		t.Fatalf("got %d fields, want 2", len(l.Fields))
	}
	a := l.Fields[0]
	b := l.Fields[1]

	if a.Name != "A" || b.Name != "B" {
		t.Errorf("got field names %q and %q, want A and B",
			a.Name, b.Name)
	}
	if a.StartPos != 0 || b.StartPos != 8 {
		t.Errorf("got byte index %d and %d, want 0 and 8",
			a.StartPos, b.StartPos)
	}
	if a.DataSize != 8 || b.DataSize != 8 {
		t.Errorf("got byte size %d and %d, want 8 and 8",
			a.DataSize, b.DataSize)
	}
}
