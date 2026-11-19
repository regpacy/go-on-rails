package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "go-on-rails",
		Short: "project and files generated for go web servers",
		Run: func(cmd *cobra.Command, args []string) {
			println(mainTemplate)
		},
		SilenceUsage: true,
	}
	initCmd := &cobra.Command{
		Use:  "init <module_path>",
		Args: cobra.ExactArgs(1),
		RunE: initProject,
	}
	pathCmd := &cobra.Command{
		Use:  "path <path> <function name>",
		Args: cobra.ExactArgs(2),
		RunE: createPath,
	}

	domainCmd := &cobra.Command{
		Use:  "domain <struct name>",
		Args: cobra.ExactArgs(1),
		RunE: createDomain,
	}

	repoCmd := &cobra.Command{
		Use:  "repo_interface <struct name for interface>",
		Args: cobra.ExactArgs(1),
		RunE: createRepoInterface,
	}

	repoImplementationCmd := &cobra.Command{
		Use:  "repo_implementation <struct name for interface> <implementation name>",
		Args: cobra.ExactArgs(2),
		RunE: createRepoImplementation,
	}
	rootCmd.AddCommand(initCmd, pathCmd, domainCmd, repoCmd, repoImplementationCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
	output, err := exec.Command("go", "mod", "tidy").CombinedOutput()
	if err != nil {
		fmt.Printf("go mod tidy failed: %v\n%s", err, output)
	}
}
