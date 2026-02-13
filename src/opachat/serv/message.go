// Package serv
package serv

import (
	"encoding/json"
)

const (
	JOINROOM  = "joinroom"
	CANDIDATE = "candidate"
	OFFER     = "offer"
	ANSWER    = "answer"
	WHOCO     = "whoco"
	WHOCOINV  = "whocoinv"
	TCON      = "tcon"
	AVC       = "avc"
	AVCD      = "avcd"
	SCRE      = "screen"
	SCRECD    = "screencd"
	BREC      = "beginrecord"
	EREC      = "endrecord"
	RREC      = "remrec"
	CLBREC    = "clbeginrecord"
	CLEREC    = "clendrecord"
	CHAT      = "chat"
	TALKERST  = "talkerstopped"
)

// Message for hub select
type HubMessage struct {
	uquser string
	msg    []byte
}

// Message define our message object
type Message struct {
	Tp      string `json:"tp"`
	Content string `json:"content"`
}

type AVConfig struct {
	Sound    bool `json:"sound"`
	Video    bool `json:"video"`
	Invis    bool `json:"invis"`
	ScreenOn bool `json:"screen_on"`
}

type WConnected struct {
	StrID      string `json:"strid"`
	Uquser     string `json:"uquser"`
	Nik        string `json:"nik"`
	Mic        bool   `json:"mic"`
	Cam        bool   `json:"cam"`
	ScreenOn   bool   `json:"screen_on"`
	Recording  bool   `json:"recording"`
	Crecording bool   `json:"crecording"`
	Vili       string `json:"vili"`
	ChatMsg    string `json:"chat_message"`
}

type ListConnected struct {
	List map[string]WConnected `json:"list"`
}

func decMessage(txt []byte) (msg *Message) {
	msg = new(Message)
	json.Unmarshal(txt, msg)
	return
}

func decAV(txt string) (av *AVConfig) {
	av = new(AVConfig)
	json.Unmarshal([]byte(txt), av)
	return
}
