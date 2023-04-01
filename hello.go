package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var (
	colonSpace = []byte(": ")
	newline = []byte("\n")
)

var clients []websocket.Conn

func main() {
	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity

		clients = append(clients, *conn)

		for {
			// Read message from browser
			_, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}

			// Print the message to the console
			fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

			for _, client := range clients {
				// Write message back to browser
				var connectionAddr = []byte(conn.RemoteAddr().String())

				w, err := client.NextWriter(websocket.TextMessage)
				if err != nil {
					return
				}
				w.Write(connectionAddr)
				w.Write(colonSpace)
				w.Write(msg)
				w.Write(newline)
				if err := w.Close(); err != nil {
					return
				}
				client.SetWriteDeadline(time.Now().Add(2 * time.Second))
			}

		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.ListenAndServe("localhost:8080", nil)
	fmt.Print("Started server and listening at 8080")
}
