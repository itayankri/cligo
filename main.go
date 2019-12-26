package main

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

func main() {
	fset := token.NewFileSet() // positions are relative to fset

	// TODO: fix later
	goPath := os.Getenv("GOPATH")

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

	// Phase 2 - Call a function that creates a cli tool
	err = generateCLITool(packageName, relativePathToPackage, commands)
	if err != nil {
		panic(err)
	}
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
						if isExported(funcDecl.Name.Name) {
							command := &command{
								strings.ToLower(funcDecl.Name.Name),
								funcDecl.Name.Name,
								make([]*argument, 0),
							}

							for _, argList := range funcDecl.Type.Params.List {
								if _type, ok := argList.Type.(*ast.Ident); ok {
									for _, arg := range argList.Names {
										//fmt.Println(arg)
										argument := &argument{
											arg.Name,
											_type.Name,
										}
										command.arguments = append(command.arguments, argument)
									}
								} else {
									return nil, errors.New("cannot create a sub-command based on a function that " +
										"requires a non-atomic argument. function name: " + funcDecl.Name.Name)
								}
							}

							commands = append(commands, command)
						} else {
							return nil, errors.New("cannot create a sub-command based on an unexported function. " +
								"function name: " + funcDecl.Name.Name)
						}
					}
				}
			}
		}
	}

	return commands, nil
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
func isExported(functionName string) bool {
	return functionName[0] >= 'A' && functionName[0] <= 'Z'
}
