package service

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/Deo-Mugabe/GOLANG/internal/platform/config"
	sysconfigdomain "github.com/Deo-Mugabe/GOLANG/internal/sysconfig/domain"
	"github.com/Deo-Mugabe/GOLANG/internal/transfer/domain"
	"github.com/rs/zerolog"
)

// TransferService handles file transfers
type TransferService struct {
	ftpService  domain.FTPService
	sftpService domain.SFTPService
	fileManager domain.FileManagerService
	configRepo  sysconfigdomain.ConfigRepository
	config      *config.Config
	logger      zerolog.Logger
}

// NewTransferService creates a new transfer service
func NewTransferService(
	ftpService domain.FTPService,
	sftpService domain.SFTPService,
	fileManager domain.FileManagerService,
	configRepo sysconfigdomain.ConfigRepository,
	cfg *config.Config,
	logger zerolog.Logger,
) *TransferService {
	return &TransferService{
		ftpService:  ftpService,
		sftpService: sftpService,
		fileManager: fileManager,
		configRepo:  configRepo,
		config:      cfg,
		logger:      logger,
	}
}

// TransferAll transfers data file and mugshots
func (s *TransferService) TransferAll(ctx context.Context) error {
	s.logger.Info().Msg("starting complete file transfer")

	// Get file config
	fileConfig, err := s.configRepo.GetFileConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to get file config: %w", err)
	}

	// Build data file path
	dataFilePath := filepath.Join(fileConfig.NewVineFilePath, fileConfig.InterFile)

	// Check if data file exists
	exists, err := s.fileManager.Exists(ctx, dataFilePath)
	if err != nil {
		return fmt.Errorf("failed to check data file: %w", err)
	}
	if !exists {
		return fmt.Errorf("data file not found: %s", dataFilePath)
	}

	// Transfer data file
	if err := s.TransferDataFile(ctx, dataFilePath); err != nil {
		return fmt.Errorf("failed to transfer data file: %w", err)
	}

	// Transfer mugshots
	if err := s.TransferMugshots(ctx, fileConfig.NewMugshotDir); err != nil {
		return fmt.Errorf("failed to transfer mugshots: %w", err)
	}

	s.logger.Info().Msg("complete file transfer successful")
	return nil
}

// TransferDataFile transfers the VINE data file
func (s *TransferService) TransferDataFile(ctx context.Context, localPath string) error {
	s.logger.Info().Str("file", localPath).Msg("transferring data file")

	// Get FTP config
	ftpConfig, err := s.configRepo.GetFTPConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to get FTP config: %w", err)
	}

	// Build remote path
	fileName := filepath.Base(localPath)
	remotePath := filepath.Join(ftpConfig.DatFolder, fileName)

	// Upload file
	if err := s.UploadFile(ctx, localPath, remotePath); err != nil {
		return fmt.Errorf("failed to upload data file: %w", err)
	}

	s.logger.Info().
		Str("local", localPath).
		Str("remote", remotePath).
		Msg("data file transferred successfully")

	return nil
}

// TransferMugshots transfers all mugshot files
func (s *TransferService) TransferMugshots(ctx context.Context, mugshotDir string) error {
	s.logger.Info().Str("dir", mugshotDir).Msg("transferring mugshots")

	// Get FTP config
	ftpConfig, err := s.configRepo.GetFTPConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to get FTP config: %w", err)
	}

	// List mugshot files
	files, err := s.fileManager.ListFiles(ctx, mugshotDir)
	if err != nil {
		return fmt.Errorf("failed to list mugshots: %w", err)
	}

	if len(files) == 0 {
		s.logger.Warn().Msg("no mugshots found to transfer")
		return nil
	}

	// Transfer each mugshot
	successCount := 0
	for _, file := range files {
		localPath := filepath.Join(mugshotDir, file)
		remotePath := filepath.Join(ftpConfig.MugshotFolder, file)

		if err := s.UploadFile(ctx, localPath, remotePath); err != nil {
			s.logger.Error().
				Err(err).
				Str("file", file).
				Msg("failed to transfer mugshot")
			continue
		}

		successCount++
	}

	s.logger.Info().
		Int("total", len(files)).
		Int("successful", successCount).
		Msg("mugshot transfer completed")

	if successCount < len(files) {
		return fmt.Errorf("failed to transfer %d out of %d mugshots", len(files)-successCount, len(files))
	}

	return nil
}

