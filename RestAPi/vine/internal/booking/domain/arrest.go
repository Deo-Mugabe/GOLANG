package domain

import "time"

// Arrest represents arrest information (armain table)
type Arrest struct {
	ID         int64     // armainid
	BookingID  int64     // book_id
	CaseID     string    // case_id
	ArrestDate time.Time // date_arr
}

// IsValid checks if arrest has required data
func (a *Arrest) IsValid() bool {
	return a.ID > 0 && a.BookingID > 0 && !a.ArrestDate.IsZero()
}
