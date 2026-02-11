package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Deo-Mugabe/GOLANG/internal/booking/domain"
	"github.com/Deo-Mugabe/GOLANG/internal/platform/database"
)

// bookingDB is the database model for jmmain table
type bookingDB struct {
	BookID   int64         `db:"book_id"`
	NameID   int64         `db:"name_id"`
	BookDate time.Time     `db:"bookdate"`
	Agency   string        `db:"agency"`
	AddTime  time.Time     `db:"addtime"`
	Status   string        `db:"bkstatus"`
	FacilID  sql.NullInt64 `db:"faci_id"`
}

type bookingRepo struct {
	db *database.DB
}

// NewBookingRepository creates a new booking repository
func NewBookingRepository(db *database.DB) domain.BookingRepository {
	return &bookingRepo{db: db}
}

// GetByID retrieves a booking by ID
func (r *bookingRepo) GetByID(ctx context.Context, id int64) (*domain.Booking, error) {
	query := `
        SELECT book_id, name_id, bookdate, agency, addtime, bkstatus, faci_id
        FROM jmmain
        WHERE book_id = @p1
    `

	var dbModel bookingDB
	err := r.db.GetContext(ctx, &dbModel, query, id)
	if err == sql.ErrNoRows {
		return nil, domain.ErrBookingNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get booking: %w", err)
	}

	return toBookingDomain(&dbModel), nil
}

// GetWithDetails retrieves booking with all related data
func (r *bookingRepo) GetWithDetails(ctx context.Context, id int64) (*domain.BookingWithDetails, error) {
	booking, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Note: Related data would be fetched by other repositories
	// This is just the structure
	return &domain.BookingWithDetails{
		Booking: booking,
		// Other fields populated by service layer
	}, nil
}

func (r *bookingRepo) ListActive(ctx context.Context, limit, offset int) ([]*domain.Booking, error) {
	query := `
		SELECT book_id, name_id, bookdate, agency, addtime, bkstatus, faci_id
        FROM jmmain
        WHERE bkstatus = 'A'
        ORDER BY book_id DESC
        OFFSET @p1 ROWS FETCH NEXT @p2 ROWS ONLY
		`
	var dbModels []bookingDB

	err := r.db.SelectContext(ctx, &dbModels, query, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("Failed to List active bookings: w", err)
	}

	bookings := make([]*domain.Booking, len(dbModels))
	for i, db := range dbModels {
		bookings[i] = toBookingDomain(&db)
	}
	return bookings, nil
}

// Count returns total booking count
func (r *bookingRepo) Count(ctx context.Context) (int64, error) {
	query := `
		select count(*) from jmain
		`
	var count int64
	err := r.db.GetContext(ctx, &count, query)
	return count, err
}

// FetchPairsForProcessing retrieves booking-name pairs for processing
func (r *bookingRepo) FetchPairsForProcessing(ctx context.Context, since time.Time) ([]domain.BookingNamePair, error) {
	query := `
        SELECT DISTINCT 
            jmmain.book_id AS book_id, 
            jmmain.name_id AS name_id
        FROM jmmain
        LEFT OUTER JOIN jrelease ON jmmain.book_id = jrelease.book_id
        WHERE jmmain.bkstatus = 'A'
           OR jrelease.releasetime >= @p1
        ORDER BY jmmain.book_id
    `
	var pairs []domain.BookingNamePair

	err := r.db.SelectContext(ctx, &pairs, query, since)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch booking pairs: %w", err)
	}

	return pairs, nil
}

// GetProcessingCandidates retrieves bookings that need processing
func (r *bookingRepo) GetProcessingCandidates(ctx context.Context, since time.Time, limit int) ([]*domain.Booking, error) {
	query := `
        SELECT TOP (@p2)
            book_id, name_id, bookdate, agency, addtime, bkstatus, faci_id
        FROM jmmain
        WHERE bkstatus = 'A' OR addtime >= @p1
        ORDER BY addtime DESC
    `
	var dbModels []bookingDB

	err := r.db.SelectContext(ctx, &dbModels, query, since, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get processing candidates: %w", err)
	}
	bookings := make([]*domain.Booking, len(dbModels))
	for i, db := range dbModels {
		bookings[i] = toBookingDomain(&db)
	}
	return bookings, nil
}

// UpdateStatus updates booking status
func (r *bookingRepo) UpdateStatus(ctx context.Context, id int64, status string) error {
	query := `
        UPDATE jmmain 
        SET bkstatus = @p1 
        WHERE book_id = @p2
    `
	result, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update booking status: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrBookingNotFound
	}

	return nil

}
