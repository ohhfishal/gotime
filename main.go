package main

import (
	"fmt"
	"os"

	"github.com/ohhfishal/gotime/cmd"
)

func main() {
	cfg := cmd.Config{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	if err := cmd.Run(os.Args[1:], cfg); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
