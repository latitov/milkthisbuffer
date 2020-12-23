package main

import (
	"github.com/latitov/milkthisbuffer"
	"time"
)

// Ctrl-D to end the input (Ctrl-Z on Windows)

func main() {
	pipe1 := milkthisbuffer.New(500)
	
	go pipe1.StdoutAsync()
	
	co1 := &milkthisbuffer.CommandObject{
		Stdout: pipe1,
		Stderr: pipe1,
	}
	
	co1.Execf("echo", "hey-hi-hello, 1-2-3, repeat")
	co1.Execf("echo", "hey-hi-hello, 1-2-3-5, repeat")
	co1.Execf("echo", "hey-hi-hello, 1-2-3-5-7, repeat")
	co1.Execf("echo", "hey-hi-hello, 1-2-3-5-7-9, repeat")
	
	co1.Execf("git", "add", "-A")
	co1.Execf("git", "commit", "-m", "'.'")
	
	pipe1.Close()	// signal the EOF
	
	time.Sleep(1 * time.Second)
}
