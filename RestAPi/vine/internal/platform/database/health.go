package database

import (
	"context"
	"fmt"
	"time"
)

// HealthStatus represents database health status
type HealthStatus struct {
	Healthy         bool          `json:"healthy"`
	Message         string        `json:"message,omitempty"`
	ResponseTime    time.Duration `json:"response_time"`
	OpenConnections int           `json:"open_connections"`
	IdleConnections int           `json:"idle_connections"`
	MaxOpenConns    int           `json:"max_open_connections"`
	Timestamp       time.Time     `json:"timestamp"`
}

// HealthCheck performs a comprehensive health check
func (db *DB) HealthCheck(ctx context.Context) *HealthStatus {
	start := time.Now()

	status := &HealthStatus{
		Healthy:   false,
		Timestamp: start,
	}

	// Check database connectivity
	if err := db.Ping(ctx); err != nil {
		status.Message = fmt.Sprintf("ping failed: %v", err)
		status.ResponseTime = time.Since(start)
		return status
	}

	// Get connection pool stats
	stats := db.Stats()
	status.OpenConnections = stats.OpenConnections
	status.IdleConnections = stats.Idle
	status.MaxOpenConns = db.config.MaxOpenConns
	status.ResponseTime = time.Since(start)

	// Check if we can execute a simple query
	var result int
	err := db.DB.GetContext(ctx, &result, "SELECT 1")
	if err != nil {
		status.Message = fmt.Sprintf("query failed: %v", err)
		return status
	}

	// All checks passed
	status.Healthy = true
	status.Message = "database is healthy"

	return status
}

// IsHealthy performs a quick health check
func (db *DB) IsHealthy(ctx context.Context) bool {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	return db.Ping(ctx) == nil
}

// WaitForHealthy waits for database to become healthy
func (db *DB) WaitForHealthy(ctx context.Context, maxWait time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, maxWait)
	defer cancel()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for database to become healthy")
		case <-ticker.C:
			if db.IsHealthy(ctx) {
				return nil
			}
		}
	}
}
