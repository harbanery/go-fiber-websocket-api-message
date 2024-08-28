package services

import (
	"gofiber-chat-api/src/models"
	"sync"

	"github.com/gofiber/contrib/websocket"
)

type Hub struct {
	Rooms map[string]*Room
	Mutex sync.Mutex
}

type Room struct {
	Clients    map[*websocket.Conn]bool
	Register   chan *websocket.Conn
	Unregister chan *websocket.Conn
	Broadcast  chan models.Message
}

var hub *Hub

func NewHub() *Hub {
	return &Hub{
		Rooms: make(map[string]*Room),
	}
}

func InitHub() {
	hub = NewHub() // Inisialisasi hanya sekali di awal
}

func GetHub() *Hub {
	return hub
}

func (r *Room) Run() {
	for {
		select {
		case conn := <-r.Register:
			r.Clients[conn] = true
		case conn := <-r.Unregister:
			delete(r.Clients, conn)
		case msg := <-r.Broadcast:
			for conn := range r.Clients {
				_ = conn.WriteJSON(msg)
			}
		}
	}
}

func (h *Hub) GetRoom(roomID string) *Room {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()

	room, exists := h.Rooms[roomID]
	if !exists {
		room = &Room{
			Clients:    make(map[*websocket.Conn]bool),
			Register:   make(chan *websocket.Conn),
			Unregister: make(chan *websocket.Conn),
			Broadcast:  make(chan models.Message),
		}
		h.Rooms[roomID] = room
		go room.Run() // Memastikan Room menjalankan loop utamanya
	}
	return room
}
