package hero

import (
	"bytes"
	"testing"
)

var buffer *bytes.Buffer

func init() {
	buffer = new(bytes.Buffer)
}

func TestNewPool(t *testing.T) {
	pool := newPool()
	if pool == nil || pool.pool == nil || len(pool.ch) != buffSize {
		t.Fail()
	}
}

func TestGetBuffer(t *testing.T) {
	if GetBuffer() == nil {
		t.Fail()
	}
}

func TestPutBubber(t *testing.T) {
	// test for panic
	PutBuffer(buffer)
}
