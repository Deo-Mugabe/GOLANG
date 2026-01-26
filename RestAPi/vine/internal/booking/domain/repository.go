package domain

import (
	"context"
	"time"
)

// BookingRepository defines booking data access
type BookingRepository interface {
	GetByID(ctx context.Context, id int64) (*Booking, error)
	GetWithDetails(ctx context.Context, id int64) (*BookingWithDetails, error)
	ListActive(ctx context.Context, id int64) ([]*Booking, error)
	Count(ctx context.Context) (int64, error)

	// Processing queries
	FetchPairsForProcessing(ctx context.Context, since time.Time) ([]BookingNamePair, error)
	GetProcessingCandidates(ctx context.Context, since time.Time, limit int) ([]*Booking, error)

	// Write operations (if needed)
	UpdateStatus(ctx context.Context, id int64, status string) error
}

// PrisonerRepository defines prisoner data access
type PrisonerRepository interface {
	// Read operations
	GetByID(ctx context.Context, id int64) (*Prisoner, error)
	GetByBookingID(ctx context.Context, bookingID int64) (*Prisoner, error)
	GetAlias(ctx context.Context, aliasID string, nameType string) (*Prisoner, error)
	List(ctx context.Context, limit, offset int) ([]*Prisoner, error)

	// Search operations
	SearchByName(ctx context.Context, firstName, lastName string) ([]*Prisoner, error)
	SearchByStateID(ctx context.Context, stateID string) (*Prisoner, error)
}

// ChargeRepository defines charge data access
type ChargeRepository interface {
	// Read operations
	GetByID(ctx context.Context, id int64) (*Charge, error)
	GetByBookingID(ctx context.Context, bookingID int64) ([]*Charge, error)
	GetByArrestID(ctx context.Context, arrestID int64) ([]*Charge, error)
	GetFirstByBookingID(ctx context.Context, bookingID int64) (*Charge, error)

	// Aggregate queries
	CountByBookingID(ctx context.Context, bookingID int64) (int, error)
}

// ArrestRepository handles arrest data access
type ArrestRepository interface {
	// Read operations
	GetByID(ctx context.Context, id int64) (*Arrest, error)
	GetByBookingID(ctx context.Context, bookingID int64) (*Arrest, error)
	GetFirstByBookingID(ctx context.Context, bookingID int64) (*Arrest, error)
	ListByBookingID(ctx context.Context, bookingID int64) ([]*Arrest, error)
}

// ReleaseRepository handles release data access
type ReleaseRepository interface {
	// Read operations
	GetByID(ctx context.Context, id int) (*Release, error)
	GetByBookingID(ctx context.Context, bookingID int64) (*Release, error)
	ListRecent(ctx context.Context, since time.Time, limit int) ([]*Release, error)
}

// FacilityRepository handles facility history data access
type FacilityRepository interface {
	// Read operations
	GetByID(ctx context.Context, id int) (*Facility, error)
	GetLatestByBookingID(ctx context.Context, bookID int64) (*Facility, error)
	ListByBookingID(ctx context.Context, bookingID int64) ([]*Facility, error)
}

// MugshotRepository handles mugshot metadata access
type MugshotRepository interface {
	// Read operations
	GetByID(ctx context.Context, id int) (*Mugshot, error)
	GetByNameID(ctx context.Context, nameID int64) ([]*Mugshot, error)
	GetLatestByNameID(ctx context.Context, nameID int64) (*Mugshot, error)
	ListByNameID(ctx context.Context, nameID int64) ([]*Mugshot, error)
}
