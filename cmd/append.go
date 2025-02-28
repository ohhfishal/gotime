package cmd

import (
	"errors"
	"fmt"
	"github.com/ohhfishal/gotime/entry"
	"strings"
	"time"
)

var ErrFileEmpty = errors.New("no entries in file")

type AppendCmd struct {
	Note []string `arg: "" optional: ""`
}

func (cmd *AppendCmd) Run(entrySet EntrySet, now func() time.Time) error {
	entries, err := entrySet.GetAll()
	if err != nil {
		return fmt.Errorf(`reading entries: %w`, err)
	}

	if len(entries) == 0 {
		return ErrFileEmpty
	}

	last := entries[len(entries)-1]
	newEntry := entry.Entry{
		Category: last.Category,
		Note:     strings.Join(cmd.Note, " "),
		Time:     time.Now(),
	}
	return entrySet.Append(newEntry)
}
