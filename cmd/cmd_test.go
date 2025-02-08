package cmd

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ohhfishal/gotime/entry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
	file, err := os.CreateTemp(dir, "gotime-test-*.log")
	defer file.Close()
	require.Nil(err)

	for _, newEntry := range entries {
		err := entry.Write(file, newEntry)
		require.Nil(err)
	}

	return file.Name()
}

func TestCommand(t *testing.T) {
}
