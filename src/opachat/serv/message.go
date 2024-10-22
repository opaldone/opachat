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
	ROOMCON   = "roomcon"
	TCON      = "tcon"
	AVC       = "avc"
	AVCD      = "avcd"
	SCRE      = "screen"
	SCRECD    = "screencd"
	BREC      = "beginrecord"
	EREC      = "endrecord"
	AREC      = "anotherrecord"
	RREC      = "remrec"
)

// Message define our message object
type Message struct {
	Tp      string `json:"tp"`
	Content string `json:"content"`
}

type AVConfig struct {
	Sound    bool `json:"sound"`
	Video    bool `json:"video"`
	ScreenOn bool `json:"screen_on"`
}

type WConnected struct {
	StrID     string `json:"strid"`
	Uquser    string `json:"uquser"`
	Nik       string `json:"nik"`
	Mic       bool   `json:"mic"`
	Cam       bool   `json:"cam"`
	ScreenOn  bool   `json:"screen_on"`
	Recording bool   `json:"recording"`
	Vili      string `json:"vili"`
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
