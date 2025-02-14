package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/alecthomas/kong"
	"github.com/ohhfishal/gotime/entry"
)

const DefaultLogPath = "$HOME/.config/gotime.log"

// TODO: Make GetAllEntries and OpenWriter/Reader testable
type Config struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
	Getenv func(string) string
}

type RootCmd struct {
	Log    LogCmd    `cmd:"" help:"Manage journal entries"`
	Report ReportCmd `cmd:"" help:"Print summary report"`
}

func Run(args []string, config ...Config) error {
	cfg := DefaultConfig(config...)
	var cmd RootCmd
	return RunCmd(cfg, &cmd, args...)
}

func RunCmd(cfg Config, cmd any, args ...string) error {
	var exit bool
	parser, err := kong.New(
		cmd,
		kong.Exit(func(_ int) { exit = true }),
		kong.Bind(cfg),
	)
	if err != nil {
		return err
	}
	parser.Stdout = cfg.Stdout
	parser.Stderr = cfg.Stderr

	context, err := parser.Parse(args)
	if err != nil || exit {
		return err
	}

	err = context.Run()
	if err != nil {
		return err
	}
	return nil
}

func (c Config) LogPath() string {
	unexpanded := c.GetenvDefault("GOTIME_LOG", DefaultLogPath)
	return os.Expand(unexpanded, c.Getenv)
}

func (c Config) GetAllEntries() ([]entry.Entry, error) {
	var zero []entry.Entry
	path := c.LogPath()
	file, err := os.Open(path)
	if err != nil {
		return zero, fmt.Errorf(`can not open "%s": %w`, path, err)
	}
	defer file.Close()

	entries, err := entry.ReadAll(file)
	if err != nil {
		return zero, fmt.Errorf(`parsing file "%s": %w`, path, err)
	}
	return entries, nil
}

func (c Config) GetenvDefault(key, value string) string {
	if v := c.Getenv(key); v != `` {
		return v
	}
	return value
}

func (c Config) OpenWriteCloser() (io.WriteCloser, error) {
	path := c.LogPath()

	perms := os.O_APPEND | os.O_WRONLY
	if c.CanCreateFile() {
		perms = os.O_APPEND | os.O_CREATE | os.O_WRONLY
	}

	file, err := os.OpenFile(path, perms, 0644)
	if err != nil {
		return nil, fmt.Errorf(`can not open "%s": %w`, path, err)
	}
	return file, nil
}

func (c Config) CanCreateFile() bool {
	// TODO: Allow control of this via an env
	return true
}

func DefaultConfig(config ...Config) Config {
	if len(config) > 0 {
		return sanitizeConfig(config[0])
	}
	return sanitizeConfig(Config{})
}

func sanitizeConfig(cfg Config) Config {
	if cfg.Stdin == nil {
		cfg.Stdin = os.Stdin
	}

	if cfg.Stdout == nil {
		cfg.Stdout = os.Stdout
	}

	if cfg.Stderr == nil {
		cfg.Stderr = os.Stderr
	}

	if cfg.Getenv == nil {
		cfg.Getenv = os.Getenv
	}
	return cfg
}
