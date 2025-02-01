package main

import (
	"fmt"
	"github.com/ohhfishal/gotime/cmd"
	"os"
)

func main() {
	if err := cmd.Run(cmd.Config{}, os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
