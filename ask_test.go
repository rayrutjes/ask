package ask

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type WriterMock struct {
	mock.Mock
}

func (m *WriterMock) Write(p []byte) (n int, err error) {
	args := m.Called(p)
	return args.Int(0), args.Error(1)
}

type ReaderMock struct {
	mock.Mock
	Answer []byte
}

func (m *ReaderMock) Read(p []byte) (n int, err error) {
	args := m.Called(p)
	err = args.Error(1)
	if err == nil {
		n = copy(p, m.Answer)
		return n, nil
	}
	return args.Int(0), err
}

type ValidateHelperFuncs struct {
	FailsCountBeforeValidate int
	FailsCount               int
	Message                  string
}

func (v *ValidateHelperFuncs) counterReached(answer string) (bool, string) {
	if v.FailsCount == v.FailsCountBeforeValidate {
		return true, ""
	}
	v.FailsCount++
	return false, v.Message
}

func TestFask(t *testing.T) {
	question := "question"
	expectedQuestion := "question: "

	answer := "answer\n"
	expectedAnswer := "answer"

	w := new(bytes.Buffer)
	r := bytes.NewBufferString(answer)

	a, err := Fask(question, w, r)
	assert.Nil(t, err)
	assert.Equal(t, expectedQuestion, w.String())
	assert.Equal(t, expectedAnswer, a)

	// Error while reading answer
	a, err = Fask(question, ioutil.Discard, new(bytes.Buffer))
	assert.Empty(t, a)
	assert.Error(t, err)

	// Error while writing question
	brokenWriter := new(WriterMock)
	brokenWriter.On("Write", []byte(expectedQuestion)).Return(0, io.EOF)
	a, err = Fask(question, brokenWriter, r)
	assert.Empty(t, a)
	assert.Error(t, err)
}

func TestAsk(t *testing.T) {
	question := "question"
	expectedAnswer := "answer"

	out = new(bytes.Buffer)
	in = bytes.NewBufferString("answer\n")

	answer, err := Ask(question)
	assert.Nil(t, err)
	assert.Equal(t, expectedAnswer, answer)
}

func TestFaskWhile(t *testing.T) {
	question := []byte("question: ")
	validationMsg := []byte("failure\n")

	writerMock := new(WriterMock)
	writerMock.On("Write", question).Return(len(question), nil).Once()
	writerMock.On("Write", validationMsg).Return(len(validationMsg), nil).Once()
	writerMock.On("Write", question).Return(len(question), nil).Once()

	readerMock := new(ReaderMock)
	readerMock.Answer = []byte("answer\n")

	buffer := make([]byte, 4096)
	readerMock.On("Read", buffer).Return(len(readerMock.Answer), nil).Twice()
	readerMock.On("Read", buffer).Return(0, io.EOF).Once()

	helper := &ValidateHelperFuncs{
		FailsCount:               0,
		FailsCountBeforeValidate: 1,
		Message:                  "failure",
	}

	answer, err := FaskWhile("question", helper.counterReached, writerMock, readerMock)
	assert.Nil(t, err)

	assert.Equal(t, "answer", answer)

	writerMock.AssertExpectations(t)

	// if write fails
	helper = &ValidateHelperFuncs{
		FailsCount:               0,
		FailsCountBeforeValidate: 1,
		Message:                  "failure",
	}
	brokenWriter := new(WriterMock)
	brokenWriter.On("Write", question).Return(0, io.EOF)

	answer, err = FaskWhile("question", helper.counterReached, brokenWriter, readerMock)
	assert.Error(t, err)
	assert.Empty(t, answer)
}

func TestAskWhile(t *testing.T) {
	question := []byte("question: ")
	validationMsg := []byte("failure\n")

	writerMock := new(WriterMock)
	writerMock.On("Write", question).Return(len(question), nil).Once()
	writerMock.On("Write", validationMsg).Return(len(validationMsg), nil).Once()
	writerMock.On("Write", question).Return(len(question), nil).Once()

	out = writerMock

	readerMock := new(ReaderMock)
	readerMock.Answer = []byte("answer\n")

	buffer := make([]byte, 4096)
	readerMock.On("Read", buffer).Return(len(readerMock.Answer), nil).Twice()
	readerMock.On("Read", buffer).Return(0, io.EOF).Once()

	in = readerMock

	helper := &ValidateHelperFuncs{
		FailsCount:               0,
		FailsCountBeforeValidate: 1,
		Message:                  "failure",
	}

	answer, err := AskWhile("question", helper.counterReached)
	assert.Nil(t, err)
	assert.Equal(t, "answer", answer)
}
