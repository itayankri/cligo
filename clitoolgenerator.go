package main

import (
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"os/exec"
)

func generateCLITool(packageName, pathToPackage string, commands []*command) error {
	file := fmt.Sprintf(`
// Package declaration
%s

// Imports
%s

// Main func
%s
	`, generatePackage(packageName), generateImports(pathToPackage), generateMain(packageName, commands))

	err := ioutil.WriteFile("executables/"+packageName+"_cligo.go", []byte(file), 0644)
	if err != nil {
		return errors.Wrap(err, "failed to write go file")
	}

	cmd := exec.Command("go", "build", "-o", "executables/"+packageName+"_cligo.exe", "executables/"+packageName+"_cligo.go")

	err = cmd.Run()
	if err != nil {
		return errors.Wrap(err, "failed to compile go file")
	}

	return nil
}

func generatePackage(packageName string) string {
	return fmt.Sprintf("package main")
}

func generateImports(pathToPackage string) string {
	return fmt.Sprintf(`
import (
	"errors"
	"flag"
	"os"
	"%s"
)
	`, pathToPackage)
}

func generateMain(packageName string, commands []*command) string {
	return fmt.Sprintf(`
func main() {
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
	case "Int":
		{
			return fmt.Sprintf("%s := \"%s\".Int(\"%s\", 0, \"Explanation here\")\n",
				cmdName,
				arg.name,
				arg.name)
		}
	case "Int64":
		{
			return fmt.Sprintf("%s := \"%s\".Int64(\"%s\", 0, \"Explanation here\")\n",
				cmdName,
				arg.name,
				arg.name)
		}
	case "Uint":
		{
			return fmt.Sprintf("%s := \"%s\".Uint(\"%s\", 0, \"Explanation here\")\n",
				cmdName,
				arg.name,
				arg.name)
		}
	case "Uint64":
		{
			return fmt.Sprintf("%s := \"%s\".Uint64(\"%s\", 0, \"Explanation here\")\n",
				cmdName,
				arg.name,
				arg.name)
		}
	case "Float64":
		{
			return fmt.Sprintf("%s := \"%s\".Float64(\"%s\", 0, \"Explanation here\")\n",
				cmdName,
				arg.name,
				arg.name)
		}
	case "String":
		{
			return fmt.Sprintf("%s := \"%s\".String(\"%s\", \"\", \"Explanation here\")\n",
				cmdName,
				arg.name,
				arg.name)
		}
	case "Bool":
		{
			return fmt.Sprintf("%s := \"%s\".Bool(\"%s\", false, \"Explanation here\")\n",
				cmdName,
				arg.name,
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

	cases += generateDefaultCase()

	return cases
}

func generateCase(packageName string, command command) string {
	return fmt.Sprintf(`
case "%s":
	{
		%s.Parse(os.Args[2:])
		%s.%s()
	}
	`, command.name, command.name, packageName, command.name)
}

func generateDefaultCase() string {
	return fmt.Sprintf(`
default:
	{
		panic(errors.New("unrecognized command - " + os.Args[1]))
	}
	`)
}

func generateManPage(commands []*command) string {
	return "Man page"
}

func generateParseCheckBlock(command command) string {
	return ""
}
