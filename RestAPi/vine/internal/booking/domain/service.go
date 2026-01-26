package domain

import (
	"context"
	"time"
)

// ProcessorService handles booking processing business logic
type ProcessorService interface {
	// Main processing
	ProcessBookings(ctx context.Context, since time.Time) (ProcessingResult, error)
	ProcessSingleBooking(ctx context.Context, bookingID int64) error

	// Processing workflow
	FetchBookingsForProcessing(ctx context.Context, since time.Time) ([]BookingNamePair, error)
	GenerateVINEData(ctx context.Context, pairs []BookingNamePair) (*VINEFileData, error)
	ValidateBookingData(ctx context.Context, bookingID, nameID int64) error
}

// FileGeneratorService handles VINE file generation
type FileGeneratorService interface {
	// File generation
	GeneratePrisonerRecord(ctx context.Context, bookingID, nameID int64) (string, error)
	GenerateChargesRecords(ctx context.Context, bookingID, nameID int64) (string, error)
	GenerateMugshotRecords(ctx context.Context, bookingID, nameID int64) (string, error)

	// File operations
	WriteVINEFile(ctx context.Context, data *VINEFileData, filepath string) error
	ValidateVINEFile(ctx context.Context, filepath string) error
}

// MugshotService handles mugshot processing
type MugshotService interface {
	// Mugshot operations
	CopyMugshot(ctx context.Context, nameID, bookingID int64) error
	ClearOutputDirectory(ctx context.Context) error
	GetMugshotPath(ctx context.Context, nameID int64) (string, error)
	ListMugshotsForBooking(ctx context.Context, bookingID int64) ([]string, error)
}

// ProcessingResult contains processing results
type ProcessingResult struct {
	TotalProcessed   int64
	SuccessCount     int64
	FailureCount     int64
	StartTime        time.Time
	EndTime          time.Time
	ProcessedIDs     []int64
	FailedIDs        []int64
	GeneratedFile    string
	TransferredFiles []string
}
