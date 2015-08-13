package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	var command string

	// Check arg for appname to load
	if len(os.Args) > 1 {
		command = os.Args[1]
	} else {
		fmt.Println("Must pass appname as argument")
		os.Exit(1)
	}

	helpMessage := "usage: binpath (command_name)"

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

func checkForBin(appname string) {
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

				fmt.Println("Found Bin in: ", dir)

				// set env for this bin path
				v := fmt.Sprintf("$PATH:%v/bin", dir)
				fmt.Println(v)

				os.Setenv("PATH", v)

				// Make sure to pass any arguments to the app
				args := os.Args[2:len(os.Args)]

				// Call command
				cmd := exec.Command(appname, args...)
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

	checkForBin(appname)
}
