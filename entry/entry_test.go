package entry

import (
	"strings"
	"testing"
	"time"

	"github.com/ohhfishal/gotime/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const TimeFormat = time.DateTime

func TestEncodeDecode(t *testing.T) {
	type EncodeDecodeTest struct {
		Name  string
		Entry Entry
	}

	tests := []EncodeDecodeTest{
		{
			Name: "empty", Entry: Entry{},
		},
		{
			Name: "normal",
			Entry: Entry{
				Time:     time.Now(),
				Category: "Project A",
				Note:     "Important Task",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			require := require.New(t)

			var buffer strings.Builder
			err := file.Write(&buffer, test.Entry)
			require.Nil(err)

			t.Log(buffer.String())
			entry, err := Decode(buffer.String())
			require.Nil(err)
			RequireEntryEqual(t, entry, test.Entry)
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
					Time:     time.Now(),
					Category: "Project A",
					Note:     "Important Task",
				},
				{
					Category: "Project B",
					Note:     "Important Task B",
				},
				{
					Time:     time.Now(),
					Category: "Project C",
					Note:     "Important Task C",
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
				err := file.Write(&buffer, entry)
				require.Nil(err)
			}

			reader := strings.NewReader(buffer.String())
			t.Logf("Wrote:\n%s", buffer.String())
			entries, err := file.ReadAll(reader, Decode)
			// entries, err := ReadAll(reader)
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

func TestFilter(t *testing.T) {
	entries := []Entry{}
	for i := 0; i < 10; i++ {
		entries = append(entries, Entry{Time: time.Now()})
	}
	filtered := Filter(entries, entries[1].Time, entries[8].Time)
	RequireEntrySliceEqual(t, entries[1:9], filtered)
}

func TestDurationMaps(t *testing.T) {
	type Test struct {
		Name       string
		Times      []string
		Categories []string
		Cutoff     string
		Expected   map[string]string
	}

	tests := []Test{
		{
			Name: `basic example works`,
			Times: []string{
				`2025-01-01 08:00:00`,
			},
			Categories: []string{
				`test/a`,
			},
			Cutoff: `2025-01-01 08:30:00`,
			Expected: map[string]string{
				`test/a`: `30m0s`,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			require := require.New(t)
			entries := []Entry{}
			require.Equal(len(test.Times), len(test.Categories))
			for i, strTime := range test.Times {
				parsed, err := time.Parse(time.DateTime, strTime)
				require.Nil(err)
				entry := Entry{
					Time:     parsed,
					Category: test.Categories[i],
				}
				entries = append(entries, entry)
			}
			parsed, err := time.Parse(time.DateTime, test.Cutoff)
			require.Nil(err)
			result := DurationMap(entries, parsed)
			for key, expected := range test.Expected {
				unparsed, ok := result[key]
				require.True(ok)
				require.Equal(expected, unparsed.String())
			}
		})

	}
}

func TestTotal(t *testing.T) {
	t.Fail()
}

func TestCompare(t *testing.T) {
	a := Entry{Time: time.Now()}
	require := require.New(t)
	b := Entry{Time: time.Now()}
	require.Equal(Compare(a, b), -1, " a  b")
	require.Equal(Compare(a, a), 0, " a  a")
	require.Equal(Compare(b, a), 1, " b  a")
}

func RequireEntrySliceEqual(tb testing.TB, expected, actual []Entry) {
	require := require.New(tb)
	require.Equal(len(expected), len(actual))
	for i, entry := range actual {
		require.LessOrEqual(i, len(expected))
		RequireEntryEqual(tb, expected[i], entry)
	}

}

func RequireEntryEqual(tb testing.TB, expected, actual Entry) {
	require := require.New(tb)
	require.Equal(expected.Category, actual.Category)
	require.Equal(expected.Note, actual.Note)
	require.Equal(expected.Time.Format(TimeFormat), actual.Time.Format(TimeFormat))
}
