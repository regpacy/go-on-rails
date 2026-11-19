package main

import (
	_ "embed"
	"os"
	"text/template"

	"github.com/spf13/cobra"
)

//go:embed templates/repo_implementation.go.tmp
var repoImplementationTemplate string

func createRepoImplementation(cmd *cobra.Command, args []string) error {
	structName := args[0]
	packageName := args[1]
	implementationName := firstLetterToUpper(packageName)
	filename := structToFilename(structName)
	if err := os.MkdirAll("repo/"+implementationName, 0o755); err != nil {
		return err
	}
	file, err := os.Create("repo/" + implementationName + "/" + filename)
	defer file.Close()
	if err != nil {
		return err
	}
	tmpl, err := template.New("repo_implmentation_tmp").Parse(repoImplementationTemplate)
	if err != nil {
		return err
	}
	err = tmpl.Execute(file, struct {
		StructName        string
		ImplementaionName string
		PackageName       string
	}{
		StructName:        structName,
		ImplementaionName: implementationName,
		PackageName:       packageName,
	})
	return err
}
