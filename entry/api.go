package entry

import (
	"fmt"
	"time"
)

type FileHandler struct {
	filename string
}

type Option func([]Entry) ([]Entry, error)

func NewFileHandler(filename string) (*FileHandler, error) {
	return &FileHandler{
		filename: filename,
	}, nil
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
		filtered := []Entry{}
		for _, entry := range entries {
			if entry.Time.After(today) && entry.Time.Before(tomorrow) {
				filtered = append(filtered, entry)
			}
		}
		return filtered, nil
	}
}
