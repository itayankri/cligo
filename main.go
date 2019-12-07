package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func main() {
	fset := token.NewFileSet() // positions are relative to fset

	packages, err := parser.ParseDir(
		fset,
		"C:\\Users\\Itay\\Projects\\src\\github.com\\itayankri\\caf\\internal\\pkg\\validators\\jsonvalidator",
		nil,
		parser.ParseComments)

	if err != nil {
		panic(err)
	}

	exportedFunctions := getAllExportedFunctions(packages["jsonvalidator"])

	fmt.Println(exportedFunctions)
}

func getAllExportedFunctions(pkg *ast.Package) []*ast.FuncDecl {
	exportedFunctions := make([]*ast.FuncDecl, 0)

	for _, file := range pkg.Files {
		for _, decl := range file.Decls {
			if v, ok := decl.(*ast.FuncDecl); ok {
				exportedFunctions = append(exportedFunctions, v)
			}
		}
	}

	return exportedFunctions
}
