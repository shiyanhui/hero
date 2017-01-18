package hero

import (
	"testing"
)

func TestExecCommand(t *testing.T) {
	// test for whether panic
	execCommand("")
	execCommand("ls")
}
