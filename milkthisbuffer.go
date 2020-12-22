// Copyright (c) 2020 Leonid Titov. All rights reserved.
// MIT licence.
// Version 2020-12-23

package milkthisbuffer

// Mutex/Lock-Free Thread-Safe Buffer

// The same as bytes.Buffer, but thread-safe.

// Wraps (implements) the io.Reader and io.Writer with a channel.

// Can be used as a "pipe" between goroutines.

// Etymology:
//	Mutex/Lock-Free Thread-Safe Buffer
//	MxLkF_ThS_Buffer
//	MilkThisBuffer

import (
	"io"
	"fmt"
)

type MilkThisBuffer struct {
	ch	chan byte
}

func New(len int) (b *MilkThisBuffer) {
	b = &MilkThisBuffer{
		ch:	make(chan byte, len),
	}
	return
}

// blocking behaviour
func (b *MilkThisBuffer) Write(p []byte) (n int, err error) {

	// here's why:
	// https://stackoverflow.com/a/34899098/11729048
	defer func() {
		if err2 := recover(); err2 != nil {
			err = io.ErrClosedPipe
		}
	}()
	
	for n = 0; n < len(p); n++ {
		b.ch <- p[n]
	}
	return
}

// blocking behaviour until at least one byte read, then behaves non-blockingly
func (b *MilkThisBuffer) Read(p []byte) (n int, err error) {
L:
	for n = 0; n < len(p); {
		if n == 0 {
			b1, ok := <-b.ch
			if !ok {
				err = io.EOF
				break L
			}
			p[n] = b1
			n++
		} else {
			select {
			case b1, ok := <-b.ch:
				if !ok {
					err = io.EOF
					break L
				}
				p[n] = b1
				n++
			default:
				break L
			}
		}
	}
	return
}

func (b *MilkThisBuffer) Close() error {
	close(b.ch)
	return nil
}

func (b *MilkThisBuffer) StdoutAsync() {
	buf := make([]byte, 100)
	for {
		n, err := b.Read(buf)
		if err != nil {
			fmt.Printf("milkthisbuffer.StdoutAsync: %v\n", err)
			break
		}
		fmt.Printf("%v", string(buf[:n]))
	}
}
