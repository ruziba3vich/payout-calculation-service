package payout

import "time"

// ParsePeriod turns "YYYY-MM" into the first day of that month in UTC.
func ParsePeriod(s string) (time.Time, error) {
	t, err := time.Parse("2006-01", s)
	if err != nil {
		return time.Time{}, ErrInvalidPeriod
	}
	return t.UTC(), nil
}

// PeriodOf returns the first day of the month that t falls in, in UTC.
func PeriodOf(t time.Time) time.Time {
	t = t.UTC()
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
}

// PeriodBounds returns the half-open [start, end) range for a period.
func PeriodBounds(period time.Time) (start, end time.Time) {
	start = PeriodOf(period)
	end = start.AddDate(0, 1, 0)
	return start, end
}
