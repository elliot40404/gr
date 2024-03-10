package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	args := os.Args[1:]
	bin := flag.String("bin", "", "Binary to run")
	flag.Parse()
	// check if main.go exists
	if _, err := os.Stat("main.go"); !os.IsNotExist(err) {
		run([]string{"go", "run", "main.go"}, args...)
	} else {
		// look for cmd directory
		if _, err := os.Stat("cmd"); !os.IsNotExist(err) {
			//list all directories in cmd
			files, err := os.ReadDir("cmd")
			if err != nil {
				fmt.Println(err)
				return
			}
			binaries := make([]string, 0)
			for _, file := range files {
				if file.IsDir() {
					binaries = append(binaries, file.Name())
				}
			}
			if len(binaries) == 0 {
				fmt.Println("No binaries found in cmd directory")
				return
			}
			if len(binaries) == 1 {
				run([]string{"go", "run", fmt.Sprintf("cmd/%s/main.go", binaries[0])}, args...)
				return
			}
			switch len(binaries) {
			case 0:
				fmt.Println("No binaries found in cmd directory")
			case 1:
				run([]string{"go", "run", fmt.Sprintf("cmd/%s/main.go", binaries[0])}, args...)
			default:
				if *bin != "" {
					found := false
					for _, binary := range binaries {
						if binary == *bin {
							found = true
							run([]string{"go", "run", fmt.Sprintf("cmd/%s/main.go", binary)}, args...)
							break
						}
					}
					if !found {
						fmt.Println("Binary not found in cmd directory")
					}
				} else {
					fmt.Println("Multiple binaries found in cmd directory, please specify which one to run")
					fmt.Println("Available binaries:")
					for _, binary := range binaries {
						fmt.Println(binary)
					}
					fmt.Println("Usage: gr --bin <binary_name> [args]")
				}
			}
		} else {
			fmt.Println("Could not find main.go or cmd directory")
		}
	}
}

func run(inp []string, args ...string) {
	if len(inp) == 0 {
		fmt.Println("No command provided")
		return
	}
	inp = append(inp, args...)
	cmd := exec.Command(inp[0], inp[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}
