package domain

import "time"

// SchedulerConfig represents scheduler configuration (scheduler_config table)
type SchedulerConfig struct {
	ConfigName    string     // config_name (PK)
	Enabled       bool       // enabled
	IntervalMin   int        // interval_minutes
	LastStartTime *time.Time // last_start_time
	LastStopTime  *time.Time // last_stop_time
	LastRunTime   *time.Time // last_run_time
	NextRunTime   *time.Time // next_run_time
	StartFromTime time.Time  // start_from_time
	CreatedAt     time.Time  // created_at
	UpdatedAt     *time.Time // updated_at
}

// IsEnabled checks if scheduler is enabled
func (s *SchedulerConfig) IsEnabled() bool {
	return s.Enabled
}

// IsRunning checks if scheduler is currently running
func (s *SchedulerConfig) IsRunning() bool {
	return s.Enabled && s.LastStartTime != nil &&
		(s.LastStopTime == nil || s.LastStartTime.After(*s.LastStopTime))
}

// ShouldRun checks if job should run now
func (s *SchedulerConfig) ShouldRun() bool {
	if !s.Enabled {
		return false
	}
	if s.NextRunTime == nil {
		return false
	}
	return time.Now().After(*s.NextRunTime)
}
