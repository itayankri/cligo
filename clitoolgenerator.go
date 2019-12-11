package main

import (
	"fmt"
)

func generateCLITool(commands []*command) error {
	return nil
}

func generatePackage(packageName string) string {
	return fmt.Sprintf("package %s\n", packageName)
}

func generateImports(pathToPackage string) string {
	return fmt.Sprintf(`
		import (
			"flag"
			"%s"
		)
	`, pathToPackage)
}

func generateMain(commands []*command) string {
	return fmt.Sprintf(`
		func main() {
			// Commands

			// Flag-Pointers

			// Switch-Case
		}
	`)
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

func generateSwitch(commands []*command) string {
	return ""
}

func generateCase(commands []*command) string {
	return ""
}

func generateDefaultCase() string {
	return ""
}

func generateManPage(commands []*command) string {
	return ""
}
