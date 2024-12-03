[![API Documentation](https://godoc.org/github.com/pascaldekloe/gomel?status.svg)](https://godoc.org/github.com/pascaldekloe/gomel)
[![Build Status](https://github.com/pascaldekloe/gomel/actions/workflows/go.yml/badge.svg)](https://github.com/pascaldekloe/gomel/actions/workflows/go.yml)

# Gomel

Gomel provides insights on the memory layout of data structures in Go. The
gomel(1) command searches for types within a Go module, including all of its
dependencies. The target architecture defaults to `go env GOARCH`.

For example, command `gomel net.TCPAddr` prints the following table. It shows
how field Zone takes 16 bytes, starting at byte-index 32 within the struct.

```
Name	Type	Size	Offset
IP	net.IP	24	0
Port	int	8	24
Zone	string	16	32
```

Generic types need all of their type parameters specified with extra arguments.
The following table is the output from `gomel sync/atomic.Pointer float64`. Note
how the first two fields have no size. The generic parameter (`float64`) has no
effect on the outcome in this example.

```
Name	Type	Size	Offset
_	[0]*float64	0	0
_	sync/atomic.noCopy	0	0
v	unsafe.Pointer	8	0
```

[Padding](https://en.wikipedia.org/wiki/Data_structure_alignment#Data_structure_padding)
is shown as a hypen ('-') in both the name and the type column.
