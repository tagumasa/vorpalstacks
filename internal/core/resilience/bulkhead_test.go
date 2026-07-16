package resilience

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestBulkhead_Execute_Success(t *testing.T) {
	config := &BulkheadConfig{
		Name:          "test",
		MaxConcurrent: 2,
		MaxWait:       100 * time.Millisecond,
	}
	bh := NewBulkhead(config)

	err := bh.Execute(context.Background(), func() error {
		return nil
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if bh.ActiveCount() != 0 {
		t.Fatalf("expected active count 0, got %d", bh.ActiveCount())
	}
}

func TestBulkhead_Execute_BlocksWhenFull(t *testing.T) {
	config := &BulkheadConfig{
		Name:          "test",
		MaxConcurrent: 1,
		MaxWait:       50 * time.Millisecond,
	}
	bh := NewBulkhead(config)

	var wg sync.WaitGroup
	wg.Add(2)

	errCh := make(chan error, 2)

	go func() {
		defer wg.Done()
		err := bh.Execute(context.Background(), func() error {
			time.Sleep(100 * time.Millisecond)
			return nil
		})
		errCh <- err
	}()

	time.Sleep(10 * time.Millisecond)

	go func() {
		defer wg.Done()
		err := bh.Execute(context.Background(), func() error {
			return nil
		})
		errCh <- err
	}()

	wg.Wait()
	close(errCh)

	errors := 0
	for err := range errCh {
		if err != nil {
			errors++
		}
	}

	if errors != 1 {
		t.Fatalf("expected 1 error (rate limit), got %d", errors)
	}
}

func TestBulkhead_Execute_ContextCancel(t *testing.T) {
	config := &BulkheadConfig{
		Name:          "test",
		MaxConcurrent: 1,
		MaxWait:       1 * time.Second,
	}
	bh := NewBulkhead(config)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := bh.Execute(ctx, func() error {
		return nil
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestBulkhead_TryExecute_Success(t *testing.T) {
	config := &BulkheadConfig{
		Name:          "test",
		MaxConcurrent: 2,
	}
	bh := NewBulkhead(config)

	err := bh.TryExecute(func() error {
		return nil
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBulkhead_TryExecute_NoSlotAvailable(t *testing.T) {
	config := &BulkheadConfig{
		Name:          "test",
		MaxConcurrent: 1,
	}
	bh := NewBulkhead(config)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		bh.Execute(context.Background(), func() error {
			time.Sleep(1000 * time.Millisecond)
			return nil
		})
	}()

	time.Sleep(500 * time.Millisecond)

	err2 := bh.TryExecute(func() error {
		return nil
	})

	if err2 == nil {
		t.Fatal("expected error when no slot available")
	}

	if !IsRateLimit(err2) {
		t.Fatalf("expected rate limit error, got %v", err2)
	}

	wg.Wait()
}

func TestBulkhead_ActiveCount(t *testing.T) {
	config := &BulkheadConfig{
		Name:          "test",
		MaxConcurrent: 5,
	}
	bh := NewBulkhead(config)

	if bh.ActiveCount() != 0 {
		t.Fatalf("expected 0, got %d", bh.ActiveCount())
	}

	var wg sync.WaitGroup
	wg.Add(3)

	for i := 0; i < 3; i++ {
		go func() {
			defer wg.Done()
			bh.Execute(context.Background(), func() error {
				time.Sleep(2000 * time.Millisecond)
				return nil
			})
		}()
	}

	time.Sleep(500 * time.Millisecond)

	count := bh.ActiveCount()
	if count != 3 {
		t.Fatalf("expected active count 3, got %d", count)
	}

	wg.Wait()

	time.Sleep(50 * time.Millisecond)

	if bh.ActiveCount() != 0 {
		t.Fatalf("expected active count 0 after completion, got %d", bh.ActiveCount())
	}
}

func TestBulkhead_AvailableSlots(t *testing.T) {
	config := &BulkheadConfig{
		Name:          "test",
		MaxConcurrent: 2,
	}
	bh := NewBulkhead(config)

	if bh.AvailableSlots() != 2 {
		t.Fatalf("expected 2 slots, got %d", bh.AvailableSlots())
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		bh.Execute(context.Background(), func() error {
			time.Sleep(500 * time.Millisecond)
			return nil
		})
	}()

	time.Sleep(200 * time.Millisecond)

	if bh.AvailableSlots() != 1 {
		t.Fatalf("expected 1 slot, got %d", bh.AvailableSlots())
	}

	wg.Wait()

	if bh.AvailableSlots() != 2 {
		t.Fatalf("expected 2 slots after completion, got %d", bh.AvailableSlots())
	}
}

func TestBulkhead_Stats(t *testing.T) {
	config := &BulkheadConfig{
		Name:          "test-bulkhead",
		MaxConcurrent: 10,
	}
	bh := NewBulkhead(config)

	stats := bh.Stats()

	if stats.Name != "test-bulkhead" {
		t.Errorf("expected name 'test-bulkhead', got '%s'", stats.Name)
	}
	if stats.MaxConcurrent != 10 {
		t.Errorf("expected MaxConcurrent 10, got %d", stats.MaxConcurrent)
	}
	if stats.AvailableSlots != 10 {
		t.Errorf("expected AvailableSlots 10, got %d", stats.AvailableSlots)
	}
}

func TestBulkhead_MaxConcurrent(t *testing.T) {
	config := &BulkheadConfig{
		MaxConcurrent: 5,
	}
	bh := NewBulkhead(config)

	if bh.MaxConcurrent() != 5 {
		t.Fatalf("expected MaxConcurrent 5, got %d", bh.MaxConcurrent())
	}
}

func TestBulkhead_ConcurrentLimit(t *testing.T) {
	config := &BulkheadConfig{
		Name:          "concurrent-test",
		MaxConcurrent: 2,
		MaxWait:       500 * time.Millisecond,
	}
	bh := NewBulkhead(config)

	var mu sync.Mutex
	activeCount := 0
	maxObserved := 0

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			bh.Execute(context.Background(), func() error {
				mu.Lock()
				activeCount++
				if activeCount > maxObserved {
					maxObserved = activeCount
				}
				mu.Unlock()

				time.Sleep(100 * time.Millisecond)

				mu.Lock()
				activeCount--
				mu.Unlock()
				return nil
			})
		}()
	}

	wg.Wait()

	if maxObserved > 2 {
		t.Fatalf("expected max 2 concurrent, observed %d", maxObserved)
	}
}

func TestBulkheadStats_String(t *testing.T) {
	s := BulkheadStats{
		Name:           "test",
		MaxConcurrent:  10,
		ActiveCount:    3,
		AvailableSlots: 7,
	}
	got := s.String()
	if !strings.Contains(got, "test") {
		t.Fatalf("BulkheadStats.String() should contain name, got %q", got)
	}
	if !strings.Contains(got, "MaxConcurrent: 10") {
		t.Fatalf("BulkheadStats.String() should contain MaxConcurrent: 10, got %q", got)
	}
	if !strings.Contains(got, "ActiveCount: 3") {
		t.Fatalf("BulkheadStats.String() should contain ActiveCount: 3, got %q", got)
	}
}
