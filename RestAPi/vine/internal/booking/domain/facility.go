package domain

import "time"

// Facility represents facility/housing information (jfachist table)
type Facility struct {
	ID        int64     // jfachistid
	BookingID int64     // book_id
	Facility  string    // facility
	Section   string    // section
	Unit      string    // unit
	Bed       string    // bed
	EventDate time.Time // eventdate
}

// Location returns formatted location string
func (f *Facility) Location() string {
	return f.Facility + "/" + f.Section + "/" + f.Unit + "/" + f.Bed
}
