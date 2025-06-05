package entry

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"slices"
	"time"
)

const NUM_RECORDS = 3
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
	slices.SortFunc(filtered, Compare)
	return filtered
}

// TODO: This function is never used...
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

// TODO: This function is not used either...
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

func Append(file io.Writer, entry Entry) error {
	writer := csv.NewWriter(file)
	if err := writer.Write([]string{
		entry.Time.Format(timeLayout),
		entry.Category,
		entry.Note,
	}); err != nil {
		return fmt.Errorf(`writing to file: %w`, err)
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return fmt.Errorf(`flushing writer: %w`, err)
	}
	return nil
}

func ReadAll(file io.Reader) ([]Entry, error) {
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf(`could not read file: %w`, err)
	}

	var entries []Entry
	for _, record := range records {
		if len(record) != NUM_RECORDS {
			return nil, fmt.Errorf(
				`record contained %d columns but expected %d`, len(record), NUM_RECORDS,
			)
		}

		var entry Entry
		noLoc, err := time.Parse(timeLayout, record[0])
		if err != nil {
			return nil, fmt.Errorf(`invalid time "%s": %w`, record[0], err)
		}
		entry.Time = time.Date(
			noLoc.Year(), noLoc.Month(), noLoc.Day(),
			noLoc.Hour(), noLoc.Minute(), noLoc.Second(), noLoc.Nanosecond(),
			time.Now().Location(),
		)
		entry.Category = record[1]
		entry.Note = record[2]
		entries = append(entries, entry)
	}
	return entries, nil
}

func AppendFile(path string, entry Entry) error {
	file, err := os.OpenFile(
		path,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0666,
	)
	if err != nil {
		return fmt.Errorf(`could not open file: %w`, err)
	}
	defer file.Close()
	return Append(file, entry)
}

func ReadAllFromFile(path string) ([]Entry, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf(`can not open file: %s`, path)
	}
	defer file.Close()
	return ReadAll(file)
}
