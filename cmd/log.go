package cmd

import (
	"strings"
	"time"

	"github.com/ohhfishal/gotime/entry"
)

type LogCmd struct {
	Category string    `arg:"" required:"" help:"Category to log the new entry under"`
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

func (cmd *LogCmd) Run(entrySet EntrySet) error {
	newEntry := entry.Entry{
		Category: cmd.Category,
		Note:     strings.Join(cmd.Note, " "),
		Time:     cmd.Time,
	}
	return entrySet.Append(newEntry)
}
