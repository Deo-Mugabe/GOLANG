package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/Deo-Mugabe/GOLANG/internal/platform/database"
	"github.com/Deo-Mugabe/GOLANG/internal/sysconfig/domain"
)

// configDB is the database model for sys_cfg table
type configDB struct {
	ID    int64  `db:"sys_cfgid"`
	Name  string `db:"sysname"`
	Value string `db:"defavalue"`
}

type configRepo struct {
	db *database.DB
}

// NewConfigRepository creates a new config repository
func NewConfigRepository(db *database.DB) domain.ConfigRepository {
	return &configRepo{db: db}
}

// GetByID retrieves config by ID
func (r *configRepo) GetByID(ctx context.Context, id int64) (*domain.SystemConfig, error) {
	query := `
        SELECT sys_cfgid, sysname, defavalue
        FROM sys_cfg
        WHERE sys_cfgid = @p1
    `

	var dbModel configDB
	err := r.db.GetContext(ctx, &dbModel, query, id)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("config not found: %d", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	return toConfigDomain(&dbModel), nil
}

// GetByName retrieves config by name
func (r *configRepo) GetByName(ctx context.Context, name string) (*domain.SystemConfig, error) {
	query := `
        SELECT sys_cfgid, sysname, defavalue
        FROM sys_cfg
        WHERE sysname = @p1
    `

	var dbModel configDB
	err := r.db.GetContext(ctx, &dbModel, query, name)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("config not found: %s", name)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get config by name: %w", err)
	}

	return toConfigDomain(&dbModel), nil
}

// GetAll retrieves all configs
func (r *configRepo) GetAll(ctx context.Context) ([]*domain.SystemConfig, error) {
	query := `
        SELECT sys_cfgid, sysname, defavalue
        FROM sys_cfg
        ORDER BY sysname
    `

	var dbModels []configDB
	err := r.db.SelectContext(ctx, &dbModels, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all configs: %w", err)
	}

	configs := make([]*domain.SystemConfig, len(dbModels))
	for i, db := range dbModels {
		configs[i] = toConfigDomain(&db)
	}

	return configs, nil
}

// GetByPrefix retrieves configs by name prefix
func (r *configRepo) GetByPrefix(ctx context.Context, prefix string) ([]*domain.SystemConfig, error) {
	query := `
        SELECT sys_cfgid, sysname, defavalue
        FROM sys_cfg
        WHERE sysname LIKE @p1
        ORDER BY sysname
    `

	var dbModels []configDB
	err := r.db.SelectContext(ctx, &dbModels, query, prefix+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to get configs by prefix: %w", err)
	}

	configs := make([]*domain.SystemConfig, len(dbModels))
	for i, db := range dbModels {
		configs[i] = toConfigDomain(&db)
	}

	return configs, nil
}

// Create creates a new config
func (r *configRepo) Create(ctx context.Context, config *domain.SystemConfig) error {
	query := `
        INSERT INTO sys_cfg (sysname, defavalue)
        OUTPUT INSERTED.sys_cfgid
        VALUES (@p1, @p2)
    `

	var id int64
	err := r.db.GetContext(ctx, &id, query, config.Name, config.Value)
	if err != nil {
		return fmt.Errorf("failed to create config: %w", err)
	}

	config.ID = id
	return nil
}

// Update updates a config
func (r *configRepo) Update(ctx context.Context, config *domain.SystemConfig) error {
	query := `
        UPDATE sys_cfg
        SET sysname = @p1, defavalue = @p2
        WHERE sys_cfgid = @p3
    `

	result, err := r.db.ExecContext(ctx, query, config.Name, config.Value, config.ID)
	if err != nil {
		return fmt.Errorf("failed to update config: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("config not found: %d", config.ID)
	}

	return nil
}

// Delete deletes a config
func (r *configRepo) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx,
		"DELETE FROM sys_cfg WHERE sys_cfgid = @p1", id)
	if err != nil {
		return fmt.Errorf("failed to delete config: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("config not found: %d", id)
	}

	return nil
}

// GetMultipleByNames retrieves multiple configs by names
func (r *configRepo) GetMultipleByNames(ctx context.Context, names []string) (map[string]*domain.SystemConfig, error) {
	if len(names) == 0 {
		return make(map[string]*domain.SystemConfig), nil
	}

	// Build IN clause
	placeholders := make([]string, len(names))
	args := make([]interface{}, len(names))
	for i, name := range names {
		placeholders[i] = fmt.Sprintf("@p%d", i+1)
		args[i] = name
	}

	query := fmt.Sprintf(`
        SELECT sys_cfgid, sysname, defavalue
        FROM sys_cfg
        WHERE sysname IN (%s)
    `, strings.Join(placeholders, ","))

	var dbModels []configDB
	err := r.db.SelectContext(ctx, &dbModels, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get multiple configs: %w", err)
	}

	result := make(map[string]*domain.SystemConfig)
	for _, db := range dbModels {
		result[db.Name] = toConfigDomain(&db)
	}

	return result, nil
}

