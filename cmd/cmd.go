package cmd

import (
	"io"
	"os"
	"time"

	"github.com/alecthomas/kong"
	"github.com/ohhfishal/gotime/entry"
	"github.com/ohhfishal/gotime/file"
	"github.com/ohhfishal/gotime/cmd/resume"
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
	Log    LogCmd    `default:"withargs" cmd:"" help:"Manage journal entries"`
	Report ReportCmd `cmd:"" help:"Print summary report"`
  Resume resume.CMD `cmd:"" help:"Resume given category in log"`
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
	return file.ReadAllFrom(c.LogPath(), entry.Decode)
}

func (c Config) Today() (*time.Time, error) {
	now := time.Now().Format(time.DateOnly)
	today, err := time.Parse(time.DateOnly, now)
	if err != nil {
		return nil, err
	}
	return &today, nil
}

func (c Config) Now() (*time.Time, error) {
	now := time.Now().Format(time.DateTime)
	today, err := time.Parse(time.DateTime, now)
	if err != nil {
		return nil, err
	}
	return &today, nil
}

func (c Config) GetenvDefault(key, value string) string {
	if v := c.Getenv(key); v != `` {
		return v
	}
	return value
}

func (c Config) FilePerms() int {
	return os.O_APPEND | os.O_CREATE | os.O_WRONLY
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
