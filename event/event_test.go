package event

import (
	"strings"
	"testing"
	"time"

	"github.com/ohhfishal/gotime/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


func AssertEqual(t *testing.T, expected, result Event) {
  assert := assert.New(t)
  assert.Equal(expected.Title, result.Title)
  assert.Equal(expected.Description, result.Description)
  // TODO: If this ever yells make it more information
  assert.True(expected.Date.Equal(result.Date))
  assert.Equal(expected.Tags, result.Tags)
}

func TestEncodeDecode(t *testing.T) {
	type EncodeDecodeTest struct {
		Name  string
		Event Event
	}

	tests := []EncodeDecodeTest{
		{
			Name:  "empty",
			Event: Event{},
		},
		{
			Name:  "example",
			Event: Event{
        Title: "example",
        Description: "some cool event",
        Date: time.Now(),
        Tags: []Tag{"A", "B"},
      },
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			require := require.New(t)

			var buffer strings.Builder
			err := file.Write(&buffer, test.Event)
			require.Nil(err)

			t.Log(buffer.String())
			entry, err := Decode(buffer.String())
			require.Nil(err)
			AssertEqual(t, entry, test.Event)
		})
	}
}
