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
