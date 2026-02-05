package sqlutil

import (
	"database/sql"
	"time"
)

// StringPtr converts string to pointer (nil if empty)
func StringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// Int64Ptr converts int64 to pointer (nil if zero)
func Int64Ptr(i int64) *int64 {
	if i == 0 {
		return nil
	}
	return &i
}

// TimePtr converts time to pointer (nil if zero)
func TimePtr(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	return &t
}

// StringValue returns string from pointer (empty if nil)
func StringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// Int64Value returns int64 from pointer (zero if nil)
func Int64Value(i *int64) int64 {
	if i == nil {
		return 0
	}
	return *i
}

// TimeValue returns time from pointer (zero if nil)
func TimeValue(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}

// NullStringToPtr converts NullString to pointer
func NullStringToPtr(ns sql.NullString) *string {
	if !ns.Valid {
		return nil
	}
	return &ns.String
}

// NullInt64ToPtr converts NullInt64 to pointer
func NullInt64ToPtr(ni sql.NullInt64) *int64 {
	if !ni.Valid {
		return nil
	}
	return &ni.Int64
}

// NullTimeToPtr converts NullTime to pointer
func NullTimeToPtr(nt sql.NullTime) *time.Time {
	if !nt.Valid {
		return nil
	}
	return &nt.Time
}
