package hero

import (
	"bytes"
	"fmt"
	"html"
	"os"
	"os/exec"
	"strings"
	"unsafe"
)

// Convert converts value to string and write it to the buffer. If escaped is
// true, the strings will be escaped.
func Convert(v interface{}, escaped bool, buffer *bytes.Buffer) {
	var content string

	switch value := v.(type) {
	case string:
		content = value
	case []byte:
		content = *(*string)(unsafe.Pointer(&value))
	default:
		content = fmt.Sprintf("%v", v)
	}

	if escaped {
		content = html.EscapeString(content)
	}

	buffer.WriteString(content)
}

func execCommand(command string) {
	parts := strings.Split(command, " ")
	if len(parts) == 0 {
		return
	}

	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Run()
}
