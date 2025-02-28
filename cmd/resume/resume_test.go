package resume

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/ohhfishal/gotime/entry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var Error = errors.New("Test Error")
var Now = time.Now()

type Test struct {
	TimeNow   time.Time
	Entry     entry.Entry
	LastOfErr error
	AppendErr error

	// Runtime values
	appended []entry.Entry
}

func (test Test) Name() string {
	return fmt.Sprintf(
		`"%s" @ %s (%s, %s)`,
		test.Entry.Category,
		test.TimeNow.Format(time.DateOnly),
		ErrString(test.LastOfErr),
		ErrString(test.AppendErr),
	)
}

func (test Test) Run(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	cmd := CMD{
		Category: test.Entry.Category,
	}

	t.Run(test.Name(), func(t *testing.T) {
		test.appended = []entry.Entry{}
		err := cmd.Run(&test, test)
		if expectedErr := test.ExpectedErr(); expectedErr != nil {
			assert.ErrorIs(expectedErr, err)
			assert.Equal(0, len(test.appended), "entry was appended despite error")
			return
		}
		assert.Nil(err)
		require.Equal(1, len(test.appended), "single entry check")
		assert.Equal(test.ExpectedEntry(), test.appended[0])
	})
}

func TestResume(t *testing.T) {
	entries := []entry.Entry{
		entry.Entry{
			Time:     Now.Add(-1 * time.Hour),
			Category: "category",
			Note:     "Some note",
		},
		entry.Entry{},
	}
	nows := []time.Time{Now, time.Time{}}
	lastOfErrs := []error{nil, Error}
	appendErrs := lastOfErrs

	tests := []Test{}
	for _, entry := range entries {
		for _, now := range nows {
			for _, err1 := range lastOfErrs {
				for _, err2 := range appendErrs {
					tests = append(tests, Test{
						Entry:     entry,
						TimeNow:   now,
						LastOfErr: err1,
						AppendErr: err2,
					})
				}
			}
		}
	}

	for _, test := range tests {
		test.Run(t)
	}
}

func (test Test) ExpectedEntry() entry.Entry {
	expected := test.Entry
	expected.Note = "Cont: " + test.Entry.Note
	expected.Time = test.Now()
	return expected

}

func (test Test) ExpectedErr() error {
	return errors.Join(test.LastOfErr, test.AppendErr)
}

func (test Test) LastOf(_ string) (entry.Entry, error) {
	return test.Entry, test.LastOfErr
}

func (test *Test) Append(e entry.Entry) error {
	if test.AppendErr != nil {
		return test.AppendErr
	}
	test.appended = append(test.appended, e)
	return nil
}

func (test Test) Now() time.Time {
	return test.TimeNow
}

func ErrString(err error) string {
	if err == nil {
		return `nil`
	}
	return `err`
}
