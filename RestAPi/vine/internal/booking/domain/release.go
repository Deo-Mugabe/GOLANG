package domain

import "time"

// Release represents release information (jrelease table)
type Release struct {
	ID          int       // jreleaseid
	BookingID   int64     // book_id
	ReleaseTime time.Time // releasetime
	Reason      string    // relsreason
}

// IsReleased checks if release is valid
func (r *Release) IsReleased() bool {
	return !r.ReleaseTime.IsZero()
}
