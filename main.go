package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

const space = "  "

func main() {
	var command string

	helpMessage := "binpath command [arguments]"
	helpMessage = fmt.Sprintf("Example uasge: \n%s%v", space, helpMessage)

	// Check arg for appname to load
	if len(os.Args) > 1 {
		command = os.Args[1]
	} else {
		fmt.Println(helpMessage)
		os.Exit(1)
	}

	// Look for help flag before executing command
	switch command {
	case "--help":
		fallthrough
	case "-help":
		fallthrough
	case "--h":
		fallthrough
	case "-h":
		fmt.Println(helpMessage)
		os.Exit(1)
	default:
		checkForBin(command)
	}

}

func checkForBin(command string) {
	dirname := "." + string(filepath.Separator)

	d, err := os.Open(dirname)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer d.Close()

	files, err := d.Readdir(-1)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}

	for _, file := range files {
		if file.Mode().IsDir() {
			if file.Name() == "bin" {
				dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
				if err != nil {
					log.Fatal(err)
				}

				// Set bin path
				binPath := filepath.Join(dir, file.Name())

				// Set path to command
				commandPath := filepath.Join(binPath, command)

				// Make to see if command exist in bin path
				_, err = os.Stat(commandPath)
				if err != nil {
					fmt.Printf("-binpath: %v: command not found in: %v \n", command, binPath)
					os.Exit(1)
				}

				// Create separator for easier reading
				separator := "-----------------------"
				fmt.Println(separator)

				// Notify if a bin is found
				fmt.Printf("Exec command: %v\n", commandPath)

				// Print separator for separation of binpath and command
				fmt.Println(separator)

				// And binPath to $PATH env
				v := fmt.Sprintf("$PATH:%v", binPath)
				os.Setenv("PATH", v)

				// Make sure to pass any arguments to the app
				args := os.Args[2:len(os.Args)]

				// Call command
				cmd := exec.Command(command, args...)
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				err = cmd.Run()
				if _, ok := err.(*exec.ExitError); ok {
					log.Fatal(err)
				}

				os.Exit(0)
			}
		}
	}

	// if it doesn't move up a directory and test
	os.Chdir(".." + string(filepath.Separator))

	checkForBin(command)
}
