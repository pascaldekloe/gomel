// Package gomel provides type information with insights on Go's memory layout.
package gomel

import (
	"errors"
	"fmt"
	"go/types"
	"slices"
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
	StartPos int64 // byte offset within struct
}

// StructLayout reads the memory structure t for a specific target.
func StructLayout(t *types.Struct, target types.Sizes) Layout {
	fields := make([]*types.Var, t.NumFields())
	for i := range fields {
		fields[i] = t.Field(i)
	}
	l := Layout{
		DataType: t,
		Fields:   make([]Field, len(fields)),
	}

	offsets := target.Offsetsof(fields)
	if len(offsets) != len(fields) {
		panic("number of offsets doesn't match requested")
	}

	for i := range l.Fields {
		f := &l.Fields[i]
		f.Name = fields[i].Name()
		f.DataType = fields[i].Type()
		f.DataSize = target.Sizeof(f.DataType)
		f.StartPos = offsets[i]
	}
	return l
}

type query struct {
	pkg string
	typ string
}

func parseQuery(s string) query {
	i := strings.LastIndexByte(s, '.')
	if i < 0 {
		return query{
			pkg: "builtin",
			typ: s,
		}
	}
	return query{
		pkg: s[:i],
		typ: s[i+1:],
	}
}

func packagesOf(queries []query) []string {
	var list []string
	for i := range queries {
		p := queries[i].pkg
		if p != "builtin" && !slices.Contains(list, p) {
			list = append(list, p)
		}
	}
	return list
}

// ErrNotFound signals lookup failure.
var ErrNotFound = errors.New("no such type")

// Find returns a type match for mainQuery or an ErrNotFound. Generic types also
// need paramQueries for each type parameter.
func Find(mainQuery string, paramQueries ...string) (types.Type, error) {
	queries := make([]query, 1+len(paramQueries))
	queries[0] = parseQuery(mainQuery)
	for i := range paramQueries {
		queries[i+1] = parseQuery(paramQueries[i])
	}

	// lookup
	found, err := findTypes(queries)
	if err != nil {
		return nil, err
	}
	mainType := found[0]
	paramTypes := found[1:]

	var generics *types.TypeParamList
	switch t := mainType.(type) {
	case *types.Named:
		generics = t.TypeParams()
	case *types.Signature:
		generics = t.TypeParams()

	default:
		if len(paramTypes) != 0 {
			return nil, fmt.Errorf("found %s while queried with %d type paramaters",
				mainType, len(paramTypes))
		}
		return mainType, nil
	}

	// pass non-generic types as is
	if generics == nil {
		// mainType is not generic
		if len(paramTypes) == 0 {
			return mainType, nil
		}
		return nil, fmt.Errorf("found non-generic type %s while queried with %d type parameters",
			mainType, len(paramTypes))
	}

	// match generics with the paramTypes found
	if genericN, queryN := generics.Len(), len(paramTypes); genericN != queryN {
		return nil, fmt.Errorf("found %s with %d type parameters while queried with %d",
			mainType, genericN, queryN)
	}
	for i, param := range paramTypes {
		t := paramTypes[i] // *types.Named
		u := t.Underlying()
		if u == types.Typ[types.Invalid] {
			// should not happen ™️
			return nil, fmt.Errorf("found invalid %s for paramameter query № %d",
				paramTypes[i], i+1)
		}

		// Underlying of types.TypeParam always returns an interface
		constraint := generics.At(i).Underlying().(*types.Interface)
		// interfaces can match types by name or by the underlying type
		if !types.Satisfies(u, constraint) && !types.Satisfies(t, constraint) {
			return nil, fmt.Errorf("found %s for parameter query № %d does not satisfy generic %s",
				param, i+1, mainType)
		}
	}

	return types.Instantiate(types.NewContext(), mainType, paramTypes, false)
}

func findTypes(queries []query) ([]types.Type, error) {
	config := packages.Config{
		Mode: packages.NeedImports | packages.NeedTypes,
	}
	loaded, err := packages.Load(&config, packagesOf(queries)...)
	if err != nil {
		return nil, err
	}

	found := make([]types.Type, len(queries))
MapQuery:
	for i := range queries {
		if queries[i].pkg == "builtin" {
			for _, basic := range types.Typ {
				if basic.Name() == queries[i].typ {
					found[i] = basic
					continue MapQuery
				}
			}
			return nil, fmt.Errorf("%w: %q does not match any of the basic types",
				ErrNotFound, queries[i].typ)
		}

		for _, p := range loaded {
			// match type name in package
			if p.Types.Path() != queries[i].pkg {
				continue
			}

			// Load returns non-existing packages with no Name
			if p.Types.Name() == "" {
				return nil, fmt.Errorf("%w: package %q for %q not found",
					ErrNotFound, queries[i].pkg, queries[i].typ)
			}

			hit := p.Types.Scope().Lookup(queries[i].typ)
			if hit != nil {
				found[i] = hit.Type()
				continue MapQuery
			}
		}
		return nil, fmt.Errorf("%w: %q not in package %q",
			ErrNotFound, queries[i].typ, queries[i].pkg)
	}
	return found, nil
}
