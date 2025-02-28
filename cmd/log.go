package cmd

import (
	"errors"
	"strings"
	"time"

	"github.com/ohhfishal/gotime/entry"
)

var ErrFileEmpty = errors.New("no entries in file")

type LogCmd struct {
	New    NewCmd    `default:"withargs" cmd: "" help:"Log using a custom category and note."`
	Append AppendCmd `cmd: "" help:"Log using the last category in the log."`
	Until  UntilCmd  `cmd: "" help:"How long until a certain amount of time is logged."`
}

type NewCmd struct {
	Category string   `arg: "" required: ""`
	Note     []string `arg: "" optional: ""`
}

type AppendCmd struct {
	Note []string `arg: "" optional: ""`
}

func (cmd *NewCmd) Run(cfg Config) error {
	newEntry := entry.Entry{
		Category: cmd.Category,
		Note:     strings.Join(cmd.Note, " "),
		Time:     time.Now(),
	}
	return cfg.WriteToLog(newEntry)
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
	return cfg.WriteToLog(newEntry)
}

// TODO: REMOVE
func Log(cfg Config, args ...string) error {
	var cmd LogCmd
	return RunCmd(cfg, &cmd, args...)
}
