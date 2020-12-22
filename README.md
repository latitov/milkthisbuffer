# milkthisbuffer
Mutex/Lock-Free Thread-Safe Buffer in Go

## Mutex/Lock-Free Thread-Safe Buffer

The same as bytes.Buffer, but thread-safe.

Wraps (implements) the io.Reader and io.Writer with a channel.

Can also be used as a "pipe" between goroutines.

Etymology:
- Mutex/Lock-Free Thread-Safe Buffer
- MxLkF_ThS_Buffer
- MilkThisBuffer

## What's the use of it

It's an absolutely useful tool for a shell automation. Basically, the standard lib, specifically the os/exec, allows for quite convenient shell commands execution, nothing else is needed.... Except one thing. That thing is the ability to see the output of long-running command in real time, and/or process that output programmatically, and/or pipe that output to another command. All these cases are covered by this simple pipe-buffer.

## Example usage

	package main
	
	import (
		"github.com/latitov/milkthisbuffer"
		"log"
		"os"
		"os/exec"
	)
	
	// Ctrl-D to end the input (Ctrl-Z on Windows)
	
	func main() {
		mtb1 := milkthisbuffer.New(500)
		go mtb1.StdoutAsync()
	
		cmd := exec.Command("tr", "a-z", "A-Z")
		cmd.Stdin = os.Stdin
		cmd.Stdout = mtb1
	
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}

## One more example, piping two commands in a chain:

	package main
	
	import (
		"github.com/latitov/milkthisbuffer"
		"os/exec"
		"time"
	)
	
	// Ctrl-D to end the input (Ctrl-Z on Windows)
	
	func main() {
		pipe1 := milkthisbuffer.New(500)
		pipe2 := milkthisbuffer.New(500)
	
		go pipe2.StdoutAsync()
	
		go func() {
			cmd := exec.Command("tr", "a-z", "A-Z")
			cmd.Stdin = pipe1
			cmd.Stdout = pipe2
			cmd.Run()
		}()
		
		cmd := exec.Command("echo", "hey-hi-hello, 1-2-3, repeat")
		cmd.Stdout = pipe1
		cmd.Run()
		pipe1.Close()	// signal the EOF
		
		time.Sleep(1 * time.Second)
	}

The output will be:

	HEY-HI-HELLO, 1-2-3, REPEAT

Basically, you are free to build any graph out of pipes (this buffer as a pipe), and any commands. This opens the road to use the Go language for a Shell automation (which is faster than Bash or Python or Perl).
