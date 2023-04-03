package main

import (
	"fmt"
)

type Hub struct {

	// Registered clients KV(name -> client).
	clients map[string]*Client

	// Register requests from clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
	
	// Inbound messages from clients.
	unicast chan *Message
}

func newHub() *Hub {
	return &Hub{
		clients: make(map[string]*Client),
		register: make(chan *Client),
		unregister: make(chan *Client),
		unicast: make(chan *Message),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.email] = client
		case client := <-h.unregister:
			if _, ok := h.clients[client.email]; ok {
				delete(h.clients, client.email)
				close(client.send)
			}
		case m := <-h.unicast:
			if client, ok := h.clients[m.from]; ok {
				select {
				case client.send <- m.message:
				default:
					close(client.send)
					delete(h.clients, client.email)
				}
			}
		}
	}
}