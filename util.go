package main

import (
	"bufio"
	"io"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func MustReadLine(r *bufio.Reader) []byte {
	line, isPrefix, err := r.ReadLine()
	if err == io.EOF {
		return nil
	}
	check(err)
	if isPrefix {
		panic("overlong line")
	}
	return line
}
