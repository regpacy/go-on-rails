package main

import "github.com/spf13/cobra"

func newRootCommand() *cobra.Command {
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
		RunE: func(_ *cobra.Command, args []string) error {
			return initProject(args[0])
		},
	}
	pathCmd := &cobra.Command{
		Use:  "path <path> <function name>",
		Args: cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			return createPath(args[0], args[1])
		},
	}

	domainCmd := &cobra.Command{
		Use:  "domain <struct name>",
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return createDomain(args[0])
		},
	}

	repoCmd := &cobra.Command{
		Use:  "repo_interface <struct name for interface>",
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return createRepoInterface(args[0])
		},
	}

	repoImplementationCmd := &cobra.Command{
		Use:  "repo_implementation <struct name for interface> <implementation name>",
		Args: cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			return createRepoImplementation(args[0], args[1])
		},
	}
	rootCmd.AddCommand(initCmd, pathCmd, domainCmd, repoCmd, repoImplementationCmd)
	return rootCmd
}