// UploadFile uploads a single file
func (s *TransferService) UploadFile(ctx context.Context, localPath, remotePath string) error {
	// Get FTP config
	ftpConfig, err := s.configRepo.GetFTPConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to get FTP config: %w", err)
	}

	// Choose appropriate service
	if ftpConfig.UseSFTP {
		// Connect SFTP
		if !s.sftpService.IsConnected() {
			if err := s.sftpService.Connect(ctx, ftpConfig); err != nil {
				return fmt.Errorf("failed to connect SFTP: %w", err)
			}
			defer s.sftpService.Disconnect(ctx)
		}

		return s.sftpService.Upload(ctx, localPath, remotePath)
	}

	// Connect FTP
	if !s.ftpService.IsConnected() {
		if err := s.ftpService.Connect(ctx, ftpConfig); err != nil {
			return fmt.Errorf("failed to connect FTP: %w", err)
		}
		defer s.ftpService.Disconnect(ctx)
	}

	return s.ftpService.Upload(ctx, localPath, remotePath)
}

// DownloadFile downloads a single file
func (s *TransferService) DownloadFile(ctx context.Context, remotePath, localPath string) error {
	// Get FTP config
	ftpConfig, err := s.configRepo.GetFTPConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to get FTP config: %w", err)
	}

	// Choose appropriate service
	if ftpConfig.UseSFTP {
		if !s.sftpService.IsConnected() {
			if err := s.sftpService.Connect(ctx, ftpConfig); err != nil {
				return fmt.Errorf("failed to connect SFTP: %w", err)
			}
			defer s.sftpService.Disconnect(ctx)
		}

		return s.sftpService.Download(ctx, remotePath, localPath)
	}

	if !s.ftpService.IsConnected() {
		if err := s.ftpService.Connect(ctx, ftpConfig); err != nil {
			return fmt.Errorf("failed to connect FTP: %w", err)
		}
		defer s.ftpService.Disconnect(ctx)
	}

	return s.ftpService.Download(ctx, remotePath, localPath)
}

// TestConnection tests FTP/SFTP connection
func (s *TransferService) TestConnection(ctx context.Context) error {
	s.logger.Info().Msg("testing connection")

	ftpConfig, err := s.configRepo.GetFTPConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to get FTP config: %w", err)
	}

	if ftpConfig.UseSFTP {
		if err := s.sftpService.Connect(ctx, ftpConfig); err != nil {
			return fmt.Errorf("SFTP connection failed: %w", err)
		}
		defer s.sftpService.Disconnect(ctx)
	} else {
		if err := s.ftpService.Connect(ctx, ftpConfig); err != nil {
			return fmt.Errorf("FTP connection failed: %w", err)
		}
		defer s.ftpService.Disconnect(ctx)
	}

	s.logger.Info().Msg("connection test successful")
	return nil
}

// Reconnect reconnects to FTP/SFTP
func (s *TransferService) Reconnect(ctx context.Context) error {
	s.logger.Info().Msg("reconnecting")

	ftpConfig, err := s.configRepo.GetFTPConfig(ctx)
	if err != nil {
		return err
	}

	if ftpConfig.UseSFTP {
		s.sftpService.Disconnect(ctx)
		return s.sftpService.Connect(ctx, ftpConfig)
	}

	s.ftpService.Disconnect(ctx)
	return s.ftpService.Connect(ctx, ftpConfig)
}
