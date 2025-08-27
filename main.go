package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
)

func findCmdBinaries() ([]string, error) {
	binaries := []string{}
	cmdPath := "./cmd"

	if _, err := os.Stat(cmdPath); os.IsNotExist(err) {
		return binaries, nil
	}

	files, err := os.ReadDir(cmdPath)
	if err != nil {
		return nil, fmt.Errorf("error: could not read 'cmd' directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() {
			if _, err := os.Stat(filepath.Join(cmdPath, file.Name(), "main.go")); err == nil {
				binaries = append(binaries, file.Name())
			}
		}
	}
	return binaries, nil
}

func run(path string, args ...string) {
	cmdArgs := []string{"run", path}
	cmdArgs = append(cmdArgs, args...)

	cmd := exec.Command("go", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			os.Exit(exitError.ExitCode())
		}
		fmt.Printf("Error: Failed to execute command: %v\n", err)
		os.Exit(1)
	}
}

func printUsage(binaries []string) {
	fmt.Println("Multiple binaries found. Please specify which one to run.")
	fmt.Println("\nAvailable binaries:")
	for _, bin := range binaries {
		fmt.Printf("  %s\n", bin)
	}
	fmt.Println("\nUsage:")
	fmt.Println("  gr <binary_name> [args_for_program...]")
	fmt.Println("  gr --bin <binary_name> [args_for_program...]")
}

func main() {
	args := os.Args[1:]

	if _, err := os.Stat("main.go"); err == nil {
		run(".", args...)
		return
	}

	binaries, err := findCmdBinaries()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(binaries) == 0 {
		fmt.Println("Error: No 'main.go' file or runnable 'cmd/' packages found.")
		os.Exit(1)
	}

	var binToRun string
	var programArgs []string

	if len(args) > 0 {
		firstArg := args[0]
		isBinaryName := slices.Contains(binaries, firstArg)

		if firstArg == "--bin" {
			if len(args) > 1 {
				binToRun = args[1]
				programArgs = args[2:]
			}
		} else if isBinaryName {
			binToRun = firstArg
			programArgs = args[1:]
		} else {
			programArgs = args
		}
	} else {
		programArgs = []string{}
	}

	if binToRun == "" {
		if len(binaries) == 1 {
			binToRun = binaries[0]
		} else {
			printUsage(binaries)
			os.Exit(1)
		}
	}

	isValid := slices.Contains(binaries, binToRun)

	if !isValid {
		fmt.Printf("Error: Binary '%s' not found in 'cmd' directory.\n", binToRun)
		printUsage(binaries)
		os.Exit(1)
	}

	run(fmt.Sprintf("./cmd/%s/", binToRun), programArgs...)
}
