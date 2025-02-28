package resume

import (
	"errors"
	"testing"
	"time"

	"github.com/ohhfishal/gotime/entry"
	"github.com/stretchr/testify/assert"
)

var Error = errors.New("Test Error")
var now = time.Now()

func Now() time.Time { return now }

var FullLog = []entry.Entry{
	{Time: now, Category: "First", Note: "First"},
	{Time: now, Category: "Second", Note: "Second"},
	{Time: now, Category: "Third", Note: "Third"},
}

type Mock struct {
	Entries   []entry.Entry
	GetAllErr error
	AppendErr error
}

func TestResumeErr(t *testing.T) {
	ValidCMD := CMD{
		Category: "First",
	}

	tests := []struct {
		Name string
		CMD  CMD
		Mock Mock
	}{
		{
			Name: "empty",
			CMD:  ValidCMD,
		},
		{
			Name: "GetAll error",
			CMD:  ValidCMD,
			Mock: Mock{
				Entries:   FullLog,
				GetAllErr: Error,
			},
		},
		{
			Name: "Append error",
			CMD:  ValidCMD,
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
			err := test.CMD.Run(&test.Mock, Now)
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

			cmd := CMD{
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
