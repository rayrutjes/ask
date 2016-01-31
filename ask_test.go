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
}

func (m *ReaderMock) Read(p []byte) (n int, err error) {
	args := m.Called(p)
	return args.Int(0), args.Error(1)
}

func runReadFunc(p []byte) func(mock.Arguments) {
	return func(args mock.Arguments) {
		copy(args.Get(0).([]uint8), p)
	}
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

func TestFaskWhile(t *testing.T) {
	question := []byte("question: ")
	validationMsg := []byte("failure\n")

	writerMock := new(WriterMock)
	writerMock.On("Write", question).Return(len(question), nil).Once()
	writerMock.On("Write", validationMsg).Return(len(validationMsg), nil).Once()
	writerMock.On("Write", question).Return(len(question), nil).Once()

	readAnswer := []byte("answer\n")
	readerMock := new(ReaderMock)
	readerMock.On("Read", mock.AnythingOfType("[]uint8")).Return(len(readAnswer), nil).Run(runReadFunc(readAnswer)).Twice()

	helper := &ValidateHelperFuncs{
		FailsCount:               0,
		FailsCountBeforeValidate: 1,
		Message:                  "failure",
	}

	answer, err := FaskWhile("question", helper.counterReached, writerMock, readerMock)
	writerMock.AssertExpectations(t)
	readerMock.AssertExpectations(t)

	assert.Nil(t, err)
	assert.Equal(t, "answer", answer)

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

func TestFconfirm(t *testing.T) {
	// invalid answer then confirmed.
	question := []byte("question (y/n): ")
	validationMsg := []byte("please confirm or infirm by typing 'y' or 'n'\n")

	writerMock := new(WriterMock)
	writerMock.On("Write", question).Return(len(question), nil).Once()
	writerMock.On("Write", validationMsg).Return(len(validationMsg), nil).Once()
	writerMock.On("Write", question).Return(len(question), nil).Once()

	invalidAnswer := []byte("invalid\n")
	yesAnswer := []byte("y\n")
	readerMock := new(ReaderMock)
	readerMock.On("Read", mock.AnythingOfType("[]uint8")).Return(len(invalidAnswer), nil).Run(runReadFunc(invalidAnswer)).Once()
	readerMock.On("Read", mock.AnythingOfType("[]uint8")).Return(len(yesAnswer), nil).Run(runReadFunc(yesAnswer)).Once()

	confirmed, err := Fconfirm("question", writerMock, readerMock)
	assert.Nil(t, err)
	assert.True(t, confirmed)
	writerMock.AssertExpectations(t)
	readerMock.AssertExpectations(t)

	// infirmed.
	writerMock = new(WriterMock)
	writerMock.On("Write", question).Return(len(question), nil).Once()

	noAnswer := []byte("n\n")
	readerMock = new(ReaderMock)
	readerMock.On("Read", mock.AnythingOfType("[]uint8")).Return(len(noAnswer), nil).Run(runReadFunc(noAnswer)).Once()

	confirmed, err = Fconfirm("question", writerMock, readerMock)
	assert.Nil(t, err)
	assert.False(t, confirmed)
	writerMock.AssertExpectations(t)
	readerMock.AssertExpectations(t)

	// if write fails.
	brokenWriter := new(WriterMock)
	brokenWriter.On("Write", question).Return(0, io.EOF)

	confirmed, err = Fconfirm("question", brokenWriter, readerMock)
	assert.Error(t, err)
	assert.False(t, confirmed)
}
