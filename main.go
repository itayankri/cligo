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

	// TODO: fix later
	goPath := strings.Split(os.Getenv("GOPATH"), ";")[1]

	//
	relativePathToPackage := strings.ReplaceAll(os.Args[1], "\\", "/")

	//
	pathToPackageDir := strings.ReplaceAll(goPath+"\\src\\"+os.Args[1], "\\", "/") + "/"

	//
	endpoints := strings.Split(pathToPackageDir, "/")

	//
	packageName := endpoints[len(endpoints)-2]

	packages, err := parser.ParseDir(
		fset,
		pathToPackageDir,
		nil,
		parser.ParseComments)

	if err != nil {
		panic(err)
	}

	// Phase 1 - Extract all the data from the given package
	commands, err := parseAnnotations(packages[packageName])
	if err != nil {
		panic(err)
	}

	// TODO: Phase 2 - Call a function that creates a cli tool
	err = generateCLITool(packageName, relativePathToPackage, commands)
	if err != nil {
		panic(err)
	}

	// TODO: Remove this block when Phase 2 is done.
	//for i, command := range commands {
	//	fmt.Printf("command #%d - %s\n", i, command.name)
	//	for j, argument := range command.arguments {
	//		fmt.Printf("\targument #%d - %s\n", j, argument.name)
	//	}
	//
	//	for j, option := range command.options {
	//		fmt.Printf("\targument #%d - --%s\n", j, option.name)
	//	}
	//}
}

func parseAnnotations(pkg *ast.Package) ([]*command, error) {
	commands := make([]*command, 0)

	for _, file := range pkg.Files {
		for _, decl := range file.Decls {
			if funcDecl, ok := decl.(*ast.FuncDecl); ok {
				tokens, err := lex(funcDecl.Doc.Text())
				if err != nil {
					return nil, err
				}

				for _, tok := range tokens {
					if tok.value == string(CLIGO_COMMAND) {
						//alias, err := parseCommandAnnotation(tokens, index)
						//if err != nil {
						//	return nil, err
						//}

						command := &command{
							funcDecl.Name.Name,
							"",
							make([]*argument, 0),
							make([]*option, 0),
						}

						commands = append(commands, command)
					}
				}
			}
		}
	}

	return commands, nil
}

func parseCommandAnnotation(tokens []*Token, pos int) (string, error) {
	// If the annotation contains parentheses, set the value in them a an alias of the command.
	if tokens[pos+1].value == "(" && tokens[pos+3].value == ")" {
		return tokens[pos+2].value, nil
	}

	return "", nil
}

func parseArgumentAnnotation(comment string) (*argument, error) {
	return nil, nil
}

func parseOptionAnnotation(comment string) (*option, error) {
	return nil, nil
}

func parseCommandAnnotation2(functionName, comment string) (*command, error) {
	tokens, err := lex(comment)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(tokens); i++ {
		fmt.Println(tokens[i].value)
	}

	return nil, nil
}

//func getAllExportedFunctions(pkg *ast.Package) []*ast.FuncDecl {
//	exportedFunctions := make([]*ast.FuncDecl, 0)
//
//	for _, file := range pkg.Files {
//		for _, decl := range file.Decls {
//			if funcDecl, ok := decl.(*ast.FuncDecl); ok {
//				if isExported(funcDecl.Name.Name) {
//					exportedFunctions = append(exportedFunctions, funcDecl)
//				}
//
//				annotations := getAnnotations(funcDecl.Name.Name, funcDecl.Doc)
//				for index, ann := range annotations {
//					fmt.Printf("annotation #%d - %s\n", index, *ann)
//				}
//			}
//		}
//	}
//
//	return exportedFunctions
//}
//
//func getAnnotations(functionName string, doc *ast.CommentGroup) []*annotation {
//	if doc == nil {
//		return nil
//	}
//
//	annotations := make([]*annotation, 0)
//
//	for _, comment := range doc.List {
//		switch {
//		case strings.Contains(comment.Text, string(CLIGO_COMMAND)):
//			{
//				annotations = append(annotations, &annotation{CLIGO_COMMAND, functionName, ""})
//				command, err := parseCommandAnnotation(functionName, comment.Text)
//				if err != nil {
//					return nil
//				}
//			}
//		case strings.Contains(comment.Text, string(CLIGO_ARGUMENT)):
//			{
//				annotations = append(annotations, &annotation{CLIGO_ARGUMENT, "", ""})
//			}
//		}
//	}
//
//	return annotations
//}
//
//func isExported(functionName string) bool {
//	return functionName[0] >= 'A' && functionName[0] <= 'Z'
//}
