package cmd

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/alecthomas/kong"
	"github.com/ohhfishal/gotime/entry"
)

var ErrFileEmpty = errors.New("no entries in file")

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
	last, err := cfg.LastEntry()
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

func (c Config) LastEntry() (*entry.Entry, error) {
	entries, err := c.GetAllEntries()
	if err != nil {
		return nil, fmt.Errorf(`reading entries: %w`, err)
	}

	if len(entries) == 0 {
		return nil, ErrFileEmpty
	}
	return &entries[0], nil
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
	var exit bool
	parser, err := kong.New(
		&cmd,
		kong.Name("gotime log"),
		kong.Description("Log a new entry in timesheet file."),
		kong.Exit(func(_ int) { exit = true }),
		kong.Bind(cfg),
	)
	if err != nil {
		return err
	}
	parser.Stdout = cfg.Stdout
	parser.Stderr = cfg.Stderr

	context, err := parser.Parse(args)
	if err != nil || exit {
		return err
	}

	err = context.Run()
	if err != nil {
		return err
	}
	return nil
}
