package main

import (
	"flag"
	"fmt"
	"go/build"
	"go/types"
	"io"
	"log"
	"os"
	"strings"

	"github.com/pascaldekloe/gomel"
)

var (
	compFlag = flag.String("comp", build.Default.Compiler, "Select the applicable compiler by `name`.")
	archFlag = flag.String("arch", build.Default.GOARCH, "Select the target architecture by `name`.")

	quietFlag   = flag.Bool("quiet", false, "Disable standard reporting.")
	verboseFlag = flag.Bool("verbose", false, "Enable detailed reporting.")
)

func main() {
	flag.Parse()

	// standard logging
	if *quietFlag {
		log.SetOutput(io.Discard)
	}
	log.SetFlags(0)

	// detailed logging
	verboseOut := io.Discard
	if *verboseFlag {
		verboseOut = os.Stderr
	}
	verbose := log.New(verboseOut, "gomel: ", 0)

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
	verbose.Printf("found type %s as %T",
		hit, hit)

	asStruct, ok := hit.(*types.Struct)
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
	fmt.Println("name\ttype\tstart\tsize")

	var pass int64
	for i := range l.Fields {
		f := &l.Fields[i]

		if pass < f.StartPos {
			// last padded
			fmt.Printf("-\t-\t%d\t%d\t\n",
				pass, f.StartPos-pass)
		}

		fmt.Printf("%s\t%s\t%d\t%d\n",
			f.Name, typeName(f.DataType), f.StartPos, f.DataSize)

		pass = f.StartPos + f.DataSize
	}

	remain := sizes.Sizeof(l.DataType) - pass
	if remain != 0 {
		// end padded
		fmt.Printf("-\t-\t%d\t%d\n",
			pass, remain)
	}
}

func typeName(t types.Type) string {
	s := t.String()

	// omit builtin package
	i := strings.LastIndexByte(s, '.')
	if s[:i+1] == "builtin." {
		s = s[8:]
	}

	return s
}
