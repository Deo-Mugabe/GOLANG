package domain

import (
	"context"
	"time"
)

// TransferRepository handles transfer logging (optional - in-memory or DB)
type TransferRepository interface {
	// Read operations
	GetByID(ctx context.Context, id string) (*Transfer, error)
	ListRecent(ctx context.Context, limit int) ([]*Transfer, error)
	ListByStatus(ctx context.Context, status TransferStatus, limit int) ([]*Transfer, error)
	ListByDateRange(ctx context.Context, from, to time.Time) ([]*Transfer, error)

	// Write operations
	Create(ctx context.Context, transfer *Transfer) error
	Update(ctx context.Context, transfer *Transfer) error
	UpdateStatus(ctx context.Context, id string, status TransferStatus) error

	// Statistics
	CountByStatus(ctx context.Context, status TransferStatus) (int64, error)
	GetTotalBytesTransferred(ctx context.Context, since time.Time) (int64, error)
}
