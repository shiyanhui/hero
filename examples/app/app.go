package main

import (
	"bytes"
	"net/http"

	"github.com/shiyanhui/hero/examples/app/template"
)

func main() {
	http.HandleFunc("/users", func(w http.ResponseWriter, req *http.Request) {
		var userList = []string{
			"Alice",
			"Bob",
			"Tom",
		}

		// Had better use buffer sync.Pool.
		// Hero exports GetBuffer and PutBuffer for this.
		//
		// buffer := hero.GetBuffer()
		// defer hero.PutBuffer(buffer)
		buffer := new(bytes.Buffer)
		template.UserList(userList, buffer)

		w.Write(buffer.Bytes())
	})

	http.ListenAndServe(":8080", nil)
}
