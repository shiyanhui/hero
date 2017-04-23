package hero

import (
	"bytes"
	"runtime"
	"testing"
)

func TestExecCommand(t *testing.T) {
	// test for panic
	if runtime.GOOS != "windows" {
		execCommand("")
		execCommand("ls")
	} else {
		execCommand("dir")
	}
}

func TestFormatUint(t *testing.T) {
	cases := []struct {
		in  uint64
		out string
	}{
		{in: 0, out: "0"},
		{in: 1, out: "1"},
		{in: 100, out: "100"},
		{in: 101, out: "101"},
	}

	buffer := new(bytes.Buffer)

	for _, c := range cases {
		FormatUint(c.in, buffer)
		if buffer.String() != c.out {
			t.Fail()
		}
		buffer.Reset()
	}
}
