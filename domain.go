package main

import (
	_ "embed"
	"os"
	"text/template"
)

//go:embed templates/domain.go.tmp
var domainTemplate string

func createDomain(structName string) error {
	filename := structToFilename(structName)

	// Creates the directory if it doesn't exist
	if err := os.MkdirAll("domain", 0o755); err != nil {
		return err
	}

	file, err := os.Create("domain/" + filename)
	if err != nil {
		return err
	}
	defer file.Close()

	tmpl, err := template.New("domain_tmp").Parse(domainTemplate)
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
