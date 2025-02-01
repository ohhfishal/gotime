package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/ohhfishal/gotime/entry"
)

var ErrLogUse = WrapInvalidUse("log category [note...]")

func Log(cfg Config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("Missing category\n%w", ErrLogUse)
	}

	newEntry := entry.Entry{
		Category: args[0],
		Note:     strings.Join(args[1:], " "),
		Time:     time.Now(),
	}

	writer, err := OpenWriter(cfg)
	if err != nil {
		return err
	}

	// TODO: Pass in an entry.Config to allow different file encodings EX timefmt
	if err := entry.Write(writer, newEntry); err != nil {
		return fmt.Errorf(`failed writing: %w`, err)
	}
	return nil
}

func OpenWriter(cfg Config) (io.Writer, error) {
	if cfg.Getenv("GOTIME_NOFILE") == "stdout" {
		return cfg.Stdout, nil
	}

	path := cfg.GetenvDefault("GOTIME_LOG", "gotime.log")
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf(`can not open "%s": %w`, path, err)
	}
	return file, nil
}
