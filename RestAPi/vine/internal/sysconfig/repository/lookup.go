package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Deo-Mugabe/GOLANG/internal/platform/database"
	"github.com/Deo-Mugabe/GOLANG/internal/sysconfig/domain"
)

// lookupDB is the database model for systab1 table
type lookupDB struct {
	ID         int    `db:"systab1id"`
	CodeAgency string `db:"codeAgcy"`
	CodeKey    string `db:"code_key"`
	Message    string `db:"sys_msg"`
}

type lookupRepo struct {
	db *database.DB
}

// NewLookupRepository creates a new lookup repository
func NewLookupRepository(db *database.DB) domain.LookupRepository {
	return &lookupRepo{db: db}
}

// GetByID retrieves a lookup by ID
func (r *lookupRepo) GetByID(ctx context.Context, id int) (*domain.SystemLookup, error) {
	query := `
        SELECT systab1id, codeAgcy, code_key, sys_msg
        FROM systab1
        WHERE systab1id = @p1
    `

	var dbModel lookupDB
	err := r.db.GetContext(ctx, &dbModel, query, id)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("lookup not found: %d", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get lookup: %w", err)
	}

	return toLookupDomain(&dbModel), nil
}

// GetByKey retrieves a lookup by agency and key
func (r *lookupRepo) GetByKey(ctx context.Context, codeAgency, codeKey string) (*domain.SystemLookup, error) {
	query := `
        SELECT systab1id, codeAgcy, code_key, sys_msg
        FROM systab1
        WHERE codeAgcy = @p1 AND code_key = @p2
    `

	var dbModel lookupDB
	err := r.db.GetContext(ctx, &dbModel, query, codeAgency, codeKey)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("lookup not found: %s:%s", codeAgency, codeKey)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get lookup by key: %w", err)
	}

	return toLookupDomain(&dbModel), nil
}

// GetByAgency retrieves all lookups for an agency
func (r *lookupRepo) GetByAgency(ctx context.Context, codeAgency string) ([]*domain.SystemLookup, error) {
	query := `
        SELECT systab1id, codeAgcy, code_key, sys_msg
        FROM systab1
        WHERE codeAgcy = @p1
        ORDER BY code_key
    `

	var dbModels []lookupDB
	err := r.db.SelectContext(ctx, &dbModels, query, codeAgency)
	if err != nil {
		return nil, fmt.Errorf("failed to get lookups by agency: %w", err)
	}

	lookups := make([]*domain.SystemLookup, len(dbModels))
	for i, db := range dbModels {
		lookups[i] = toLookupDomain(&db)
	}

	return lookups, nil
}

// List retrieves all lookups with pagination
func (r *lookupRepo) List(ctx context.Context, limit, offset int) ([]*domain.SystemLookup, error) {
	query := `
        SELECT systab1id, codeAgcy, code_key, sys_msg
        FROM systab1
        ORDER BY codeAgcy, code_key
        OFFSET @p1 ROWS FETCH NEXT @p2 ROWS ONLY
    `

	var dbModels []lookupDB
	err := r.db.SelectContext(ctx, &dbModels, query, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list lookups: %w", err)
	}

	lookups := make([]*domain.SystemLookup, len(dbModels))
	for i, db := range dbModels {
		lookups[i] = toLookupDomain(&db)
	}

	return lookups, nil
}

// Create creates a new lookup
func (r *lookupRepo) Create(ctx context.Context, lookup *domain.SystemLookup) error {
	query := `
        INSERT INTO systab1 (codeAgcy, code_key, sys_msg)
        OUTPUT INSERTED.systab1id
        VALUES (@p1, @p2, @p3)
    `

	var id int
	err := r.db.GetContext(ctx, &id, query, lookup.CodeAgency, lookup.CodeKey, lookup.Message)
	if err != nil {
		return fmt.Errorf("failed to create lookup: %w", err)
	}

	lookup.ID = id
	return nil
}

// Update updates a lookup
func (r *lookupRepo) Update(ctx context.Context, lookup *domain.SystemLookup) error {
	query := `
        UPDATE systab1
        SET codeAgcy = @p1, code_key = @p2, sys_msg = @p3
        WHERE systab1id = @p4
    `

	result, err := r.db.ExecContext(ctx, query,
		lookup.CodeAgency, lookup.CodeKey, lookup.Message, lookup.ID)
	if err != nil {
		return fmt.Errorf("failed to update lookup: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("lookup not found: %d", lookup.ID)
	}

	return nil
}

// Delete deletes a lookup
func (r *lookupRepo) Delete(ctx context.Context, id int) error {
	result, err := r.db.ExecContext(ctx,
		"DELETE FROM systab1 WHERE systab1id = @p1", id)
	if err != nil {
		return fmt.Errorf("failed to delete lookup: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("lookup not found: %d", id)
	}

	return nil
}

// toLookupDomain converts database model to domain model
func toLookupDomain(db *lookupDB) *domain.SystemLookup {
	return &domain.SystemLookup{
		ID:         db.ID,
		CodeAgency: db.CodeAgency,
		CodeKey:    db.CodeKey,
		Message:    db.Message,
	}
}
