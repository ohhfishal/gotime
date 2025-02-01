package entry

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"time"
)

type Entry struct {
	Time     time.Time
	Category string
	Note     string
}

func (entry Entry) String() string {
	// TODO: REMOVE this is config dependant
	msg := entry.Time.Format(time.Kitchen) + " " + entry.Category
	if entry.Note == "" {
		return msg
	}
	return msg + fmt.Sprintf(` "%s"`, entry.Note)
}

func ReadAll(reader io.Reader, config ...Config) ([]Entry, error) {
	entries := []Entry{}
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		entry, err := Read(strings.NewReader(scanner.Text()), config...)
		if err != nil {
			return []Entry{}, fmt.Errorf("reading line: %w", err)
		}
		entries = append(entries, *entry)
	}

	if err := scanner.Err(); err != nil {
		return []Entry{}, fmt.Errorf("reading: %w", err)
	}

	return entries, nil
}

func Write(writer io.Writer, entry Entry, config ...Config) error {
	cfg := DefaultConfig(config...)

	line, err := cfg.Encode(entry)
	if err != nil {
		return fmt.Errorf("encoding entry: %w", err)
	}

	_, err = fmt.Fprintln(writer, line)
	return err
}

func Read(reader io.Reader, config ...Config) (*Entry, error) {
	cfg := DefaultConfig(config...)
	bytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("reading: %w", err)
	}
	entry, err := cfg.Decode(string(bytes))
	if err != nil {
		return nil, fmt.Errorf("parsing: %w", err)
	}
	return entry, nil
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
