package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Deo-Mugabe/GOLANG/internal/booking/domain"
	"github.com/Deo-Mugabe/GOLANG/internal/platform/database"
)

// facilityDB is the database model for jfachist table
type facilityDB struct {
	FacHistID int64     `db:"jfachistid"`
	BookingID int64     `db:"book_id"`
	Facility  string    `db:"facility"`
	Section   string    `db:"section"`
	Unit      string    `db:"unit"`
	Bed       string    `db:"bed"`
	EventDate time.Time `db:"eventdate"`
}

type facilityRepo struct {
	db *database.DB
}

// NewFacilityRepository creates a new facility repository
func NewFacilityRepository(db *database.DB) domain.FacilityRepository {
	return &facilityRepo{db: db}
}

// GetByID retrieves a facility record by ID
func (r *facilityRepo) GetByID(ctx context.Context, id int64) (*domain.Facility, error) {
	query := `
        SELECT jfachistid, book_id, facility, section, unit, bed, eventdate
        FROM jfachist
        WHERE jfachistid = @p1
    `

	var dbModel facilityDB
	err := r.db.GetContext(ctx, &dbModel, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get facility: %w", err)
	}

	return toFacilityDomain(&dbModel), nil
}

// GetLatestByBookingID retrieves the most recent facility assignment
func (r *facilityRepo) GetLatestByBookingID(ctx context.Context, bookingID int64) (*domain.Facility, error) {
	query := `
        SELECT TOP 1 jfachistid, book_id, facility, section, unit, bed, eventdate
        FROM jfachist
        WHERE book_id = @p1
        ORDER BY eventdate DESC
    `

	var dbModel facilityDB
	err := r.db.GetContext(ctx, &dbModel, query, bookingID)
	if err == sql.ErrNoRows {
		return nil, nil // No facility assignment is valid
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get latest facility: %w", err)
	}

	return toFacilityDomain(&dbModel), nil
}

// ListByBookingID retrieves all facility history for a booking
func (r *facilityRepo) ListByBookingID(ctx context.Context, bookingID int64) ([]*domain.Facility, error) {
	query := `
        SELECT jfachistid, book_id, facility, section, unit, bed, eventdate
        FROM jfachist
        WHERE book_id = @p1
        ORDER BY eventdate DESC
    `

	var dbModels []facilityDB
	err := r.db.SelectContext(ctx, &dbModels, query, bookingID)
	if err != nil {
		return nil, fmt.Errorf("failed to list facilities: %w", err)
	}

	facilities := make([]*domain.Facility, len(dbModels))
	for i, db := range dbModels {
		facilities[i] = toFacilityDomain(&db)
	}

	return facilities, nil
}

// toFacilityDomain converts database model to domain model
func toFacilityDomain(db *facilityDB) *domain.Facility {
	return &domain.Facility{
		ID:        db.FacHistID,
		BookingID: db.BookingID,
		Facility:  db.Facility,
		Section:   db.Section,
		Unit:      db.Unit,
		Bed:       db.Bed,
		EventDate: db.EventDate,
	}
}
