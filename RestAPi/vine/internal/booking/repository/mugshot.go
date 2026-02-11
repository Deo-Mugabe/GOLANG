package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Deo-Mugabe/GOLANG/internal/booking/domain"
	"github.com/Deo-Mugabe/GOLANG/internal/platform/database"
)

// mugshotDB is the database model for sys_img table
type mugshotDB struct {
	SysImgID  int64     `db:"sys_imgid"`
	SystemID  int64     `db:"sysid"`
	SystemKey string    `db:"syskey"`
	Ext1      int       `db:"ext1"`
	Ext2      int       `db:"ext2"`
	AddTime   time.Time `db:"addtime"`
}

type mugshotRepo struct {
	db *database.DB
}

// NewMugshotRepository creates a new mugshot repository
func NewMugshotRepository(db *database.DB) domain.MugshotRepository {
	return &mugshotRepo{db: db}
}

// GetByID retrieves a mugshot by ID
func (r *mugshotRepo) GetByID(ctx context.Context, id int64) (*domain.Mugshot, error) {
	query := `
        SELECT sys_imgid, sysid, syskey, ext1, ext2, addtime
        FROM sys_img
        WHERE sys_imgid = @p1
    `

	var dbModel mugshotDB
	err := r.db.GetContext(ctx, &dbModel, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get mugshot: %w", err)
	}

	return toMugshotDomain(&dbModel), nil
}

// GetByNameID retrieves all mugshots for a name ID
func (r *mugshotRepo) GetByNameID(ctx context.Context, nameID int64) ([]*domain.Mugshot, error) {
	query := `
        SELECT sys_imgid, sysid, syskey, ext1, ext2, addtime
        FROM sys_img
        WHERE syskey = 'N' AND sysid = @p1
        ORDER BY addtime DESC
    `

	var dbModels []mugshotDB
	err := r.db.SelectContext(ctx, &dbModels, query, nameID)
	if err != nil {
		return nil, fmt.Errorf("failed to get mugshots by name: %w", err)
	}

	mugshots := make([]*domain.Mugshot, len(dbModels))
	for i, db := range dbModels {
		mugshots[i] = toMugshotDomain(&db)
	}

	return mugshots, nil
}

// GetLatestByNameID retrieves the most recent mugshot
func (r *mugshotRepo) GetLatestByNameID(ctx context.Context, nameID int64) (*domain.Mugshot, error) {
	query := `
        SELECT TOP 1 sys_imgid, sysid, syskey, ext1, ext2, addtime
        FROM sys_img
        WHERE syskey = 'N' AND sysid = @p1
        ORDER BY addtime DESC
    `

	var dbModel mugshotDB
	err := r.db.GetContext(ctx, &dbModel, query, nameID)
	if err == sql.ErrNoRows {
		return nil, nil // No mugshot is valid
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get latest mugshot: %w", err)
	}

	return toMugshotDomain(&dbModel), nil
}

// ListByNameID retrieves mugshots with limit
func (r *mugshotRepo) ListByNameID(ctx context.Context, nameID int64, limit int) ([]*domain.Mugshot, error) {
	query := `
        SELECT TOP (@p2) sys_imgid, sysid, syskey, ext1, ext2, addtime
        FROM sys_img
        WHERE syskey = 'N' AND sysid = @p1
        ORDER BY addtime DESC
    `

	var dbModels []mugshotDB
	err := r.db.SelectContext(ctx, &dbModels, query, nameID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list mugshots: %w", err)
	}

	mugshots := make([]*domain.Mugshot, len(dbModels))
	for i, db := range dbModels {
		mugshots[i] = toMugshotDomain(&db)
	}

	return mugshots, nil
}

// toMugshotDomain converts database model to domain model
func toMugshotDomain(db *mugshotDB) *domain.Mugshot {
	return &domain.Mugshot{
		ID:        db.SysImgID,
		SystemID:  db.SystemID,
		SystemKey: db.SystemKey,
		Ext1:      db.Ext1,
		Ext2:      db.Ext2,
		AddTime:   db.AddTime,
	}
}
