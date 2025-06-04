package cmd

import (
	"io"
	"time"

	"github.com/alecthomas/kong"
)

// TODO: Make GetAllEntries and OpenWriter/Reader testable
type Config struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

type RootCmd struct {
	Log     LogCmd    `default:"withargs" cmd:"" help:"Log using a custom category and note."`
	Report  ReportCmd `cmd:"" help:"Print summary report"`
	LogFile string    `default:"~/.config/gotime.csv" help:"Path to log file." env:"GOTIME_LOG"`
}

func Run(args []string, config Config) error {
	cmd := &RootCmd{}

	var exit bool
	parser, err := kong.New(
		cmd,
		kong.Exit(func(_ int) { exit = true }),
		kong.Bind(config),
		kong.Bind(time.Now),
		kong.BindTo(config.Stdout, (*io.Writer)(nil)),
	)
	if err != nil {
		return err
	}
	parser.Stdout = config.Stdout
	parser.Stderr = config.Stderr

	context, err := parser.Parse(args)
	if err != nil || exit {
		return err
	}

	err = context.Run(cmd.LogFile)
	if err != nil {
		return err
	}
	return nil
}
