package entry

import (
	"time"
)

const timeLayout = time.DateTime

type Entry struct {
	Time     time.Time
	Category string // TODO: Make its own type with sub categories ex: "foo/bar"
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

func DurationMap(entries []Entry, cutoff time.Time) map[string]time.Duration {
	totals := map[string]time.Duration{}
	// Ensure we always have an end time
	entries = append(entries, Entry{Time: cutoff})
	for i, entry := range entries {
		if i != len(entries)-1 {
			next := entries[i+1]
			duration := next.Time.Sub(entry.Time)
			if _, ok := totals[entry.Category]; ok {
				totals[entry.Category] += duration
			} else {
				totals[entry.Category] = duration
			}
		}
	}
	return totals
}

func Total(totals map[string]time.Duration) time.Duration {
	var duration time.Duration
	for _, value := range totals {
		duration += value
	}
	return duration
}

func Compare(a, b Entry) int {
	return a.Time.Compare(b.Time)
}
