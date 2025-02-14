package cmd

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ohhfishal/gotime/entry"
)

var ErrFileEmpty = errors.New("no entries in file")
var ErrCategoryNotFound = errors.New("category not found")

type LogCmd struct {
	New      NewCmd      `default:"withargs" cmd: "" help:"Log using a custom category and note."`
	Append   AppendCmd   `cmd: "" help:"Log using the last category in the log."`
	Continue ContinueCmd `cmd: "" help:"Log using the last message of a given category."`
}

type NewCmd struct {
	Category string   `arg: "" required: ""`
	Note     []string `arg: "" optional: ""`
}

type AppendCmd struct {
	Note []string `arg: "" optional: ""`
}

type ContinueCmd struct {
	Category string `arg: "" required: ""`
}

func (cmd *NewCmd) Run(cfg Config) error {
	newEntry := entry.Entry{
		Category: cmd.Category,
		Note:     strings.Join(cmd.Note, " "),
		Time:     time.Now(),
	}
	return cfg.Write(newEntry)
}

func (cmd *AppendCmd) Run(cfg Config) error {
	last, err := cfg.LastEntry()
	if err != nil {
		return err
	}

	newEntry := entry.Entry{
		Category: last.Category,
		Note:     strings.Join(cmd.Note, " "),
		Time:     time.Now(),
	}
	return cfg.Write(newEntry)
}

func (cmd *ContinueCmd) Run(cfg Config) error {
	last, err := cfg.LastEntryOf(cmd.Category)
	if err != nil {
		return err
	}

	newEntry := entry.Entry{
		Category: cmd.Category,
		Note:     "Cont: " + last.Note,
		Time:     time.Now(),
	}
	return cfg.Write(newEntry)
}

func (c Config) LastEntryOf(category string) (*entry.Entry, error) {
	entries, err := c.GetAllEntries()
	if err != nil {
		return nil, fmt.Errorf(`reading entries: %w`, err)
	}

	var last *entry.Entry
	for _, entry := range entries {
		if entry.Category == category {
			last = &entry
		}
	}

	if last == nil {
		return nil, ErrCategoryNotFound
	}
	return last, nil
}

func (c Config) LastEntry() (*entry.Entry, error) {
	entries, err := c.GetAllEntries()
	if err != nil {
		return nil, fmt.Errorf(`reading entries: %w`, err)
	}

	if len(entries) == 0 {
		return nil, ErrFileEmpty
	}
	return &entries[len(entries)-1], nil
}

func (c Config) Write(newEntry entry.Entry) error {
	writer, err := c.OpenWriteCloser()
	if err != nil {
		return fmt.Errorf(`opening file: %w`, err)
	}
	defer writer.Close()

	// TODO: Pass in an entry.Config to allow different file encodings EX timefmt
	if err := entry.Write(writer, newEntry); err != nil {
		return fmt.Errorf(`writing to file: %w`, err)
	}
	return nil
}

func Log(cfg Config, args ...string) error {
	var cmd LogCmd
	return RunCmd(cfg, &cmd, args...)
}
