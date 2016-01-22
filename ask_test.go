package ask

import (
	"bytes"
	"errors"
	"io/ioutil"
	"testing"
)

type BrokenWriter struct{}

func (w *BrokenWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("broken")
}

func TestAsk(t *testing.T) {
	question := "question"
	expectedAnswer := "answer"

	out = new(bytes.Buffer)
	in = bytes.NewBufferString("answer\n")

	answer, err := Ask(question)
	if err != nil {
		t.Fatalf("Ask unexpected error %v", err)
	}

	if answer != expectedAnswer {
		t.Fatal("Answer output error, Got: %s Expected: %s", answer, expectedAnswer)
	}
}

func TestFask(t *testing.T) {
	question := "question"
	expectedQuestion := "question: "

	answer := "answer\n"
	expectedAnswer := "answer"

	w := new(bytes.Buffer)
	r := bytes.NewBufferString(answer)

	a, err := Fask(question, w, r)
	if err != nil {
		t.Fatalf("Fask unexpected error %v", err)
	}

	if w.String() != expectedQuestion {
		t.Fatalf("Question output error, Got: %s Expected: %s", question, expectedQuestion)
	}

	if a != expectedAnswer {
		t.Fatal("Answer output error, Got: %s Expected: %s", answer, expectedAnswer)
	}

	// Error while reading answer
	a, err = Fask(question, ioutil.Discard, new(bytes.Buffer))
	if a != "" {
		t.Fatalf("Answer expected to be empty, Got: %s")
	}
	if err == nil {
		t.Fatalf("Error expected, Got: nil")
	}

	// Error while writing question
	a, err = Fask(question, &BrokenWriter{}, r)
	if a != "" {
		t.Fatalf("Answer expected to be empty, Got: %s", a)
	}
	if err == nil {
		t.Fatalf("Error expected, Got: nil")
	}
}
