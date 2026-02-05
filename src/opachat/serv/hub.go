package serv

import (
	"sync"
)

// Hub with clients
type Hub struct {
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
	lockHub    sync.RWMutex
}

// NewHub create new hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) addClient(cl *Client) {
	h.lockHub.Lock()
	h.clients[cl.uquser] = cl
	h.lockHub.Unlock()
}

func (h *Hub) removeClient(cl *Client) {
	h.lockHub.Lock()
	delete(h.clients, cl.uquser)
	h.lockHub.Unlock()
}

// Run hub
func (h *Hub) Run() {
	for {
		select {
		case uqcl := <-h.register:
			h.addClient(uqcl)
		case uqcl := <-h.unregister:
			h.removeClient(uqcl)
		}
	}
}
