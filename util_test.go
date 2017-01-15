package hero

import (
	"bytes"
	"testing"
)

type dummy struct{}

func (d dummy) String() string {
	return "this is dummy"
}

func TestConvert(t *testing.T) {
	cases := []struct {
		v       interface{}
		escaped bool
		out     string
	}{
		{100, true, "100"},
		{1.23, true, "1.23"},
		{"<div>hello</div>", false, "<div>hello</div>"},
		{"<div>hello</div>", true, "&lt;div&gt;hello&lt;/div&gt;"},
		{dummy{}, true, "this is dummy"},
	}

	for _, c := range cases {
		buffer := new(bytes.Buffer)
		if Convert(c.v, c.escaped, buffer); buffer.String() != c.out {
			t.Fail()
		}
	}
}

func TestExecCommand(t *testing.T) {
	// test for whether panic
	execCommand("")
	execCommand("ls")
}
