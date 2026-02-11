package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Deo-Mugabe/GOLANG/internal/booking/domain"
	"github.com/Deo-Mugabe/GOLANG/internal/platform/database"
)

// releaseDB is the database model for jrelease table
type releaseDB struct {
	ReleaseID   int       `db:"jreleaseid"`
	BookingID   int64     `db:"book_id"`
	ReleaseTime time.Time `db:"releasetime"`
	Reason      string    `db:"relsreason"`
}

type releaseRepo struct {
	db *database.DB
}

// NewReleaseRepository creates a new release repository
func NewReleaseRepository(db *database.DB) domain.ReleaseRepository {
	return &releaseRepo{db: db}
}

// GetByID retrieves a release by ID
func (r *releaseRepo) GetByID(ctx context.Context, id int) (*domain.Release, error) {
	query := `
        SELECT jreleaseid, book_id, releasetime, relsreason
        FROM jrelease
        WHERE jreleaseid = @p1
    `

	var dbModel releaseDB
	err := r.db.GetContext(ctx, &dbModel, query, id)
	if err == sql.ErrNoRows {
		return nil, nil // Release not found is not an error
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get release: %w", err)
	}

	return toReleaseDomain(&dbModel), nil
}

// GetByBookingID retrieves release by booking ID
func (r *releaseRepo) GetByBookingID(ctx context.Context, bookingID int64) (*domain.Release, error) {
	query := `
        SELECT jreleaseid, book_id, releasetime, relsreason
        FROM jrelease
        WHERE book_id = @p1
    `

	var dbModel releaseDB
	err := r.db.GetContext(ctx, &dbModel, query, bookingID)
	if err == sql.ErrNoRows {
		return nil, nil // No release record is valid for active bookings
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get release by booking: %w", err)
	}

	return toReleaseDomain(&dbModel), nil
}

// ListRecent retrieves recent releases
func (r *releaseRepo) ListRecent(ctx context.Context, since time.Time, limit int) ([]*domain.Release, error) {
	query := `
        SELECT TOP (@p2) jreleaseid, book_id, releasetime, relsreason
        FROM jrelease
        WHERE releasetime >= @p1
        ORDER BY releasetime DESC
    `

	var dbModels []releaseDB
	err := r.db.SelectContext(ctx, &dbModels, query, since, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list recent releases: %w", err)
	}

	releases := make([]*domain.Release, len(dbModels))
	for i, db := range dbModels {
		releases[i] = toReleaseDomain(&db)
	}

	return releases, nil
}

// toReleaseDomain converts database model to domain model
func toReleaseDomain(db *releaseDB) *domain.Release {
	return &domain.Release{
		ID:          db.ReleaseID,
		BookingID:   db.BookingID,
		ReleaseTime: db.ReleaseTime,
		Reason:      db.Reason,
	}
}
