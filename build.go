package main

import (
	"log"
	"os"
	"io"
	"os/exec"
)

func CopyFile(srcName, dstName string) error {
	srcFile, err := os.Open(srcName)
	if err != nil {
		return err
	}

	defer srcFile.Close()

	dstFile, err := os.Create(dstName)
	if err != nil {
		return err
	}

	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)

	return err
}

func main() {
	//copy help.txt to be under bin folder so that the program can access the help file
	err := CopyFile("./misc/help.txt", "./bin/help.txt")
	if err != nil {
		log.Println(err)
	}

	cmd := exec.Command("go", "install", "-v", "main")
	err = cmd.Run()
	if err != nil {
		log.Println(err)
	}
}
