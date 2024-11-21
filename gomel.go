// Package gomel provides insights about memory layout.
package gomel

import (
	"errors"
	"go/types"
	"log"
	"strings"

	"golang.org/x/tools/go/packages"
)

// Layout defines a memory structure.
type Layout struct {
	DataType types.Type
	Fields   []Field
}

// Field defines memory from a struct member.
type Field struct {
	Name     string // label in source code
	DataType types.Type
	DataSize int64 // number of bytes
	StartPos int64 // index of first byte within struct
}

func LayoutOf(t *types.Struct, sizes types.Sizes) Layout {
	fields := make([]*types.Var, t.NumFields())
	for i := range fields {
		fields[i] = t.Field(i)
	}
	l := Layout{
		DataType: t,
		Fields:   make([]Field, len(fields)),
	}

	offsets := sizes.Offsetsof(fields)
	if len(offsets) != len(fields) {
		panic("number of offsets doesn't match requested")
	}

	for i := range l.Fields {
		f := &l.Fields[i]
		f.Name = fields[i].Name()
		f.DataType = fields[i].Type()
		f.DataSize = sizes.Sizeof(f.DataType)
		f.StartPos = offsets[i]
	}
	return l
}

var ErrNotFound = errors.New("type not found")

func Find(typeQuery string, report *log.Logger) (types.Type, error) {
	i := strings.LastIndexByte(typeQuery, '.')
	if i < 0 {
		return FindInPackage("builtin", typeQuery, report)
	}
	return FindInPackage(typeQuery[:i], typeQuery[i+1:], report)
}

func FindInPackage(packageQuery, typeQuery string, report *log.Logger) (types.Type, error) {
	config := packages.Config{
		Mode: packages.NeedImports | packages.NeedExportFile | packages.NeedTypes | packages.NeedSyntax,
	}
	pkgs, err := packages.Load(&config, packageQuery)
	if err != nil {
		return nil, err
	}

	for _, p := range pkgs {
		hit := p.Types.Scope().Lookup(typeQuery)
		if hit == nil {
			report.Printf("type %q not in package path %q",
				typeQuery, p.PkgPath)
		} else {
			report.Printf("type %q found in package path %q",
				typeQuery, p.PkgPath)
			return hit.Type().Underlying(), nil
		}
	}

	return nil, ErrNotFound
}
