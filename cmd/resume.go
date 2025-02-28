package cmd

import (
	"errors"
	"fmt"
	"time"

	"github.com/ohhfishal/gotime/entry"
)

var ErrCategoryNotFound = errors.New("category not found")

type Resume struct {
	Category string `arg: "" required: ""`
}

func (cmd *Resume) Run(entrySet EntrySet, now func() time.Time) error {
	entries, err := entrySet.GetAll()
	if err != nil {
		return fmt.Errorf(`reading entries: %w`, err)
	}

	var matches []entry.Entry
	for _, entry := range entries {
		if entry.Category == cmd.Category {
			matches = append(matches, entry)
		}
	}

	if len(matches) == 0 {
		return fmt.Errorf(`%w: "%s"`, ErrCategoryNotFound, cmd.Category)
	}

	last := matches[len(matches)-1]
	last.Note = "Cont: " + last.Note
	last.Time = now()
	return entrySet.Append(last)
}
