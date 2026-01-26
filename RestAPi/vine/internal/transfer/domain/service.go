package domain

import (
	"context"

	sysconfigdomain "github.com/Deo-Mugabe/GOLANG/internal/sysconfig/domain"
)

// TransferService handles file transfers
type TransferService interface {
	// Transfer operations
	TransferDataFile(ctx context.Context, localPath string) error
	TransferMugshots(ctx context.Context, mugshotDir string) error
	TransferAll(ctx context.Context) error

	// Single file operations
	UploadFile(ctx context.Context, localPath, remotePath string) error
	DownloadFile(ctx context.Context, remotePath, localPath string) error

	// Connection management
	TestConnection(ctx context.Context) error
	Reconnect(ctx context.Context) error
}

// FTPService handles FTP operations
type FTPService interface {
	// Connection
	Connect(ctx context.Context, config *sysconfigdomain.FTPConfig) error
	Disconnect(ctx context.Context) error
	IsConnected() bool

	// File operations
	Upload(ctx context.Context, localPath, remotePath string) error
	Download(ctx context.Context, remotePath, localPath string) error
	List(ctx context.Context, remotePath string) ([]string, error)
	Delete(ctx context.Context, remotePath string) error
}

// SFTPService handles SFTP operations
type SFTPService interface {
	// Connection
	Connect(ctx context.Context, config *sysconfigdomain.FTPConfig) error
	Disconnect(ctx context.Context) error
	IsConnected() bool

	// File operations
	Upload(ctx context.Context, localPath, remotePath string) error
	Download(ctx context.Context, remotePath, localPath string) error
	List(ctx context.Context, remotePath string) ([]string, error)
	Delete(ctx context.Context, remotePath string) error
}

// FileManagerService handles local file operations
type FileManagerService interface {
	// File operations
	Copy(ctx context.Context, src, dst string) error
	Move(ctx context.Context, src, dst string) error
	Delete(ctx context.Context, path string) error
	Exists(ctx context.Context, path string) (bool, error)

	// Directory operations
	CreateDir(ctx context.Context, path string) error
	ClearDir(ctx context.Context, path string) error
	ListFiles(ctx context.Context, dir string) ([]string, error)

	// Validation
	ValidateFilePath(ctx context.Context, path string) error
	GetFileSize(ctx context.Context, path string) (int64, error)
}
