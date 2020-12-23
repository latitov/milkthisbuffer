package main

import (
	"github.com/latitov/milkthisbuffer"
	"os/exec"
	"log"
	"time"
)

// Ctrl-D to end the input (Ctrl-Z on Windows)

func main() {
	pipe1 := milkthisbuffer.New(500)
	pipe2 := milkthisbuffer.New(500)

	go pipe2.StdoutAsync()
	
	var err error

	go func() {
		cmd := exec.Command("tr", "a-z", "A-Z")
		cmd.Stdin = pipe1
		cmd.Stdout = pipe2
		cmd.Stderr = pipe2
		cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}()
	
	cmd := exec.Command("echo", "hey-hi-hello, 1-2-3, repeat")
	cmd.Stdout = pipe1
	cmd.Stderr = pipe1
	cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	
	pipe1.Close()	// signal the EOF
	
	time.Sleep(1 * time.Second)
}