// GetFTPConfig retrieves grouped FTP configuration
func (r *configRepo) GetFTPConfig(ctx context.Context) (*domain.FTPConfig, error) {
	names := []string{
		"gcvineftpusername",
		"gcvineftppassword",
		"gcvineprimaryftpservername",
		"gcvineftpdatfoldername",
		"gnvineftpfirewalloutport",
		"gcvineftpmugshotfoldername",
		"glvineusesftp",
	}

	configs, err := r.GetMultipleByNames(ctx, names)
	if err != nil {
		return nil, err
	}

	ftpConfig := &domain.FTPConfig{}
	if c, ok := configs["gcvineftpusername"]; ok {
		ftpConfig.Username = c.Value
	}
	if c, ok := configs["gcvineftppassword"]; ok {
		ftpConfig.Password = c.Value
	}
	if c, ok := configs["gcvineprimaryftpservername"]; ok {
		ftpConfig.ServerName = c.Value
	}
	if c, ok := configs["gcvineftpdatfoldername"]; ok {
		ftpConfig.DatFolder = c.Value
	}
	if c, ok := configs["gnvineftpfirewalloutport"]; ok {
		ftpConfig.FirewallPort = c.Value
	}
	if c, ok := configs["gcvineftpmugshotfoldername"]; ok {
		ftpConfig.MugshotFolder = c.Value
	}
	if c, ok := configs["glvineusesftp"]; ok {
		ftpConfig.UseSFTP = strings.ToLower(c.Value) == "true" || c.Value == ".t."
	}

	return ftpConfig, nil
}

// GetFileConfig retrieves grouped file configuration
func (r *configRepo) GetFileConfig(ctx context.Context) (*domain.FileConfig, error) {
	names := []string{
		"gcvinechargesfileheader",
		"gcvineprisonerfileheader",
		"gcvinejailidnumber",
		"gcvinenewmugshotdirectory",
		"gcvinemugshotdirectory",
		"gcvinenewvinefilepath",
		"gcvineinterfile",
	}

	configs, err := r.GetMultipleByNames(ctx, names)
	if err != nil {
		return nil, err
	}

	fileConfig := &domain.FileConfig{}
	if c, ok := configs["gcvinechargesfileheader"]; ok {
		fileConfig.ChargesFileHeader = c.Value
	}
	if c, ok := configs["gcvineprisonerfileheader"]; ok {
		fileConfig.PrisonerFileHeader = c.Value
	}
	if c, ok := configs["gcvinejailidnumber"]; ok {
		fileConfig.JailIDNumber = c.Value
	}
	if c, ok := configs["gcvinenewmugshotdirectory"]; ok {
		fileConfig.NewMugshotDir = c.Value
	}
	if c, ok := configs["gcvinemugshotdirectory"]; ok {
		fileConfig.MugshotDir = c.Value
	}
	if c, ok := configs["gcvinenewvinefilepath"]; ok {
		fileConfig.NewVineFilePath = c.Value
	}
	if c, ok := configs["gcvineinterfile"]; ok {
		fileConfig.InterFile = c.Value
	}

	return fileConfig, nil
}

// UpdateFTPConfig updates FTP configuration
func (r *configRepo) UpdateFTPConfig(ctx context.Context, config *domain.FTPConfig) error {
	updates := map[string]string{
		"gcvineftpusername":          config.Username,
		"gcvineftppassword":          config.Password,
		"gcvineprimaryftpservername": config.ServerName,
		"gcvineftpdatfoldername":     config.DatFolder,
		"gnvineftpfirewalloutport":   config.FirewallPort,
		"gcvineftpmugshotfoldername": config.MugshotFolder,
		"glvineusesftp":              boolToString(config.UseSFTP),
	}

	for name, value := range updates {
		query := "UPDATE sys_cfg SET defavalue = @p1 WHERE sysname = @p2"
		_, err := r.db.ExecContext(ctx, query, value, name)
		if err != nil {
			return fmt.Errorf("failed to update %s: %w", name, err)
		}
	}

	return nil
}

// UpdateFileConfig updates file configuration
func (r *configRepo) UpdateFileConfig(ctx context.Context, config *domain.FileConfig) error {
	updates := map[string]string{
		"gcvinechargesfileheader":   config.ChargesFileHeader,
		"gcvineprisonerfileheader":  config.PrisonerFileHeader,
		"gcvinejailidnumber":        config.JailIDNumber,
		"gcvinenewmugshotdirectory": config.NewMugshotDir,
		"gcvinemugshotdirectory":    config.MugshotDir,
		"gcvinenewvinefilepath":     config.NewVineFilePath,
		"gcvineinterfile":           config.InterFile,
	}

	for name, value := range updates {
		query := "UPDATE sys_cfg SET defavalue = @p1 WHERE sysname = @p2"
		_, err := r.db.ExecContext(ctx, query, value, name)
		if err != nil {
			return fmt.Errorf("failed to update %s: %w", name, err)
		}
	}

	return nil
}

func toConfigDomain(db *configDB) *domain.SystemConfig {
	return &domain.SystemConfig{
		ID:    db.ID,
		Name:  db.Name,
		Value: db.Value,
	}
}

func boolToString(b bool) string {
	if b {
		return ".t."
	}
	return ".f."
}
