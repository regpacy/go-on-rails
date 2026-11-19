package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

//go:embed templates/main.go.tmp
var mainTemplate string

//go:embed templates/go.mod.tmp
var goModTemplate string

//go:embed templates/server.go.tmp
var serverTemplate string

//go:embed templates/config.go.tmp
var configTemplate string

//go:embed templates/gitignore.tmp
var gitignoreTemplate string

func initProject(modulePath string) error {
	mainFile, err := os.Create("main.go")
	if err != nil {
		return err
	}
	defer mainFile.Close()
	mainFile.WriteString(mainTemplate)

	serverFile, err := os.Create("server.go")
	if err != nil {
		return err
	}
	defer serverFile.Close()
	serverFile.WriteString(serverTemplate)

	configFile, err := os.Create("config.go")
	if err != nil {
		return err
	}
	defer configFile.Close()
	configFile.WriteString(configTemplate)

	gitignoreFile, err := os.Create(".gitignore")
	if err != nil {
		return err
	}
	defer gitignoreFile.Close()
	gitignoreFile.WriteString(gitignoreTemplate)

	dotenv, err := os.Create(".env")
	if err != nil {
		return err
	}
	dotenv.Close()
	tmpl, err := template.New("go_mod").Parse(goModTemplate)
	if err != nil {
		return err
	}
	modFile, err := os.Create("go.mod")
	if err != nil {
		return err
	}
	defer modFile.Close() // bruh do not check git history to see what mistake was here
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
