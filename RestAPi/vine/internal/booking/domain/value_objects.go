package domain

import "time"

// BookingNamePair represents a booking-prisoner relationship
type BookingNamePair struct {
	BookingID int64
	NameID    int64
}

// ProcessingPeriod represents a time range for processing
type ProcessingPeriod struct {
	From time.Time
	To   time.Time
}

// IsValid checks if period is valid
func (p ProcessingPeriod) IsValid() bool {
	return !p.From.IsZero() && !p.To.IsZero() && p.To.After(p.From)
}

// VINEFileData represents generated VINE file content
type VINEFileData struct {
	PrisonerData string
	ChargesData  string
	MugshotData  string
}

// Combine returns complete file content
func (v *VINEFileData) Combine() string {
	return v.PrisonerData + v.ChargesData + v.MugshotData
}

// IsEmpty checks if any data exists
func (v *VINEFileData) IsEmpty() bool {
	return v.PrisonerData == "" && v.ChargesData == "" && v.MugshotData == ""
}
