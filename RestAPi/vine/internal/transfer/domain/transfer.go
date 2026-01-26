package domain

import "time"

// Transfer represents a file transfer operation
type Transfer struct {
	ID         string // UUID
	Type       TransferType
	FileName   string
	FileSize   int64
	RemotePath string
	Status     TransferStatus
	StartTime  time.Time
	EndTime    *time.Time
	Error      string
	RetryCount int
}

// TransferType represents type of transfer
type TransferType string

const (
	TransferTypeData    TransferType = "DATA"    // .dat file
	TransferTypeMugshot TransferType = "MUGSHOT" // image files
)

// TransferStatus represents transfer status
type TransferStatus string

const (
	TransferStatusPending   TransferStatus = "PENDING"
	TransferStatusRunning   TransferStatus = "RUNNING"
	TransferStatusCompleted TransferStatus = "COMPLETED"
	TransferStatusFailed    TransferStatus = "FAILED"
)

// IsCompleted checks if transfer completed
func (t *Transfer) IsCompleted() bool {
	return t.Status == TransferStatusCompleted
}

// CanRetry checks if transfer can be retried
func (t *Transfer) CanRetry() bool {
	return t.Status == TransferStatusFailed && t.RetryCount < 3
}
