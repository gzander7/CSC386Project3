// Package main provides a simple interactive shell program in Go.
package main

import (
	"Project2Demo/FileSystem"
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// make func to mkdir using our FileSys.go file system
func mkdir(name string) {
	newDirectoryInode, newInodeNum := FileSystem.Open(FileSystem.CREATE, name, FileSystem.RootFolder)
	directoryBlock, newDirectoryInode := FileSystem.CreateDirectoryFile(FileSystem.ReadSuperBlock().RootDirInode, newInodeNum)
	bytesForDirectoryBlock := FileSystem.EncodeToBytes(directoryBlock)
	FileSystem.Write(&newDirectoryInode, newInodeNum, bytesForDirectoryBlock)
}

func rm(name string) {

}

// make func to mv using our FileSys.go file system

func ls(dirName string) {
	dirInode, _ := FileSystem.Open(FileSystem.READ, dirName, FileSystem.RootFolder)

	if !dirInode.IsDirectory {
		fmt.Println("Error:", dirName, "is not a directory")
		return
	}

	files := FileSystem.Ls(&dirInode)
	fmt.Println("Files in directory:", dirName)
	for _, file := range files {
		fmt.Println(file)
	}
}

// main is the entry point of the program.
func main() {
	FileSystem.InitializeFileSystem()
	fmt.Println("Welcome to GagesGoShell!")
	// Create a reader to read input from standard input
	reader := bufio.NewReader(os.Stdin)

	// Loop indefinitely to continuously accept user input
	for {
		// Print the shell prompt.
		fmt.Print("GagesGoShell> ")

		// Read input from the user until a newline character
		input, err := reader.ReadString('\n')
		if err != nil {
			// If there is an error reading input print the error message and continue to the next iteration
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			continue
		}

		// Remove newline character from input
		input = strings.TrimSpace(input)

		// Split input into command and arguments
		parts := strings.Fields(input)
		if len(parts) == 0 {
			continue // Skip empty input lines
		}

		// Check for commands
		switch parts[0] {
		case "ls":
			if len(parts) < 2 {
				fmt.Println("Usage: ls [directory]")
				continue
			}
			ls(parts[1]) // Call ls function with the directory name
		case "wc":
			// Execute wc command with arguments
			cmd := exec.Command("wc", parts[1:]...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
				fmt.Println("Error:", err)
			}
		case "mkdir":
			// Execute mkdir command with arguments
			mkdir(parts[1])
			fmt.Println("Directory created:", parts[1])
		case "cp":
			// Execute cp command with arguments
			cmd := exec.Command("cp", parts[1:]...)
			err := cmd.Run()
			if err != nil {
				fmt.Println("Error:", err)
			}
		case "mv":
			// Execute mv command with arguments
			cmd := exec.Command("mv", parts[1:]...)
			err := cmd.Run()
			if err != nil {
				fmt.Println("Error:", err)
			}
		case "cd":
			// Change directory to the specified path
			if len(parts) < 2 {
				fmt.Println("Usage: cd [directory]")
				continue
			}
			err := os.Chdir(parts[1])
			if err != nil {
				fmt.Println("Error:", err)
			}
		case "whoami":
			// Print users name and user ID
			//fmt.Println("User:", os.Getenv("USER"))
			fmt.Println("User:", "Gage Ross", "User ID:", os.Getuid()) // Hardcoded user name for demonstration
		case "rm":
			// remove file
			if len(parts) < 2 {
				fmt.Println("Error: missing file name")
			}
			rm(parts[1])
		case "exit":
			// Exit the shell
			os.Exit(0)
		default:
			// Print message for unsupported commands
			fmt.Println("Command not supported:", parts[0])
		}
	}
}
