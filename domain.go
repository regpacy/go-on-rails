package main

import (
	_ "embed"
	"os"
	"strings"
	"text/template"
	"unicode/utf8"

	"github.com/spf13/cobra"
)

//go:embed templates/domain.go.tmp
var domain_tmp string

func createDomain(cmd *cobra.Command, args []string) error {
	structName := args[0]
	fileName := StructToFilename(structName)

	// Creates the directory if it doesn't exist
	if err := os.MkdirAll("domain", 0o755); err != nil {
		return err
	}

	file, err := os.Create("domain/" + fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	tmpl, err := template.New("domain_tmp").Parse(domain_tmp)
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

func StructToFilename(s string) string {
	if s == "" {
		return s
	}

	_, size := utf8.DecodeRuneInString(s)
	return strings.ToLower(s[:size]) + s[size:] + ".go"
}
