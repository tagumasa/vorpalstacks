package buffer

import (
	"bytes"
	"sync"
	"testing"
)

func TestPoolGetPut(t *testing.T) {
	pool := NewPool()

	buf := pool.Get(512)
	if buf == nil {
		t.Fatal("expected non-nil buffer")
	}
	if buf.Cap() != smallSize {
		t.Errorf("expected small buffer cap %d, got %d", smallSize, buf.Cap())
	}
	pool.Put(buf)

	buf = pool.Get(2048)
	if buf == nil {
		t.Fatal("expected non-nil buffer")
	}
	if buf.Cap() != mediumSize {
		t.Errorf("expected medium buffer cap %d, got %d", mediumSize, buf.Cap())
	}
	pool.Put(buf)

	buf = pool.Get(8192)
	if buf == nil {
		t.Fatal("expected non-nil buffer")
	}
	if buf.Cap() != largeSize {
		t.Errorf("expected large buffer cap %d, got %d", largeSize, buf.Cap())
	}
	pool.Put(buf)
}

func TestPoolReset(t *testing.T) {
	pool := NewPool()

	buf := pool.Get(512)
	buf.WriteString("test data")
	if buf.Len() != 9 {
		t.Errorf("expected len 9, got %d", buf.Len())
	}

	pool.Put(buf)

	if buf.Len() != 0 {
		t.Errorf("expected buffer to be reset, got len %d", buf.Len())
	}
}

func TestPoolPutNil(t *testing.T) {
	pool := NewPool()
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Put(nil) panicked: %v", r)
		}
	}()
	pool.Put(nil)
}

func TestGlobalPool(t *testing.T) {
	buf := GlobalPool.Get(1024)
	if buf == nil {
		t.Fatal("expected non-nil buffer")
	}

	buf.WriteString("hello")
	GlobalPool.Put(buf)
}

func TestPoolConcurrent(t *testing.T) {
	pool := NewPool()
	var wg sync.WaitGroup
	const goroutines = 100

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			buf := pool.Get(1024)
			if buf == nil {
				t.Error("expected non-nil buffer from pool")
				return
			}
			buf.WriteString("test")
			if buf.Len() != 4 {
				t.Errorf("expected buffer len 4, got %d", buf.Len())
				return
			}
			pool.Put(buf)
		}()
	}

	wg.Wait()
}

func BenchmarkPoolGetPut(b *testing.B) {
	pool := NewPool()
	for i := 0; i < b.N; i++ {
		buf := pool.Get(1024)
		buf.WriteString("test data")
		pool.Put(buf)
	}
}

func BenchmarkBytesNewBuffer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := bytes.NewBuffer(make([]byte, 0, 1024))
		buf.WriteString("test data")
		_ = buf
	}
}
