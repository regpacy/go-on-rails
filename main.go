package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"strings"
	"text/template"

	_ "embed"

	"github.com/spf13/cobra"
)

//go:embed templates/main.go
var main_tmp string

//go:embed templates/go.mod.tmp
var go_mod_tmp string

//go:embed templates/server.go.tmp
var server_tmp string

func main() {
	rootCmd := &cobra.Command{
		Use:   "go-on-rails",
		Short: "project and files generated for go web servers",
		Run: func(cmd *cobra.Command, args []string) {
			println(main_tmp)
		},
		SilenceUsage: true,
	}
	initCmd := &cobra.Command{
		Use:  "init <module_path>",
		Args: cobra.ExactArgs(1),
		RunE: initInit,
	}
	pathCmd := &cobra.Command{
		Use:  "path <path> <function name>",
		Args: cobra.ExactArgs(2),
		RunE: createPath,
	}

	rootCmd.AddCommand(initCmd, pathCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
	output, err := exec.Command("go", "mod", "tidy").CombinedOutput()
	if err != nil {
		fmt.Printf("go mod tidy failed: %w\n%s", err, output)
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

func createPath(cmd *cobra.Command, args []string) error {
	path := args[0]
	functionName := args[1]

	err := exec.Command("go", "build", "-o", "/dev/null", "./...").Run()
	if err != nil {
		return errors.New("Project Should be compilable before adding path")
	}
	src, err := os.ReadFile("server.go")
	if err != nil {
		return err
	}

	src, err = addRouter(path, functionName, src)
	if err != nil {
		return err
	}
	err = addFunc(functionName, src)
	return err
}

func addFunc(functionName string, src []byte) error {
	newFunc := fmt.Sprintf(`
	func %s(c *gin.Context) {
    c.JSON(200, gin.H{})
	}
	`, functionName)

	updated := append(src, []byte("\n"+newFunc)...)

	formatted, err := format.Source(updated)
	if err != nil {
		return err
	}

	return os.WriteFile("server.go", formatted, 0o644)
}

func addRouter(path string, functionName string, src []byte) ([]byte, error) {
	fset := token.NewFileSet()

	file, err := parser.ParseFile(fset, "server.go", src, 0)
	if err != nil {
		return nil, err
	}
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Name.Name != "registerRoutes" {
			continue
		}

		insertPos := fset.Position(fn.Body.Rbrace).Offset
		updated := string(src[:insertPos]) +
			"\n\ts.r.GET(\"/" + path + "\", " + functionName + ")\n" +
			string(src[insertPos:])
		return []byte(updated), nil
	}

	return nil, fmt.Errorf("registerRoutes function not found")
}
