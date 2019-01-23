package routes

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Hub is the websocket controller
type Hub struct {
	// Registered clients.
	clients map[*Client]bool
	// Inbound messages from the clients.
	broadcast chan []byte
	// Register requests from the clients.
	register chan *Client
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub
	// The websocket connection.
	conn *websocket.Conn
}

// NewHub comment
func NewHub() *Hub {
	h := &Hub{
		broadcast: make(chan []byte),
		register:  make(chan *Client),
		clients:   make(map[*Client]bool),
	}

	go func() {
		for {
			select {
			case client := <-h.register:
				h.clients[client] = true
			case msg := <-h.broadcast:
				for client := range h.clients {
					if err := client.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
						// Remove the client if there's an error writing
						delete(h.clients, client)
					}
				}
			}
		}
	}()

	return h
}

// WSHandler handles websocket requests from the peer.
func WSHandler(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn}
	client.hub.register <- client

	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		client.hub.broadcast <- message
	}
}
