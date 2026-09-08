# timesheet-week-total

I log time to a plain CSV file — one line per block of work, project and
hours and a note. What I actually want to know, most days, is a single
thing: how many hours have I logged this week, split out by day and by
project, and how far that is from 40 hours. Spreadsheets can answer this
but re-deriving the pivot every time is annoying. `twt` just answers that
one question, from the command line.

## timesheet format

A CSV file with a header row. Required columns: `date` (YYYY-MM-DD),
`project`, `hours`. An optional `note` column is ignored by the tool today
but is a reasonable place to keep context for yourself. Column order and
case don't matter.

```csv
date,project,hours,note
2026-09-07,acme-website,3.5,homepage redesign
2026-09-07,acme-website,1.0,bugfix triage
2026-09-08,internal-tools,4.0,timesheet cli prototype
2026-09-08,acme-website,2.0,client call
2026-09-09,internal-tools,6.5,report formatting
2026-09-10,acme-website,8.0,launch prep
2026-09-11,internal-tools,3.0,code review
```

See `testdata/sample.csv` for this same file.

## usage

```
go build -o twt .
./twt -file testdata/sample.csv -week 2026-09-09
```

`-week` takes any date that falls inside the week you want (the week is
always Monday through Sunday); it defaults to today. `-target` overrides
the 40 hour default weekly target.

Human-readable output:

```
week 2026-09-07 to 2026-09-13

  2026-09-07   4.50h
  2026-09-08   6.00h
  2026-09-09   6.50h
  2026-09-10   8.00h
  2026-09-11   3.00h
  2026-09-12   0.00h
  2026-09-13   0.00h

by project:
  acme-website          14.50h
  internal-tools         13.50h

total:  28.00h
target: 40.00h
12.00h under target
```

Same query, machine-readable:

```
./twt -file testdata/sample.csv -week 2026-09-09 --json
```

```json
{
  "week_start": "2026-09-07",
  "week_end": "2026-09-13",
  "days": [
    { "date": "2026-09-07", "hours": 4.5 },
    { "date": "2026-09-08", "hours": 6 },
    { "date": "2026-09-09", "hours": 6.5 },
    { "date": "2026-09-10", "hours": 8 },
    { "date": "2026-09-11", "hours": 3 },
    { "date": "2026-09-12", "hours": 0 },
    { "date": "2026-09-13", "hours": 0 }
  ],
  "projects": [
    { "project": "acme-website", "hours": 14.5 },
    { "project": "internal-tools", "hours": 13.5 }
  ],
  "total_hours": 28,
  "target_hours": 40,
  "difference_hours": -12
}
```

The JSON mode is the point of the tool as much as the CSV parsing is — it's
meant to be piped into `jq` or read by some other script, not just eyeballed.

## status

First pass. No dependencies, standard library only.
