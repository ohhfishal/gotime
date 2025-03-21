package cmd

import (
	"strings"
	"time"

	"github.com/ohhfishal/gotime/entry"
)

type LogCmd struct {
	Category string        `arg:"" required:"" help:"Category to log the new entry under"`
	Note     []string      `arg:"" optional:"" help:"Description for the entry"`
	Offset   time.Duration `short:"O" default:"0m" help:"Offset from now to set the entry as"`
}

func (cmd *LogCmd) Run(entrySet EntrySet, now func() time.Time) error {
	newEntry := entry.Entry{
		Category: cmd.Category,
		Note:     strings.Join(cmd.Note, " "),
		Time:     now().Add(cmd.Offset),
	}
	return entrySet.Append(newEntry)
}
