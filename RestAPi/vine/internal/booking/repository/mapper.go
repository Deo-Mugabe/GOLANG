package repository

import (
	"database/sql"

	"github.com/Deo-Mugabe/GOLANG/internal/booking/domain"
)

// toBookingDomain converts database model to domain model
func toBookingDomain(db *bookingDB) *domain.Booking {
	var facilityID *int64
	if db.FacilID.Valid {
		facilityID = &db.FacilID.Int64
	}

	return &domain.Booking{
		ID:         db.BookID,
		NameID:     db.NameID,
		BookDate:   db.BookDate,
		Agency:     db.Agency,
		AddTime:    db.AddTime,
		Status:     db.Status,
		FacilityID: facilityID,
		CreatedAt:  db.AddTime,
		UpdatedAt:  db.AddTime,
	}
}

// toBookingDB converts domain model to database model
func toBookingDB(d *domain.Booking) *bookingDB {
	db := &bookingDB{
		BookID:   d.ID,
		NameID:   d.NameID,
		BookDate: d.BookDate,
		Agency:   d.Agency,
		AddTime:  d.AddTime,
		Status:   d.Status,
	}

	if d.FacilityID != nil {
		db.FacilID = sql.NullInt64{Int64: *d.FacilityID, Valid: true}
	}

	return db
}

func toPrisonerDomain(db *prisonerDB) *domain.Prisoner {
	return &domain.Prisoner{
		ID:          db.NameID,
		StateID:     db.StateID,
		NameType:    db.NameType,
		AliasID:     db.AliasID,
		FirstName:   db.FirstName,
		MiddleName:  db.MiddleName,
		LastName:    db.LastName,
		DOB:         db.DOB,
		Sex:         db.Sex,
		Race:        db.Race,
		Height:      db.Height,
		Weight:      db.Weight,
		Eye:         db.Eye,
		Hair:        db.Hair,
		SSN:         db.SSN,
		DrLicense:   db.DrLicense,
		DLState:     db.DLState,
		Birthplace:  db.Birthplace,
		StreetNbr:   db.StreetNbr,
		Street:      db.Street,
		City:        db.City,
		State:       db.State,
		Zip:         db.Zip,
		HomePhone:   db.HomePhone,
		WorkPhone:   db.WorkPhone,
		MobilePhone: db.MobilePhone,
		Marital:     db.Marital,
		Occupation:  db.Occupation,
		Employer:    db.Employer,
	}
}
