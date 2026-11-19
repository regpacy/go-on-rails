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

	"github.com/spf13/cobra"
)

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
    c.JSON(http.StatusOK, gin.H{})
	}
	`, functionName)

	updated := append(src, []byte("\n"+newFunc)...)

	formatted, err := format.Source(updated)
	if err != nil {
		return err
	}

	err = os.WriteFile("server.go", formatted, 0o644)
	if err != nil {
		return err
	}
	err = organizeImports("server.go")
	return err
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

func organizeImports(path string) error {
	output, err := exec.Command("gopls", "codeaction", "-kind=source.organizeImports", "-exec", "-write", path).CombinedOutput()
	if err != nil && !strings.HasPrefix(string(output), "gopls: no matching code action at ") {
		return fmt.Errorf("gopls organize imports failed: %w\n%s", err, output)
	}
	return nil
}
