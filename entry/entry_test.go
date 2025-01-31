package entry

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"time"
)

const TimeFormat = time.DateTime

func TestEncodeDecode(t *testing.T) {
	type EncodeDecodeTest struct {
		Name  string
		Entry Entry
	}

	tests := []EncodeDecodeTest{
		{
			Name:  "empty",
			Entry: Entry{},
		},
		{
			Name: "normal",
			Entry: Entry{
				Time:        time.Now(),
				Category:    "Project A",
				Note: "Important Task",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			require := require.New(t)

			var buffer strings.Builder
			err := Write(&buffer, test.Entry)
			require.Nil(err)

			t.Log(buffer.String())
			reader := strings.NewReader(buffer.String())
			entry, err := Read(reader)
			require.Nil(err)
			RequireEntryEqual(t, *entry, test.Entry)
		})
	}
}

func TestReadAll(t *testing.T) {
	type Test struct {
		Name    string
		Entries []Entry
	}
	tests := []Test{
		{
			Name: "empty",
		},
		{
			Name: "happy case",
			Entries: []Entry{
				{
					Time:        time.Now(),
					Category:    "Project A",
					Note: "Important Task",
				},
				{
					Category:    "Project B",
					Note: "Important Task B",
				},
				{
					Time:        time.Now(),
					Category:    "Project C",
					Note: "Important Task C",
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)
			var buffer strings.Builder

			for _, entry := range test.Entries {
				err := Write(&buffer, entry)
				require.Nil(err)
			}

			reader := strings.NewReader(buffer.String())
			t.Logf("Wrote:\n%s", buffer.String())
			entries, err := ReadAll(reader)
			require.Nil(err)

			assert.Equal(len(test.Entries), len(entries), "entries length")
			for i, entry := range entries {
				require.LessOrEqual(i, len(test.Entries), "out of bounds")
				RequireEntryEqual(t, entry, test.Entries[i])
				t.Log("Correct:", entry)
			}
		})
	}
}

func RequireEntryEqual(tb testing.TB, expected, actual Entry) {
	require := require.New(tb)
	require.Equal(expected.Category, actual.Category)
	require.Equal(expected.Note, actual.Note)
	require.Equal(expected.Time.Format(TimeFormat), actual.Time.Format(TimeFormat))
}
