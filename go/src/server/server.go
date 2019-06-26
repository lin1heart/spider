package server

import (
	"io"
	"log"
	"net/http"
	"fmt"
)

func StartServer () {
	fmt.Printf("start server")

	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Hello, world!\n")
	}
	http.HandleFunc("/hello", helloHandler)

	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to my website!")
	})

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.ListenAndServe(":80", nil)
	log.Fatal(http.ListenAndServe(":8080", nil))
}