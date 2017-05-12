// Code generated by hero.
// source: /Users/Lime/Documents/workspace/GoProject/src/github.com/shiyanhui/hero/examples/app/template/userlistbytearray.html
// DO NOT EDIT!
package template

import (
	"bytes"

	"github.com/shiyanhui/hero"
)

func UserListReturnsByteArray(userList []string) []byte {
	_buffer := new(bytes.Buffer)
	_buffer.WriteString(`<!DOCTYPE html>
<html>
    <head>
        <meta charset="utf-8">
    </head>

    <body>
        `)
	for _, user := range userList {
		_buffer.WriteString(`
        <ul>
            `)
		_buffer.WriteString(`<li>
    `)
		hero.EscapeHTML(user, _buffer)
		_buffer.WriteString(`
</li>
`)

		_buffer.WriteString(`
        </ul>
    `)
	}

	_buffer.WriteString(`
    </body>
</html>
`)
	return _buffer.Bytes()
}
