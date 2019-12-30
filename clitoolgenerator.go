package main

import (
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

func generateCLITool(outputPath, packageName, pathToPackage string, verbose bool, commands []*command) error {
	fmt.Println("--> Generating Cli Tool code")
	file := fmt.Sprintf(`
// Package declaration
%s

// Imports
%s

// Main func
%s
	`, generatePackage(packageName), generateImports(pathToPackage), generateMain(packageName, commands))

	goFilePath := filepath.Join(outputPath, packageName+"_cligo.go")
	cliToolPath := filepath.Join(outputPath, packageName+"-cli.exe")

	if verbose {
		fmt.Println("--> Writing Cli Tool code to ", goFilePath)
	}

	err := ioutil.WriteFile(goFilePath, []byte(file), 0644)
	if err != nil {
		return errors.Wrap(err, "failed to write go file")
	}

	cmd := exec.Command(
		"go",
		"build",
		"-o",
		cliToolPath,
		goFilePath)

	if verbose {
		fmt.Println("--> Compiling Cli Tool executable at ", outputPath)
	}

	err = cmd.Run()
	if err != nil {
		return errors.Wrap(err, "failed to compile go file")
	}

	if verbose {
		fmt.Println("--> Removing go file")
	}

	err = os.Remove(goFilePath)
	if err != nil {
		return errors.Wrap(err, "failed to remove go code")
	}

	if verbose {
		fmt.Println("--> Cli Tool generated successfully")
	}

	return nil
}

func generatePackage(packageName string) string {
	return fmt.Sprintf("package main")
}

func generateImports(pathToPackage string) string {
	return fmt.Sprintf(`
import (
	"flag"
	"fmt"
	"os"
	"%s"
)
	`, pathToPackage)
}

func generateMain(packageName string, commands []*command) string {
	return fmt.Sprintf(`
func main() {
	if len(os.Args) < 2 {
        fmt.Println("subcommand is required")
        os.Exit(1)
    }

	// Commands
	%s

	// Flag-Pointers
	%s

	// Switch-Case
	%s
}
	`, generateCommands(commands), generateFlags(commands), generateSwitch(packageName, commands))
}

func generateCommands(commands []*command) string {
	cmds := ""

	for _, cmd := range commands {
		cmds += fmt.Sprintf("%s := flag.NewFlagSet(\"%s\", flag.ExitOnError)\n", cmd.name, cmd.name)
	}

	return cmds
}

func generateFlags(commands []*command) string {
	flags := ""

	for _, cmd := range commands {
		for _, arg := range cmd.arguments {
			flags += generateFlag(cmd.name, *arg)
		}
	}

	return flags
}

func generateFlag(cmdName string, arg argument) string {
	switch arg._type {
	case "int":
		{
			return fmt.Sprintf("%s_%s := %s.Int(\"%s\", 0, \"Explanation here\")\n",
				cmdName,
				arg.name,
				cmdName,
				arg.name)
		}
	case "int64":
		{
			return fmt.Sprintf("%s_%s := %s.Int64(\"%s\", 0, \"Explanation here\")\n",
				cmdName,
				arg.name,
				cmdName,
				arg.name)
		}
	case "uint":
		{
			return fmt.Sprintf("%s_%s := %s.Uint(\"%s\", 0, \"Explanation here\")\n",
				cmdName,
				arg.name,
				cmdName,
				arg.name)
		}
	case "uint64":
		{
			return fmt.Sprintf("%s_%s := %s.Uint64(\"%s\", 0, \"Explanation here\")\n",
				cmdName,
				arg.name,
				cmdName,
				arg.name)
		}
	case "float64":
		{
			return fmt.Sprintf("%s_%s := %s.Float64(\"%s\", 0, \"Explanation here\")\n",
				cmdName,
				arg.name,
				cmdName,
				arg.name)
		}
	case "string":
		{
			return fmt.Sprintf("%s_%s := %s.String(\"%s\", \"\", \"Explanation here\")\n",
				cmdName,
				arg.name,
				cmdName,
				arg.name)
		}
	case "bool":
		{
			return fmt.Sprintf("%s_%s := %s.Bool(\"%s\", false, \"Explanation here\")\n",
				cmdName,
				arg.name,
				cmdName,
				arg.name)
		}
	}

	return ""
}

func generateSwitch(packageName string, commands []*command) string {
	return fmt.Sprintf("switch os.Args[1] {\n" + generateCases(packageName, commands) + "\n}")
}

func generateCases(packageName string, commands []*command) string {
	cases := ""

	for _, cmd := range commands {
		cases += generateCase(packageName, *cmd)
	}

	cases += generateHelpCase(packageName, commands)
	cases += generateDefaultCase()

	return cases
}

func generateCase(packageName string, command command) string {
	return fmt.Sprintf(`
case "%s":
	{
		err := %s.Parse(os.Args[2:])
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		%s.%s(`+generateArguments(command)+`)
	}
	`, command.name, command.name, packageName, command.funcName)
}

func generateArguments(command command) string {
	if command.arguments == nil || len(command.arguments) == 0 {
		return ""
	}

	arguments := ""

	for _, arg := range command.arguments {
		arguments += "*" + command.name + "_" + arg.name + ","
	}

	return arguments[:len(arguments)-1]
}

func generateDefaultCase() string {
	return fmt.Sprintf(`
default:
	{
		//panic(errors.New("unrecognized command - " + os.Args[1]))
		fmt.Println("unrecognized command - " + os.Args[1] + ", use --help for more information")
		os.Exit(1)
	}
	`)
}

func generateHelpCase(packageName string, commands []*command) string {
	return fmt.Sprintf(`
case "--help", "-h":
{
	fmt.Println(`+"`%s`"+`)
}
`, generateManPage(packageName, commands))
}

func generateManPage(packageName string, commands []*command) string {
	man := fmt.Sprintf("\n [This CLI Tool has been generated by CLIGO]\n\n %s:\n", packageName)

	for index, cmd := range commands {
		man += fmt.Sprintf("    %s\t%s %d\n", cmd.name, "description", index)
	}

	return man
}

func generateParseCheckBlock(command command) string {
	return ""
}
