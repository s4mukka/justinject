package utils

import "os"

var (
	osReadFile = os.ReadFile
)

func (u *Utils) ReadFile(name string) ([]byte, error) {
	return osReadFile(name)
}
