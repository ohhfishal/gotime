package cmd

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ohhfishal/gotime/entry"
)

type LogCmd struct {
	Category string    `arg:"" required:"" help:"Category to log the new entry under. (Use '-' to use the last category logged)"`
	Note     []string  `arg:"" optional:"" help:"Description for the entry"`
	Date     time.Time `short:"d" optional:"" format:"2006-01-02" help:"Date to log (default: today)"`
	Time     time.Time `short:"t" optional:"" format:"15:04" help:"Time to log entry at (default: now)"`
}

func (cmd *LogCmd) AfterApply() error {
	if cmd.Date.IsZero() {
		now := time.Now()
		cmd.Date = time.Date(
			now.Year(), now.Month(), now.Day(),
			0, 0, 0, 0,
			now.Location(),
		)
	}

	if cmd.Time.IsZero() {
		cmd.Time = time.Now()
	}

	// Combine time and date
	cmd.Time = time.Date(
		cmd.Date.Year(), cmd.Date.Month(), cmd.Date.Day(),
		cmd.Time.Hour(), cmd.Time.Minute(), cmd.Time.Second(), 0,
		cmd.Date.Location(),
	)
	return nil
}

func (cmd *LogCmd) Run(log string) error {
	if cmd.Category == "-" {
		entries, err := entry.ReadAllFromFile(log)
		if err != nil {
			return fmt.Errorf(`could not replace "-": %w`, err)
		} else if len(entries) == 0 {
			return errors.New(`could not replace "-": no entries in log`)
		}
		cmd.Category = entries[len(entries)-1].Category
	}

	newEntry := entry.Entry{
		Category: cmd.Category,
		Note:     strings.Join(cmd.Note, " "),
		Time:     cmd.Time,
	}
	return entry.AppendFile(log, newEntry)
}
