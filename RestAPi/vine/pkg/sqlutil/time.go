package sqlutil

import (
	"fmt"
	"time"
)

// ParseFlexibleDate parses various date formats
func ParseFlexibleDate(s string) (time.Time, error) {
	if s == "" || s == "null" || s == "N/A" {
		return time.Time{}, nil
	}

	// Common date formats to try
	formats := []string{
		"2006-01-02",          // YYYY-MM-DD
		"20060102",            // YYYYMMDD
		"01/02/2006",          // MM/DD/YYYY
		"1/2/2006",            // M/D/YYYY
		"2006-01-02 15:04:05", // YYYY-MM-DD HH:MM:SS
		"01/02/2006 15:04:05", // MM/DD/YYYY HH:MM:SS
	}

	for _, format := range formats {
		if t, err := time.Parse(format, s); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", s)
}

// FormatDateForDB formats time for database
func FormatDateForDB(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

// FormatDateYYYYMMDD formats time as YYYYMMDD
func FormatDateYYYYMMDD(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("20060102")
}
