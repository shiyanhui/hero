package hero

import (
	"bytes"
	"sync"
)

const buffSize = 10000

var defaultPool *pool

func init() {
	defaultPool = newPool()
}

type pool struct {
	pool *sync.Pool
	ch   chan *bytes.Buffer
}

func newPool() *pool {
	p := &pool{
		pool: new(sync.Pool),
		ch:   make(chan *bytes.Buffer, buffSize),
	}

	for i := 0; i < buffSize; i++ {
		p.ch <- new(bytes.Buffer)
	}

	return p
}

// GetBuffer returns a *bytes.Buffer from sync.Pool.
func GetBuffer() (buffer *bytes.Buffer) {
	v := defaultPool.pool.Get()
	if v == nil {
		buffer = new(bytes.Buffer)
	} else {
		buffer = v.(*bytes.Buffer)
	}
	return
}

// PutBuffer puts a *bytes.Buffer to the sync.Pool.
func PutBuffer(buffer *bytes.Buffer) {
	buffer.Reset()
	defaultPool.pool.Put(buffer)
}
