package main

import (
	_ "embed"
	"os"
	"text/template"
)

//go:embed templates/repo_implementation.go.tmp
var repoImplementationTemplate string

func createRepoImplementation(structName, packageName string) error {
	implementationName := firstLetterToUpper(packageName)
	filename := structToFilename(structName)
	if err := os.MkdirAll("repo/"+implementationName, 0o755); err != nil {
		return err
	}
	file, err := os.Create("repo/" + implementationName + "/" + filename)
	if err != nil {
		return err
	}
	defer file.Close()
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
