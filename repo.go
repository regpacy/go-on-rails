package main

import (
	_ "embed"
	"os"
	"text/template"

	"github.com/spf13/cobra"
)

//go:embed templates/repo_interface.go.tmp
var repoInterfaceTmp string

func createRepoInterface(cmd *cobra.Command, args []string) error {
	structName := args[0]
	fileName := StructToFilename(structName)
	if err := os.MkdirAll("repo", 0o755); err != nil {
		return err
	}

	file, err := os.Create("repo/" + fileName)
	defer file.Close()
	if err != nil {
		return err
	}
	tmpl, err := template.New("repo_interface_tmp").Parse(repoInterfaceTmp)
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
