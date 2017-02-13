package hero

import (
	"testing"
)

func TestExecCommand(t *testing.T) {
	// test for panic
	execCommand("")
	execCommand("ls")
}
