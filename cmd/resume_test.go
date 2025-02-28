package cmd

import (
	"errors"
	"testing"
	"time"

	"github.com/ohhfishal/gotime/entry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var Error = errors.New("Test Error")

func Now() time.Time { return now }

type Mock struct {
	Entries   []entry.Entry
	GetAllErr error
	AppendErr error
}

func TestResume(t *testing.T) {
	aliases := []string{"resume", "continue", "cont"}
	entries := FullLog
	for _, alias := range aliases {
		for _, e := range entries {
			testResumeValidArgs(t, alias, e)

			t.Run(alias+"/no args", func(t *testing.T) {
				tc := NewTestConfig(t)
				config := tc.Config()
				err := Run([]string{alias}, config)
				require.NotNil(t, err)
			})
		}
	}
}

func testResumeValidArgs(t *testing.T, alias string, match entry.Entry) {
	args := []string{alias, match.Category}
	tests := []Test{
		{
			Name:     "valid input/valid file",
			LogState: FullLog,
			Expected: Expected{
				Entry: entry.Entry{
					Time: now, Category: match.Category, Note: "Cont: " + match.Note,
				},
			},
		},
		{
			Name:     "empty file",
			LogState: EmptyLog,
			Expected: Expected{
				Err: assert.AnError,
			},
		},
		{
			Name: "missing file",
			Expected: Expected{
				Err: assert.AnError,
			},
		},
	}
	for _, test := range tests {
		test.Args = args
		testCmd(t, test)
	}
}

func TestResumeErr(t *testing.T) {
	ValidResume := Resume{
		Category: "First",
	}

	tests := []struct {
		Name   string
		Resume Resume
		Mock   Mock
	}{
		{
			Name:   "empty",
			Resume: ValidResume,
		},
		{
			Name:   "GetAll error",
			Resume: ValidResume,
			Mock: Mock{
				Entries:   FullLog,
				GetAllErr: Error,
			},
		},
		{
			Name:   "Append error",
			Resume: ValidResume,
			Mock: Mock{
				Entries:   FullLog,
				AppendErr: Error,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			assert := assert.New(t)
			entries := test.Mock.Entries
			err := test.Resume.Run(&test.Mock, Now)
			assert.NotNil(err)
			// Ensure nothing changed
			assert.Equal(entries, test.Mock.Entries)
		})
	}

}

func TestResumeHappyPath(t *testing.T) {
	mock := Mock{
		Entries: FullLog,
	}
	for _, entry := range FullLog {
		t.Run(entry.Category, func(t *testing.T) {
			assert := assert.New(t)
			expected := entry
			expected.Note = "Cont: " + entry.Note

			cmd := Resume{
				Category: entry.Category,
			}
			err := cmd.Run(&mock, Now)

			assert.Nil(err)
			assert.Equal(expected, mock.Entries[len(mock.Entries)-1])
		})
	}
}

func (m Mock) GetAll() ([]entry.Entry, error) {
	return m.Entries, m.GetAllErr
}

func (m *Mock) Append(e entry.Entry) error {
	if m.AppendErr != nil {
		return m.AppendErr
	}
	m.Entries = append(m.Entries, e)
	return nil
}
