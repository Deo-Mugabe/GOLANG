package database

import (
	"fmt"
	"time"
)

// Config holds database connection configuration
type Config struct {
	// Connection details
	Host     string
	Port     int
	Database string
	User     string
	Password string
	Instance string // SQL Server instance name (e.g., "ADTEMPUS")

	// Connection pool settings
	MaxOpenConns    int           // Maximum open connections
	MaxIdleConns    int           // Maximum idle connections
	ConnMaxLifetime time.Duration // Maximum connection lifetime
	ConnMaxIdleTime time.Duration // Maximum idle time

	// Timeouts
	ConnectTimeout time.Duration // Connection timeout
	QueryTimeout   time.Duration // Query timeout

	// SSL/TLS
	Encrypt         bool // Enable encryption
	TrustServerCert bool // Trust server certificate

	// Additional options
	AppName       string        // Application name
	EnableRetry   bool          // Enable automatic retry
	RetryAttempts int           // Number of retry attempts
	RetryInterval time.Duration // Interval between retries
}

// DefaultConfig returns default database configuration
func DefaultConfig() *Config {
	return &Config{
		Host:            "localhost",
		Port:            1433,
		Database:        "vine-test",
		User:            "vine-user",
		Password:        "",
		Instance:        "",
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 2 * time.Minute,
		ConnectTimeout:  10 * time.Second,
		QueryTimeout:    30 * time.Second,
		Encrypt:         true,
		TrustServerCert: true,
		AppName:         "vine-automation",
		EnableRetry:     true,
		RetryAttempts:   3,
		RetryInterval:   2 * time.Second,
	}
}

// DSN generates SQL Server connection string
func (c *Config) DSN() string {
	// Format: sqlserver://user:password@host:port?database=dbname&param=value
	host := c.Host
	if c.Instance != "" {
		host = fmt.Sprintf("%s\\%s", c.Host, c.Instance)
	}

	dsn := fmt.Sprintf(
		"sqlserver://%s:%s@%s:%d?database=%s&connection+timeout=%d&encrypt=%s&TrustServerCertificate=%s&app+name=%s",
		c.User,
		c.Password,
		host,
		c.Port,
		c.Database,
		int(c.ConnectTimeout.Seconds()),
		boolToString(c.Encrypt),
		boolToString(c.TrustServerCert),
		c.AppName,
	)

	return dsn
}

// Validate checks if configuration is valid
func (c *Config) Validate() error {
	if c.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("invalid database port: %d", c.Port)
	}
	if c.Database == "" {
		return fmt.Errorf("database name is required")
	}
	if c.User == "" {
		return fmt.Errorf("database user is required")
	}
	if c.MaxOpenConns < c.MaxIdleConns {
		return fmt.Errorf("max open connections must be >= max idle connections")
	}
	return nil
}

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
