package service

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	sysconfigdomain "github.com/Deo-Mugabe/GOLANG/internal/sysconfig/domain"
	"github.com/pkg/sftp"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/ssh"
)

// SFTPService handles SFTP operations
type SFTPService struct {
	client *sftp.Client
	conn   *ssh.Client
	logger zerolog.Logger
}

// NewSFTPService creates a new SFTP service
func NewSFTPService(logger zerolog.Logger) *SFTPService {
	return &SFTPService{
		logger: logger,
	}
}

// Connect establishes SFTP connection
func (s *SFTPService) Connect(ctx context.Context, config *sysconfigdomain.FTPConfig) error {
	s.logger.Info().
		Str("server", config.ServerName).
		Str("username", config.Username).
		Msg("connecting to SFTP server")

	// Parse port
	port := 22
	if config.FirewallPort != "" {
		if p, err := strconv.Atoi(config.FirewallPort); err == nil {
			port = p
		}
	}

	// Create SSH client config
	sshConfig := &ssh.ClientConfig{
		User: config.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(config.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Note: In production, verify host key
	}

	// Connect to SSH server
	addr := fmt.Sprintf("%s:%d", config.ServerName, port)
	conn, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to SSH server: %w", err)
	}

	// Create SFTP client
	client, err := sftp.NewClient(conn)
	if err != nil {
		conn.Close()
		return fmt.Errorf("failed to create SFTP client: %w", err)
	}

	s.conn = conn
	s.client = client

	s.logger.Info().Msg("SFTP connection established")
	return nil
}

// Disconnect closes SFTP connection
func (s *SFTPService) Disconnect(ctx context.Context) error {
	if s.client != nil {
		s.client.Close()
		s.client = nil
	}
	if s.conn != nil {
		s.conn.Close()
		s.conn = nil
	}

	s.logger.Info().Msg("SFTP connection closed")
	return nil
}

// IsConnected checks if connected
func (s *SFTPService) IsConnected() bool {
	return s.client != nil
}

// Upload uploads a file
func (s *SFTPService) Upload(ctx context.Context, localPath, remotePath string) error {
	if !s.IsConnected() {
		return fmt.Errorf("not connected to SFTP server")
	}

	s.logger.Debug().
		Str("local", localPath).
		Str("remote", remotePath).
		Msg("uploading file via SFTP")

	// Open local file
	localFile, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("failed to open local file: %w", err)
	}
	defer localFile.Close()

	// Get file info
	fileInfo, err := localFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	// Create remote directory if needed
	remoteDir := filepath.Dir(remotePath)
	if err := s.client.MkdirAll(remoteDir); err != nil {
		s.logger.Warn().Err(err).Str("dir", remoteDir).Msg("failed to create remote directory")
	}

	// Create remote file
	remoteFile, err := s.client.Create(remotePath)
	if err != nil {
		return fmt.Errorf("failed to create remote file: %w", err)
	}
	defer remoteFile.Close()

	// Copy file content
	written, err := io.Copy(remoteFile, localFile)
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	s.logger.Info().
		Str("file", filepath.Base(localPath)).
		Int64("size", written).
		Msg("file uploaded successfully via SFTP")

	if written != fileInfo.Size() {
		return fmt.Errorf("incomplete upload: wrote %d of %d bytes", written, fileInfo.Size())
	}

	return nil
}

// Download downloads a file
func (s *SFTPService) Download(ctx context.Context, remotePath, localPath string) error {
	if !s.IsConnected() {
		return fmt.Errorf("not connected to SFTP server")
	}

	s.logger.Debug().
		Str("remote", remotePath).
		Str("local", localPath).
		Msg("downloading file via SFTP")

	// Open remote file
	remoteFile, err := s.client.Open(remotePath)
	if err != nil {
		return fmt.Errorf("failed to open remote file: %w", err)
	}
	defer remoteFile.Close()

	// Create local directory if needed
	localDir := filepath.Dir(localPath)
	if err := os.MkdirAll(localDir, 0755); err != nil {
		return fmt.Errorf("failed to create local directory: %w", err)
	}

	// Create local file
	localFile, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create local file: %w", err)
	}
	defer localFile.Close()

	// Copy file content
	written, err := io.Copy(localFile, remoteFile)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}

	s.logger.Info().
		Str("file", filepath.Base(remotePath)).
		Int64("size", written).
		Msg("file downloaded successfully via SFTP")

	return nil
}

// List lists files in remote directory
func (s *SFTPService) List(ctx context.Context, remotePath string) ([]string, error) {
	if !s.IsConnected() {
		return nil, fmt.Errorf("not connected to SFTP server")
	}

	files, err := s.client.ReadDir(remotePath)
	if err != nil {
		return nil, fmt.Errorf("failed to list remote directory: %w", err)
	}

	fileNames := make([]string, 0, len(files))
	for _, file := range files {
		if !file.IsDir() {
			fileNames = append(fileNames, file.Name())
		}
	}

	return fileNames, nil
}

// Delete deletes a remote file
func (s *SFTPService) Delete(ctx context.Context, remotePath string) error {
	if !s.IsConnected() {
		return fmt.Errorf("not connected to SFTP server")
	}

	if err := s.client.Remove(remotePath); err != nil {
		return fmt.Errorf("failed to delete remote file: %w", err)
	}

	s.logger.Info().Str("file", remotePath).Msg("remote file deleted")
	return nil
}
