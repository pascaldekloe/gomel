package main

import (
	"flag"
	"fmt"
	"go/build"
	"go/types"
	"log"

	"github.com/pascaldekloe/gomel"
)

var (
	compFlag = flag.String("comp", build.Default.Compiler, "Select the applicable compiler by `name`.")
	archFlag = flag.String("arch", build.Default.GOARCH, "Select the target architecture by `name`.")
)

func main() {
	flag.Parse()

	log.SetFlags(0)

	// establish compilation target
	sizes := types.SizesFor(*compFlag, *archFlag)
	if sizes == nil {
		log.Fatalf("gomel: unknown compiler/architecture pair %q and %q",
			*compFlag, *archFlag)
	}

	// read types from arguments
	args := flag.Args()
	if len(args) == 0 {
		log.Fatal("gomel: need type argument, as in <package>.<type>")
	}
	hit, err := gomel.Find(args[0], args[1:]...)
	if err != nil {
		log.Fatal(err)
	}

	asStruct, ok := hit.Underlying().(*types.Struct)
	if !ok {
		// TODO(pascaldekloe): deal with non-struct
		log.Fatalf("gomel: type %s as %T is not a struct",
			hit, hit)
	}

	l := gomel.StructLayout(asStruct, sizes)
	print(l, sizes)
}

func print(l gomel.Layout, sizes types.Sizes) {
	// header
	fmt.Println("Name\tType\tSize\tOffset")

	var pass int64
	for i := range l.Fields {
		f := &l.Fields[i]

		if pass < f.Offset {
			// padding between previous field
			fmt.Printf("-\t-\t%d\t%d\t\n",
				f.Offset-pass, pass)
		}

		fmt.Printf("%s\t%s\t%d\t%d\n",
			f.Name, f.DataType, f.DataSize, f.Offset)

		pass = f.Offset + f.DataSize
	}

	remain := sizes.Sizeof(l.DataType) - pass
	if remain != 0 {
		// padding at struct end
		fmt.Printf("-\t-\t%d\t%d\n",
			remain, pass)
	}
}
