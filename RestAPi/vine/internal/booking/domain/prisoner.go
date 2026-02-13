package domain

import (
	"fmt"
	"strings"
	"time"
)

type Prisoner struct {
	// Identity
	ID       int64  // name_id
	StateID  string // state_id
	NameType string // nametype (M=Main, AKA=Alias)
	AliasID  string // alias_id

	// Personal Information
	FirstName  string // firstname
	MiddleName string // middlename
	LastName   string // lastname
	DOB        string // dob (stored as string in DB)
	Sex        string // sex
	Race       string // race

	// Physical Characteristics
	Height string // height
	Weight string // weight
	Eye    string // eye
	Hair   string // hair

	// Identification
	SSN        string // ssn
	DrLicense  string // dr_lic
	DLState    string // dl_state
	Birthplace string // birthplace

	// Address
	StreetNbr string // streetnbr
	Street    string // street
	City      string // city
	State     string // state
	Zip       string // zip

	// Contact
	HomePhone   string // hphone
	WorkPhone   string // wphone
	MobilePhone string // mphone

	// Other
	Marital    string // marital
	Occupation string // occupation
	Employer   string // employer
}

// NameType constants
const (
	NameTypeMain  = "M"
	NameTypeAlias = "AKA"
)

// FullName returns formatted full name
func (p *Prisoner) FullName() string {
	if p.MiddleName != "" {
		return p.FirstName + " " + p.MiddleName + " " + p.LastName
	}
	return p.FirstName + " " + p.LastName
}

// FullAddress returns formatted address
func (p *Prisoner) FullAddress() string {
	addr := ""
	if p.StreetNbr != "" {
		addr = p.StreetNbr + " "
	}
	if p.Street != "" {
		addr += p.Street
	}
	return addr
}

// IsMainRecord checks if this is the main name record
func (p *Prisoner) IsMainRecord() bool {
	return p.NameType == NameTypeMain
}

// ParseDOB attempts to parse DOB string to time.Time
func (p *Prisoner) ParseDOB() (*time.Time, error) {
	if p.DOB == "" || p.DOB == "null" || p.DOB == "N/A" {
		return &time.Time{}, nil
	}

	cleanDOB := strings.TrimSpace(p.DOB)

	// Handle empty or known null representations
	switch cleanDOB {
	case "", "null", "NULL", "N/A", "n/a":
		return nil, nil
	case "00/00/0000", "0000-00-00":
		return nil, nil
	}

	// Try different date formats
	formats := []string{
		"2006-01-02", // YYYY-MM-DD
		"20060102",   // YYYYMMDD
		"01/02/2006", // MM/DD/YYYY
		"1/2/2006",   // M/D/YYYY
		"2006",       // Year only
	}
	// Attempt parsing with each format
	for _, layout := range formats {
		parsed, err := time.Parse(layout, cleanDOB)
		if err == nil {
			// Normalize to UTC midnight (prevents timezone drift issues)
			normalized := parsed.UTC()
			return &normalized, nil
		}
	}
	// If we reach here, DOB was provided but invalid
	return nil, fmt.Errorf("invalid DOB format: %q", cleanDOB)
}
