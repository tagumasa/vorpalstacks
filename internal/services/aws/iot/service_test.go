package iot

import (
	"sync"
	"testing"
)

func TestIoTServiceInitIdempotent(t *testing.T) {
	s := NewIoTService("test-account")
	deps := IoTServiceDeps{}

	s.Init(deps)
	s.Init(deps)

	if !s.initialised {
		t.Error("expected initialised=true after Init")
	}
}

func TestIoTServiceInitConcurrent(t *testing.T) {
	s := NewIoTService("test-account")
	deps := IoTServiceDeps{
		StorageManager: nil,
	}

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.Init(deps)
		}()
	}
	wg.Wait()

	if !s.initialised {
		t.Error("expected initialised=true after concurrent Init")
	}
}
