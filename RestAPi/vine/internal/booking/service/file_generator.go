package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Deo-Mugabe/GOLANG/internal/booking/domain"
	"github.com/Deo-Mugabe/GOLANG/internal/platform/config"
	syscfg "github.com/Deo-Mugabe/GOLANG/internal/sysconfig/domain"
	"github.com/rs/zerolog"
)

// FileGeneratorService generates VINE format files
type FileGeneratorService struct {
	bookingRepo  domain.BookingRepository
	prisonerRepo domain.PrisonerRepository
	chargeRepo   domain.ChargeRepository
	arrestRepo   domain.ArrestRepository
	releaseRepo  domain.ReleaseRepository
	facilityRepo domain.FacilityRepository
	lookupRepo   syscfg.LookupRepository
	config       *config.Config
	logger       zerolog.Logger
}

// NewFileGeneratorService creates a new file generator service
func NewFileGeneratorService(
	bookingRepo domain.BookingRepository,
	prisonerRepo domain.PrisonerRepository,
	chargeRepo domain.ChargeRepository,
	arrestRepo domain.ArrestRepository,
	releaseRepo domain.ReleaseRepository,
	facilityRepo domain.FacilityRepository,
	lookupRepo syscfg.LookupRepository,
	config *config.Config,
	logger zerolog.Logger,
) *FileGeneratorService {
	return &FileGeneratorService{
		bookingRepo:  bookingRepo,
		prisonerRepo: prisonerRepo,
		chargeRepo:   chargeRepo,
		arrestRepo:   arrestRepo,
		releaseRepo:  releaseRepo,
		facilityRepo: facilityRepo,
		lookupRepo:   lookupRepo,
		config:       config,
		logger:       logger,
	}
}

// GeneratePrisonerRecord generates a prisoner record in VINE format
func (s *FileGeneratorService) GeneratePrisonerRecord(ctx context.Context, bookingID, nameID int64) (string, error) {
	var sb strings.Builder

	// Get all required data
	prisoner, err := s.prisonerRepo.GetByID(ctx, nameID)
	if err != nil {
		return "", fmt.Errorf("failed to get prisoner: %w", err)
	}

	arrest, err := s.arrestRepo.GetFirstByBookingID(ctx, bookingID)
	if err != nil {
		return "", fmt.Errorf("failed to get arrest: %w", err)
	}

	booking, err := s.bookingRepo.GetByID(ctx, bookingID)
	if err != nil {
		return "", fmt.Errorf("failed to get booking: %w", err)
	}

	facility, _ := s.facilityRepo.GetLatestByBookingID(ctx, bookingID)
	release, _ := s.releaseRepo.GetByBookingID(ctx, bookingID)
	alias, _ := s.prisonerRepo.GetAlias(ctx, fmt.Sprintf("%d", nameID), "AKA")

	// Get agency lookup
	agencyLookup, _ := s.lookupRepo.GetByKey(ctx, booking.Agency, "AGCY")
	agencyMsg := ""
	if agencyLookup != nil {
		agencyMsg = agencyLookup.Message
	}

	// Build VINE prisoner record
	sb.WriteString(padRight(s.config.VINE.PrisonerFileHeader, 10))
	sb.WriteString(padRight(s.config.VINE.JailIDNumber, 12))
	sb.WriteString(padRight(prisoner.StateID, 25))
	sb.WriteString(padRight(fmt.Sprintf("%d", prisoner.ID), 25))
	sb.WriteString(padRight(fmt.Sprintf("%d", booking.ID), 25))
	sb.WriteString(padRight(arrest.CaseID, 25))
	sb.WriteString(padRight(prisoner.FirstName, 20))
	sb.WriteString(padRight(prisoner.MiddleName, 20))
	sb.WriteString(padRight(prisoner.LastName, 20))

	// DOB
	dob, _ := prisoner.ParseDOB()
	if !dob.IsZero() {
		sb.WriteString(dob.Format("20060102"))
	} else {
		sb.WriteString(strings.Repeat(" ", 8))
	}

	sb.WriteString(padRight(prisoner.Race, 1))
	sb.WriteString(padRight(prisoner.Sex, 1))
	sb.WriteString(padRight(prisoner.Height, 4))
	sb.WriteString(padRight(prisoner.Weight, 4))
	sb.WriteString(padRight(prisoner.SSN, 9))
	sb.WriteString(padRight(agencyMsg, 12))

	// Arrest date
	if !arrest.ArrestDate.IsZero() {
		sb.WriteString(arrest.ArrestDate.Format("20060102"))
		sb.WriteString(strings.Repeat(" ", 8)) // Time placeholder
	} else {
		sb.WriteString(strings.Repeat(" ", 16))
	}

	sb.WriteString(padRight(fmt.Sprintf("%d", booking.FacilityID), 12))

	// Booking date
	sb.WriteString(booking.BookDate.Format("200601021504"))

	// Release info
	if release != nil {
		sb.WriteString(padRight(release.Reason, 12))
		sb.WriteString(release.ReleaseTime.Format("20060102"))
		sb.WriteString(release.ReleaseTime.Format("1504"))
	} else {
		sb.WriteString(strings.Repeat(" ", 28))
	}

	// Arrest ID
	sb.WriteString(padRight(fmt.Sprintf("%d", arrest.ID), 14))

	// Address
	address := prisoner.FullAddress()
	sb.WriteString(padRight(address, 58))
	sb.WriteString(padRight(prisoner.City, 20))
	sb.WriteString(padRight(prisoner.State, 2))
	sb.WriteString(padRight(prisoner.Zip, 10))
	sb.WriteString(padRight(prisoner.Birthplace, 20))
	sb.WriteString(padRight(prisoner.DrLicense, 25))
	sb.WriteString(padRight(prisoner.DLState, 2))
	sb.WriteString(padRight(prisoner.Marital, 1))
	sb.WriteString(padRight(prisoner.Occupation, 15))
	sb.WriteString(padRight(prisoner.Eye, 10))
	sb.WriteString(padRight(prisoner.Hair, 10))
	sb.WriteString(padRight(prisoner.Employer, 30))
	sb.WriteString(padRight(prisoner.HomePhone, 10))
	sb.WriteString(padRight(prisoner.WorkPhone, 15))
	sb.WriteString(padRight(prisoner.MobilePhone, 10))

	// Alias
	if alias != nil && (alias.FirstName != "" || alias.LastName != "") {
		sb.WriteString("N")
		sb.WriteString(padRight(alias.FirstName, 20))
		sb.WriteString(padRight(alias.LastName, 20))
	} else {
		sb.WriteString(strings.Repeat(" ", 41))
	}

	// Facility
	if facility != nil {
		sb.WriteString(padRight(facility.Facility, 10))
		sb.WriteString(padRight(facility.Section, 10))
		sb.WriteString(padRight(facility.Unit, 10))
		sb.WriteString(padRight(facility.Bed, 10))
	} else {
		sb.WriteString(strings.Repeat(" ", 40))
	}

	sb.WriteString("\n")
	return sb.String(), nil
}

