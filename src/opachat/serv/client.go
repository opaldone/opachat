package serv

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"opachat/tools"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 1048576
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	uqroom     string
	uquser     string
	nik        string
	ke         string
	recording  bool
	screen     bool
	talker     *Talker
	hub        *Hub
	conn       *websocket.Conn
	send       chan []byte
	lockClient sync.RWMutex
}

func (c *Client) sendMeCandidate(cand string) {
	if len(cand) == 0 {
		return
	}
	msg := new(Message)
	msg.Tp = CANDIDATE
	msg.Content = cand
	bts, _ := json.Marshal(msg)
	c.send <- bts
}

func (c *Client) sendMeOffer(of string) {
	msg := new(Message)
	msg.Tp = OFFER
	msg.Content = of
	bts, _ := json.Marshal(msg)
	c.send <- bts
}

func (c *Client) sendMeWhoConnected() {
	str := whoConnectedRoom(c.uqroom, c.uquser)

	if len(str) == 0 {
		return
	}

	msg := new(Message)
	msg.Tp = TCON
	msg.Content = str
	bts, _ := json.Marshal(msg)
	c.send <- bts
}

func (c *Client) sendMeAvcChanged(str string) {
	if len(str) == 0 {
		return
	}

	msg := new(Message)
	msg.Tp = AVCD
	msg.Content = str
	bts, _ := json.Marshal(msg)
	c.send <- bts
}

func (c *Client) sendMeScreenChanged(str string) {
	if len(str) == 0 {
		return
	}

	msg := new(Message)
	msg.Tp = SCRECD
	msg.Content = str
	bts, _ := json.Marshal(msg)
	c.send <- bts
}

func (c *Client) sendMeStartedRecord(str string) {
	if len(str) == 0 {
		return
	}

	msg := new(Message)
	msg.Tp = BREC
	msg.Content = str
	bts, _ := json.Marshal(msg)
	c.send <- bts
}

func (c *Client) sendMeStoppedRecord(str string) {
	if len(str) == 0 {
		return
	}

	msg := new(Message)
	msg.Tp = EREC
	msg.Content = str
	bts, _ := json.Marshal(msg)
	c.send <- bts
}

func (c *Client) sendMeAnotherRecord(str string) {
	if len(str) == 0 {
		return
	}

	msg := new(Message)
	msg.Tp = AREC
	msg.Content = str
	bts, _ := json.Marshal(msg)
	c.send <- bts
}

func (c *Client) sendMeChat(str string) {
	if len(str) == 0 {
		return
	}

	msg := new(Message)
	msg.Tp = CHAT
	msg.Content = str
	bts, _ := json.Marshal(msg)
	c.send <- bts
}

func (c *Client) processMessage(msg *Message) {
	switch msg.Tp {
	case JOINROOM:
		av := decAV(msg.Content)
		t := joinRoom(c, av)
		if t == nil {
			return
		}
		c.talker = t
	case CANDIDATE:
		c.talker.setCandidate(msg.Content)
	case ANSWER:
		c.talker.setAnswer(msg.Content)
	case WHOCO:
		c.sendMeWhoConnected()
	case AVC:
		av := decAV(msg.Content)

		if c.talker != nil {
			c.talker.changeOpts(av)
		}

		talkerChangedOpts(c)
	case SCRE:
		sv := decAV(msg.Content)

		if c.talker != nil && sv.ScreenOn {
			sv.Video = true

			c.lockClient.Lock()
			c.screen = true
			c.lockClient.Unlock()

			c.talker.changeOpts(sv)
		}

		if !sv.ScreenOn && c.screen {
			c.lockClient.Lock()
			c.screen = false
			c.lockClient.Unlock()

			c.talker.changeOpts(sv)
		}

		talkerChangedScreen(c, sv)
	case BREC:
		startRecord(c)
	case EREC:
		stopRecord(c)
	case RREC:
		removeRecord(c)
	case CHAT:
		chatMessage(c, msg.Content)
	}
}

// Gets some message
func (c *Client) readPump() {
	defer func() {
		c.stopClient()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				tools.Danger("readPump", err)
			}
			break
		}

		message = bytes.TrimSpace(bytes.ReplaceAll(message, newline, space))

		msg := decMessage(message)

		c.processMessage(msg)
	}
}

// Message to the outside
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			w.Write(message)

			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) stopClient() {
	c.hub.unregister <- c
	c.conn.Close()

	if c.talker == nil {
		return
	}

	c.talker.stopTalker()
}

// ServeWs handles websocket requests from the peer.
func ServeWs(roomuq_in string, useruq_in string,
	perroom int, nik_in string, ke_in string,
	hub_in *Hub, w_in http.ResponseWriter, r_in *http.Request,
) {
	start_virt := len(ke_in) > 0

	if start_virt && !CheckKeRoom(roomuq_in, ke_in) {
		tools.Log("ServeWs", "virt was not checked ke_in =", ke_in)
		return
	}

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	conn_in, err := upgrader.Upgrade(w_in, r_in, nil)
	if err != nil {
		tools.Danger("ServeWs", err)
		return
	}

	nc := &Client{
		uqroom:    roomuq_in,
		uquser:    useruq_in,
		nik:       nik_in,
		ke:        ke_in,
		recording: false,
		hub:       hub_in,
		conn:      conn_in,
		send:      make(chan []byte, 256),
	}

	if !start_virt {
		createRoom(nc.uqroom, perroom)
	}

	nc.hub.register <- nc
	go nc.writePump()
	go nc.readPump()
}
