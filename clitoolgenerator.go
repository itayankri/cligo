package main

import (
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

type CliToolGeneratorConfig struct {
	outputPath   string
	packageName  string
	importString string
	verbose      bool
}

func generateCLITool(config CliToolGeneratorConfig, commands []*command) error {
	if config.verbose {
		fmt.Println("--> Generating Cli Tool code")
	}

	file := fmt.Sprintf(`
// Package declaration
%s

// Imports
%s

// Main func
%s
	`, generatePackage(config.packageName), generateImports(config.importString), generateMain(config.packageName, commands))

	goFilePath := filepath.Join(config.outputPath, config.packageName+"_cligo.go")
	cliToolPath := filepath.Join(config.outputPath, config.packageName+"-cli.exe")

	if config.verbose {
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

	if config.verbose {
		fmt.Println("--> Compiling Cli Tool executable at ", config.outputPath)
	}

	err = cmd.Run()
	if err != nil {
		return errors.Wrap(err, "failed to compile go file")
	}

	if config.verbose {
		fmt.Println("--> Removing go file")
	}

	err = os.Remove(goFilePath)
	if err != nil {
		return errors.Wrap(err, "failed to remove go code")
	}

	if config.verbose {
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
	"strconv"
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
		for _, opt := range cmd.options {
			flags += generateFlag(cmd.name, *opt)
		}
	}

	return flags
}

func generateCommandArguments(command command) string {
	commandArguments := ""

	for index, arg := range command.arguments {
		commandArguments += generateCommandArgument(index, command.name, *arg)
	}

	return commandArguments
}

func generateCommandArgument(argIndex int, cmdName string, arg argument) string {
	switch arg._type {
	case "int64", "int":
		{
			return fmt.Sprintf(`%s_%s, err := strconv.ParseInt(%s.Arg(%d), 10, 64)
	if err != nil {
		fmt.Println("arg %s expected to be of type int")
		os.Exit(2)
	}
`,
				cmdName, arg.name, cmdName, argIndex, arg.name)
		}
	case "float64":
		{
			return fmt.Sprintf(`%s_%s, err := strconv.ParseFloat(%s.Arg(%d), 64)
	if err != nil {
		fmt.Println("arg %s expected to be of type float")
		os.Exit(2)
	}
`,
				cmdName, arg.name, cmdName, argIndex, arg.name)
		}
	case "string":
		{
			return fmt.Sprintf(`%s_%s := %s.Arg(%d)`, cmdName, arg.name, cmdName, argIndex)
		}
	case "bool":
		{
			return fmt.Sprintf(`%s_%s, err := strconv.ParseBool(%s.Arg(%d))
	if err != nil {
		fmt.Println("arg %s expected to be of type bool")
		os.Exit(2)
	}
`,
				cmdName, arg.name, cmdName, argIndex, arg.name)
		}
	}

	return ""
}

func generateFlag(cmdName string, opt option) string {
	switch opt._type {
	case "int":
		{
			return fmt.Sprintf("%s_%s := %s.Int(\"%s\", 0, \"%s\")\n",
				cmdName,
				opt.name,
				cmdName,
				opt.name,
				opt.description)
		}
	case "int64":
		{
			return fmt.Sprintf("%s_%s := %s.Int64(\"%s\", 0, \"%s\")\n",
				cmdName,
				opt.name,
				cmdName,
				opt.name,
				opt.description)
		}
	case "uint":
		{
			return fmt.Sprintf("%s_%s := %s.Uint(\"%s\", 0, \"%s\")\n",
				cmdName,
				opt.name,
				cmdName,
				opt.name,
				opt.description)
		}
	case "uint64":
		{
			return fmt.Sprintf("%s_%s := %s.Uint64(\"%s\", 0, \"%s\")\n",
				cmdName,
				opt.name,
				cmdName,
				opt.name,
				opt.description)
		}
	case "float64":
		{
			return fmt.Sprintf("%s_%s := %s.Float64(\"%s\", 0, \"%s\")\n",
				cmdName,
				opt.name,
				cmdName,
				opt.name,
				opt.description)
		}
	case "string":
		{
			return fmt.Sprintf("%s_%s := %s.String(\"%s\", \"\", \"%s\")\n",
				cmdName,
				opt.name,
				cmdName,
				opt.name,
				opt.description)
		}
	case "bool":
		{
			return fmt.Sprintf("%s_%s := %s.Bool(\"%s\", false, \"%s\")\n",
				cmdName,
				opt.name,
				cmdName,
				opt.name,
				opt.description)
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

		if %s.Parsed() {
			%s
			%s.%s(`+generateArguments(command)+`)
		}
	}
	`, command.name, command.name, command.name, generateCommandArguments(command), packageName, command.funcName)
}

func generateArguments(command command) string {
	if command.options == nil && command.arguments == nil {
		return ""
	}

	arguments := ""

	for _, opt := range command.options {
		arguments += "*" + command.name + "_" + opt.name + ","
	}

	for _, arg := range command.arguments {
		arguments += command.name + "_" + arg.name + ","
	}

	return arguments[:len(arguments)-1]
}

func generateDefaultCase() string {
	return fmt.Sprintf(`
default:
	{
		fmt.Println("unrecognized command - " + os.Args[1] + ", use --help for more information")
		os.Exit(1)
	}
	`)
}

func generateHelpCase(packageName string, commands []*command) string {
	return fmt.Sprintf(`
case "--help", "-h":
{
	%s
}
`, generateManPage(packageName, commands))
}

func generateManPage(packageName string, commands []*command) string {
	man := fmt.Sprintf("fmt.Println(\"\\n [This CLI Tool has been generated by CLIGO]\\n\\n Usage of %s:\\n\")",
		packageName)

	for _, cmd := range commands {
		man += fmt.Sprintf("\nfmt.Println(\"\\n%s\\t%s\")", cmd.name, cmd.description)
		man += fmt.Sprintf("\n%s.PrintDefaults()", cmd.name)
	}

	return man
}

func generateParseCheckBlock(command command) string {
	return ""
}
