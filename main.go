package main

import (
	"fmt"
	"os"

	"github.com/ohhfishal/gotime/cmd"
)

func main() {
	if err := cmd.Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
