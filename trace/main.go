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
	"strings"
)

var depth int

// TimedNode is a replacement for any CallExpr that performs the original
// CallExpr but wraps it in a timing block and returns the value of the
// CallExpr. It also prints the original source code of the CallExpr.
func replace(n ast.CallExpr, typ interface{}, fset *token.FileSet) {

	var b bytes.Buffer

	printer.Fprint(&b, fset, n)
	s := b.String()

	if len(s) > 0 {
		fmt.Printf("%s\t%s:\t%s\n", reflect.TypeOf(typ), s, fset.Position(n.Pos()))
	}

	//n = &ast.FuncLit{
	//	Fun: n,
	//	Lparen   n.Pos,
	//	Args     []Expr
	//	Ellipsis token.Pos
	//	Rparen   token.Pos
	//}
}

var fset *token.FileSet

func main() {
	var file = flag.String("file", "", "A file to trace")
	flag.Parse()

	if *file == "" {
		fmt.Printf("Specify a filename to run and trace")
		flag.PrintDefaults()
		return
	}

	fset = token.NewFileSet()
	tree, err := parser.ParseFile(fset, *file, nil, 0)

	if err != nil {
		panic(err)
	}
	//ast.Inspect(tree, func(n ast.Node) bool {
	//	switch typ := n.(type) {
	//	default:
	//		replace(n.(*ast.CallExpr), typ, fset)
	//	}
	//})
	depth = 0
	walk(tree)
}

func printLineNode(n ast.Node) {
	fmt.Printf("id: %v\n", reflect.TypeOf(&ast.Ident{}))
	t := reflect.TypeOf(n)
	switch t {

	case reflect.TypeOf(&ast.Ident{}):
		if n.(*ast.Ident).Obj != nil {
			var b bytes.Buffer
			printer.Fprint(&b, fset, n)
			fmt.Printf("CallExpr: %s\n", b.String())
		}

	case reflect.TypeOf(&ast.CallExpr{}),
		reflect.TypeOf(&ast.SelectorExpr{}),
		reflect.TypeOf(&ast.CommClause{}),
		reflect.TypeOf(&ast.ExprStmt{}):
		var b bytes.Buffer
		printer.Fprint(&b, fset, n)
		fmt.Printf("CallExpr: %s\n", b.String())
	}
}
func walk(n ast.Node) {
	depth++
	printLineNode(n)
	switch tt := n.(type) {
	case ast.Node:

		v := reflect.ValueOf(n)
		x := reflect.Indirect(v)
		if !v.IsNil() {
			for i := 0; i < x.NumField(); i++ {
				field := x.Field(i)
				fmt.Printf(
					"%s%s : %s : %s : %s\n",
					strings.Repeat("  ", depth),
					field.Type(),
					reflect.TypeOf(tt),
					field.Type().Kind(),
					fset.Position(n.Pos()),
				)

				if field.Type().Kind() == reflect.Slice {
					if field.CanInterface() && field.Interface() != nil {
						s := reflect.ValueOf(field.Interface())
						for j := 0; j < s.Len(); j++ {
							walk(s.Index(j).Interface().(ast.Node))
						}
					}
				}
				if field.CanInterface() && field.Interface() != nil {
					switch field.Interface().(type) {
					case ast.Node:
						walk(field.Interface().(ast.Node))
					}
				}
			}
		}
	}
	depth--
}
