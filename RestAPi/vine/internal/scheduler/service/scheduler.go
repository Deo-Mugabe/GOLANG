package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Deo-Mugabe/GOLANG/internal/scheduler/domain"
	"github.com/rs/zerolog"
)

// SchedulerService manages job scheduling
type SchedulerService struct {
	configRepo    domain.ConfigRepository
	executionRepo domain.ExecutionRepository
	executor      *ExecutorService
	logger        zerolog.Logger
}

// NewSchedulerService creates a new scheduler service
func NewSchedulerService(
	configRepo domain.ConfigRepository,
	executionRepo domain.ExecutionRepository,
	executor *ExecutorService,
	logger zerolog.Logger,
) *SchedulerService {
	return &SchedulerService{
		configRepo:    configRepo,
		executionRepo: executionRepo,
		executor:      executor,
		logger:        logger,
	}
}

// Start starts the scheduler
func (s *SchedulerService) Start(ctx context.Context, configName string, intervalMin int) error {
	s.logger.Info().
		Str("config", configName).
		Int("interval_minutes", intervalMin).
		Msg("starting scheduler")

	// Validate interval
	if intervalMin < 1 {
		return fmt.Errorf("interval must be at least 1 minute")
	}

	// Get or create config
	config, err := s.configRepo.Get(ctx, configName)
	if err != nil {
		// Config doesn't exist, create it
		config = &domain.SchedulerConfig{
			ConfigName:    configName,
			Enabled:       false,
			IntervalMin:   intervalMin,
			StartFromTime: time.Now().AddDate(0, 0, -30), // Default: 30 days ago
			CreatedAt:     time.Now(),
		}
		if err := s.configRepo.Create(ctx, config); err != nil {
			return fmt.Errorf("failed to create config: %w", err)
		}
	}

	// Check if already running
	if config.IsRunning() {
		return fmt.Errorf("scheduler is already running")
	}

	// Update config
	now := time.Now()
	config.Enabled = true
	config.IntervalMin = intervalMin
	config.LastStartTime = &now
	config.LastStopTime = nil
	nextRun := now.Add(time.Duration(intervalMin) * time.Minute)
	config.NextRunTime = &nextRun

	if err := s.configRepo.Update(ctx, config); err != nil {
		return fmt.Errorf("failed to update config: %w", err)
	}

	s.logger.Info().
		Str("config", configName).
		Time("next_run", nextRun).
		Msg("scheduler started")

	return nil
}

// Stop stops the scheduler
func (s *SchedulerService) Stop(ctx context.Context, configName string) error {
	s.logger.Info().Str("config", configName).Msg("stopping scheduler")

	config, err := s.configRepo.Get(ctx, configName)
	if err != nil {
		return fmt.Errorf("config not found: %w", err)
	}

	if !config.IsRunning() {
		return fmt.Errorf("scheduler is not running")
	}

	// Update config
	now := time.Now()
	config.Enabled = false
	config.LastStopTime = &now
	config.NextRunTime = nil

	if err := s.configRepo.Update(ctx, config); err != nil {
		return fmt.Errorf("failed to update config: %w", err)
	}

	s.logger.Info().Str("config", configName).Msg("scheduler stopped")
	return nil
}

// Restart restarts the scheduler
func (s *SchedulerService) Restart(ctx context.Context, configName string) error {
	config, err := s.configRepo.Get(ctx, configName)
	if err != nil {
		return err
	}

	if err := s.Stop(ctx, configName); err != nil {
		return err
	}

	return s.Start(ctx, configName, config.IntervalMin)
}

// GetStatus returns scheduler status
func (s *SchedulerService) GetStatus(ctx context.Context, configName string) (*domain.SchedulerStatus, error) {
	config, err := s.configRepo.Get(ctx, configName)
	if err != nil {
		return nil, err
	}

	// Get execution statistics
	totalExecs, _ := s.executionRepo.CountByStatus(ctx, domain.ExecutionStatusStarted)
	successExecs, _ := s.executionRepo.CountByStatus(ctx, domain.ExecutionStatusCompleted)
	failedExecs, _ := s.executionRepo.CountByStatus(ctx, domain.ExecutionStatusFailed)

	// Get last error
	lastError := ""
	failedExecutions, err := s.executionRepo.ListByStatus(ctx, domain.ExecutionStatusFailed, 1)
	if err == nil && len(failedExecutions) > 0 {
		lastError = failedExecutions[0].ErrorMessage
	}

	status := &domain.SchedulerStatus{
		ConfigName:           config.ConfigName,
		Enabled:              config.Enabled,
		Running:              config.IsRunning(),
		IntervalMinutes:      config.IntervalMin,
		LastStartTime:        config.LastStartTime,
		LastStopTime:         config.LastStopTime,
		LastRunTime:          config.LastRunTime,
		NextRunTime:          config.NextRunTime,
		TotalExecutions:      totalExecs + successExecs + failedExecs,
		SuccessfulExecutions: successExecs,
		FailedExecutions:     failedExecs,
		LastError:            lastError,
	}

	return status, nil
}

// IsRunning checks if scheduler is running
func (s *SchedulerService) IsRunning(ctx context.Context, configName string) (bool, error) {
	config, err := s.configRepo.Get(ctx, configName)
	if err != nil {
		return false, err
	}
	return config.IsRunning(), nil
}

// GetNextRunTime returns next scheduled run time
func (s *SchedulerService) GetNextRunTime(ctx context.Context, configName string) (*time.Time, error) {
	config, err := s.configRepo.Get(ctx, configName)
	if err != nil {
		return nil, err
	}
	return config.NextRunTime, nil
}

// UpdateInterval updates scheduler interval
func (s *SchedulerService) UpdateInterval(ctx context.Context, configName string, intervalMin int) error {
	if intervalMin < 1 {
		return fmt.Errorf("interval must be at least 1 minute")
	}

	return s.configRepo.UpdateInterval(ctx, configName, intervalMin)
}

// UpdateStartFromTime updates the processing start time
func (s *SchedulerService) UpdateStartFromTime(ctx context.Context, configName string, t time.Time) error {
	config, err := s.configRepo.Get(ctx, configName)
	if err != nil {
		return err
	}

	config.StartFromTime = t
	return s.configRepo.Update(ctx, config)
}

// TriggerNow manually triggers job execution
func (s *SchedulerService) TriggerNow(ctx context.Context, configName string) error {
	s.logger.Info().Str("config", configName).Msg("manually triggering job")

	config, err := s.configRepo.Get(ctx, configName)
	if err != nil {
		return err
	}

	// Execute the job
	jobName := "bookingProcessorJob"
	jobGroup := "vine-group"

	if err := s.executor.Execute(ctx, jobName, jobGroup); err != nil {
		return fmt.Errorf("failed to execute job: %w", err)
	}

	// Update last run time
	now := time.Now()
	config.LastRunTime = &now
	if config.Enabled {
		nextRun := now.Add(time.Duration(config.IntervalMin) * time.Minute)
		config.NextRunTime = &nextRun
	}

	if err := s.configRepo.Update(ctx, config); err != nil {
		return fmt.Errorf("failed to update config: %w", err)
	}

	return nil
}
