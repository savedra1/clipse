package main

import (
	"fmt"
	"os"
	"strconv"
)

const TerminalPPIDEnvVar = "TERMINAL_PPID"

var TerminalPPID int

func init() {
	// Retrieve the PPID passed as an argument
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run test.go <terminal_ppid>")
		os.Exit(1)
	}
	fmt.Println("ppid:", os.Args[1])
	ppidStr := os.Args[1]

	// Convert the PPID string to an integer
	ppid, err := strconv.Atoi(ppidStr)
	if err != nil {
		fmt.Println("Invalid PPID:", err)
		os.Exit(1)
	}

	// Set the global variable
	TerminalPPID = ppid

	// Optionally, you can also export it as an environment variable
	os.Setenv(TerminalPPIDEnvVar, ppidStr)
}

func main() {
	// Your main function logic here
	// You can access the TerminalPPID global variable throughout your program
	fmt.Println("Terminal PPID:", TerminalPPID)
}
