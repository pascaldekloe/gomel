package main

import (
	"flag"
	"fmt"
	"go/build"
	"go/types"
	"io"
	"log"
	"os"
	"text/tabwriter"

	"github.com/pascaldekloe/gomel"
)

var (
	compFlag = flag.String("comp", build.Default.Compiler, "Select the applicable compiler by `name`.")
	archFlag = flag.String("arch", build.Default.GOARCH, "Select the target architecture by `name`.")

	sepFlag = flag.String("sep", "tab", "Print output separated by either tab, spaces, or comma.")
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

	var w io.Writer = os.Stdout
	var sep rune = '\t'
	switch *sepFlag {
	case "tab":
		break
	case "spaces":
		tw := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)
		defer tw.Flush()
		w = tw
	case "comma":
		sep = ','
	default:
		log.Fatalf("gomel: unknown separator %q; need tab, spaces or comma",
			*sepFlag)
	}
	print(w, l, sizes, sep)
}

func print(w io.Writer, l gomel.Layout, sizes types.Sizes, sep rune) {
	// header
	fmt.Fprintf(w, "Name%cType%cSize%cOffset\n",
		sep, sep, sep)

	var pass int64
	for i := range l.Fields {
		f := &l.Fields[i]

		if pass < f.Offset {
			// padding between previous field
			fmt.Fprintf(w, "-%c-%c%d%c%d\n",
				sep, sep, f.Offset-pass, sep, pass)
		}

		typeName := types.TypeString(f.DataType, packageNameOnly)
		fmt.Fprintf(w, "%s%c%s%c%d%c%d\n",
			f.Name, sep, typeName, sep, f.DataSize, sep, f.Offset)

		pass = f.Offset + f.DataSize
	}

	remain := sizes.Sizeof(l.DataType) - pass
	if remain != 0 {
		// padding at struct end
		fmt.Fprintf(w, "-%c-%c%d%c%d\n",
			sep, sep, remain, sep, pass)
	}
}

// PackageNameOnly skips/trims the path component, if any.
func packageNameOnly(p *types.Package) string { return p.Name() }
