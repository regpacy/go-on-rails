package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"

	_ "embed"

	"github.com/spf13/cobra"
)

//go:embed templates/main.go.tmp
var main_tmp string

//go:embed templates/go.mod.tmp
var go_mod_tmp string

func main() {
	rootCmd := &cobra.Command{
		Use:   "go-on-rails",
		Short: "project and files generated for go web servers",
		Run: func(cmd *cobra.Command, args []string) {
			println(main_tmp)
		},
	}
	initCmd := &cobra.Command{
		Use:  "init <module_path>",
		Args: cobra.ExactArgs(1),
		RunE: initInit,
	}

	rootCmd.AddCommand(initCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func initInit(cmd *cobra.Command, args []string) error {
	modulePath := args[0]
	file, err := os.Create("main.go")
	if err != nil {
		return err
	}
	file.WriteString(main_tmp)
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
	return tmpl.Execute(file, struct {
		Version    string
		ModulePath string
	}{
		Version:    goVersion,
		ModulePath: modulePath,
	})
	return nil
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
