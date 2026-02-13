package service

import (
	"context"
	"time"

	"github.com/Deo-Mugabe/GOLANG/internal/scheduler/domain"
	"github.com/rs/zerolog"
)

// HistoryService manages job execution history
type HistoryService struct {
	executionRepo domain.ExecutionRepository
	logger        zerolog.Logger
}

// NewHistoryService creates a new history service
func NewHistoryService(
	executionRepo domain.ExecutionRepository,
	logger zerolog.Logger,
) *HistoryService {
	return &HistoryService{
		executionRepo: executionRepo,
		logger:        logger,
	}
}

// GetExecutionHistory retrieves execution history for a job
func (s *HistoryService) GetExecutionHistory(
	ctx context.Context,
	jobName, jobGroup string,
	limit, offset int,
) ([]*domain.JobExecution, error) {
	executions, err := s.executionRepo.List(ctx, jobName, jobGroup, limit, offset)
	if err != nil {
		return nil, err
	}

	s.logger.Debug().
		Str("job_name", jobName).
		Int("count", len(executions)).
		Msg("retrieved execution history")

	return executions, nil
}

// GetLatestExecutions retrieves latest executions across all jobs
func (s *HistoryService) GetLatestExecutions(ctx context.Context, limit int) ([]*domain.JobExecution, error) {
	since := time.Now().AddDate(0, 0, -30) // Last 30 days
	executions, err := s.executionRepo.ListSince(ctx, since, limit, 0)
	if err != nil {
		return nil, err
	}

	s.logger.Debug().
		Int("count", len(executions)).
		Msg("retrieved latest executions")

	return executions, nil
}

// GetFailedExecutions retrieves failed executions
func (s *HistoryService) GetFailedExecutions(ctx context.Context, limit int) ([]*domain.JobExecution, error) {
	executions, err := s.executionRepo.ListByStatus(ctx, domain.ExecutionStatusFailed, limit)
	if err != nil {
		return nil, err
	}

	s.logger.Debug().
		Int("count", len(executions)).
		Msg("retrieved failed executions")

	return executions, nil
}

// GetExecutionStats retrieves execution statistics for a job
func (s *HistoryService) GetExecutionStats(
	ctx context.Context,
	jobName, jobGroup string,
) (*domain.ExecutionStats, error) {
	// Get total executions count
	totalCount, err := s.executionRepo.CountByStatus(ctx, domain.ExecutionStatusStarted)
	if err != nil {
		return nil, err
	}

	// Get successful count
	successCount, err := s.executionRepo.CountSuccessful(ctx, jobName, jobGroup)
	if err != nil {
		return nil, err
	}

	// Get failed count
	failedCount, err := s.executionRepo.CountFailed(ctx, jobName, jobGroup)
	if err != nil {
		return nil, err
	}

	// Get recent executions for calculating average duration
	executions, err := s.executionRepo.List(ctx, jobName, jobGroup, 100, 0)
	if err != nil {
		return nil, err
	}

	var totalDuration int64
	var totalRecords int64
	var lastSuccessful *time.Time
	var lastFailed *time.Time

	for _, exec := range executions {
		totalDuration += exec.DurationMs
		totalRecords += exec.RecordsProcessed

		if exec.IsCompleted() && (lastSuccessful == nil || exec.EndTime.After(*lastSuccessful)) {
			lastSuccessful = exec.EndTime
		}

		if exec.IsFailed() && (lastFailed == nil || exec.EndTime.After(*lastFailed)) {
			lastFailed = exec.EndTime
		}
	}

	avgDuration := int64(0)
	if len(executions) > 0 {
		avgDuration = totalDuration / int64(len(executions))
	}

	stats := &domain.ExecutionStats{
		TotalExecutions:       totalCount + successCount + failedCount,
		SuccessfulExecutions:  successCount,
		FailedExecutions:      failedCount,
		AverageDurationMs:     avgDuration,
		TotalRecordsProcessed: totalRecords,
		LastSuccessfulRun:     lastSuccessful,
		LastFailedRun:         lastFailed,
	}

	s.logger.Debug().
		Str("job_name", jobName).
		Int64("total", stats.TotalExecutions).
		Int64("successful", stats.SuccessfulExecutions).
		Int64("failed", stats.FailedExecutions).
		Msg("calculated execution stats")

	return stats, nil
}

// GetSuccessRate calculates success rate for a time period
func (s *HistoryService) GetSuccessRate(
	ctx context.Context,
	jobName, jobGroup string,
	since time.Time,
) (float64, error) {
	// Get executions since the specified time
	executions, err := s.executionRepo.ListSince(ctx, since, 1000, 0)
	if err != nil {
		return 0, err
	}

	if len(executions) == 0 {
		return 0, nil
	}

	successCount := 0
	for _, exec := range executions {
		if exec.JobName == jobName && exec.JobGroup == jobGroup && exec.IsCompleted() {
			successCount++
		}
	}

	successRate := float64(successCount) / float64(len(executions)) * 100

	s.logger.Debug().
		Str("job_name", jobName).
		Float64("success_rate", successRate).
		Int("total", len(executions)).
		Int("successful", successCount).
		Msg("calculated success rate")

	return successRate, nil
}

// CleanupOldHistory deletes old execution history
func (s *HistoryService) CleanupOldHistory(ctx context.Context, olderThan time.Duration) (int64, error) {
	cutoff := time.Now().Add(-olderThan)

	deleted, err := s.executionRepo.DeleteOlderThan(ctx, cutoff)
	if err != nil {
		return 0, err
	}

	s.logger.Info().
		Time("cutoff", cutoff).
		Int64("deleted", deleted).
		Msg("cleaned up old execution history")

	return deleted, nil
}
