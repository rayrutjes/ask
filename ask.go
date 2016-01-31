package ask

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
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

// While writes a question and returns the answer once the validation function
// returns true. The questions is written again on each failed validation pass.
func While(q string, v ValidateFunc) (string, error) {
	return Fwhile(q, v, out, in)
}

// Fwhile writes a question and returns the answer once the validation function
// returns true. The questions is written again on each failed validation pass.
func Fwhile(q string, v ValidateFunc, w io.Writer, r io.Reader) (string, error) {
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

// confirmValidateFunc function that validates that the answer is y or n.
var confirmValidateFunc = func(answer string) (bool, string) {
	if !confirmed(answer) && !infirmed(answer) {
		return false, "please confirm or infirm by typing 'y' or 'n'"
	}
	return true, ""
}

// confirmed returns true if the answer is Y or y.
func confirmed(answer string) bool {
	a := strings.ToLower(answer)
	return strings.EqualFold(a, "y")
}

// infirmed returns true if the answer is N or n.
func infirmed(answer string) bool {
	a := strings.ToLower(answer)
	return strings.EqualFold(a, "n")
}

// Fconfirm prompts the user with a yes/no question. Keeps prompting the question until
// the users answers with "y" or "n". Returns a boolean value indicating if the user confirmed
// or infirmed.
func Fconfirm(question string, w io.Writer, r io.Reader) (bool, error) {
	answer, err := Fwhile(question+" (y/n)", confirmValidateFunc, w, r)
	if err != nil {
		return false, err
	}

	return confirmed(answer), nil
}

// Confirm prompts the user with a yes/no question. Keeps prompting the question until
// the users answers with "y" or "n". Returns a boolean value indicating if the user confirmed
// or infirmed. Question and answers are written and read back from standard output/input.
func Confirm(question string) (bool, error) {
	return Fconfirm(question, out, in)
}
