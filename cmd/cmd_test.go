package cmd

import (
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ohhfishal/gotime/entry"
	"github.com/ohhfishal/gotime/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var now = time.Now()
var FullLog = []entry.Entry{
	{Time: now, Category: "category", Note: "First"},
	{Time: now, Category: "category", Note: "Original"},
	{Time: now, Category: "junk", Note: "Bad Note"},
}

var EmptyLog = []entry.Entry{}

func TestResume(t *testing.T) {
	aliases := []string{"resume", "continue", "cont"}
	for _, alias := range aliases {
		testResumeValidArgs(t, alias)
	}
}

func testResumeValidArgs(t *testing.T, alias string) {
	args := []string{alias, "category"}
	type Expected struct {
		Entry  entry.Entry
		Err    error
		Stdout string
		Stderr string
	}
	tests := []struct {
		Name     string
		Args     []string
		Envs     map[string]string
		LogState []entry.Entry
		Expected Expected
	}{
		{
			Name:     "valid input/valid file",
			LogState: FullLog,
			Expected: Expected{
				Entry: entry.Entry{
					Time: now, Category: "category", Note: "Cont: Original",
				},
			},
		},
		{
			Name:     "empty file",
			LogState: EmptyLog,
			Expected: Expected{
				Err: ErrCategoryNotFound,
			},
		},
		{
			Name: "missing file",
			Expected: Expected{
				Err: ErrCategoryNotFound,
			},
		},
	}
	for _, test := range tests {
		t.Run(alias+":"+test.Name, func(t *testing.T) {
			tc := NewTestConfig(t).
				WithEnvs(test.Envs).
				WithLogFile(test.LogState)
			config := tc.Config()

			if test.Args != nil {
				args = test.Args
			}
			err := Run(args, config)

			require := require.New(t)
			if errors.Is(test.Expected.Err, assert.AnError) {
				require.Error(err)
				return
			}
			if test.Expected.Err != nil {
				require.NotNil(err)
				return
			}

			require.Nil(err)
			assert.Equal(t, test.Expected.Stdout, tc.Stdout())
			assert.Equal(t, test.Expected.Stderr, tc.Stderr())

			entries, err := config.GetAllEntries()
			require.Nil(err)
			if len(entries) != 0 {
				entries = entries[:len(entries)-1]
			}
			AssertEqualSlice(t, test.LogState, entries)

			last, err := config.LastEntry()
			require.Nil(err)
			AssertEqual(t, test.Expected.Entry, *last)

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
		err := file.Write(writer, newEntry)
		require.Nil(err)
	}

	return writer.Name()
}

func TestCommand(t *testing.T) {
}
