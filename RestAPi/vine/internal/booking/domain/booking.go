package domain

import (
	"time"
)

// Booking represents a jail booking (jmmain table)
type Booking struct {
	// Core fields
	ID         int64     // book_id
	NameID     int64     // name_id
	BookDate   time.Time // bookdate
	Agency     string    // agency
	AddTime    time.Time // addtime
	Status     string    // bkstatus (A=Active, R=Released)
	FacilityID *int64    // faci_id (nullable)

	// Metadata
	CreatedAt time.Time
	UpdatedAt time.Time
}

// BookingStatus represents valid booking statuses
type BookingStatus string

const (
	BookingStatusActive   BookingStatus = "A"
	BookingStatusReleased BookingStatus = "R"
)

// IsActive checks if booking is currently active
func (b *Booking) IsActive() bool {
	return b.Status == string(BookingStatusActive)
}

// CanProcess checks if booking can be processed for VINE
func (b *Booking) CanProcess() bool {
	return b.ID > 0 && b.NameID > 0 && b.BookDate.Year() > 1900
}

// BookingWithDetails aggregates booking with related data
type BookingWithDetails struct {
	Booking  *Booking
	Prisoner *Prisoner
	Arrest   *Arrest
	Charges  []*Charge
	Release  *Release
	Facility *Facility
	Mugshots []*Mugshot
}
