package mock

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

func CaptureError(fn func()) string {
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	fn()

	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	os.Stderr = old
	fmt.Printf("ERR %v\n", buf)
	return buf.String()
}

func CaptureOutput(fn func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	os.Stdout = old
	fmt.Printf("OUT %v\n", buf)
	return buf.String()
}
