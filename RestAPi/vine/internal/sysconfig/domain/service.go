package domain

import "context"

// ConfigService manages system configuration
type ConfigService interface {
	// Configuration retrieval
	GetConfig(ctx context.Context, name string) (*SystemConfig, error)
	GetAllConfigs(ctx context.Context) ([]*SystemConfig, error)

	// Specialized configs
	GetFTPConfig(ctx context.Context) (*FTPConfig, error)
	GetFileConfig(ctx context.Context) (*FileConfig, error)

	// Configuration updates
	UpdateConfig(ctx context.Context, name, value string) error
	UpdateFTPConfig(ctx context.Context, config *FTPConfig) error
	UpdateFileConfig(ctx context.Context, config *FileConfig) error

	// Validation
	ValidateConfig(ctx context.Context, name, value string) error
	ValidateFTPConfig(ctx context.Context, config *FTPConfig) error

	// Reload
	ReloadConfiguration(ctx context.Context) error
}

// LookupService manages system lookups
type LookupService interface {
	// Lookup retrieval
	GetLookup(ctx context.Context, codeAgency, codeKey string) (*SystemLookup, error)
	GetAgencyLookups(ctx context.Context, codeAgency string) ([]*SystemLookup, error)

	// Translation
	TranslateCode(ctx context.Context, codeAgency, codeKey string) (string, error)

	// Cache management
	InvalidateCache(ctx context.Context) error
	PreloadCache(ctx context.Context) error
}
