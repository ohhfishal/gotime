package resume

import (
	"fmt"
	"time"

	"github.com/ohhfishal/gotime/entry"
)

type CMD struct {
	Category string `arg: "" required: ""`
}

type EntrySet interface {
	GetAll() ([]entry.Entry, error)
	Append(entry.Entry) error
}

func (cmd *CMD) Run(entrySet EntrySet, now func() time.Time) error {
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
		return fmt.Errorf(`no entry matches category: "%s"`, cmd.Category)
	}

	last := matches[len(matches)-1]
	last.Note = "Cont: " + last.Note
	last.Time = now()
	return entrySet.Append(last)
}
