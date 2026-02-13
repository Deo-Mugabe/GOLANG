package service

import (
	"context"
	"fmt"
	"time"

	bookingdomain "github.com/Deo-Mugabe/GOLANG/internal/booking/domain"
	"github.com/Deo-Mugabe/GOLANG/internal/scheduler/domain"
	"github.com/rs/zerolog"
)

// ExecutorService handles job execution
type ExecutorService struct {
	configRepo    domain.ConfigRepository
	executionRepo domain.ExecutionRepository
	processor     bookingdomain.ProcessorService
	logger        zerolog.Logger
}

// NewExecutorService creates a new executor service
func NewExecutorService(
	configRepo domain.ConfigRepository,
	executionRepo domain.ExecutionRepository,
	processor bookingdomain.ProcessorService,
	logger zerolog.Logger,
) *ExecutorService {
	return &ExecutorService{
		configRepo:    configRepo,
		executionRepo: executionRepo,
		processor:     processor,
		logger:        logger,
	}
}

// Execute executes a job
func (s *ExecutorService) Execute(ctx context.Context, jobName, jobGroup string) error {
	s.logger.Info().
		Str("job_name", jobName).
		Str("job_group", jobGroup).
		Msg("executing job")

	// Record job start
	executionID, err := s.RecordStart(ctx, jobName, jobGroup)
	if err != nil {
		return fmt.Errorf("failed to record job start: %w", err)
	}

	// Get processing start time
	processFromTime, err := s.GetLastProcessingTime(ctx, jobName)
	if err != nil {
		s.logger.Warn().Err(err).Msg("failed to get last processing time, using default")
		processFromTime = time.Now().AddDate(0, 0, -30) // Default: 30 days ago
	}

	// Execute booking processing
	result, err := s.processor.ProcessBookings(ctx, processFromTime)
	if err != nil {
		recordErr := s.RecordFailure(ctx, executionID, err)
		if recordErr != nil {
			s.logger.Error().Err(recordErr).Msg("failed to record job failure")
		}
		return fmt.Errorf("job execution failed: %w", err)
	}

	// Record successful completion
	if err := s.RecordCompletion(ctx, executionID, result.TotalProcessed); err != nil {
		s.logger.Error().Err(err).Msg("failed to record job completion")
		return err
	}

	// Update processing time
	if err := s.UpdateProcessingTime(ctx, jobName, result.EndTime); err != nil {
		s.logger.Warn().Err(err).Msg("failed to update processing time")
	}

	s.logger.Info().
		Str("job_name", jobName).
		Int64("records_processed", result.TotalProcessed).
		Dur("duration", result.EndTime.Sub(result.StartTime)).
		Msg("job execution completed")

	return nil
}

// RecordStart records job execution start
func (s *ExecutorService) RecordStart(ctx context.Context, jobName, jobGroup string) (int64, error) {
	execution := &domain.JobExecution{
		JobName:      jobName,
		JobGroup:     jobGroup,
		TriggerName:  "manual",
		TriggerGroup: jobGroup,
		StartTime:    time.Now(),
		Status:       domain.ExecutionStatusStarted,
	}

	executionID, err := s.executionRepo.Create(ctx, execution)
	if err != nil {
		return 0, fmt.Errorf("failed to create execution record: %w", err)
	}

	s.logger.Debug().
		Int64("execution_id", executionID).
		Str("job_name", jobName).
		Msg("recorded job start")

	return executionID, nil
}

// RecordCompletion records successful job completion
func (s *ExecutorService) RecordCompletion(ctx context.Context, executionID int64, recordsProcessed int64) error {
	endTime := time.Now()

	if err := s.executionRepo.Complete(ctx, executionID, endTime, recordsProcessed); err != nil {
		return fmt.Errorf("failed to mark execution as complete: %w", err)
	}

	s.logger.Debug().
		Int64("execution_id", executionID).
		Int64("records_processed", recordsProcessed).
		Msg("recorded job completion")

	return nil
}

// RecordFailure records job execution failure
func (s *ExecutorService) RecordFailure(ctx context.Context, executionID int64, execErr error) error {
	endTime := time.Now()
	errorMsg := execErr.Error()

	if err := s.executionRepo.Fail(ctx, executionID, endTime, errorMsg); err != nil {
		return fmt.Errorf("failed to mark execution as failed: %w", err)
	}

	s.logger.Debug().
		Int64("execution_id", executionID).
		Str("error", errorMsg).
		Msg("recorded job failure")

	return nil
}

// GetLastProcessingTime retrieves the last successful processing time
func (s *ExecutorService) GetLastProcessingTime(ctx context.Context, jobName string) (time.Time, error) {
	// Try to get last successful execution
	lastExec, err := s.executionRepo.GetLastSuccessful(ctx, jobName, "vine-group")
	if err != nil {
		return time.Time{}, err
	}

	if lastExec != nil && lastExec.ProcessToTime != nil {
		s.logger.Debug().
			Time("last_process_time", *lastExec.ProcessToTime).
			Msg("using last successful execution time")
		return *lastExec.ProcessToTime, nil
	}

	// Fall back to config's start_from_time
	config, err := s.configRepo.Get(ctx, "booking-processor")
	if err == nil && !config.StartFromTime.IsZero() {
		s.logger.Debug().
			Time("start_from_time", config.StartFromTime).
			Msg("using configured start time")
		return config.StartFromTime, nil
	}

	// Default to 30 days ago
	defaultTime := time.Now().AddDate(0, 0, -30)
	s.logger.Debug().
		Time("default_time", defaultTime).
		Msg("using default start time (30 days ago)")
	return defaultTime, nil
}

// UpdateProcessingTime updates the last processing time
func (s *ExecutorService) UpdateProcessingTime(ctx context.Context, jobName string, t time.Time) error {
	config, err := s.configRepo.Get(ctx, "booking-processor")
	if err != nil {
		return err
	}

	config.LastRunTime = &t
	return s.configRepo.Update(ctx, config)
}
