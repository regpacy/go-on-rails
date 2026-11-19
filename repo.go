package main

import (
	_ "embed"
	"os"
	"text/template"
)

//go:embed templates/repo_interface.go.tmp
var repoInterfaceTemplate string

func createRepoInterface(structName string) error {
	filename := structToFilename(structName)
	if err := os.MkdirAll("repo", 0o755); err != nil {
		return err
	}

	file, err := os.Create("repo/" + filename)
	if err != nil {
		return err
	}
	defer file.Close()
	tmpl, err := template.New("repo_interface_tmp").Parse(repoInterfaceTemplate)
	if err != nil {
		return err
	}
	err = tmpl.Execute(file, struct {
		StructName string
	}{
		StructName: structName,
	})
	return err
}
