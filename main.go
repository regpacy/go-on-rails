package main

import (
	"fmt"
	"os"

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
		Use:  "init",
		RunE: initInit,
	}

	rootCmd.AddCommand(initCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func initInit(cmd *cobra.Command, args []string) error {
	file, err := os.Create("main.go")
	if err != nil {
		return err
	}
	file.WriteString(main_tmp)
	file.Close()

	file, err = os.Create("go.mod")
	if err != nil {
		return err
	}
	file.WriteString(go_mod_tmp)
	file.Close()
	return nil
}
