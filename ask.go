package ask

import (
	"bufio"
	"io"
	"os"
)

var (
	out io.Writer = os.Stdout
	in  io.Reader = os.Stdin
)

// Ask writes a question into the standard output and reads the answer from the
// standard input.
func Ask(q string) (string, error) {
	return Fask(q, out, in)
}

// Fask writes the question into to given writer and returns the answer from the reader.
func Fask(q string, w io.Writer, r io.Reader) (string, error) {
	q += ": "
	if _, err := w.Write([]byte(q)); err != nil {
		return "", err
	}
	a, err := bufio.NewReader(r).ReadString('\n')
	if err != nil {
		return "", err
	}

	return a[:len(a)-1], err
}