// GenerateChargesRecords generates charge records in VINE format
func (s *FileGeneratorService) GenerateChargesRecords(ctx context.Context, bookingID, nameID int64) (string, error) {
	var sb strings.Builder

	prisoner, err := s.prisonerRepo.GetByID(ctx, nameID)
	if err != nil {
		return "", err
	}

	arrest, err := s.arrestRepo.GetFirstByBookingID(ctx, bookingID)
	if err != nil {
		return "", err
	}

	charges, err := s.chargeRepo.GetByArrestID(ctx, arrest.ID)
	if err != nil {
		return "", err
	}

	for _, charge := range charges {
		sb.WriteString(padRight(s.config.VINE.ChargesFileHeader, 10))
		sb.WriteString(padRight(s.config.VINE.JailIDNumber, 12))
		sb.WriteString(padRight(prisoner.StateID, 25))
		sb.WriteString(padRight(fmt.Sprintf("%d", prisoner.ID), 25))
		sb.WriteString(padRight(fmt.Sprintf("%d", charge.BookingID), 25))
		sb.WriteString(padRight(charge.ChargeCode, 25))
		sb.WriteString(padRight(charge.FelonyMisd, 10))
		sb.WriteString(padRight(charge.Count, 4))
		sb.WriteString(padRight(charge.Sequence, 4))
		sb.WriteString(padRight(charge.BondAmount, 15))
		sb.WriteString(padRight(charge.BondType, 4))
		sb.WriteString(padRight(fmt.Sprintf("%d", charge.ArrestID), 14))
		sb.WriteString(padRight(charge.Description, 60))
		sb.WriteString("\n")
	}

	return sb.String(), nil
}

// GenerateMugshotRecords generates mugshot records in VINE format
func (s *FileGeneratorService) GenerateMugshotRecords(ctx context.Context, bookingID, nameID int64) (string, error) {
	// This will be implemented by MugshotService
	// Returns the mugshot metadata lines for the VINE file
	return "", nil
}

// WriteVINEFile writes VINE data to file
func (s *FileGeneratorService) WriteVINEFile(ctx context.Context, data *domain.VINEFileData, filePath string) error {
	if data.IsEmpty() {
		return fmt.Errorf("no data to write")
	}

	content := data.Combine()

	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write file
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	s.logger.Info().
		Str("file", filePath).
		Int("size", len(content)).
		Msg("wrote VINE file")

	return nil
}

// ValidateVINEFile validates VINE file format
func (s *FileGeneratorService) ValidateVINEFile(ctx context.Context, filepath string) error {
	// Check file exists
	if _, err := os.Stat(filepath); err != nil {
		return fmt.Errorf("file not found: %w", err)
	}

	// Read file
	content, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Basic validation
	if len(content) == 0 {
		return fmt.Errorf("file is empty")
	}

	lines := strings.Split(string(content), "\n")
	if len(lines) == 0 {
		return fmt.Errorf("file has no lines")
	}

	s.logger.Debug().
		Str("file", filepath).
		Int("lines", len(lines)).
		Msg("validated VINE file")

	return nil
}

// padRight pads string to width with spaces
func padRight(s string, width int) string {
	if len(s) >= width {
		return s[:width]
	}
	return s + strings.Repeat(" ", width-len(s))
}
