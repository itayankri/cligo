package main

import "fmt"

func generatePackage(packageName string) string {
	return "package " + packageName
}

func generateImports(pathToPackage string) string {
	return fmt.Sprintf(`
		import (
			"fmt"
			%s
		)
	`, pathToPackage)
}
