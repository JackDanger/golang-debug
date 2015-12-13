package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"reflect"
)

// A replacement for any CallExpr
type TimedCallExpr ast.CallExpr

func main() {
	var file = flag.String("file", "", "A file to trace")
	flag.Parse()

	if *file == "" {
		fmt.Printf("Specify a filename to run and trace")
		flag.PrintDefaults()
	}

	fset := token.NewFileSet()
	tree, err := parser.ParseFile(fset, *file, nil, 0)

	if err != nil {
		panic(err)
	}

	ast.Inspect(tree, func(n ast.Node) bool {
		var b bytes.Buffer

		printer.Fprint(&b, fset, n)

		s := b.String()

		if len(s) > 0 {
			fmt.Printf("%s\t%s:\t%s\n", reflect.TypeOf(n), s, fset.Position(n.Pos()))
		}
		return true
	})
}
