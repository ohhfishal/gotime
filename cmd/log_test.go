package cmd

import (
	"errors"
	"testing"
	"time"

	"github.com/ohhfishal/gotime/entry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var now = time.Now()

func TestLog(t *testing.T) {
	type Expected struct {
		Entry  entry.Entry
		Err    error
		Stdout string
		Stderr string
	}
	tests := []struct {
		Name string
		Args []string
		// GOTIME_LOG is set **always** to contain LogState
		Envs     map[string]string
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
			Name: "new/valid input/empty file",
			Args: []string{"new", "happy", "Cool", "note"},
			Expected: Expected{
				Entry: entry.Entry{Category: "happy", Note: "Cool note"},
			},
		},
		{
			Name: "new/valid input/existing file",
			Args: []string{"new", "second", "Note"},
			LogState: []entry.Entry{
				{Time: now, Category: "first", Note: "Should remain unchanged"},
			},
			Expected: Expected{
				Entry: entry.Entry{Category: "second", Note: "Note"},
			},
		},
		{
			Name: "append/valid input/valid file",
			Args: []string{"append", "Cool", "note"},
			LogState: []entry.Entry{
				{Time: now, Category: "current", Note: "Should remain unchanged"},
			},
			Expected: Expected{
				Entry: entry.Entry{Category: "current", Note: "Cool note"},
			},
		},
		{
			Name: "append/empty file",
			Args: []string{"append", "Cool", "note"},
			Expected: Expected{
				Err: ErrFileEmpty,
			},
		},
		{
			Name: "continue/valid input/valid file",
			Args: []string{"continue", "current"},
			LogState: []entry.Entry{
				{Time: now, Category: "current", Note: "Original"},
				{Time: now, Category: "junk", Note: "Bad Note"},
			},
			Expected: Expected{
				Entry: entry.Entry{
					Time: now, Category: "current", Note: "Cont: Original",
				},
			},
		},
		{
			Name: "continue/empty file",
			Args: []string{"continue", "category"},
			Expected: Expected{
				Err: ErrFileEmpty,
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
			tc := NewTestConfig(t).
				WithEnvs(test.Envs).
				WithLogFile(test.LogState)
			config := tc.Config()

			err := Log(config, test.Args...)

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

			entries, err := config.GetAllEntries()
			require.Nil(err)
			expected := test.LogState
			var zero entry.Entry
			if test.Expected.Entry != zero {
				expected = append(expected, test.Expected.Entry)
			}
			AssertEqualSlice(t, expected, entries)
		})
	}
}
