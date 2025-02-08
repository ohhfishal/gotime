package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/ohhfishal/gotime/entry"
)

type Command func(Config, ...string) error

var commands = map[string]Command{
	"report": Report,
	"log":    Log,
	"help":   Help,
}

const DefaultLogPath = "gotime.log"

var ErrInvalidUse = errors.New(`usage: gotime`)

func WrapInvalidUse(msg string) error {
	return fmt.Errorf("%w %s\nTry `gotime help [command]` for more information.", ErrInvalidUse, msg)
}

// TODO: Make GetAllEntries and OpenWriter/Reader testable
type Config struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
	Getenv func(string) string
	// TODO: Private struct that unmarshals from ENVs
}

func Run(args []string, config ...Config) error {
	cfg := DefaultConfig(config...)
	if len(args) == 0 {
		return Help(cfg, args...)
	}

	if command, ok := commands[args[0]]; ok {
		return command(cfg, args[1:]...)
	}
	// TODO: Parse some flags to see if they are trying to do something else?
	return Log(cfg, args...)
}

func (c Config) LogPath() string {
	return c.GetenvDefault("GOTIME_LOG", DefaultLogPath)
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
