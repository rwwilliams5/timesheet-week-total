package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const dateLayout = "2006-01-02"

// Entry is a single logged block of time.
type Entry struct {
	Date    time.Time
	Project string
	Hours   float64
	Note    string
}

// loadEntries reads a CSV timesheet log. The file must have a header row
// with at least "date", "project" and "hours" columns (any order); a
// "note" column is optional. Column names are matched case-insensitively.
func loadEntries(path string) ([]Entry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening timesheet: %w", err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.FieldsPerRecord = -1 // rows may be shorter when trailing columns are blank

	rows, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("reading csv: %w", err)
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("%s: empty file", path)
	}

	col := make(map[string]int, len(rows[0]))
	for i, name := range rows[0] {
		col[strings.ToLower(strings.TrimSpace(name))] = i
	}
	for _, want := range []string{"date", "project", "hours"} {
		if _, ok := col[want]; !ok {
			return nil, fmt.Errorf("%s: missing required column %q", path, want)
		}
	}
	noteIdx, hasNote := col["note"]

	entries := make([]Entry, 0, len(rows)-1)
	for i, row := range rows[1:] {
		lineNum := i + 2 // account for the header and 1-indexing

		date, err := time.Parse(dateLayout, field(row, col["date"]))
		if err != nil {
			return nil, fmt.Errorf("%s: line %d: invalid date: %w", path, lineNum, err)
		}
		hours, err := strconv.ParseFloat(field(row, col["hours"]), 64)
		if err != nil {
			return nil, fmt.Errorf("%s: line %d: invalid hours: %w", path, lineNum, err)
		}
		project := field(row, col["project"])
		if project == "" {
			return nil, fmt.Errorf("%s: line %d: project is empty", path, lineNum)
		}
		note := ""
		if hasNote {
			note = field(row, noteIdx)
		}

		entries = append(entries, Entry{Date: date, Project: project, Hours: hours, Note: note})
	}
	return entries, nil
}

// field returns the trimmed value at idx, or "" if the row is too short to
// have it. Rows are allowed to omit trailing optional columns.
func field(row []string, idx int) string {
	if idx < 0 || idx >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[idx])
}
