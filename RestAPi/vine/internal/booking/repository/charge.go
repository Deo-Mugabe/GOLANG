package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Deo-Mugabe/GOLANG/internal/booking/domain"
	"github.com/Deo-Mugabe/GOLANG/internal/platform/database"
)

// chargeDB is the database model for archrg table
type chargeDB struct {
	ChargeID    int64  `db:"archrgid"`
	BookingID   int64  `db:"book_id"`
	ArrestID    int64  `db:"armainid"`
	ChargeCode  string `db:"arr_chrg"`
	Description string `db:"chrgdesc"`
	FelonyMisd  string `db:"fel_misd"`
	Count       string `db:"chrg_cnt"`
	Sequence    string `db:"chrg_seq"`
	BondAmount  string `db:"bondamt"`
	BondType    string `db:"bondtype"`
}

type chargeRepo struct {
	db *database.DB
}

// NewChargeRepository creates a new charge repository
func NewChargeRepository(db *database.DB) chargeDB {
	return &chargeRepo{db: db}
}

// GetByID retrieves a charge by ID
func (r *chargeRepo) GetByID(ctx context.Context, id int64) (*domain.Charge, error) {
	query := `
        SELECT archrgid, book_id, armainid, arr_chrg, chrgdesc, fel_misd, 
               chrg_cnt, chrg_seq, bondamt, bondtype
        FROM archrg
        WHERE archrgid = @p1
    `

	var dbModel chargeDB
	err := r.db.GetContext(ctx, &dbModel, query, id)
	if err == sql.ErrNoRows {
		return nil, domain.ErrChargeNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get charge: %w", err)
	}

	return toChargeDomain(&dbModel), nil
}

// GetByBookingID retrieves all charges for a booking
func (r *chargeRepo) GetByBookingID(ctx context.Context, bookingID int64) ([]*domain.Charge, error) {
	query := `
        SELECT archrgid, book_id, armainid, arr_chrg, chrgdesc, fel_misd, 
               chrg_cnt, chrg_seq, bondamt, bondtype
        FROM archrg
        WHERE book_id = @p1
        ORDER BY armainid ASC, chrg_seq ASC
    `

	var dbModels []chargeDB
	err := r.db.SelectContext(ctx, &dbModels, query, bookingID)
	if err != nil {
		return nil, fmt.Errorf("failed to get charges by booking: %w", err)
	}

	charges := make([]*domain.Charge, len(dbModels))
	for i, db := range dbModels {
		charges[i] = toChargeDomain(&db)
	}

	return charges, nil
}

// GetByArrestID retrieves charges for a specific arrest
func (r *chargeRepo) GetByArrestID(ctx context.Context, arrestID int64) ([]*domain.Charge, error) {
	query := `
        SELECT archrgid, book_id, armainid, arr_chrg, chrgdesc, fel_misd, 
               chrg_cnt, chrg_seq, bondamt, bondtype
        FROM archrg
        WHERE armainid = @p1
        ORDER BY chrg_seq ASC
    `

	var dbModels []chargeDB
	err := r.db.SelectContext(ctx, &dbModels, query, arrestID)
	if err != nil {
		return nil, fmt.Errorf("failed to get charges by arrest: %w", err)
	}

	charges := make([]*domain.Charge, len(dbModels))
	for i, db := range dbModels {
		charges[i] = toChargeDomain(&db)
	}

	return charges, nil
}

// GetFirstByBookingID retrieves the first charge for a booking
func (r *chargeRepo) GetFirstByBookingID(ctx context.Context, bookingID int64) (*domain.Charge, error) {
	query := `
        SELECT TOP 1 
            archrgid, book_id, armainid, arr_chrg, chrgdesc, fel_misd, 
            chrg_cnt, chrg_seq, bondamt, bondtype
        FROM archrg
        WHERE book_id = @p1
        ORDER BY armainid ASC, chrg_seq ASC
    `

	var dbModel chargeDB
	err := r.db.GetContext(ctx, &dbModel, query, bookingID)
	if err == sql.ErrNoRows {
		return nil, domain.ErrChargeNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get first charge: %w", err)
	}

	return toChargeDomain(&dbModel), nil
}

// CountByBookingID counts charges for a booking
func (r *chargeRepo) CountByBookingID(ctx context.Context, bookingID int64) (int, error) {
	var count int
	err := r.db.GetContext(ctx, &count,
		"SELECT COUNT(*) FROM archrg WHERE book_id = @p1", bookingID)
	return count, err
}

// toChargeDomain converts database model to domain model
func toChargeDomain(db *chargeDB) *domain.Charge {
	return &domain.Charge{
		ID:          db.ChargeID,
		BookingID:   db.BookingID,
		ArrestID:    db.ArrestID,
		ChargeCode:  db.ChargeCode,
		Description: db.Description,
		FelonyMisd:  db.FelonyMisd,
		Count:       db.Count,
		Sequence:    db.Sequence,
		BondAmount:  db.BondAmount,
		BondType:    db.BondType,
	}
}
