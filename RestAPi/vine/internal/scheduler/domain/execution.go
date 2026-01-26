package domain

import "time"

// JobExecution represents job execution history (job_execution_history table)
type JobExecution struct {
	ID               int64           // execution_id
	JobName          string          // job_name
	JobGroup         string          // job_group
	TriggerName      string          // trigger_name
	TriggerGroup     string          // trigger_group
	StartTime        time.Time       // start_time
	EndTime          *time.Time      // end_time
	Status           ExecutionStatus // status
	ErrorMessage     string          // error_message
	RecordsProcessed int64           // records_processed
	DurationMs       int64           // duration_ms
	ProcessFromTime  *time.Time      // process_from_time
	ProcessToTime    *time.Time      // process_to_time
}

// ExecutionStatus represents job execution status
type ExecutionStatus string

const (
	ExecutionStatusStarted     ExecutionStatus = "STARTED"
	ExecutionStatusCompleted   ExecutionStatus = "COMPLETED"
	ExecutionStatusFailed      ExecutionStatus = "FAILED"
	ExecutionStatusInterrupted ExecutionStatus = "INTERRUPTED"
)

// IsCompleted checks if execution completed successfully
func (e *JobExecution) IsCompleted() bool {
	return e.Status == ExecutionStatusCompleted
}

// IsFailed checks if execution failed
func (e *JobExecution) IsFailed() bool {
	return e.Status == ExecutionStatusFailed
}

// Duration returns execution duration
func (e *JobExecution) Duration() time.Duration {
	if e.EndTime == nil {
		return time.Since(e.StartTime)
	}
	return e.EndTime.Sub(e.StartTime)
}
