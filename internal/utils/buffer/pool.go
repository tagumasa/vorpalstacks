// Package buffer provides utilities for managing byte buffers.
package buffer

import (
	"bytes"
	"sync"
)

const (
	smallSize  = 1024
	mediumSize = 4096
	largeSize  = 16384
)

// Pool is a pool of byte buffers with different size tiers.
type Pool struct {
	small  *sync.Pool
	medium *sync.Pool
	large  *sync.Pool
}

// NewPool creates a new buffer pool with three size tiers.
func NewPool() *Pool {
	return &Pool{
		small: &sync.Pool{
			New: func() interface{} {
				return bytes.NewBuffer(make([]byte, 0, smallSize))
			},
		},
		medium: &sync.Pool{
			New: func() interface{} {
				return bytes.NewBuffer(make([]byte, 0, mediumSize))
			},
		},
		large: &sync.Pool{
			New: func() interface{} {
				return bytes.NewBuffer(make([]byte, 0, largeSize))
			},
		},
	}
}

// Get retrieves a buffer from the pool suitable for the given size.
func (p *Pool) Get(size int) *bytes.Buffer {
	switch {
	case size <= smallSize:
		return p.small.Get().(*bytes.Buffer)
	case size <= mediumSize:
		return p.medium.Get().(*bytes.Buffer)
	default:
		return p.large.Get().(*bytes.Buffer)
	}
}

// Put returns a buffer to the pool.
func (p *Pool) Put(buf *bytes.Buffer) {
	if buf == nil {
		return
	}

	buf.Reset()

	switch buf.Cap() {
	case smallSize:
		p.small.Put(buf)
	case mediumSize:
		p.medium.Put(buf)
	case largeSize:
		p.large.Put(buf)
	}
}

// GlobalPool is a global buffer pool instance.
var GlobalPool = NewPool()
