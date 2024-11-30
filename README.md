# Gomel

Gomel provides insights on the memory layout of data structures in Go. The
gomel(1) command searches for types within a Go module, including all of its
dependencies. The target architecture defaults to `go env GOARCH`.

For example, command `gomel net.TCPAddr` prints the following table. It shows
how field Zone takes 16 bytes, starting at byte-index 32 within the struct.

```
name	type	start	size
IP	net.IP	0	24
Port	int	24	8
Zone	string	32	16
```

Generic types need all of their type parameters specified with extra arguments.
The following table is the output from `gomel sync/atomic.Pointer float64`. Note
how the first two fields have no size. The generic parameter (`float64`) has no
effect on the outcome in this example.

```
name	type	start	size
_	[0]*float64	0	0
_	sync/atomic.noCopy	0	0
v	unsafe.Pointer	0	8
```

[Padding](https://en.wikipedia.org/wiki/Data_structure_alignment#Data_structure_padding)
is shown as a hypen ('-') in both the name and the type column.
