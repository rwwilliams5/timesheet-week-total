package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "timesheet.csv")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writing temp csv: %v", err)
	}
	return path
}

func TestLoadEntries(t *testing.T) {
	path := writeTemp(t, `date,project,hours,note
2026-09-07,acme-website,3.5,homepage redesign
2026-09-08,internal-tools,4,
`)

	entries, err := loadEntries(path)
	if err != nil {
		t.Fatalf("loadEntries: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("got %d entries, want 2", len(entries))
	}

	want := Entry{
		Date:    time.Date(2026, 9, 7, 0, 0, 0, 0, time.UTC),
		Project: "acme-website",
		Hours:   3.5,
		Note:    "homepage redesign",
	}
	if got := entries[0]; !got.Date.Equal(want.Date) || got.Project != want.Project || got.Hours != want.Hours || got.Note != want.Note {
		t.Errorf("entries[0] = %+v, want %+v", got, want)
	}
	if entries[1].Note != "" {
		t.Errorf("entries[1].Note = %q, want empty", entries[1].Note)
	}
}

func TestLoadEntriesHeaderCaseAndOrder(t *testing.T) {
	path := writeTemp(t, `Hours,Project,Date
2,acme-website,2026-09-07
`)

	entries, err := loadEntries(path)
	if err != nil {
		t.Fatalf("loadEntries: %v", err)
	}
	if len(entries) != 1 || entries[0].Hours != 2 || entries[0].Project != "acme-website" {
		t.Errorf("entries = %+v", entries)
	}
}

func TestLoadEntriesMissingFile(t *testing.T) {
	if _, err := loadEntries(filepath.Join(t.TempDir(), "does-not-exist.csv")); err == nil {
		t.Fatal("expected an error for a missing file, got nil")
	}
}

func TestLoadEntriesEmptyFile(t *testing.T) {
	path := writeTemp(t, "")
	if _, err := loadEntries(path); err == nil {
		t.Fatal("expected an error for an empty file, got nil")
	}
}

func TestLoadEntriesMissingColumn(t *testing.T) {
	path := writeTemp(t, `date,project
2026-09-07,acme-website
`)
	if _, err := loadEntries(path); err == nil {
		t.Fatal("expected an error for a missing hours column, got nil")
	}
}

func TestLoadEntriesInvalidDate(t *testing.T) {
	path := writeTemp(t, `date,project,hours
09/07/2026,acme-website,3
`)
	if _, err := loadEntries(path); err == nil {
		t.Fatal("expected an error for an invalid date, got nil")
	}
}

func TestLoadEntriesInvalidHours(t *testing.T) {
	path := writeTemp(t, `date,project,hours
2026-09-07,acme-website,not-a-number
`)
	if _, err := loadEntries(path); err == nil {
		t.Fatal("expected an error for invalid hours, got nil")
	}
}

func TestLoadEntriesEmptyProject(t *testing.T) {
	path := writeTemp(t, `date,project,hours
2026-09-07,,3
`)
	if _, err := loadEntries(path); err == nil {
		t.Fatal("expected an error for an empty project, got nil")
	}
}
