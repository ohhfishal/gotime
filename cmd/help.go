package cmd

import (
	"fmt"
)

func Help(cfg Config, args ...string) error {
	fmt.Fprintln(cfg.Stdout, "HELP")
	return nil
}
