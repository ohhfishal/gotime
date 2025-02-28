package cmd

import (
	"strings"
	"time"

	"github.com/ohhfishal/gotime/entry"
)

type LogCmd struct {
	Category string   `arg: "" required: ""`
	Note     []string `arg: "" optional: ""`
}

func (cmd *LogCmd) Run(entrySet EntrySet, now func() time.Time) error {
	newEntry := entry.Entry{
		Category: cmd.Category,
		Note:     strings.Join(cmd.Note, " "),
		Time:     time.Now(),
	}
	return entrySet.Append(newEntry)
}
