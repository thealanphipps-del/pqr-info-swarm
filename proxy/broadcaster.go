package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Broadcaster struct {
	clients map[*websocket.Conn]bool
	mu      sync.Mutex
}

func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		clients: make(map[*websocket.Conn]bool),
	}
}

func (b *Broadcaster) HandleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}
	
	b.mu.Lock()
	b.clients[conn] = true
	b.mu.Unlock()

	defer func() {
		b.mu.Lock()
		delete(b.clients, conn)
		b.mu.Unlock()
		conn.Close()
	}()

	// keep connection alive
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
}

func (b *Broadcaster) Broadcast(payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		return
	}
	
	b.mu.Lock()
	defer b.mu.Unlock()
	for client := range b.clients {
		err := client.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			client.Close()
			delete(b.clients, client)
		}
	}
}
