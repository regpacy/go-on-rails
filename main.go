package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	if err := newRootCommand().Execute(); err != nil {
		os.Exit(1)
	}
	output, err := exec.Command("go", "mod", "tidy").CombinedOutput()
	if err != nil {
		fmt.Printf("go mod tidy failed: %v\n%s", err, output)
	}
}
