package domain

import (
	"context"
	"time"
)

// BookingRepository defines booking data access
type BookingRepository interface {
	GetByID(ctx context.Context, id int64) (*Booking, error)
	GetWithDetails(ctx context.Context, id int64) (*BookingWithDetails, error)
	FetchPairsForProcessing(ctx context.Context, since time.Time) ([]BookingNamePair, error)
}

// PrisonerRepository defines prisoner data access
type PrisonerRepository interface {
	GetByID(ctx context.Context, id int64) (*Prisoner, error)
	GetAlias(ctx context.Context, aliasID string) (*Prisoner, error)
}

// ChargeRepository defines charge data access
type ChargeRepository interface {
	GetByBookingID(ctx context.Context, bookingID int64) ([]*Charge, error)
	GetFirstByBookingID(ctx context.Context, bookingID int64) (*Charge, error)
}

// ... similar for Arrest, Release, Facility, Mugshot repositories
