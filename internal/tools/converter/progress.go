// Package converter provides tools for converting Smithy models.
package converter

import (
	"fmt"
	"sync"
	"time"
)

// ProgressTracker tracks and displays progress for service processing
type ProgressTracker struct {
	total     int
	completed int
	current   map[string]bool // Currently processing services
	startTime time.Time
	mu        sync.Mutex
}

// NewProgressTracker creates a new ProgressTracker
func NewProgressTracker() *ProgressTracker {
	return &ProgressTracker{
		current: make(map[string]bool),
	}
}

// SetTotal sets the total number of services to process
func (p *ProgressTracker) SetTotal(total int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.total = total
	p.startTime = time.Now()
}

// Start marks a service as currently being processed
func (p *ProgressTracker) Start(serviceName string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.current[serviceName] = true
	p.print()
}

// Complete marks a service as completed
func (p *ProgressTracker) Complete(serviceName string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.current, serviceName)
	p.completed++
	p.print()
}

// print prints the current progress
func (p *ProgressTracker) print() {
	if p.total == 0 {
		return
	}

	// Build current services list
	var currentList []string
	for name := range p.current {
		currentList = append(currentList, name)
	}

	// Calculate progress percentage
	percentage := float64(p.completed) / float64(p.total) * 100

	// Calculate elapsed time
	elapsed := time.Since(p.startTime)
	elapsedStr := elapsed.Round(time.Second).String()

	// Print progress
	if len(currentList) > 0 {
		fmt.Printf("\r[%d/%d] (%.1f%%) Processing: %s | Elapsed: %s",
			p.completed, p.total, percentage, currentList[0], elapsedStr)
		if len(currentList) > 1 {
			fmt.Printf(" (+%d more)", len(currentList)-1)
		}
	} else {
		fmt.Printf("\r[%d/%d] (%.1f%%) | Elapsed: %s",
			p.completed, p.total, percentage, elapsedStr)
	}
}

// Done prints the final completion message
func (p *ProgressTracker) Done() {
	p.mu.Lock()
	defer p.mu.Unlock()

	elapsed := time.Since(p.startTime)
	elapsedStr := elapsed.Round(time.Second).String()

	fmt.Printf("\n✓ All %d services processed successfully in %s\n", p.completed, elapsedStr)
}

// Error prints an error message
func (p *ProgressTracker) Error(serviceName string, err error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	fmt.Printf("\n✗ Error processing service '%s': %v\n", serviceName, err)
}
