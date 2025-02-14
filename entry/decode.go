package entry

import (
	"fmt"
	"strings"
	"time"
)

func Decode(line string) (Entry, error) {
	var entry Entry
	var err error

	// Parse time
	timeLength := len(timeLayout)
	entry.Time, err = time.Parse(timeLayout, line[:timeLength])
	if err != nil {
		return Entry{}, fmt.Errorf("parsing time: %w", err)
	}

	// Parse Category
	line = line[timeLength:]
	if line == `` {
		return Entry{}, fmt.Errorf("unexpected eof")
	}

	if line[0] != ',' {
		return Entry{}, fmt.Errorf(`expected :",": found "%b"`, line[0])
	}
	cutIndex := strings.Index(line, ":")
	if cutIndex == -1 {
		return Entry{}, fmt.Errorf("unexpected eof")
	}
	entry.Category = line[1:cutIndex]
	entry.Note = strings.TrimSpace(line[cutIndex+1:])
	return entry, nil
}
