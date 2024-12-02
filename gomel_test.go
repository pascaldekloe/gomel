package gomel

import (
	"go/types"
	"testing"
)

func TestStructLayout(t *testing.T) {
	target := types.StdSizes{WordSize: 8, MaxAlign: 8}
	tests := []struct {
		mainQ   string
		paramsQ []string
		fields  []Field
	}{
		{
			mainQ: "github.com/pascaldekloe/gomel/internal/testset.Bytes",
			fields: []Field{
				{Name: "A", DataSize: 0, StartPos: 0},
				{Name: "B", DataSize: 1, StartPos: 0},
				{Name: "C", DataSize: 2, StartPos: 1},
			},
		}, {
			mainQ: "github.com/pascaldekloe/gomel/internal/testset.BytesAlias",
			fields: []Field{
				{Name: "A", DataSize: 0, StartPos: 0},
				{Name: "B", DataSize: 1, StartPos: 0},
				{Name: "C", DataSize: 2, StartPos: 1},
			},
		},

		{
			mainQ:   "github.com/pascaldekloe/gomel/internal/testset.GenericInts",
			paramsQ: []string{"int32"},
			fields: []Field{
				{Name: "A", DataSize: 4, StartPos: 0},
				{Name: "B", DataSize: 4, StartPos: 4},
			},
		},

		{
			mainQ: "github.com/pascaldekloe/gomel/internal/testset.Nested",
			fields: []Field{
				{Name: "Sub", DataSize: 17, StartPos: 0},
			},
		},

		{
			mainQ:   "github.com/pascaldekloe/gomel/internal/testset.GenericNested",
			paramsQ: []string{"github.com/pascaldekloe/gomel/internal/testset/other.Nested"},
			fields: []Field{
				{Name: "Sub", DataSize: 17, StartPos: 0},
			},
		},

		{
			mainQ:   "github.com/pascaldekloe/gomel/internal/testset.InheritGeneric",
			paramsQ: []string{"int64"},
			fields: []Field{
				{Name: "GenericInts", DataSize: 16, StartPos: 0},
			},
		},
	}

	for _, test := range tests {
		hit, err := Find(test.mainQ, test.paramsQ...)
		if err != nil {
			// no context; the error must be descriptive
			t.Error("Find error:", err)
			continue
		}

		asStruct, ok := hit.Underlying().(*types.Struct)
		if !ok {
			t.Errorf("Find %q got type %T, want a *types.Struct",
				test.mainQ, hit)
			continue
		}
		l := StructLayout(asStruct, &target)

		if len(l.Fields) != len(test.fields) {
			t.Fatalf("Find %q got %d fields, want %d",
				test.mainQ, len(l.Fields), len(test.fields))
			continue
		}
		for i := range l.Fields {
			got, want := &l.Fields[i], &test.fields[i]

			if got.Name != want.Name {
				t.Errorf("Find %q got field %q, want field %q",
					test.mainQ, got.Name, want.Name)
				continue
			}

			if got.DataSize != want.DataSize {
				t.Errorf("Find %q field %q got a %d B data size, want %d B",
					test.mainQ, got.Name, got.DataSize, want.DataSize)
			}
			if got.StartPos != want.StartPos {
				t.Errorf("Find %q field %q got a %d B offset, want %d B",
					test.mainQ, got.Name, got.StartPos, want.StartPos)
			}
		}
	}
}

func TestFind_types(t *testing.T) {
	queries := []string{
		"float32",
		"string",

		"github.com/pascaldekloe/gomel/internal/testset.FloatAlias",
	}
	for _, q := range queries {
		hit, err := Find(q)
		switch {
		case err != nil:
			t.Errorf("Find %q got error: %s", q, err)
		case hit.String() != q:
			t.Errorf("Find %q got %q", q, hit)
		}
	}
}

func TestFind_errors(t *testing.T) {
	tests := []struct {
		mainQ   string
		paramsQ []string
		want    string
	}{
		// not found in package
		{
			mainQ: "github.com/pascaldekloe/gomel.DoesNotExist",
			want:  `no such type: "DoesNotExist" not in package "github.com/pascaldekloe/gomel"`,
		}, {
			mainQ:   "github.com/pascaldekloe/gomel/internal/testset.GenericInts",
			paramsQ: []string{"github.com/pascaldekloe/gomel.DoesNotExist"},
			want:    `no such type: "DoesNotExist" not in package "github.com/pascaldekloe/gomel"`,
		}, {
			mainQ: "github.com/pascaldekloe/gomel/doesnotexist.Arbitrary",
			want:  `no such type: package "github.com/pascaldekloe/gomel/doesnotexist" for "Arbitrary" not found`,
		}, {
			mainQ:   "github.com/pascaldekloe/gomel/internal/testset.GenericInts",
			paramsQ: []string{"github.com/pascaldekloe/gomel/doesnotexist.Arbitrary"},
			want:    `no such type: package "github.com/pascaldekloe/gomel/doesnotexist" for "Arbitrary" not found`,
		}, {
			mainQ: "builtin.Unknown",
			want:  `no such type: "Unknown" does not match any of the basic types`,
		},

		// generics mismatch
		{
			mainQ:   "github.com/pascaldekloe/gomel/internal/testset.GenericInts",
			paramsQ: nil,
			want:    `found github.com/pascaldekloe/gomel/internal/testset.GenericInts[T int32 | int64] with 1 type parameters while queried with 0`,
		}, {
			mainQ:   "github.com/pascaldekloe/gomel/internal/testset.GenericInts",
			paramsQ: []string{"builtin.int64", "builtin.int64"},
			want:    `found github.com/pascaldekloe/gomel/internal/testset.GenericInts[T int32 | int64] with 1 type parameters while queried with 2`,
		}, {
			mainQ:   "github.com/pascaldekloe/gomel/internal/testset.GenericInts",
			paramsQ: []string{"builtin.bool"},
			want:    `found bool for parameter query № 1 does not satisfy generic github.com/pascaldekloe/gomel/internal/testset.GenericInts[T int32 | int64]`,
		},
	}

	for _, test := range tests {
		hit, err := Find(test.mainQ, test.paramsQ...)
		if err == nil {
			t.Errorf("lookup of %q with %q got %T, want error",
				test.mainQ, test.paramsQ, hit)
			continue
		}
		if got := err.Error(); got != test.want {
			t.Errorf("lookup of %q with %q got error %q, want %q",
				test.mainQ, test.paramsQ, got, test.want)
		}
	}
}
