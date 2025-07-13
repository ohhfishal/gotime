package entry

import (
	"fmt"
	"slices"
	"time"
)

type FileHandler struct {
	filename string
}

type Summary struct {
	Schedule []ScheduleEntry
	// TODO: Include cagtegory map
	Total time.Duration
}

type ScheduleEntry struct {
	Entry    Entry
	Duration time.Duration
}

type Option func([]Entry) ([]Entry, error)

func NewFileHandler(filename string) (*FileHandler, error) {
	return &FileHandler{
		filename: filename,
	}, nil
}

func Summarize(entries []Entry) Summary {
	var summary Summary
	if len(entries) == 0 {
		return summary
	}

	entries = append(
		entries,
		Entry{Time: time.Now()},
	)

	for i, entry := range entries[:len(entries)-1] {
		entry := ScheduleEntry{
			Entry:    entry,
			Duration: entries[i+1].Time.Sub(entry.Time),
		}
		// TODO: Implement here instead
		// if _, ok := args.CategoryBreakdown[entry.Category]; ok {
		// 	args.CategoryBreakdown[entry.Category] += tuple.Duration
		// } else {
		// 	args.CategoryBreakdown[entry.Category] = tuple.Duration
		// }
		summary.Schedule = append(summary.Schedule, entry)
		summary.Total += entry.Duration
	}
	return summary
}

func (handler *FileHandler) CreateEntry(newEntry Entry) error {
	return AppendFile(handler.filename, newEntry)
}

func (handler *FileHandler) GetAllEntries(options ...Option) ([]Entry, error) {
	var entries []Entry
	var err error

	// TODO: Inline this function, should not be used outside of this struct now
	entries, err = ReadAllFromFile(handler.filename)
	if err != nil {
		return nil, fmt.Errorf(`getting entries: %w`, err)
	}

	for _, option := range options {
		entries, err = option(entries)
		if err != nil {
			return nil, fmt.Errorf(`filtering entries: %w`, err)
		}
	}
	slices.SortFunc(entries, Compare)
	return entries, nil
}

func Today() Option {
	now := time.Now()
	today := time.Date(
		now.Year(), now.Month(), now.Day(),
		0, 0, 0, 0,
		now.Location(),
	)
	tomorrow := today.Add(24 * time.Hour)
	return func(entries []Entry) ([]Entry, error) {
		return Filter(entries, today, tomorrow), nil
	}
}
