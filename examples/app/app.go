package main

import (
	"bytes"
	"log"
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

		if _, err := w.Write(buffer.Bytes()); err != nil {
			log.Printf("ERR: %s\n", err)
		}
	})

	http.HandleFunc("/users2", func(w http.ResponseWriter, req *http.Request) {
		var userList = []string{
			"Alice",
			"Bob",
			"Tom",
		}

		// using an io.Writer for automatic buffer management (i.e. hero built-in buffer pool)
		template.UserListToWriter(userList, w)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
