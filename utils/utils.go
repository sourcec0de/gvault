package utils

import (
	"os"
)

// CWD returns a string pointing to the current working directory
// where the process was started
func CWD() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return cwd
}
