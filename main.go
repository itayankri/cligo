package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

func main() {
	fset := token.NewFileSet() // positions are relative to fset

	pathToPackageDir := strings.ReplaceAll(os.Args[1], "\\", "/") + "/"
	endpoints := strings.Split(pathToPackageDir, "/")
	packageName := endpoints[len(endpoints)-2]

	packages, err := parser.ParseDir(
		fset,
		pathToPackageDir,
		nil,
		parser.ParseComments)

	if err != nil {
		panic(err)
	}

	exportedFunctions := getAllExportedFunctions(packages[packageName])

	for index, exportedFunction := range exportedFunctions {
		fmt.Printf("function #%d - %s\n", index, exportedFunction.Name.Name)
	}
}

func getAllExportedFunctions(pkg *ast.Package) []*ast.FuncDecl {
	exportedFunctions := make([]*ast.FuncDecl, 0)

	for _, file := range pkg.Files {
		for _, decl := range file.Decls {
			if funcDecl, ok := decl.(*ast.FuncDecl); ok {
				if isExported(funcDecl.Name.Name) {
					exportedFunctions = append(exportedFunctions, funcDecl)
				}

				annotations := getAnnotations(funcDecl.Name.Name, funcDecl.Doc)
				for index, ann := range annotations {
					fmt.Printf("annotation #%d - %s\n", index, *ann)
				}
			}
		}
	}

	return exportedFunctions
}

func getAnnotations(functionName string, doc *ast.CommentGroup) []*annotation {
	if doc == nil {
		return nil
	}

	annotations := make([]*annotation, 0)

	for _, comment := range doc.List {
		switch {
		case strings.Contains(comment.Text, string(CLIGO_COMMAND)):
			{
				annotations = append(annotations, &annotation{CLIGO_COMMAND, functionName, ""})
			}
		case strings.Contains(comment.Text, string(CLIGO_ARGUMENT)):
			{
				annotations = append(annotations, &annotation{CLIGO_ARGUMENT, "", ""})
			}
		}
	}

	return annotations
}

func isExported(functionName string) bool {
	return functionName[0] >= 'A' && functionName[0] <= 'Z'
}
