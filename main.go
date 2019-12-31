package main

import (
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	cliToolGeneratorConfig := handleCli()

	// TODO: fix later
	goPath := os.Getenv("GOPATH")

	// Calculate the absolute path of the package.
	pathToPackageDir := filepath.Join(goPath, "src", cliToolGeneratorConfig.importString)

	// Check if the path to the package directory does not exist.
	if _, err := os.Stat(pathToPackageDir); os.IsNotExist(err) {
		handleError(err)
	}

	// Crete a new file set.
	fset := token.NewFileSet()

	// Use the golang parser to parse the code in the package directory.
	packages, err := parser.ParseDir(
		fset,
		pathToPackageDir,
		nil,
		parser.ParseComments)

	if err != nil {
		handleError(errors.Wrap(err, "failed to parse golang package"))
	}

	// Phase 1 - Extract all the data from the given package
	commands, err := parseAnnotations(packages[cliToolGeneratorConfig.packageName])
	if err != nil {
		handleError(err)
	}

	// Phase 2 - Call a function that creates a cli tool
	err = generateCLITool(cliToolGeneratorConfig, commands)
	if err != nil {
		handleError(err)
	}
}

func handleError(err error) {
	fmt.Println("Cligo Error: ", err.Error())
	os.Exit(2)
}

func handleCli() CliToolGeneratorConfig {
	// If the user did now provide one argument, exit the program.
	if len(os.Args) < 2 {
		fmt.Println("Path to Golang package must be provided. Use -help for more information")
		os.Exit(1)
	}

	// A variable that will be populated by the flags "o" or "output"
	var (
		isHelp     bool
		isVerbose  bool
		outputPath string
	)

	currentWorkingDirectory, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	// Parse the -output flag into outputPath
	flag.StringVar(&outputPath,
		"output",
		currentWorkingDirectory,
		"A path which the cli tool will be written to.")

	// Parse the -verbose flag into isVerbose
	flag.BoolVar(&isVerbose,
		"verbose",
		false,
		"Print to console progress messages.")

	// Parse the -help flag into isHelp
	flag.BoolVar(&isHelp,
		"help",
		false,
		"Usage Explanation.")

	// Parse the Arg list into the flag variables.
	flag.Parse()

	// Get the path to the golang package which is the input of our cli tool generator.
	importString := strings.ReplaceAll(flag.Arg(0), "\\", "/")

	// Get the last element in the path, which is the name of the package.
	packageName := filepath.Base(importString)

	// If the user used the help flag print usage page, not matter what other
	// flags he used.
	if isHelp {
		printUsage()
		os.Exit(0)
	}

	return CliToolGeneratorConfig{
		outputPath:   outputPath,
		packageName:  packageName,
		importString: importString,
		verbose:      isVerbose,
	}
}

func printUsage() {
	fmt.Println()
	fmt.Println("A CLI Tool that generates CLI Tools from golang packages.")
	fmt.Println()
	fmt.Println("Usage of CLIGO: ")
	fmt.Println("\tcligo [OPTIONS] <IMPORT STRING>")
	fmt.Println()
	fmt.Println("Options:")
	flag.PrintDefaults()
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

func isExported(functionName string) bool {
	return functionName[0] >= 'A' && functionName[0] <= 'Z'
}
