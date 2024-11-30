package gomel

import (
	"go/types"
	"strings"
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
		},

		{
			mainQ:   "github.com/pascaldekloe/gomel/internal/testset.GenericInts",
			paramsQ: []string{"int32"},
			fields: []Field{
				{Name: "A", DataSize: 4, StartPos: 0},
				{Name: "B", DataSize: 4, StartPos: 4},
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

		asStruct, ok := hit.(*types.Struct)
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

func TestFind_errors(t *testing.T) {
	tests := []struct {
		typeQ string
		argQ  []string
		want  string
	}{
		{
			typeQ: "github.com/pascaldekloe/gomel/internal/testset.GenericInts",
			argQ:  []string{"builtin.int64", "builtin.int64"},
			want:  `type github.com/pascaldekloe/gomel/internal/testset.GenericInts[T int32 | int64] has 1 generic parameters while queried with ["builtin.int64" "builtin.int64"]`,
		},

		{
			typeQ: "github.com/pascaldekloe/gomel/internal/testset.GenericInts",
			argQ:  []string{"builtin.bool"},
			want:  "generic parameter № 1 type bool does not satisfy interface int32 | int64",
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
