package main

import (
	"github.com/latitov/milkthisbuffer"
	"time"
)

// Ctrl-D to end the input (Ctrl-Z on Windows)

func main() {
	pipe1 := milkthisbuffer.New(500)
	pipe2 := milkthisbuffer.New(500)
	
	go pipe2.StdoutAsync()
	
	co1 := &milkthisbuffer.CommandObject{
		Stdout: pipe1,
		Stderr: pipe1,
	}
	co2 := &milkthisbuffer.CommandObject{
		Stdin: pipe1,
		Stdout: pipe2,
		Stderr: pipe2,
	}
	
	go func() {
		co2.Execf("tr", "a-z", "A-Z")
	}()
	
	co1.Execf("echo", "hey-hi-hello, 1-2-3, repeat")
	co1.Execf("echo", "hey-hi-hello, 1-2-3-5, repeat")
	co1.Execf("echo", "hey-hi-hello, 1-2-3-5-7, repeat")
	co1.Execf("echo", "hey-hi-hello, 1-2-3-5-7-9, repeat")
	
	pipe1.Close()	// signal the EOF
	
	time.Sleep(1 * time.Second)
}
