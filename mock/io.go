package mock

import (
	"bytes"
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
	_, err := io.Copy(&buf, r)
	if err != nil {
		panic(err)
	}
	os.Stderr = old
	return buf.String()
}

func CaptureOutput(fn func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	w.Close()
	var buf bytes.Buffer
	_, err := io.Copy(&buf, r)
	if err != nil {
		panic(err)
	}
	os.Stdout = old
	return buf.String()
}
