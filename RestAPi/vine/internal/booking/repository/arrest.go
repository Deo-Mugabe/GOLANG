package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Deo-Mugabe/GOLANG/internal/booking/domain"
	"github.com/Deo-Mugabe/GOLANG/internal/platform/database"
)

// arrestDB is the database model for armain table
type arrestDB struct {
	ArrestID   int64     `db:"armainid"`
	BookingID  int64     `db:"book_id"`
	CaseID     string    `db:"case_id"`
	ArrestDate time.Time `db:"date_arr"`
}

type arrestRepo struct {
	db *database.DB
}

// NewArrestRepository creates a new arrest repository
func NewArrestRepository(db *database.DB) domain.ArrestRepository {
	return &arrestRepo{db: db}
}

// GetByID retrieves an arrest by ID
func (r *arrestRepo) GetByID(ctx context.Context, id int64) (*domain.Arrest, error) {
	query := `
        SELECT armainid, book_id, case_id, date_arr
        FROM armain
        WHERE armainid = @p1
    `

	var dbModel arrestDB
	err := r.db.GetContext(ctx, &dbModel, query, id)
	if err == sql.ErrNoRows {
		return nil, domain.ErrArrestNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get arrest: %w", err)
	}

	return toArrestDomain(&dbModel), nil
}

// GetByBookingID retrieves arrest by booking ID
func (r *arrestRepo) GetByBookingID(ctx context.Context, bookingID int64) (*domain.Arrest, error) {
	query := `
        SELECT armainid, book_id, case_id, date_arr
        FROM armain
        WHERE book_id = @p1
    `

	var dbModel arrestDB
	err := r.db.GetContext(ctx, &dbModel, query, bookingID)
	if err == sql.ErrNoRows {
		return nil, domain.ErrArrestNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get arrest by booking: %w", err)
	}

	return toArrestDomain(&dbModel), nil
}

// GetFirstByBookingID retrieves the first arrest for a booking
func (r *arrestRepo) GetFirstByBookingID(ctx context.Context, bookingID int64) (*domain.Arrest, error) {
	query := `
        SELECT TOP 1 armainid, book_id, case_id, date_arr
        FROM armain
        WHERE book_id = @p1
        ORDER BY armainid ASC
    `

	var dbModel arrestDB
	err := r.db.GetContext(ctx, &dbModel, query, bookingID)
	if err == sql.ErrNoRows {
		return nil, domain.ErrArrestNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get first arrest: %w", err)
	}

	return toArrestDomain(&dbModel), nil
}

// ListByBookingID retrieves all arrests for a booking
func (r *arrestRepo) ListByBookingID(ctx context.Context, bookingID int64) ([]*domain.Arrest, error) {
	query := `
        SELECT armainid, book_id, case_id, date_arr
        FROM armain
        WHERE book_id = @p1
        ORDER BY armainid ASC
    `

	var dbModels []arrestDB
	err := r.db.SelectContext(ctx, &dbModels, query, bookingID)
	if err != nil {
		return nil, fmt.Errorf("failed to list arrests: %w", err)
	}

	arrests := make([]*domain.Arrest, len(dbModels))
	for i, db := range dbModels {
		arrests[i] = toArrestDomain(&db)
	}

	return arrests, nil
}

// toArrestDomain converts database model to domain model
func toArrestDomain(db *arrestDB) *domain.Arrest {
	return &domain.Arrest{
		ID:         db.ArrestID,
		BookingID:  db.BookingID,
		CaseID:     db.CaseID,
		ArrestDate: db.ArrestDate,
	}
}
