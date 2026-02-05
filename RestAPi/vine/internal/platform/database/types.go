package database

import (
	"database/sql"
	"time"
)

// NullString is a nullable string type
type NullString = sql.NullString

// NullInt64 is a nullable int64 type
type NullInt64 = sql.NullInt64

// NullBool is a nullable bool type
type NullBool = sql.NullBool

// NullTime is a nullable time type
type NullTime = sql.NullTime

// Helper functions for creating null types

// NewNullString creates a NullString
func NewNullString(s string, valid bool) NullString {
	return NullString{String: s, Valid: valid}
}

// ToNullString converts a string pointer to NullString
func ToNullString(s *string) NullString {
	if s == nil {
		return NullString{Valid: false}
	}
	return NullString{String: *s, Valid: true}
}

// NewNullInt64 creates a NullInt64
func NewNullInt64(i int64, valid bool) NullInt64 {
	return NullInt64{Int64: i, Valid: valid}
}

// ToNullInt64 converts an int64 pointer to NullInt64
func ToNullInt64(i *int64) NullInt64 {
	if i == nil {
		return NullInt64{Valid: false}
	}
	return NullInt64{Int64: *i, Valid: true}
}

// NewNullTime creates a NullTime
func NewNullTime(t time.Time, valid bool) NullTime {
	return NullTime{Time: t, Valid: valid}
}

// ToNullTime converts a time pointer to NullTime
func ToNullTime(t *time.Time) NullTime {
	if t == nil {
		return NullTime{Valid: false}
	}
	return NullTime{Time: *t, Valid: true}
}
