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
	file_main, err := os.Create("main.go")
	defer file_main.Close()
	if err != nil {
		return err
	}
	file_main.WriteString(main_tmp)

	file_server, err := os.Create("server.go")
	defer file_server.Close()
	if err != nil {
		return err
	}
	file_server.WriteString(server_tmp)

	tmpl, err := template.New("go_mod").Parse(go_mod_tmp)
	if err != nil {
		return err
	}
	file_mod, err := os.Create("go.mod")
	defer file_mod.Close() // bruh do not check git history to see what mistake was here
	if err != nil {
		return err
	}
	goVersion, err := getGoVersion()
	if err != nil {
		return err
	}
	err = tmpl.Execute(file_mod, struct {
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
