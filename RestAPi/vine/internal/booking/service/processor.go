package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Deo-Mugabe/GOLANG/internal/booking/domain"
	"github.com/rs/zerolog"
)

// ProcessorService handles booking processing business logic
type ProcessorService struct {
	bookingRepo  domain.BookingRepository
	prisonerRepo domain.PrisonerRepository
	chargeRepo   domain.ChargeRepository
	arrestRepo   domain.ArrestRepository
	releaseRepo  domain.ReleaseRepository
	facilityRepo domain.FacilityRepository
	mugshotRepo  domain.MugshotRepository
	fileGen      *domain.FileGeneratorService
	mugshotSvc   *domain.MugshotService
	logger       zerolog.Logger
}

// NewProcessorService creates a new processor service
func NewProcessorService(
	bookingRepo domain.BookingRepository,
	prisonerRepo domain.PrisonerRepository,
	chargeRepo domain.ChargeRepository,
	arrestRepo domain.ArrestRepository,
	releaseRepo domain.ReleaseRepository,
	facilityRepo domain.FacilityRepository,
	mugshotRepo domain.MugshotRepository,
	fileGen *domain.FileGeneratorService,
	mugshotSvc *domain.MugshotService,
	logger zerolog.Logger,
) *ProcessorService {
	return &ProcessorService{
		bookingRepo:  bookingRepo,
		prisonerRepo: prisonerRepo,
		chargeRepo:   chargeRepo,
		arrestRepo:   arrestRepo,
		releaseRepo:  releaseRepo,
		facilityRepo: facilityRepo,
		mugshotRepo:  mugshotRepo,
		fileGen:      fileGen,
		mugshotSvc:   mugshotSvc,
		logger:       logger,
	}
}

// ProcessBookings processes bookings and generates VINE files
func (s *ProcessorService) ProcessBookings(ctx context.Context, since time.Time) (domain.ProcessingResult, error) {
	startTime := time.Now()

	s.logger.Info().
		Time("since", since).
		Msg("starting booking processing")

	result := domain.ProcessingResult{
		StartTime:    startTime,
		ProcessedIDs: make([]int64, 0),
		FailedIDs:    make([]int64, 0),
	}

	// Clear mugshot output directory
	if err := s.mugshotSvc.ClearOutputDirectory(ctx); err != nil {
		return result, fmt.Errorf("failed to clear mugshot directory: %w", err)
	}

	// Fetch booking-name pairs for processing
	pairs, err := s.FetchBookingsForProcessing(ctx, since)
	if err != nil {
		return result, fmt.Errorf("failed to fetch bookings: %w", err)
	}

	s.logger.Info().Int("count", len(pairs)).Msg("fetched booking pairs")

	// Generate VINE data for all pairs
	vineData, err := s.GenerateVINEData(ctx, pairs)
	if err != nil {
		return result, fmt.Errorf("failed to generate VINE data: %w", err)
	}

	// Track processing results
	for _, pair := range pairs {
		result.ProcessedIDs = append(result.ProcessedIDs, pair.BookingID)
		result.SuccessCount++
	}

	result.TotalProcessed = int64(len(pairs))
	result.EndTime = time.Now()

	s.logger.Info().
		Int64("processed", result.TotalProcessed).
		Int64("success", result.SuccessCount).
		Int64("failed", result.FailureCount).
		Dur("duration", result.EndTime.Sub(result.StartTime)).
		Msg("completed booking processing")

	return result, nil
}

// ProcessSingleBooking processes a single booking
func (s *ProcessorService) ProcessSingleBooking(ctx context.Context, bookingID int64) error {
	s.logger.Info().Int64("booking_id", bookingID).Msg("processing single booking")

	// Get booking
	booking, err := s.bookingRepo.GetByID(ctx, bookingID)
	if err != nil {
		return fmt.Errorf("failed to get booking: %w", err)
	}

	// Validate
	if err := s.ValidateBookingData(ctx, booking.ID, booking.NameID); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Generate files for this booking
	pair := domain.BookingNamePair{
		BookingID: booking.ID,
		NameID:    booking.NameID,
	}

	_, err = s.GenerateVINEData(ctx, []domain.BookingNamePair{pair})
	return err
}

