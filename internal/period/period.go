// Package period computes date ranges for week/month/quarter/year,
// optionally anchored to a specific period identifier.
package period

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

type Range struct {
	Name      string `json:"name"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

const dateFmt = "2006-01-02"

// Compute returns the date range for the given period kind.
// If specific is empty, the current period (relative to now) is used.
// kind must be one of: week, month, quarter, year.
func Compute(kind, specific string) (Range, error) {
	if specific != "" {
		return parseSpecific(kind, specific)
	}
	return current(kind, time.Now())
}

func current(kind string, today time.Time) (Range, error) {
	today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	switch kind {
	case "week":
		weekday := int(today.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		start := today.AddDate(0, 0, -(weekday - 1))
		end := start.AddDate(0, 0, 6)
		year, week := start.ISOWeek()
		return Range{
			Name:      fmt.Sprintf("%04d-W%02d", year, week),
			StartDate: start.Format(dateFmt),
			EndDate:   end.Format(dateFmt),
		}, nil
	case "month":
		start := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location())
		end := start.AddDate(0, 1, -1)
		return Range{
			Name:      start.Format("2006-01"),
			StartDate: start.Format(dateFmt),
			EndDate:   end.Format(dateFmt),
		}, nil
	case "quarter":
		q := (int(today.Month()) - 1) / 3
		start := time.Date(today.Year(), time.Month(q*3+1), 1, 0, 0, 0, 0, today.Location())
		end := start.AddDate(0, 3, -1)
		return Range{
			Name:      fmt.Sprintf("%d-Q%d", start.Year(), q+1),
			StartDate: start.Format(dateFmt),
			EndDate:   end.Format(dateFmt),
		}, nil
	case "year":
		start := time.Date(today.Year(), 1, 1, 0, 0, 0, 0, today.Location())
		end := time.Date(today.Year(), 12, 31, 0, 0, 0, 0, today.Location())
		return Range{
			Name:      fmt.Sprintf("%d", today.Year()),
			StartDate: start.Format(dateFmt),
			EndDate:   end.Format(dateFmt),
		}, nil
	}
	return Range{}, fmt.Errorf("unknown period: %q (want week|month|quarter|year)", kind)
}

var (
	reWeek    = regexp.MustCompile(`^(\d{4})-W(\d{1,2})$`)
	reMonth   = regexp.MustCompile(`^(\d{4})-(\d{1,2})$`)
	reQuarter = regexp.MustCompile(`^(\d{4})-Q([1-4])$`)
	reYear    = regexp.MustCompile(`^(\d{4})$`)
)

func parseSpecific(kind, value string) (Range, error) {
	switch kind {
	case "week":
		m := reWeek.FindStringSubmatch(value)
		if m == nil {
			return Range{}, fmt.Errorf("week format: YYYY-Www (e.g. 2026-W15), got %q", value)
		}
		year, _ := strconv.Atoi(m[1])
		week, _ := strconv.Atoi(m[2])
		start := isoWeekStart(year, week)
		end := start.AddDate(0, 0, 6)
		return Range{
			Name:      fmt.Sprintf("%04d-W%02d", year, week),
			StartDate: start.Format(dateFmt),
			EndDate:   end.Format(dateFmt),
		}, nil
	case "month":
		m := reMonth.FindStringSubmatch(value)
		if m == nil {
			return Range{}, fmt.Errorf("month format: YYYY-MM (e.g. 2026-03), got %q", value)
		}
		year, _ := strconv.Atoi(m[1])
		month, _ := strconv.Atoi(m[2])
		if month < 1 || month > 12 {
			return Range{}, fmt.Errorf("invalid month: %d", month)
		}
		start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
		end := start.AddDate(0, 1, -1)
		return Range{
			Name:      fmt.Sprintf("%04d-%02d", year, month),
			StartDate: start.Format(dateFmt),
			EndDate:   end.Format(dateFmt),
		}, nil
	case "quarter":
		m := reQuarter.FindStringSubmatch(value)
		if m == nil {
			return Range{}, fmt.Errorf("quarter format: YYYY-Qn (e.g. 2026-Q1), got %q", value)
		}
		year, _ := strconv.Atoi(m[1])
		q, _ := strconv.Atoi(m[2])
		start := time.Date(year, time.Month((q-1)*3+1), 1, 0, 0, 0, 0, time.UTC)
		end := start.AddDate(0, 3, -1)
		return Range{
			Name:      fmt.Sprintf("%d-Q%d", year, q),
			StartDate: start.Format(dateFmt),
			EndDate:   end.Format(dateFmt),
		}, nil
	case "year":
		m := reYear.FindStringSubmatch(value)
		if m == nil {
			return Range{}, fmt.Errorf("year format: YYYY (e.g. 2026), got %q", value)
		}
		year, _ := strconv.Atoi(m[1])
		return Range{
			Name:      fmt.Sprintf("%d", year),
			StartDate: fmt.Sprintf("%04d-01-01", year),
			EndDate:   fmt.Sprintf("%04d-12-31", year),
		}, nil
	}
	return Range{}, fmt.Errorf("unknown period: %q", kind)
}

// isoWeekStart returns the Monday of the given ISO week.
func isoWeekStart(year, week int) time.Time {
	jan4 := time.Date(year, 1, 4, 0, 0, 0, 0, time.UTC)
	weekday := int(jan4.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	week1Monday := jan4.AddDate(0, 0, -(weekday - 1))
	return week1Monday.AddDate(0, 0, (week-1)*7)
}
