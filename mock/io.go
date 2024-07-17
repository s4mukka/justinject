package mock

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/s4mukka/justinject/domain"
	"github.com/sirupsen/logrus"
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

func CaptureLoggerOutput(logger domain.ILogger, fn func()) string {
	old := logger.(*logrus.Entry).Logger.Out
	buf := bytes.Buffer{}
	logger.(*logrus.Entry).Logger.SetOutput(&buf)
	fn()
	logger.(*logrus.Entry).Logger.SetOutput(old)
	return buf.String()
}
