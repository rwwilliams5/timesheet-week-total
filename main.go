// Command twt answers one question: how many hours did I log in a given
// week, broken down by day and by project, and how does that compare to a
// weekly target?
package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	filePath := flag.String("file", "", "path to timesheet CSV file (required)")
	weekArg := flag.String("week", "", "a date (YYYY-MM-DD) inside the week to report on; defaults to today")
	target := flag.Float64("target", 40, "weekly target hours, used to compute over/under")
	jsonOut := flag.Bool("json", false, "print the report as JSON instead of a human-readable table")
	flag.Parse()

	if *filePath == "" {
		fmt.Fprintln(os.Stderr, "twt: -file is required")
		flag.Usage()
		os.Exit(2)
	}

	ref := time.Now()
	if *weekArg != "" {
		t, err := time.Parse(dateLayout, *weekArg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "twt: invalid -week date %q: %v\n", *weekArg, err)
			os.Exit(2)
		}
		ref = t
	}

	entries, err := loadEntries(*filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "twt: %v\n", err)
		os.Exit(1)
	}

	report := buildReport(entries, weekStart(ref), *target)

	if *jsonOut {
		if err := report.writeJSON(os.Stdout); err != nil {
			fmt.Fprintf(os.Stderr, "twt: %v\n", err)
			os.Exit(1)
		}
		return
	}
	report.writeHuman(os.Stdout)
}
