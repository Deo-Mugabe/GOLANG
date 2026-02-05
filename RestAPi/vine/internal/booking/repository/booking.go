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
