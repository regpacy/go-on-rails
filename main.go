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

//go:embed templates/server.go
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
	if err != nil {
		return err
	}
	output, err := exec.Command("go", "get", "github.com/gin-gonic/gin").CombinedOutput()
	if err != nil {
		return fmt.Errorf("go get failed: %w\n%s", err, output)
	}
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

func createPath(cmd *cobra.Command, args []string) error {
	path := args[0]
	functionName := args[1]

	err := exec.Command("go", "build", "-o", "/dev/null", "./...").Run()
	if err != nil {
		return errors.New("Project Should be compilable before adding path")
	}
	src, err := os.ReadFile("main.go")

	err = addRouter(path, functionName, src)
	if err != nil {
		return err
	}
	err = addFunc(functionName, src)
	return err
}

func addFunc(functionName string, src []byte) error {
	newFunc := fmt.Sprintf(`
	func %s(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{})
	}
	`, functionName)

	updated := append(src, []byte("\n"+newFunc)...)

	formatted, err := format.Source(updated)
	if err != nil {
		return err
	}

	return os.WriteFile("main.go", formatted, 0o644)
}

func addRouter(path string, functionName string, src []byte) error {
	fset := token.NewFileSet()

	file, err := parser.ParseFile(fset, "main.go", src, 0)
	if err != nil {
		return err
	}
	// Find main()
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Name.Name != "main" {
			continue
		}

		// Find the last route registration
		insertPos := -1

		for _, stmt := range fn.Body.List {
			exprStmt, ok := stmt.(*ast.ExprStmt)
			if !ok {
				continue
			}

			call, ok := exprStmt.X.(*ast.CallExpr)
			if !ok {
				continue
			}

			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				continue
			}

			receiver, ok := sel.X.(*ast.Ident)
			if !ok || receiver.Name != "r" {
				continue
			}

			switch sel.Sel.Name {
			case "GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS":
				insertPos = fset.Position(stmt.End()).Offset
			}
		}

		if insertPos == -1 {
			return fmt.Errorf("no route handlers found")
		}

		// Insert text after the last handler
		updated := string(src[:insertPos]) +
			"\n\n\t" +
			"r.GET(\"/" + path + "\", " + functionName + ")" +
			string(src[insertPos:])

		formatted, err := format.Source([]byte(updated))
		if err != nil {
			return err
		}
		src = formatted
		return os.WriteFile("main.go", formatted, 0o644)
	}

	return fmt.Errorf("main function not found")
}
