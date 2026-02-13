package service

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
)

// FileManagerService handles local file operations
type FileManagerService struct {
	logger zerolog.Logger
}

// NewFileManagerService creates a new file manager service
func NewFileManagerService(logger zerolog.Logger) *FileManagerService {
	return &FileManagerService{
		logger: logger,
	}
}

// Copy copies a file
func (s *FileManagerService) Copy(ctx context.Context, src, dst string) error {
	s.logger.Debug().
		Str("src", src).
		Str("dst", dst).
		Msg("copying file")

	// Open source file
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	// Get source file info
	sourceInfo, err := sourceFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to get source file info: %w", err)
	}

	// Create destination directory if needed
	destDir := filepath.Dir(dst)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Create destination file
	destFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	// Copy content
	written, err := io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	if written != sourceInfo.Size() {
		return fmt.Errorf("incomplete copy: wrote %d of %d bytes", written, sourceInfo.Size())
	}

	s.logger.Info().
		Str("file", filepath.Base(src)).
		Int64("size", written).
		Msg("file copied successfully")

	return nil
}

// Move moves a file
func (s *FileManagerService) Move(ctx context.Context, src, dst string) error {
	s.logger.Debug().
		Str("src", src).
		Str("dst", dst).
		Msg("moving file")

	// Create destination directory if needed
	destDir := filepath.Dir(dst)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Try rename first (fastest if on same filesystem)
	if err := os.Rename(src, dst); err == nil {
		s.logger.Info().
			Str("file", filepath.Base(src)).
			Msg("file moved successfully")
		return nil
	}

	// Fall back to copy + delete
	if err := s.Copy(ctx, src, dst); err != nil {
		return err
	}

	if err := os.Remove(src); err != nil {
		return fmt.Errorf("failed to remove source file after copy: %w", err)
	}

	s.logger.Info().
		Str("file", filepath.Base(src)).
		Msg("file moved successfully (via copy)")

	return nil
}

// Delete deletes a file
func (s *FileManagerService) Delete(ctx context.Context, path string) error {
	s.logger.Debug().Str("path", path).Msg("deleting file")

	if err := os.Remove(path); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	s.logger.Info().Str("file", filepath.Base(path)).Msg("file deleted")
	return nil
}

// Exists checks if a file exists
func (s *FileManagerService) Exists(ctx context.Context, path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// CreateDir creates a directory
func (s *FileManagerService) CreateDir(ctx context.Context, path string) error {
	s.logger.Debug().Str("path", path).Msg("creating directory")

	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	s.logger.Info().Str("dir", path).Msg("directory created")
	return nil
}

// ClearDir removes all files from a directory
func (s *FileManagerService) ClearDir(ctx context.Context, path string) error {
	s.logger.Debug().Str("path", path).Msg("clearing directory")

	// Ensure directory exists
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Read directory contents
	entries, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	// Delete each file
	deletedCount := 0
	for _, entry := range entries {
		if !entry.IsDir() {
			filePath := filepath.Join(path, entry.Name())
			if err := os.Remove(filePath); err != nil {
				s.logger.Warn().
					Err(err).
					Str("file", entry.Name()).
					Msg("failed to delete file")
			} else {
				deletedCount++
			}
		}
	}

	s.logger.Info().
		Str("dir", path).
		Int("deleted", deletedCount).
		Msg("directory cleared")

	return nil
}

// ListFiles lists all files in a directory
func (s *FileManagerService) ListFiles(ctx context.Context, dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	files := make([]string, 0)
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}

	s.logger.Debug().
		Str("dir", dir).
		Int("count", len(files)).
		Msg("listed files")

	return files, nil
}

// ValidateFilePath validates a file path
func (s *FileManagerService) ValidateFilePath(ctx context.Context, path string) error {
	// Check if path is absolute
	if !filepath.IsAbs(path) {
		return fmt.Errorf("path must be absolute: %s", path)
	}

	// Check if path exists
	exists, err := s.Exists(ctx, path)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("file does not exist: %s", path)
	}

	// Check if it's a file (not directory)
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return fmt.Errorf("path is a directory, not a file: %s", path)
	}

	return nil
}

// GetFileSize returns file size in bytes
func (s *FileManagerService) GetFileSize(ctx context.Context, path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, fmt.Errorf("failed to get file info: %w", err)
	}

	return info.Size(), nil
}
