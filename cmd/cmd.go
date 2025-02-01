package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
)

var ErrInvalidUse = errors.New(`usage: gotime`)

func WrapInvalidUse(msg string) error {
	return fmt.Errorf("%w %s\nTry `gotime help [command]` for more information.", ErrInvalidUse, msg)
}

type Config struct {
	Stdin  io.Reader
	Stdout io.Writer
	Getenv func(string) string
}

func (c Config) GetenvDefault(key, value string) string {
	if v := c.Getenv(key); v != `` {
		return v
	}
	return value
}

type Command func(Config, ...string) error

var commands = map[string]Command{
	"report": Report,
	"log":    Log,
	"help":   Help,
}

func Run(cfg Config, args []string) error {
	cfg = SanitizeConfig(cfg)

	if len(args) == 0 {
		return Help(cfg, args...)
	}

	if command, ok := commands[args[0]]; ok {
		return command(cfg, args[1:]...)
	}
	// TODO: Parse some flags to see if they are trying to do something else?
	return Log(cfg, args...)
}

func SanitizeConfig(cfg Config) Config {
	if cfg.Stdin == nil {
		cfg.Stdin = os.Stdin
	}

	if cfg.Stdout == nil {
		cfg.Stdout = os.Stdout
	}

	if cfg.Getenv == nil {
		cfg.Getenv = os.Getenv
	}
	return cfg
}
