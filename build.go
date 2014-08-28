package main

import (
	"log"
	"os"
	"io"
	"os/exec"
	"path/filepath"
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

	//1. set GOPATH env variable
	//take the directory of build.go as the GOPATH value
	absPath, err := filepath.Abs(".")
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("set GOPATH as: %s\n", absPath)

	err = os.Setenv("GOPATH", absPath)
	if err != nil {
		log.Println(err)
		return
	}

	//2. copy help.txt to be under bin folder so that the program can access the help file
	err = CopyFile(filepath.Join(absPath, "/misc/help.txt"), 
		filepath.Join(absPath, "/bin/help.txt"))
	if err != nil {
		log.Println(err)
		return
	}

	//3. install main package
	cmd := exec.Command("go", "install", "-v", "main")
	err = cmd.Run()
	if err != nil {
		log.Println(err)
		return
	}
}
