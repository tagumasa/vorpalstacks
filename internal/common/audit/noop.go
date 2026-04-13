// Package audit provides AWS CloudTrail audit logging functionality for vorpalstacks.
package audit

// NoOpRecorder is a no-operation recorder that discards all audit events.
type NoOpRecorder struct{}

// NewNoOpRecorder creates a new no-operation recorder.
func NewNoOpRecorder() *NoOpRecorder {
	return &NoOpRecorder{}
}

// RecordEvent is a no-operation implementation that returns nil.
func (r *NoOpRecorder) RecordEvent(event *AuditEvent) error {
	return nil
}

// Record is a no-operation implementation that returns nil.
func (r *NoOpRecorder) Record(event interface{}) error {
	return nil
}
