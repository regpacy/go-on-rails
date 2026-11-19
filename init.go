package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

//go:embed templates/main.go
var mainTemplate string

//go:embed templates/go.mod.tmp
var goModTemplate string

//go:embed templates/server.go.tmp
var serverTemplate string

func initProject(modulePath string) error {
	mainFile, err := os.Create("main.go")
	defer mainFile.Close()
	if err != nil {
		return err
	}
	mainFile.WriteString(mainTemplate)

	serverFile, err := os.Create("server.go")
	defer serverFile.Close()
	if err != nil {
		return err
	}
	serverFile.WriteString(serverTemplate)

	tmpl, err := template.New("go_mod").Parse(goModTemplate)
	if err != nil {
		return err
	}
	modFile, err := os.Create("go.mod")
	defer modFile.Close() // bruh do not check git history to see what mistake was here
	if err != nil {
		return err
	}
	goVersion, err := getGoVersion()
	if err != nil {
		return err
	}
	err = tmpl.Execute(modFile, struct {
		Version    string
		ModulePath string
	}{
		Version:    goVersion,
		ModulePath: modulePath,
	})
	return err
}

func getGoVersion() (string, error) {
	out, err := exec.Command("go", "version").Output()
	if err != nil {
		return "", err
	}

	fields := strings.Fields(string(out))
	if len(fields) < 3 {
		return "", fmt.Errorf("unexpected go version output: %q", out)
	}

	return strings.TrimPrefix(fields[2], "go"), nil
}
