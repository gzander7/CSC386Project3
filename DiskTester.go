package main

import (
	"Project2Demo/FileSystem"
	"fmt"
	"log"
	"os"
)

func main() {
	// this is I think the test I promised you - except that maybe the string is too short
	FileSystem.InitializeFileSystem()
	// d
	newFileInode, _ := FileSystem.Open(FileSystem.CREATE, "Text.txt", FileSystem.RootFolder)
	stringContents, err := os.ReadFile("testInput.txt")
	if err != nil {
		log.Fatal("Oh yikes why couldn't we open the file!?!?!")
	}
	contentToWrite := []byte(stringContents)

	FileSystem.Write(&newFileInode, contentToWrite)
	testFileInode, testFileInodeNum := FileSystem.Open(FileSystem.READ, "Text.txt", FileSystem.RootFolder)
	fmt.Println(testFileInodeNum, FileSystem.Read(testFileInode))
	newDirectoryInode, newInodeNum := FileSystem.Open(FileSystem.CREATE, "NewDir",
		FileSystem.RootFolder)
	directoryBlock := FileSystem.CreateDirectoryFile(FileSystem.ReadSuperBlock().RootDirInode, newInodeNum)
	bytesForDirectoryBlock := FileSystem.EncodeToBytes(directoryBlock)
	fmt.Println("WARNING Encoded Directory block is %d bytes", len(bytesForDirectoryBlock))
	FileSystem.Write(&newDirectoryInode, bytesForDirectoryBlock)
	file2Inode, _ := FileSystem.Open(FileSystem.CREATE, "FileInSubdir", newDirectoryInode)
	dataToWrite := []byte("Help I'm stuck in a virtual file System")
	FileSystem.Write(&file2Inode, dataToWrite)
	fmt.Println(FileSystem.Read(file2Inode))
}
