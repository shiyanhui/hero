package hero

import (
	"os"
	"os/exec"
	"strings"
)

// execCommand wraps exec.Command
func execCommand(command string) {
	parts := strings.Split(command, " ")
	if len(parts) == 0 {
		return
	}

	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Run()
}
