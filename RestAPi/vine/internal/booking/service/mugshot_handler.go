package service

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Deo-Mugabe/GOLANG/internal/booking/domain"
	"github.com/Deo-Mugabe/GOLANG/internal/platform/config"
	"github.com/rs/zerolog"
)

// MugshotService handles mugshot file operations
type MugshotService struct {
	mugshotRepo domain.MugshotRepository
	config      *config.Config
	logger      zerolog.Logger
}

// NewMugshotService creates a new mugshot service
func NewMugshotService(
	mugshotRepo domain.MugshotRepository,
	config *config.Config,
	logger zerolog.Logger,
) *MugshotService {
	return &MugshotService{
		mugshotRepo: mugshotRepo,
		config:      config,
		logger:      logger,
	}
}

// CopyMugshot copies mugshot to output directory
func (s *MugshotService) CopyMugshot(ctx context.Context, nameID, bookingID int64) error {
	// Get latest mugshot
	mugshots, err := s.mugshotRepo.GetByNameID(ctx, nameID)
	if err != nil {
		return fmt.Errorf("failed to get mugshots: %w", err)
	}

	if len(mugshots) == 0 {
		s.logger.Debug().Int64("name_id", nameID).Msg("no mugshots found")
		return nil
	}

	// Use first (most recent) mugshot
	mugshot := mugshots[0]

	// Build source and destination paths
	sourceFile := mugshot.FileName()
	sourcePath := filepath.Join(s.config.Files.MugshotDir, sourceFile)
	destFile := fmt.Sprintf("%d.jpg", bookingID)
	destPath := filepath.Join(s.config.Files.NewMugshotDir, destFile)

	// Check source exists
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		s.logger.Warn().
			Str("source", sourcePath).
			Msg("mugshot file not found on disk")
		return nil
	}

	// Copy file
	if err := copyFile(sourcePath, destPath); err != nil {
		return fmt.Errorf("failed to copy mugshot: %w", err)
	}

	s.logger.Info().
		Str("source", sourceFile).
		Str("dest", destFile).
		Msg("copied mugshot")

	return nil
}

// ClearOutputDirectory clears the mugshot output directory
func (s *MugshotService) ClearOutputDirectory(ctx context.Context) error {
	dir := s.config.Files.NewMugshotDir

	// Create directory if it doesn't exist
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// List all files
	files, err := filepath.Glob(filepath.Join(dir, "*.jpg"))
	if err != nil {
		return fmt.Errorf("failed to list files: %w", err)
	}

	// Delete each file
	for _, file := range files {
		if err := os.Remove(file); err != nil {
			s.logger.Warn().
				Err(err).
				Str("file", file).
				Msg("failed to delete file")
		}
	}

	s.logger.Info().
		Int("count", len(files)).
		Msg("cleared mugshot output directory")

	return nil
}

// GetMugshotPath returns the path to a mugshot file
func (s *MugshotService) GetMugshotPath(ctx context.Context, nameID int64) (string, error) {
	mugshot, err := s.mugshotRepo.GetLatestByNameID(ctx, nameID)
	if err != nil {
		return "", err
	}
	if mugshot == nil {
		return "", fmt.Errorf("no mugshot found for name_id %d", nameID)
	}

	return filepath.Join(s.config.Files.MugshotDir, mugshot.FileName()), nil
}

// ListMugshotsForBooking lists mugshots for a booking
func (s *MugshotService) ListMugshotsForBooking(ctx context.Context, bookingID int64) ([]string, error) {
	// This would require getting the name_id from booking first
	// Then listing mugshots for that name_id
	return nil, nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
