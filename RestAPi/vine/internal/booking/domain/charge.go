package domain

// Charge represents a criminal charge (archrg table)
type Charge struct {
	ID        int64 // archrgid
	BookingID int64 // book_id
	ArrestID  int64 // armainid

	// Charge Details
	ChargeCode  string // arr_chrg
	Description string // chrgdesc
	FelonyMisd  string // fel_misd (F=Felony, M=Misdemeanor)
	Count       string // chrg_cnt
	Sequence    string // chrg_seq

	// Bond Information
	BondAmount string // bondamt
	BondType   string // bondtype
}

// ChargeType constants
const (
	ChargeTypeFelony      = "F"
	ChargeTypeMisdemeanor = "M"
)

// IsFelony checks if charge is a felony
func (c *Charge) IsFelony() bool {
	return c.FelonyMisd == ChargeTypeFelony
}

// IsMisdemeanor checks if charge is a misdemeanor
func (c *Charge) IsMisdemeanor() bool {
	return c.FelonyMisd == ChargeTypeMisdemeanor
}
