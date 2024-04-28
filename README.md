# OS Programming Project 3
# Gage Ross Simple Go Shell Using Simple Go File System

This is a simple shell implemented in Go. It provides basic file system operations and allows for the use of OS commands.

## Features

- `mkdir [directory]`: Creates a new directory in the file system.
- `mv [source] [destination]`: Moves a file from the source to the destination. This is implemented using the copy and remove logic.
- `cp [source] [destination]`: Copies a file from the source to the destination.
- `rm [file]`: Removes a file from the file system.
- `more [file]`: Reads in a file and prints it out.
- `whoami`: Prints the user's name and user ID.
- `exit`: Exits the shell.
- Redirection: The shell supports redirection using the '>>' operator. For example, `cat testinput2.txt >> new/test.txt` will put the file `test.txt` in the directory `new`.
## Issues
When using a file that is as big as as testinput.txt it seems to crash when using mv or cp commands are used.
