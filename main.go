package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const space = "  " // Used for pretty printing
const br = "\n\n"

func main() {
	var command string

	usage := "binpath command [arguments]"
	options := "-list, -ls"
	optionsDesc := "list directory contents of bin path"

	// Create help message with usage messaging
	helpMessage := fmt.Sprintf(bold("USAGE:")+"\n%s%v", space, usage)
	// Break between messages
	helpMessage += br
	// Add options messaging
	helpMessage += fmt.Sprintf(bold("OPTIONS:")+"\n%v\n%s%v", options, space, optionsDesc)

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
	case "--show-bash-completion":
		checkForBin(command, true)
	default:
		checkForBin(command, false)
	}

}

func checkForBin(command string, showBashCompletion bool) {
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

				// If --show-bash-completion flag was used then list out files in binPath
				if showBashCompletion == true {
					// List files in bin folder
					listBinPathFiles(binPath)
				} else {
					// Execute command in binPath
					executeCommand(binPath, command)
				}

				os.Exit(0)
			}
		}
	}

	// If it doesn't move up a directory and test
	os.Chdir(".." + string(filepath.Separator))

	checkForBin(command, showBashCompletion)
}

// Print out files found in binpath
func listBinPathFiles(binPath string) {
	files, err := ioutil.ReadDir(binPath)
	if err != nil {
		fmt.Printf("-binpath: error reading: %v \n", binPath)
		os.Exit(1)
	}

	for _, file := range files {
		name := file.Name()

		// If it's a hidden file don't print it
		if !strings.HasPrefix(name, ".") {
			fmt.Println(name)
		}
	}
}

// Try to execute the command if a binpath is found
func executeCommand(binPath string, command string) {
	// Set path to command
	commandPath := filepath.Join(binPath, command)

	// Make to see if command exist in bin path
	_, err := os.Stat(commandPath)
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
}

func bold(s string) string {
	return "\033[1m" + s + "\033[0m"
}
