package resume

import (
	"time"

	"github.com/ohhfishal/gotime/entry"
)

type CMD struct {
	Category string `arg: "" required: ""`
}

type EntrySet interface {
	LastOf(string) (entry.Entry, error)
	Append(entry.Entry) error
}

type Clock interface {
	Now() time.Time
}

func (cmd *CMD) Run(entries EntrySet, clock Clock) error {
	entry, err := entries.LastOf(cmd.Category)
	if err != nil {
		return err
	}
	entry.Note = "Cont: " + entry.Note
	entry.Time = clock.Now()
	return entries.Append(entry)
}
