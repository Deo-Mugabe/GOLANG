package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/Deo-Mugabe/GOLANG/internal/sysconfig/domain"
	"github.com/rs/zerolog"
)

// ConfigService manages system configuration
type ConfigService struct {
	configRepo domain.ConfigRepository
	logger     zerolog.Logger
}

// NewConfigService creates a new config service
func NewConfigService(
	configRepo domain.ConfigRepository,
	logger zerolog.Logger,
) *ConfigService {
	return &ConfigService{
		configRepo: configRepo,
		logger:     logger,
	}
}

// GetConfig retrieves a configuration value
func (s *ConfigService) GetConfig(ctx context.Context, name string) (*domain.SystemConfig, error) {
	config, err := s.configRepo.GetByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	s.logger.Debug().
		Str("name", name).
		Str("value", config.Value).
		Msg("retrieved configuration")

	return config, nil
}

// GetAllConfigs retrieves all configurations
func (s *ConfigService) GetAllConfigs(ctx context.Context) ([]*domain.SystemConfig, error) {
	configs, err := s.configRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all configs: %w", err)
	}

	s.logger.Debug().
		Int("count", len(configs)).
		Msg("retrieved all configurations")

	return configs, nil
}

// GetFTPConfig retrieves FTP configuration
func (s *ConfigService) GetFTPConfig(ctx context.Context) (*domain.FTPConfig, error) {
	ftpConfig, err := s.configRepo.GetFTPConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get FTP config: %w", err)
	}

	s.logger.Debug().
		Str("server", ftpConfig.ServerName).
		Bool("use_sftp", ftpConfig.UseSFTP).
		Msg("retrieved FTP configuration")

	return ftpConfig, nil
}

// GetFileConfig retrieves file configuration
func (s *ConfigService) GetFileConfig(ctx context.Context) (*domain.FileConfig, error) {
	fileConfig, err := s.configRepo.GetFileConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get file config: %w", err)
	}

	s.logger.Debug().
		Str("vine_file_path", fileConfig.NewVineFilePath).
		Str("mugshot_dir", fileConfig.MugshotDir).
		Msg("retrieved file configuration")

	return fileConfig, nil
}

// UpdateConfig updates a configuration value
func (s *ConfigService) UpdateConfig(ctx context.Context, name, value string) error {
	s.logger.Info().
		Str("name", name).
		Str("value", value).
		Msg("updating configuration")

	// Validate config
	if err := s.ValidateConfig(ctx, name, value); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Get existing config
	config, err := s.configRepo.GetByName(ctx, name)
	if err != nil {
		return fmt.Errorf("config not found: %w", err)
	}

	// Update value
	config.Value = value
	if err := s.configRepo.Update(ctx, config); err != nil {
		return fmt.Errorf("failed to update config: %w", err)
	}

	s.logger.Info().
		Str("name", name).
		Msg("configuration updated successfully")

	return nil
}

// UpdateFTPConfig updates FTP configuration
func (s *ConfigService) UpdateFTPConfig(ctx context.Context, config *domain.FTPConfig) error {
	s.logger.Info().
		Str("server", config.ServerName).
		Bool("use_sftp", config.UseSFTP).
		Msg("updating FTP configuration")

	// Validate FTP config
	if err := s.ValidateFTPConfig(ctx, config); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Update config
	if err := s.configRepo.UpdateFTPConfig(ctx, config); err != nil {
		return fmt.Errorf("failed to update FTP config: %w", err)
	}

	s.logger.Info().Msg("FTP configuration updated successfully")
	return nil
}

// UpdateFileConfig updates file configuration
func (s *ConfigService) UpdateFileConfig(ctx context.Context, config *domain.FileConfig) error {
	s.logger.Info().
		Str("vine_file_path", config.NewVineFilePath).
		Msg("updating file configuration")

	// Update config
	if err := s.configRepo.UpdateFileConfig(ctx, config); err != nil {
		return fmt.Errorf("failed to update file config: %w", err)
	}

	s.logger.Info().Msg("file configuration updated successfully")
	return nil
}

// ValidateConfig validates a configuration value
func (s *ConfigService) ValidateConfig(ctx context.Context, name, value string) error {
	// Basic validation
	if name == "" {
		return fmt.Errorf("config name cannot be empty")
	}

	// Special validation for specific configs
	switch strings.ToLower(name) {
	case "gcvineftpusername":
		if value == "" {
			return fmt.Errorf("FTP username cannot be empty")
		}
	case "gcvineprimaryftpservername":
		if value == "" {
			return fmt.Errorf("FTP server name cannot be empty")
		}
	case "gcvinejailidnumber":
		if value == "" {
			return fmt.Errorf("jail ID number cannot be empty")
		}
	case "gnvineftpfirewalloutport":
		// Could add port validation here
	}

	return nil
}

// ValidateFTPConfig validates FTP configuration
func (s *ConfigService) ValidateFTPConfig(ctx context.Context, config *domain.FTPConfig) error {
	if config.ServerName == "" {
		return fmt.Errorf("FTP server name is required")
	}
	if config.Username == "" {
		return fmt.Errorf("FTP username is required")
	}
	if config.Password == "" {
		return fmt.Errorf("FTP password is required")
	}
	if config.DatFolder == "" {
		return fmt.Errorf("FTP data folder is required")
	}
	if config.MugshotFolder == "" {
		return fmt.Errorf("FTP mugshot folder is required")
	}

	s.logger.Debug().Msg("FTP configuration validated")
	return nil
}

// ReloadConfiguration reloads configuration from database
func (s *ConfigService) ReloadConfiguration(ctx context.Context) error {
	s.logger.Info().Msg("reloading configuration")

	// This could be used to refresh cached configs
	// For now, we just log it since we're not caching
	configs, err := s.configRepo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to reload configuration: %w", err)
	}

	s.logger.Info().
		Int("count", len(configs)).
		Msg("configuration reloaded successfully")

	return nil
}
