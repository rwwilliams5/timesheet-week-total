package main

import (
	"testing"
	"time"
)

func date(s string) time.Time {
	t, err := time.Parse(dateLayout, s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestWeekStart(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"2026-09-07", "2026-09-07"}, // Monday itself
		{"2026-09-08", "2026-09-07"}, // Tuesday
		{"2026-09-11", "2026-09-07"}, // Friday
		{"2026-09-13", "2026-09-07"}, // Sunday, end of the same week
		{"2026-09-14", "2026-09-14"}, // Monday of the next week
	}
	for _, c := range cases {
		got := weekStart(date(c.in)).Format(dateLayout)
		if got != c.want {
			t.Errorf("weekStart(%s) = %s, want %s", c.in, got, c.want)
		}
	}
}

func TestBuildReportTotals(t *testing.T) {
	entries := []Entry{
		{Date: date("2026-09-07"), Project: "acme-website", Hours: 3.5},
		{Date: date("2026-09-07"), Project: "acme-website", Hours: 1.0},
		{Date: date("2026-09-08"), Project: "internal-tools", Hours: 4.0},
		{Date: date("2026-09-13"), Project: "internal-tools", Hours: 2.0},
		{Date: date("2026-09-14"), Project: "acme-website", Hours: 100}, // next week, must be excluded
		{Date: date("2026-09-06"), Project: "acme-website", Hours: 100}, // prior week, must be excluded
	}

	r := buildReport(entries, weekStart(date("2026-09-09")), 40)

	if r.WeekStart != "2026-09-07" || r.WeekEnd != "2026-09-13" {
		t.Errorf("week bounds = %s..%s", r.WeekStart, r.WeekEnd)
	}
	if len(r.Days) != 7 {
		t.Fatalf("got %d days, want 7", len(r.Days))
	}
	if r.Days[0].Hours != 4.5 {
		t.Errorf("day 0 hours = %v, want 4.5", r.Days[0].Hours)
	}
	if r.Days[6].Hours != 2.0 {
		t.Errorf("day 6 hours = %v, want 2.0", r.Days[6].Hours)
	}
	// a day with no entries still shows up with zero hours
	if r.Days[2].Date != "2026-09-09" || r.Days[2].Hours != 0 {
		t.Errorf("day 2 = %+v, want 2026-09-09 with 0 hours", r.Days[2])
	}

	if r.TotalHours != 10.5 {
		t.Errorf("total hours = %v, want 10.5", r.TotalHours)
	}
	if r.TargetHours != 40 {
		t.Errorf("target hours = %v, want 40", r.TargetHours)
	}
	if r.DifferenceHours != 10.5-40 {
		t.Errorf("difference hours = %v, want %v", r.DifferenceHours, 10.5-40)
	}

	if len(r.Projects) != 2 {
		t.Fatalf("got %d projects, want 2", len(r.Projects))
	}
	if r.Projects[0].Project != "acme-website" || r.Projects[0].Hours != 4.5 {
		t.Errorf("projects[0] = %+v", r.Projects[0])
	}
	if r.Projects[1].Project != "internal-tools" || r.Projects[1].Hours != 6.0 {
		t.Errorf("projects[1] = %+v", r.Projects[1])
	}
}

func TestBuildReportProjectSortTieBreak(t *testing.T) {
	entries := []Entry{
		{Date: date("2026-09-07"), Project: "zeta", Hours: 2},
		{Date: date("2026-09-07"), Project: "alpha", Hours: 2},
	}

	r := buildReport(entries, weekStart(date("2026-09-09")), 40)

	if len(r.Projects) != 2 {
		t.Fatalf("got %d projects, want 2", len(r.Projects))
	}
	// equal hours, so alphabetical order breaks the tie
	if r.Projects[0].Project != "alpha" || r.Projects[1].Project != "zeta" {
		t.Errorf("projects = %+v, want alpha before zeta", r.Projects)
	}
}

func TestBuildReportNoEntries(t *testing.T) {
	r := buildReport(nil, weekStart(date("2026-09-09")), 40)

	if r.TotalHours != 0 {
		t.Errorf("total hours = %v, want 0", r.TotalHours)
	}
	if len(r.Projects) != 0 {
		t.Errorf("got %d projects, want 0", len(r.Projects))
	}
	if r.DifferenceHours != -40 {
		t.Errorf("difference hours = %v, want -40", r.DifferenceHours)
	}
}
