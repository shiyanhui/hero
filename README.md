# Hero

Hero is a handy, fast and powerful go template engine, which pre-compiles the html templtes to go code.
It has been used in production environment in [bthub.io](http://bthub.io).

[![GoDoc](https://godoc.org/github.com/shiyanhui/hero?status.svg)](https://godoc.org/github.com/shiyanhui/hero)
[![Go Report Card](https://goreportcard.com/badge/github.com/shiyanhui/hero)](https://goreportcard.com/report/github.com/shiyanhui/hero)

- [Features](#features)
- [Install](#install)
- [Usage](#usage)
- [Quick Start](#quick-start)
- [Template Syntax](#template-syntax)
- [License](#license)

## Features

- Easy to use.
- Powerful. template `Extend` and `Include` supported.
- High performance.
- Auto compiling when files change.

## Install

    go get -u github.com/shiyanhui/hero

## Usage

```shell
hero [options]

options:
	- source:  the html template file or dir.
	- dest:    generated golang files dir, it will be the same with source if not set.
	- pkgname: the generated template package name, default is `template`.
	- watch:   whether automic compile when the source files change.

example:
	hero -source="./"
	hero -source="$GOPATH/src/app/template" -watch
```

## Quick Start

Assume that we are going to render a user list `userlist.html`. `index.html`
is the layout, and `user.html` is an item in the list.

And assumes that they are all under `$GOPATH/src/app/template`

### index.html

```html
<!DOCTYPE html>
<html>
    <head>
        <meta charset="utf-8">
    </head>

    <body>
        <%@ body { %>
        <% } %>
    </body>
</html>
```

### users.html

```html
<%: func UserList(userList []string) []byte %>

<%~ "index.html" %>

<%@ body { %>
    <% for _, user := range userList { %>
        <ul>
            <%+ "user.html" %>
        </ul>
    <% } %>
<% } %>
```

### user.html

```html
<li>
    <%= user %>
</li>
```

Then we compile the templates to go code.

```shell
hero -source="$GOPATH/src/app/template"
```

We will get three new `.go` files under `$GOPATH/src/app/template`,
i.e. `index.html.go`, `user.html.go` and `userlist.html.go`.

Then we write a http server in `$GOPATH/src/app/main.go`.

### main.go

```go
package main

import (
	"app/template"
	"net/http"
)

func main() {
	http.HandleFunc("/users", func(w http.ResponseWriter, req *http.Request) {
		var userList = []string {
          	"Alice",
			"Bob",
			"Tom",
		}
		w.Write(template.UserList(userList))
	})

	http.ListenAndServe(":8080", nil)
}
```

At last, start the server and visit `http://localhost:8080/users` in your browser, we will get what we want!

## Template syntax

There are only nine necessary kinds of statements, which are:

- Function Definition `<%: func define %>`
  - Function definition statement defines the function which represents a html file.
  - The function defined should return one and only one parameter `[]byte`.
  - Example:`<%: func UserList(userList []string) []byte %>` , which we have mentioned in quick start.

- Extend `<%~ "parent template" %>`
  - Extend statement states the parent template the current template extends.
  - The parent template should be quoted with `""`.
  - Example: `<%~ "index.html" >`, which we have mentioned in quick start, too.

- Include `<%+ "sub template" %>`
  - Include statement includes a sub-template to the current template. It works like `#include` in `C++`.
  - The sub-template should be quoted with `""`.
  - Example: `<%+ "user.html" >`, which we also have mentioned in quick start.

- Import `<%! go code %>`
  - Import statement imports the packages used in the defined function, and it also contains everything that is outside of the defined function.

  - Import statement will NOT be inherited by child template.

  - Example:

    ```go
    <%!
    	import (
          	"fmt"
        	"strings"
        )

    	var a int

    	const b = "hello, world"

    	func Add(a, b int) int {
        	return a + b
    	}

    	type S struct {
        	Name string
    	}

    	func (s S) String() string {
        	return s.Name
    	}
    %>
    ```

- Block `<%@ blockName { %> <% } %>`

  - Block statement represents a block. Child template overwrites blocks to extend parent template.

  - Example:

    ```html
    <!DOCTYPE html>
    <html>
        <head>
            <meta charset="utf-8">
        </head>

        <body>
            <%@ body { %>
            <% } %>
        </body>
    </html>
    ```

- Code `<% go code %>`

  - Code statement states all code inside the defined function. It's just go code.

  - Example:

    ```go
    <% for _, user := userList { %>
        <% if user != "Alice" { %>
        	<%= user %>
        <% } %>
    <% } %>

    <%
    	a, b := 1, 2
    	c := Add(a, b)
    %>
    ```

- Raw Value `<%= statement %>`

  - Raw Value statement will convert the variable to string.

  - Example:

    ```go
    <%== a %>
    <%== a + b %>
    <%== Add(a, b) %>
    <%== user.Name %>
    ```

- Escaped Value `<%== statement %>`

  - Escaped Value statement is similar with Raw Value statement, but after converting, it will escaped it with `html.EscapesString`.

  - Example:

    ```go
    <%= a %>
    <%= a + b %>
    <%= Add(a, b) %>
    <%= user.Name %>
    ```

- Note `<%# note %>`

  - Note statement add notes to the template.
  - It will not be added to the generated go source.
  - Example: `<# this is just a note example>`.

## License

Hero is licensed under the Apache License.
