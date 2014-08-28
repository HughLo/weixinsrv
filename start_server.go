package main

import (
	"os/exec"
	"path/filepath"
	"log"
	"os"
)

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

	//2. run server program
	cmd := exec.Command("sudo", "-E", "./bin/main")
	err = cmd.Run()
	if err != nil {
		log.Println(err)
		return
	}
}
