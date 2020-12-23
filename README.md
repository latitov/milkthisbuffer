# milkthisbuffer
Mutex/Lock-Free Thread-Safe Buffer in Go

## For Shell automation in Go:

- #### See the output of long-running Shell command in real time;
- #### Process that output programmatically by your code, in real time;
- #### Pipe that output to another Shell command;
- #### Execute shell commands really easily, like in the Bash;

## Mutex/Lock-Free Thread-Safe Buffer

The same as bytes.Buffer, but thread-safe.

Wraps (implements) the io.Reader and io.Writer with a channel.

Can also be used as a "pipe" between goroutines.

Etymology:
1. Mutex/Lock-Free Thread-Safe Buffer
1. MxLkF_ThS_Buffer
1. MilkThisBuffer

## What's the use of it

It's an absolutely useful tool for a shell automation. Basically, the standard lib, specifically the os/exec, allows for quite convenient shell commands execution out-of-the-box, nothing else is needed.... Except one thing. That thing is the ability to see the output of long-running command in real time, and/or process that output programmatically, and/or pipe that output to another command. All these cases are covered by this simple pipe-buffer.

You can:

- See the output of long-running Shell command in real time;
- Process that output programmatically by your code, in real time;
- Pipe that output to another Shell command;

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
		cmd.Stderr = mtb1
	
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
			cmd.Stderr = pipe2
			cmd.Run()
		}()
		
		cmd := exec.Command("echo", "hey-hi-hello, 1-2-3, repeat")
		cmd.Stdout = pipe1
		cmd.Stderr = pipe1
		cmd.Run()
		pipe1.Close()	// signal the EOF
		
		time.Sleep(1 * time.Second)
	}

The output will be:

	HEY-HI-HELLO, 1-2-3, REPEAT

Basically, you are free to build any graph out of pipes (this buffer as a pipe), and any commands. This opens the road to use the Go language for a Shell automation (which is faster than Bash or Python or Perl).

## And one more thing ))

There is this standard thing: https://golang.org/pkg/io/#Pipe

The example above can be written in 100% standard Go, without any third-party libs, this way:

	package main
	
	import (
		"fmt"
		"io"
		"os/exec"
		"time"
	)
	
	// Ctrl-D to end the input (Ctrl-Z on Windows)
	
	func main() {
		pipe1_reader, pipe1_writer := io.Pipe()
		pipe2_reader, pipe2_writer := io.Pipe()
	
		go StdoutAsync(pipe2_reader)
	
		go func() {
			cmd := exec.Command("tr", "a-z", "A-Z")
			cmd.Stdin = pipe1_reader
			cmd.Stdout = pipe2_writer
			cmd.Stderr = pipe2_writer
			cmd.Run()
		}()
		
		cmd := exec.Command("echo", "hey-hi-hello, 1-2-3, repeat")
		cmd.Stdout = pipe1_writer
		cmd.Stderr = pipe1_writer
		cmd.Run()
		pipe1_writer.Close() // signal the EOF
		
		time.Sleep(1 * time.Second)
	}
	
	func StdoutAsync(b io.Reader) {
		buf := make([]byte, 100)
		for {
			n, err := b.Read(buf)
			if err != nil {
				fmt.Printf("StdoutAsync: %v\n", err)
				break
			}
			fmt.Printf("%v", string(buf[:n]))
		}
	}
	
And the output will be, identically:

	HEY-HI-HELLO, 1-2-3, REPEAT

It's up to you how to proceed. I myself am in doubt, if I wasted my time. Did I?

## THE BONUS

Actually, it's more than that. I also added a wrapper for the code above, so now you can write the example this short:

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
		
		pipe1.Close()	// signal the EOF
		
		time.Sleep(1 * time.Second)
	}
One more example, without piping, illustrating the ease of calling many commands in sequence:

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

To do, maybe also:

- A wrapper with Context;
- A wrapper that return the error;

Very useful.
