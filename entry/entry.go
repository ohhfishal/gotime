package entry

import (
	"time"
)

const timeLayout = time.DateTime

type Entry struct {
	Time     time.Time
	Category string
	Note     string
}

func Filter(entries []Entry, start, end time.Time) []Entry {
	filtered := []Entry{}
	for _, entry := range entries {
		time := entry.Time
		if (time.Equal(start) || time.After(start)) &&
			(time.Equal(end) || time.Before(end)) {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}

func Compare(a, b Entry) int {
	return a.Time.Compare(b.Time)
}
