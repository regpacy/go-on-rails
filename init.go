package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

//go:embed templates/main.go
var main_tmp string

//go:embed templates/go.mod.tmp
var go_mod_tmp string

//go:embed templates/server.go.tmp
var server_tmp string

func initInit(cmd *cobra.Command, args []string) error {
	modulePath := args[0]
	file, err := os.Create("main.go")
	if err != nil {
		return err
	}
	file.WriteString(main_tmp)
	file.Close()

	file, err = os.Create("server.go")
	if err != nil {
		return err
	}
	file.WriteString(server_tmp)
	file.Close()

	tmpl, err := template.New("go_mod").Parse(go_mod_tmp)
	if err != nil {
		return err
	}
	file, err = os.Create("go.mod")
	if err != nil {
		return err
	}
	goVersion, err := getGoVersion()
	if err != nil {
		return err
	}
	defer file.Close()
	err = tmpl.Execute(file, struct {
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
