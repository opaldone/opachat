package serv

import (
	"sync"
)

// Hub with clients
type Hub struct {
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
	sender     chan *HubMessage
	lockHub    sync.RWMutex
}

// NewHub create new hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		sender:     make(chan *HubMessage),
	}
}

func (h *Hub) addClient(cl *Client) {
	h.lockHub.Lock()
	h.clients[cl.uquser] = cl
	h.lockHub.Unlock()
}

func (h *Hub) removeClient(cl *Client) {
	h.lockHub.Lock()
	close(cl.chasend)
	delete(h.clients, cl.uquser)
	h.lockHub.Unlock()
}

func (h *Hub) messageClient(sen *HubMessage) {
	defer h.lockHub.RUnlock()
	h.lockHub.RLock()

	clret, ok := h.clients[sen.uquser]

	if !ok {
		return
	}

	clret.chasend <- sen.msg
}

// Run hub
func (h *Hub) Run() {
	for {
		select {
		case uqcl := <-h.register:
			h.addClient(uqcl)
		case uqcl := <-h.unregister:
			h.removeClient(uqcl)
		case sen := <-h.sender:
			h.messageClient(sen)
		}
	}
}
