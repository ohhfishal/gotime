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

type TestConfig struct {
	T      *testing.T
	Envs   map[string]string
	Stdin  string
	stdout strings.Builder
	stderr strings.Builder
}

func NewTestConfig(t *testing.T) TestConfig {
	return TestConfig{
		T:    t,
		Envs: map[string]string{},
	}
}

func (tc *TestConfig) Config() Config {
	getEnv := func(key string) string {
		if val, ok := tc.Envs[key]; ok {
			return val
		}
		return ""
	}
	return Config{
		Stdin:  strings.NewReader(tc.Stdin),
		Stdout: &tc.stdout,
		Stderr: &tc.stderr,
		Getenv: getEnv,
	}
}

func (tc *TestConfig) Stdout() string {
	return tc.stdout.String()
}

func (tc *TestConfig) Stderr() string {
	return tc.stderr.String()
}

func (tc TestConfig) WithEnvs(envs map[string]string) TestConfig {
	for key, value := range envs {
		tc.Envs[key] = value
	}
	return tc
}

func (tc TestConfig) WithLogFile(entries []entry.Entry) TestConfig {
	file := newFile(tc.T, entries)
	tc.Envs[`GOTIME_LOG`] = file
	return tc
}

func newFile(t *testing.T, entries []entry.Entry) string {
	t.Helper()
	require := require.New(t)

	dir := t.TempDir()
	writer, err := os.CreateTemp(dir, "gotime-test-*.log")
	defer writer.Close()
	require.Nil(err)

	for _, newEntry := range entries {
		err := entry.Append(writer, newEntry)
		require.Nil(err)
	}

	return writer.Name()
}

type Expected struct {
	Entry  entry.Entry
	Err    error
	Stdout string
	Stderr string
}

type Test struct {
	Name string
	Args []string
	// GOTIME_LOG is set **always** to contain LogState
	Envs     map[string]string
	LogState []entry.Entry
	Expected Expected
}

func testCmd(t *testing.T, test Test) {
	t.Run(test.Name, func(t *testing.T) {
		tc := NewTestConfig(t).
			WithEnvs(test.Envs).
			WithLogFile(test.LogState)
		config := tc.Config()

		err := Run(test.Args, config)

		require := require.New(t)
		if errors.Is(test.Expected.Err, assert.AnError) {
			require.Error(err)
			return
		} else if test.Expected.Err != nil {
			require.ErrorIs(err, test.Expected.Err)
			return
		}

		require.Nil(err)
		assert.Equal(t, test.Expected.Stdout, tc.Stdout())
		assert.Equal(t, test.Expected.Stderr, tc.Stderr())

		entries, err := config.GetAll()
		require.Nil(err)
		if len(entries) != 0 {
			entries = entries[:len(entries)-1]
		}
		AssertEqualSlice(t, test.LogState, entries)
	})
}

func TestCmds(t *testing.T) {
	tests := []Test{
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
		testCmd(t, test)
	}
}
