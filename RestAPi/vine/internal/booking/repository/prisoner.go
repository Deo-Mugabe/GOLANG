package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Deo-Mugabe/GOLANG/internal/booking/domain"
	"github.com/Deo-Mugabe/GOLANG/internal/platform/database"
)

// prisonerDB is the database model for nmmain table
type prisonerDB struct {
	NameID      int64  `db:"name_id"`
	StateID     string `db:"state_id"`
	NameType    string `db:"nametype"`
	AliasID     string `db:"alias_id"`
	FirstName   string `db:"firstname"`
	MiddleName  string `db:"middlename"`
	LastName    string `db:"lastname"`
	DOB         string `db:"dob"`
	Sex         string `db:"sex"`
	Race        string `db:"race"`
	Height      string `db:"height"`
	Weight      string `db:"weight"`
	Eye         string `db:"eye"`
	Hair        string `db:"hair"`
	SSN         string `db:"ssn"`
	DrLicense   string `db:"dr_lic"`
	DLState     string `db:"dl_state"`
	Birthplace  string `db:"birthplace"`
	StreetNbr   string `db:"streetnbr"`
	Street      string `db:"street"`
	City        string `db:"city"`
	State       string `db:"state"`
	Zip         string `db:"zip"`
	HomePhone   string `db:"hphone"`
	WorkPhone   string `db:"wphone"`
	MobilePhone string `db:"mphone"`
	Marital     string `db:"marital"`
	Occupation  string `db:"occupation"`
	Employer    string `db:"employer"`
}

type prisonerRepo struct {
	db *database.DB
}

// NewPrisonerRepository creates a new prisoner repository
func NewPrisonerRepository(db *database.DB) domain.PrisonerRepository {
	return &prisonerRepo{db: db}
}

// GetByID retrieves a prisoner by ID
func (r *prisonerRepo) GetByID(ctx context.Context, id int64) (*domain.Prisoner, error) {
	query := `
        SELECT 
            name_id, state_id, nametype, alias_id, firstname, middlename, lastname,
            dob, sex, race, height, weight, eye, hair, ssn, dr_lic, dl_state,
            birthplace, streetnbr, street, city, state, zip, hphone, wphone, 
            mphone, marital, occupation, employer
        FROM nmmain
        WHERE name_id = @p1
    `
	var dbModel prisonerDB

	err := r.db.GetContext(ctx, &dbModel, query, id)

	if err == sql.ErrNoRows {
		return nil, domain.ErrPrisonerNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get prisoner: %w", err)
	}

	return toPrisonerDomain(&dbModel), nil

}

// GetByBookingID retrieves prisoner associated with a booking
func (r *prisonerRepo) GetByBookingID(ctx context.Context, bookingID int64) (*domain.Prisoner, error) {
	query := `
        SELECT 
            n.name_id, n.state_id, n.nametype, n.alias_id, n.firstname, n.middlename, n.lastname,
            n.dob, n.sex, n.race, n.height, n.weight, n.eye, n.hair, n.ssn, n.dr_lic, n.dl_state,
            n.birthplace, n.streetnbr, n.street, n.city, n.state, n.zip, n.hphone, n.wphone, 
            n.mphone, n.marital, n.occupation, n.employer
        FROM nmmain n
        INNER JOIN jmmain j ON n.name_id = j.name_id
        WHERE j.book_id = @p1
    `
	var dbModel prisonerDB

	err := r.db.GetContext(ctx, &dbModel, query, bookingID)
	if err == sql.ErrNoRows {
		return nil, domain.ErrPrisonerNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get prisoner by booking: %w", err)
	}
	return toPrisonerDomain(&dbModel), nil
}

// GetAlias retrieves prisoner alias
func (r *prisonerRepo) GetAlias(ctx context.Context, aliasID string, nameType string) (*domain.Prisoner, error) {
	query := `
        SELECT 
            name_id, state_id, nametype, alias_id, firstname, middlename, lastname,
            dob, sex, race, height, weight, eye, hair, ssn, dr_lic, dl_state,
            birthplace, streetnbr, street, city, state, zip, hphone, wphone, 
            mphone, marital, occupation, employer
        FROM nmmain
        WHERE alias_id = @p1 AND nametype = @p2
        ORDER BY name_id ASC
    `
	var dbModel prisonerDB

	err := r.db.GetContext(ctx, &dbModel, query, aliasID, nameType)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get alias: %w", err)
	}
	return toPrisonerDomain(&dbModel), nil
}

// List retrieves prisoners with pagination
func (r *prisonerRepo) List(ctx context.Context, limit, offset int) ([]*domain.Prisoner, error) {
	query := `
        SELECT 
            name_id, state_id, nametype, alias_id, firstname, middlename, lastname,
            dob, sex, race, height, weight, eye, hair, ssn, dr_lic, dl_state,
            birthplace, streetnbr, street, city, state, zip, hphone, wphone, 
            mphone, marital, occupation, employer
        FROM nmmain
        WHERE nametype = 'M'
        ORDER BY name_id DESC
        OFFSET @p1 ROWS FETCH NEXT @p2 ROWS ONLY
    `

	var dbModels []prisonerDB
	err := r.db.SelectContext(ctx, &dbModels, query, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list prisoners: %w", err)
	}

	prisoners := make([]*domain.Prisoner, len(dbModels))

	for i, db := range dbModels {
		prisoners[i] = toPrisonerDomain(&db)
	}

	return prisoners, nil

}

// SearchByName searches prisoners by name
func (r *prisonerRepo) SearchByName(ctx context.Context, firstName, lastName string) ([]*domain.Prisoner, error) {
	query := `
        SELECT 
            name_id, state_id, nametype, alias_id, firstname, middlename, lastname,
            dob, sex, race, height, weight, eye, hair, ssn, dr_lic, dl_state,
            birthplace, streetnbr, street, city, state, zip, hphone, wphone, 
            mphone, marital, occupation, employer
        FROM nmmain
        WHERE LOWER(firstname) LIKE LOWER(@p1) 
          AND LOWER(lastname) LIKE LOWER(@p2)
        ORDER BY name_id DESC
    `

	var dbModels []prisonerDB
	err := r.db.SelectContext(ctx, &dbModels, query, "%"+firstName+"%", "%"+lastName+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to search prisoners: %w", err)
	}

	prisoners := make([]*domain.Prisoner, len(dbModels))
	for i, db := range dbModels {
		prisoners[i] = toPrisonerDomain(&db)
	}

	return prisoners, nil
}

// SearchByStateID searches prisoner by state ID
func (r *prisonerRepo) SearchByStateID(ctx context.Context, stateID string) (*domain.Prisoner, error) {
	query := `
        SELECT 
            name_id, state_id, nametype, alias_id, firstname, middlename, lastname,
            dob, sex, race, height, weight, eye, hair, ssn, dr_lic, dl_state,
            birthplace, streetnbr, street, city, state, zip, hphone, wphone, 
            mphone, marital, occupation, employer
        FROM nmmain
        WHERE state_id = @p1
    `

	var dbModel prisonerDB
	err := r.db.GetContext(ctx, &dbModel, query, stateID)
	if err == sql.ErrNoRows {
		return nil, domain.ErrPrisonerNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to search by state ID: %w", err)
	}

	return toPrisonerDomain(&dbModel), nil
}
