package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"sync"
)

// main is the entry point of the application.
func main() {
	// Disable the default timestamp and source file prefixes from the log package.
	log.SetFlags(0)

	// The commands to run are passed as command-line arguments.
	// os.Args[0] is the program name itself, so we skip it.
	commands := os.Args[1:]
	if len(commands) == 0 {
		log.Fatal("Error: No commands to run. Usage: go run main.go \"command1\" \"command2\"")
	}

	// A WaitGroup waits for a collection of goroutines to finish.
	var wg sync.WaitGroup

	// Run each command in its own goroutine.
	for _, commandStr := range commands {
		// Increment the WaitGroup counter.
		wg.Add(1)
		// Launch a new goroutine to run the command.
		go runCommand(commandStr, &wg)
	}

	// Wait for all goroutines to finish.
	wg.Wait()
}

// runCommand executes a single command string, captures its output,
// and prefixes each line with the command string itself.
func runCommand(commandStr string, wg *sync.WaitGroup) {
	// Decrement the counter when the goroutine completes.
	defer wg.Done()

	var shell, flag string

	// Determine the correct shell and flag based on the operating system.
	if runtime.GOOS == "windows" {
		shell, flag = "cmd", "/C"
	} else {
		shell, flag = "sh", "-c"
	}

	// Create the command with the appropriate shell.
	cmd := exec.Command(shell, flag, commandStr)

	// Create pipes for stdout and stderr to capture the output.
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("Error creating stdout pipe for command '%s': %v", commandStr, err)
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Printf("Error creating stderr pipe for command '%s': %v", commandStr, err)
		return
	}

	// Start the command. This does not wait for it to complete.
	if err := cmd.Start(); err != nil {
		log.Printf("Error starting command '%s': %v", commandStr, err)
		return
	}

	// Create a prefix for logging output from this command.
	prefix := fmt.Sprintf("[%s] ", commandStr)

	// Concurrently scan both stdout and stderr and pipe them to our logger.
	go pipeOutput(stdout, prefix)
	go pipeOutput(stderr, prefix)

	// Wait for the command to exit and release its resources.
	if err := cmd.Wait(); err != nil {
		log.Printf("Command '%s' finished with error: %v", commandStr, err)
	}
}

// pipeOutput reads from an io.Reader line by line, adds a prefix,
// and prints it to the standard logger.
func pipeOutput(r io.Reader, prefix string) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		log.Printf("%s%s", prefix, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Printf("Error reading output: %v", err)
	}
}
