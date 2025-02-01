package cmd

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"slices"
	"time"

	"github.com/ohhfishal/gotime/entry"
)

var ErrReportUse = WrapInvalidUse("report [--start DAY] [--end DAY]")

type reportFlags struct {
	Start   time.Time
	End     time.Time
	back    int
	forward int
}

func parseReportFlags(stdout io.Writer, args []string) (reportFlags, error) {
	flags := reportFlags{}
	flagSet := flag.NewFlagSet("report", flag.ContinueOnError)
	flagSet.IntVar(&flags.back, "start", 0, "days back to include (inclusive)")
	flagSet.IntVar(&flags.forward, "end", 0, "days forward to include (inclusive)")
	// TODO: Flag for time format Kitchen, TimeOnly, Decimal

	flagSet.SetOutput(stdout)
	flagSet.Usage = func() {}
	err := flagSet.Parse(args)
	if err != nil {
		return reportFlags{}, err
	}
	// TODO: Calculate the time shifts
	if flags.back != 0 || flags.forward != 0 {
		fmt.Fprintln(stdout, "time flags not implemented")
		return reportFlags{}, errors.New("time flags not implemented")
	}
	flags.End = time.Now()
	return flags, nil

}

func Report(cfg Config, args ...string) error {
	flags, err := parseReportFlags(cfg.Stdout, args)
	if errors.Is(err, flag.ErrHelp) {
		return Help(cfg, "report")
	} else if err != nil {
		return ErrReportUse
	}

	path := cfg.GetenvDefault("GOTIME_LOG", "gotime.log")
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf(`can not open "%s": %w`, path, err)
	}

	// TODO: Add in config to customize the file format read
	entries, err := entry.ReadAll(file)
	if err != nil {
		return fmt.Errorf(`parsing file "%s": %w`, path, err)
	}
	filtered := entry.Filter(entries, flags.Start, flags.End)
	slices.SortFunc(filtered, entry.Compare)
	// TODO Make the output good and based on config
	return report(cfg, filtered)
}

func report(cfg Config, entries []entry.Entry) error {
	// TODO: Use `text/template` to make this expandable
	for i, entry := range entries {
		// TODO: Make this better
		fmt.Fprintf(cfg.Stdout, "%s", entry)
		if i < (len(entries) - 1) {
			diff := entries[i+1].Time.Sub(entry.Time)
			fmt.Fprintf(cfg.Stdout, " %v", diff)
		}

		fmt.Fprintln(cfg.Stdout)
	}
	return nil

}
