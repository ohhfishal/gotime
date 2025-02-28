package cmd

import (
	"io"
	"os"
	"time"

	"github.com/alecthomas/kong"
	"github.com/ohhfishal/gotime/entry"
	"github.com/ohhfishal/gotime/file"
)

const DefaultLogPath = "$HOME/.config/gotime.log"

// TODO: Make GetAllEntries and OpenWriter/Reader testable
type Config struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
	Getenv func(string) string
}

type EntrySet interface {
	GetAll() ([]entry.Entry, error)
	Append(entry.Entry) error
}

type RootCmd struct {
	Log    LogCmd    `default:"withargs" cmd: "" help:"Log using a custom category and note."`
	Append AppendCmd `cmd: "" help:"Log using the last category in the log."`
	Resume Resume    `cmd:"" aliases:"continue,cont" help:"Log an entry that continues the last entry of <category>."`
	Report ReportCmd `cmd:"" help:"Print summary report"`
	Until  UntilCmd  `cmd: "" help:"How long until a certain amount of time is logged."`
}

func Run(args []string, config ...Config) error {
	cfg := DefaultConfig(config...)
	cmd := &RootCmd{}

	var exit bool
	parser, err := kong.New(
		cmd,
		kong.Exit(func(_ int) { exit = true }),
		kong.Bind(cfg),
		kong.Bind(time.Now),
		kong.BindTo(cfg.Stdout, (*io.Writer)(nil)),
		kong.BindTo(cfg, (*EntrySet)(nil)),
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

func (c Config) GetAll() ([]entry.Entry, error) {
	return file.ReadAllFrom(c.LogPath(), entry.Decode)
}

func (c Config) Append(newEntry entry.Entry) error {
	return file.WriteTo(c.LogPath(), c.FilePerms(), newEntry)
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
