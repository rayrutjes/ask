package ask

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

var (
	out io.Writer = os.Stdout
	in  io.Reader = os.Stdin
)

// Ask writes a question to the standard output and reads the answer from the
// standard input.
func Ask(q string) (string, error) {
	return Fask(q, out, in)
}

// Fask writes the question into to given writer and returns the answer from the reader.
func Fask(q string, w io.Writer, r io.Reader) (string, error) {
	q += ": "
	if _, err := w.Write([]byte(q)); err != nil {
		return "", fmt.Errorf("Unable to write the question %v", err)
	}
	a, err := bufio.NewReader(r).ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("Unable to read the answer %v", err)
	}

	return a[:len(a)-1], err
}

// ValidateFunc is a function that is used to validate that an answer matches
// the expectations. It should return a boolean and an additional validation failure message
// when the validation fails.
type ValidateFunc func(string) (bool, string)

// AskWhile writes a question and returns the answer once the validation function
// returns true. The questions is written again on each failed validation pass.
func AskWhile(q string, v ValidateFunc) (string, error) {
	return FaskWhile(q, v, out, in)
}

// FaskWhile writes a question and returns the answer once the validation function
// returns true. The questions is written again on each failed validation pass.
func FaskWhile(q string, v ValidateFunc, w io.Writer, r io.Reader) (string, error) {
	for {
		a, err := Fask(q, w, r)
		if err != nil {
			return "", err
		}
		ok, msg := v(a)
		if ok {
			return a, nil
		}
		w.Write([]byte(msg + "\n"))
	}
}
