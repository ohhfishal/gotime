package cmd

import (
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ohhfishal/gotime/entry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var now = time.Now()
var FullLog = []entry.Entry{
	{Time: now, Category: "First", Note: "First"},
	{Time: now, Category: "Second", Note: "Second"},
	{Time: now, Category: "Third", Note: "Third"},
}

var EmptyLog = []entry.Entry{}

func TestCommands(t *testing.T) {
	dir := t.TempDir()
	type Expected struct {
		Entry entry.Entry
		Err   error
	}
	tests := []struct {
		Name     string
		Args     []string
		LogState []entry.Entry
		Expected Expected
	}{
		{
			Name: "default/valid input/empty file",
			Args: []string{"happy", "Cool", "note"},
			Expected: Expected{
				Entry: entry.Entry{Category: "happy", Note: "Cool note"},
			},
		},
		{
			Name: "default/valid input/existing file",
			Args: []string{"second", "Note"},
			LogState: []entry.Entry{
				{Time: now, Category: "first", Note: "Should remain unchanged"},
			},
			Expected: Expected{
				Entry: entry.Entry{Category: "second", Note: "Note"},
			},
		},
		{
			Name: "log/valid input/empty file",
			Args: []string{"log", "happy", "Cool", "note"},
			Expected: Expected{
				Entry: entry.Entry{Category: "happy", Note: "Cool note"},
			},
		},
		{
			Name: "log/valid input/existing file",
			Args: []string{"log", "second", "Note"},
			LogState: []entry.Entry{
				{Time: now, Category: "first", Note: "Should remain unchanged"},
			},
			Expected: Expected{
				Entry: entry.Entry{Category: "second", Note: "Note"},
			},
		},
		{
			Name: "help/short",
			Args: []string{"-h"},
			Expected: Expected{
				Err: assert.AnError,
			},
		},
		{
			Name: "help/long",
			Args: []string{"--help"},
			Expected: Expected{
				Err: assert.AnError,
			},
		},
		{
			Name: "bad flag",
			Args: []string{"--badflag123"},
			Expected: Expected{
				Err: assert.AnError,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			require := require.New(t)
			var stdout, stderr strings.Builder
			config := Config{
				Stdin:  strings.NewReader(""),
				Stdout: &stdout,
				Stderr: &stderr,
			}

			f, err := os.CreateTemp(dir, "gotime-test-*.log")
			require.Nil(err)
			defer os.Remove(f.Name())
			for _, newEntry := range test.LogState {
				err := entry.Append(f, newEntry)
				require.Nil(err)
			}
			f.Close()

			t.Log(f.Name())
			test.Args = append(test.Args, "--log-file", f.Name())

			err = Run(test.Args, config)
			if errors.Is(test.Expected.Err, assert.AnError) {
				require.Error(err)
				return
			} else if test.Expected.Err != nil {
				require.ErrorIs(err, test.Expected.Err)
				return
			}
			t.Log(err)
			require.Nil(err)

			entries, err := entry.ReadAllFromFile(f.Name())
			require.Nil(err)
			if len(entries) != 0 {
				entries = entries[:len(entries)-1]
			}
			AssertEqualSlice(t, test.LogState, entries)
		})
	}
}

func AssertEqualSlice(t *testing.T, expected, actual []entry.Entry) {
	t.Helper()
	if !assert.Equal(t, len(expected), len(actual)) {
		return
	}
	for i := range len(expected) {
		expectedEntry := expected[i]
		actualEntry := actual[i]
		AssertEqual(t, expectedEntry, actualEntry)
	}
}

func AssertEqual(t *testing.T, expected, actual entry.Entry) {
	t.Helper()
	assert := assert.New(t)
	assert.Equal(expected.Category, actual.Category)
	assert.Equal(expected.Note, actual.Note)
	if !expected.Time.IsZero() {
		expectString := expected.Time.Format(time.DateTime)
		actualString := actual.Time.Format(time.DateTime)
		assert.Equal(expectString, actualString)
	}
}