// FetchBookingsForProcessing retrieves bookings that need processing
func (s *ProcessorService) FetchBookingsForProcessing(ctx context.Context, since time.Time) ([]domain.BookingNamePair, error) {
	pairs, err := s.bookingRepo.FetchPairsForProcessing(ctx, since)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch booking pairs: %w", err)
	}

	s.logger.Debug().
		Int("count", len(pairs)).
		Time("since", since).
		Msg("fetched booking pairs")

	return pairs, nil
}

// GenerateVINEData generates VINE file data for booking pairs
func (s *ProcessorService) GenerateVINEData(ctx context.Context, pairs []domain.BookingNamePair) (*domain.VINEFileData, error) {
	if len(pairs) == 0 {
		return &domain.VINEFileData{}, nil
	}

	vineData := &domain.VINEFileData{}

	for _, pair := range pairs {
		// Generate prisoner record
		prisonerData, err := s.fileGen.GeneratePrisonerRecord(ctx, pair.BookingID, pair.NameID)
		if err != nil {
			s.logger.Error().
				Err(err).
				Int64("booking_id", pair.BookingID).
				Int64("name_id", pair.NameID).
				Msg("failed to generate prisoner record")
			continue
		}
		vineData.PrisonerData += prisonerData

		// Generate charges records
		chargesData, err := s.fileGen.GenerateChargesRecords(ctx, pair.BookingID, pair.NameID)
		if err != nil {
			s.logger.Error().
				Err(err).
				Int64("booking_id", pair.BookingID).
				Msg("failed to generate charges records")
			continue
		}
		vineData.ChargesData += chargesData

		// Generate mugshot records
		mugshotData, err := s.fileGen.GenerateMugshotRecords(ctx, pair.BookingID, pair.NameID)
		if err != nil {
			s.logger.Error().
				Err(err).
				Int64("booking_id", pair.BookingID).
				Msg("failed to generate mugshot records")
			continue
		}
		vineData.MugshotData += mugshotData
	}

	return vineData, nil
}

// ValidateBookingData validates booking data before processing
func (s *ProcessorService) ValidateBookingData(ctx context.Context, bookingID, nameID int64) error {
	// Check booking exists and is valid
	booking, err := s.bookingRepo.GetByID(ctx, bookingID)
	if err != nil {
		return fmt.Errorf("booking not found: %w", err)
	}

	if !booking.CanProcess() {
		return domain.ErrCannotProcessBooking
	}

	// Check prisoner exists
	_, err = s.prisonerRepo.GetByID(ctx, nameID)
	if err != nil {
		return fmt.Errorf("prisoner not found: %w", err)
	}

	// Check arrest exists
	_, err = s.arrestRepo.GetByBookingID(ctx, bookingID)
	if err != nil {
		return fmt.Errorf("arrest not found: %w", err)
	}

	return nil
}

// GetBookingWithDetails retrieves complete booking details
func (s *ProcessorService) GetBookingWithDetails(ctx context.Context, bookingID int64) (*domain.BookingWithDetails, error) {
	// Get booking
	booking, err := s.bookingRepo.GetByID(ctx, bookingID)
	if err != nil {
		return nil, err
	}

	details := &domain.BookingWithDetails{
		Booking: booking,
	}

	// Get prisoner
	prisoner, err := s.prisonerRepo.GetByBookingID(ctx, bookingID)
	if err == nil {
		details.Prisoner = prisoner
	}

	// Get arrest
	arrest, err := s.arrestRepo.GetByBookingID(ctx, bookingID)
	if err == nil {
		details.Arrest = arrest
	}

	// Get charges
	charges, err := s.chargeRepo.GetByBookingID(ctx, bookingID)
	if err == nil {
		details.Charges = charges
	}

	// Get release (if exists)
	release, err := s.releaseRepo.GetByBookingID(ctx, bookingID)
	if err == nil {
		details.Release = release
	}

	// Get facility
	facility, err := s.facilityRepo.GetLatestByBookingID(ctx, bookingID)
	if err == nil {
		details.Facility = facility
	}

	// Get mugshots
	if prisoner != nil {
		mugshots, err := s.mugshotRepo.GetByNameID(ctx, prisoner.ID)
		if err == nil {
			details.Mugshots = mugshots
		}
	}

	return details, nil
}
