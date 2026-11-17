package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "go-on-rails",
		Short: "project and files generated for go web servers",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Hello from Cobra!")
		},
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
