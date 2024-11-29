package gomel

import (
	"go/types"
	"strings"
	"testing"
)

func TestDuo(t *testing.T) {
	hit, err := Find("github.com/pascaldekloe/gomel/internal/testset.Duo")
	if err != nil {
		t.Fatal("lookup error:", err)
	}
	asStruct, ok := hit.Underlying().(*types.Struct)
	if !ok {
		t.Fatalf("got underlying type %T from Find, want a struct", hit.Underlying())
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

func TestGenericDuo(t *testing.T) {
	hit, err := Find("github.com/pascaldekloe/gomel/internal/testset.GenericDuo",
		"builtin.int64")
	if err != nil {
		t.Fatal("lookup error:", err)
	}
	asStruct, ok := hit.Underlying().(*types.Struct)
	if !ok {
		t.Fatalf("got underlying type %T from Find, want a struct", hit.Underlying())
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

func TestFind_errors(t *testing.T) {
	tests := []struct {
		typeQ string
		argQ  []string
		want  string
	}{
		{
			typeQ: "github.com/pascaldekloe/gomel/internal/testset.GenericDuo",
			argQ:  []string{"builtin.int", "builtin.int"},
			want:  `type github.com/pascaldekloe/gomel/internal/testset.GenericDuo[T int8 | int16 | int32 | int64] has 1 generic parameters while queried with ["builtin.int" "builtin.int"]`,
		},

		{
			typeQ: "github.com/pascaldekloe/gomel/internal/testset.GenericDuo",
			argQ:  []string{"builtin.bool"},
			want:  "generic parameter № 1 type bool does not satisfy interface int8 | int16 | int32 | int64",
		},
	}

	for _, test := range tests {
		hit, err := Find(test.typeQ, test.argQ...)
		if err == nil {
			t.Errorf("lookup of %q with %q got %T, want error",
				test.typeQ, test.argQ, hit)
			continue
		}
		s := err.Error()
		if !strings.Contains(s, test.want) {
			t.Errorf("lookup of %q with %q got error %q, want %q included",
				test.typeQ, test.argQ, s, test.want)
		}
	}
}
