package main

import (
	"fmt"
	"net/http"
)

func main() {
	hub := newHub()
	go hub.run()
	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.ListenAndServe("localhost:8080", nil)
	fmt.Print("Started server and listening at 8080")
}
