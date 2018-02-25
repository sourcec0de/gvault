package utils

import (
	"io/ioutil"
	"os"

	"github.com/chzyer/readline"
)

var sdinData []byte

// CWD returns a string pointing to the current working directory
// where the process was started
func CWD() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return cwd
}

// Ask a question as a promt
func Ask(question string, rl *readline.Instance) (string, error) {
	rl.SetPrompt(question)
	return rl.Readline()
}

// ReadAllStdin reads all data from stdin
func ReadAllStdin() []byte {
	bytes, _ := ioutil.ReadAll(os.Stdin)
	return bytes
}
