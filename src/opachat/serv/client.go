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
	invis      bool
	ke         string
	recording  bool
	screen     bool
	talker     *Talker
	hub        *Hub
	conn       *websocket.Conn
	chasend    chan []byte
	lockClient sync.RWMutex
}

func (c *Client) sendMe(str string, co string) {
	if len(str) == 0 {
		return
	}

	if c.chasend == nil {
		return
	}

	msg := new(Message)
	msg.Tp = co
	msg.Content = str
	bts, _ := json.Marshal(msg)

	c.chasend <- bts
}

func (c *Client) sendMeWhoConnected(onlyInvis bool) {
	str := whoConnectedRoom(c.uqroom, c.uquser, onlyInvis)

	if len(str) == 0 {
		return
	}

	c.sendMe(str, TCON)
}

func (c *Client) stopClient() {
	talkerStop(c)

	c.hub.unregister <- c
	c.conn.Close()
	c.chasend = nil

	if c.talker == nil {
		return
	}

	c.talker.stopTalker()
}

func (c *Client) processMessage(msg *Message) {
	switch msg.Tp {
	case JOINROOM:
		av := decAV(msg.Content)

		c.invis = av.Invis

		t := joinRoom(c, av)

		if t == nil {
			return
		}

		c.talker = t

		if c.invis {
			talkerHi(c)
		}
	case CANDIDATE:
		c.talker.setCandidate(msg.Content)
	case ANSWER:
		c.talker.setAnswer(msg.Content)
	case WHOCO:
		c.sendMeWhoConnected(false)
	case WHOCOINV:
		c.sendMeWhoConnected(true)
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
		case message, ok := <-c.chasend:
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

			n := len(c.chasend)
			for range n {
				w.Write(newline)
				w.Write(<-c.chasend)
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

// ServeWs handles websocket requests from the peer.
func ServeWs(roomuqin string, useruqin string,
	perroom int, nikin string, kein string,
	hubin *Hub, win http.ResponseWriter, rin *http.Request,
) {
	startvirt := len(kein) > 0

	if startvirt && !CheckKeRoom(roomuqin, kein) {
		tools.Log("ServeWs", "virt was not checked ke_in =", kein)
		return
	}

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	connin, err := upgrader.Upgrade(win, rin, nil)
	if err != nil {
		tools.Danger("ServeWs", err)
		return
	}

	nc := &Client{
		uqroom:    roomuqin,
		uquser:    useruqin,
		nik:       nikin,
		ke:        kein,
		recording: false,
		hub:       hubin,
		conn:      connin,
		chasend:   make(chan []byte, 256),
	}

	if !startvirt {
		createRoom(nc.uqroom, perroom)
	}

	nc.hub.register <- nc
	go nc.writePump()
	go nc.readPump()
}
