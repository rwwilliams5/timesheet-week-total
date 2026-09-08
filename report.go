package main

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"time"
)

// weekStart returns the Monday of the week containing t, at midnight UTC.
// Dates in the log are plain calendar days with no time zone attached, so
// everything here stays in UTC to avoid DST arithmetic drift.
func weekStart(t time.Time) time.Time {
	t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	wd := int(t.Weekday())
	if wd == 0 {
		wd = 7 // treat Sunday as day 7, not day 0, so the week starts Monday
	}
	return t.AddDate(0, 0, -(wd - 1))
}

type DayTotal struct {
	Date  string  `json:"date"`
	Hours float64 `json:"hours"`
}

type ProjectTotal struct {
	Project string  `json:"project"`
	Hours   float64 `json:"hours"`
}

type Report struct {
	WeekStart       string         `json:"week_start"`
	WeekEnd         string         `json:"week_end"`
	Days            []DayTotal     `json:"days"`
	Projects        []ProjectTotal `json:"projects"`
	TotalHours      float64        `json:"total_hours"`
	TargetHours     float64        `json:"target_hours"`
	DifferenceHours float64        `json:"difference_hours"`
}

// buildReport sums entries falling in the 7-day window starting at start.
// All 7 days are always present in the output, even with zero hours, so a
// missed day is visible rather than silently absent.
func buildReport(entries []Entry, start time.Time, target float64) Report {
	end := start.AddDate(0, 0, 6)

	dayTotals := make(map[string]float64, 7)
	dayOrder := make([]string, 7)
	for i := 0; i < 7; i++ {
		d := start.AddDate(0, 0, i).Format(dateLayout)
		dayOrder[i] = d
		dayTotals[d] = 0
	}

	projectTotals := map[string]float64{}
	var total float64

	for _, e := range entries {
		if e.Date.Before(start) || e.Date.After(end) {
			continue
		}
		key := e.Date.Format(dateLayout)
		dayTotals[key] += e.Hours
		projectTotals[e.Project] += e.Hours
		total += e.Hours
	}

	days := make([]DayTotal, 0, 7)
	for _, d := range dayOrder {
		days = append(days, DayTotal{Date: d, Hours: dayTotals[d]})
	}

	projects := make([]ProjectTotal, 0, len(projectTotals))
	for name, hours := range projectTotals {
		projects = append(projects, ProjectTotal{Project: name, Hours: hours})
	}
	sort.Slice(projects, func(i, j int) bool {
		if projects[i].Hours != projects[j].Hours {
			return projects[i].Hours > projects[j].Hours
		}
		return projects[i].Project < projects[j].Project
	})

	return Report{
		WeekStart:       start.Format(dateLayout),
		WeekEnd:         end.Format(dateLayout),
		Days:            days,
		Projects:        projects,
		TotalHours:      total,
		TargetHours:     target,
		DifferenceHours: total - target,
	}
}

func (r Report) writeJSON(w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}

func (r Report) writeHuman(w io.Writer) {
	fmt.Fprintf(w, "week %s to %s\n\n", r.WeekStart, r.WeekEnd)

	for _, d := range r.Days {
		fmt.Fprintf(w, "  %s  %5.2fh\n", d.Date, d.Hours)
	}
	fmt.Fprintln(w)

	if len(r.Projects) == 0 {
		fmt.Fprintln(w, "no entries this week")
	} else {
		fmt.Fprintln(w, "by project:")
		for _, p := range r.Projects {
			fmt.Fprintf(w, "  %-20s %5.2fh\n", p.Project, p.Hours)
		}
	}
	fmt.Fprintln(w)

	fmt.Fprintf(w, "total:  %.2fh\n", r.TotalHours)
	fmt.Fprintf(w, "target: %.2fh\n", r.TargetHours)
	if r.DifferenceHours >= 0 {
		fmt.Fprintf(w, "%.2fh over target\n", r.DifferenceHours)
	} else {
		fmt.Fprintf(w, "%.2fh under target\n", -r.DifferenceHours)
	}
}
