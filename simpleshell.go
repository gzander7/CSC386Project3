package main

import (
	"Project2Demo/FileSystem"
	"bufio"
	"fmt"
	"os"
	"strings"
)

// mkdir func creates a new directory in the file system and writes the directory block to the file system using the open and write functions from the FileSys.go file system
func mkdir(name string) {
	newDirectoryInode, newInodeNum := FileSystem.Open(FileSystem.CREATE, name, FileSystem.RootFolder)                          // Create a new directory inode
	directoryBlock, newDirectoryInode := FileSystem.CreateDirectoryFile(FileSystem.ReadSuperBlock().RootDirInode, newInodeNum) // Create a new directory block
	bytesForDirectoryBlock := FileSystem.EncodeToBytes(directoryBlock)                                                         // Encode the directory block to bytes
	FileSystem.Write(&newDirectoryInode, newInodeNum, bytesForDirectoryBlock)                                                  // Write the directory block to the file system
}

// logic for mv is using cp logic that is in the FileSys.go file system and using rm to remove the source file after copying it to the destination file
func mv(srcName string, destName string) {
	FileSystem.Cp(srcName, destName) // Copy the source file to the destination file
	rm(srcName)                      // Remove the source file
}

// logic for cp is in the FileSys.go file system
func cp(srcName string, destName string) {
	FileSystem.Cp(srcName, destName)
}

// using the FileSys.go Unlink fucntion to remove a file
func rm(path string) {
	parentinode, childinode, _, childinodenum := getParentandChildInodes(path)
	if childinode.DirectBlock1 == 0 {
		fmt.Println("nothing to remove")
	} else {
		FileSystem.Unlink(childinodenum, parentinode)

	} // end of if statement
}

// using the FileSys.go Read function to read in a file and print it out
func more(path string) {
	_, childinode, _, _ := getParentandChildInodes(path)
	if childinode.DirectBlock1 == 0 {
		fmt.Println("nothing to read in")
	} else {
		filecontent := FileSystem.Read(&childinode)
		fmt.Println(filecontent)
	} // end of if statement
}

// Function to get the parent and child inodes of a path
func getParentandChildInodes(path string) (parentinode FileSystem.INode, childinode FileSystem.INode, parentinodenum int, childinodenum int) {
	stringSlice := strings.Split(path, "/")         // Split the path into a slice of strings
	newDirectory := stringSlice[len(stringSlice)-1] // Get the new directory name
	stringSlice = stringSlice[:len(stringSlice)-1]  // Remove the new directory name from the slice
	var toPath string                               // Initialize the path to the parent directory
	for _, dir := range stringSlice {               // Loop through the slice of strings
		if dir != "" { // Check if the string is not empty
			toPath = toPath + "/" + dir // Add the string to the path
		} // end of if statement
	} // end of for loop
	parentinode, parentinodenum = FileSystem.FindSubdirectories(toPath)
	childinode, childinodenum = FileSystem.Open(FileSystem.CREATE, newDirectory, parentinode)
	return parentinode, childinode, parentinodenum, childinodenum
}

// ExecuteCommand function allows for the use of os commands in the shell such as cat, ls, etc. and allows for redirection into the file system
func ExecuteCommand(command string, args []string) {
	modifiedArgs := []string{}
	var nextOutputFile string
	pathInOutputFile := false

	// Check for redirection and store the next argument if found
	for i, arg := range args {
		if arg == ">>" {
			if i+1 < len(args) {
				nextOutputFile = args[i+1] // keep track of next argument if '>>' is found
				if strings.Contains(nextOutputFile, "/") {
					pathInOutputFile = true // Check if path is included in output file
				} // end of if statement
			} // end of if statement
			break // Stop checking for redirection if found
		} // end of if statement
		modifiedArgs = append(modifiedArgs, arg) // Add argument to modifiedArgs if not a redirection
	} // end of for loop

	inputFileContent, err := os.ReadFile(modifiedArgs[len(modifiedArgs)-1]) // Read in file content
	if err != nil {
		fmt.Println("couldn't read in file or not a valid command") // Print error message if file not found
		fmt.Println(err)
	} // end of if statement

	// allows to specify a path for the output file *e.g cat testinput2.txt >> new/test.txt* will put the file test.txt in the dir new
	if pathInOutputFile {
		stringSlice := strings.Split(nextOutputFile, "/")                                        // Split the output file path
		fileName := stringSlice[len(stringSlice)-1]                                              // Get the file name from the path
		parentinode, _, _, _ := getParentandChildInodes(nextOutputFile)                          // Get the parent inode
		newFileInode, firstInodeNun := FileSystem.Open(FileSystem.CREATE, fileName, parentinode) // Create the new file
		contentToWrite := []byte(inputFileContent)                                               // Convert the file content to bytes
		FileSystem.Write(&newFileInode, firstInodeNun, contentToWrite)                           // Write the content to the new file
		fmt.Println("file ", fileName, " read in")
	} else { // If path is not included in output file then create the file in the current directory of the shell(*e.g cat testinput2.txt >> test.txt* will put the file test.txt in the current directory of the shell in this case probably the root directory of the file system)
		newFileInode, firstInodeNun := FileSystem.Open(FileSystem.CREATE, nextOutputFile, FileSystem.RootFolder)
		contentToWrite := []byte(inputFileContent)
		FileSystem.Write(&newFileInode, firstInodeNun, contentToWrite)
		fmt.Println("file read into", nextOutputFile)
	} // end of if statement

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
		case "mkdir":
			// Execute mkdir command with arguments
			if len(parts) < 2 {
				fmt.Println("Usage: mkdir [directory]")
				continue
			}
			mkdir(parts[1])
			fmt.Println("Directory created:", parts[1])
		case "mv":
			if len(parts) < 3 {
				fmt.Println("Usage: mv [source] [destination]")
				continue
			}
			mv(parts[1], parts[2])
		case "cp":
			if len(parts) < 3 {
				fmt.Println("Usage: cp [source] [destination]")
				continue
			}
			// Execute cp command with arguments
			cp(parts[1], parts[2])
		case "rm":
			// remove file
			if len(parts) < 2 {
				fmt.Println("Usage rm [file]")
			}
			rm(parts[1])
		case "more":
			// Execute more command with arguments
			if len(parts) < 2 {
				fmt.Println("Usage: more [file]")
				continue
			}
			more(parts[1])
		case "whoami":
			// Print users name and user ID
			//fmt.Println("User:", os.Getenv("USER"))
			fmt.Println("User:", "Gage Ross", "User ID:", os.Getuid()) // Hardcoded user name for demonstration
		case "exit":
			// Exit the shell
			os.Exit(0)
		default:
			// Print message for unsupported commands
			ExecuteCommand(parts[0], parts[1:])
		} // end of switch statement
	} // end of for loop
} // end of main function
