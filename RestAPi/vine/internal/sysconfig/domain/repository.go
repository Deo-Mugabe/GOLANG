package domain

import "context"

// ConfigRepository handles system configuration
type ConfigRepository interface {
	// Read operations
	GetByID(ctx context.Context, id int64) (*SystemConfig, error)
	GetByName(ctx context.Context, name string) (*SystemConfig, error)
	GetAll(ctx context.Context) ([]*SystemConfig, error)
	GetByPrefix(ctx context.Context, prefix string) ([]*SystemConfig, error)

	// Write operations
	Create(ctx context.Context, config *SystemConfig) error
	Update(ctx context.Context, config *SystemConfig) error
	Delete(ctx context.Context, id int64) error

	// Batch operations
	GetMultipleByNames(ctx context.Context, names []string) (map[string]*SystemConfig, error)

	// Specialized queries
	GetFTPConfig(ctx context.Context) (*FTPConfig, error)
	GetFileConfig(ctx context.Context) (*FileConfig, error)
	UpdateFTPConfig(ctx context.Context, config *FTPConfig) error
	UpdateFileConfig(ctx context.Context, config *FileConfig) error
}

// LookupRepository handles system lookup tables
type LookupRepository interface {
	// Read operations
	GetByID(ctx context.Context, id int) (*SystemLookup, error)
	GetByKey(ctx context.Context, codeAgency, codeKey string) (*SystemLookup, error)
	GetByAgency(ctx context.Context, codeAgency string) ([]*SystemLookup, error)
	List(ctx context.Context, limit, offset int) ([]*SystemLookup, error)

	// Write operations
	Create(ctx context.Context, lookup *SystemLookup) error
	Update(ctx context.Context, lookup *SystemLookup) error
	Delete(ctx context.Context, id int) error
}
