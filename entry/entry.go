package entry

import (
	"bufio"
	// "errors"
	"fmt"
	"io"
	"strings"
	"time"
)

type Entry struct {
	Time        time.Time
	Category    string
	Note string
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
