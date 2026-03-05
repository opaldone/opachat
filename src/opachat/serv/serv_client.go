package serv

import (
	"net/http"

	"opachat/tools"
)

func ServeWsErec(roomuqin string, win http.ResponseWriter, rin *http.Request) {
	connin, err := upgrader.Upgrade(win, rin, nil)
	if err != nil {
		tools.Danger("ServeWsErec", err)
		return
	}
	defer connin.Close()

	roo := getRoom(roomuqin)

	if roo == nil {
		return
	}

	cl := roo.getRecordingClient()

	if cl == nil {
		cl = &Client{
			uqroom:    roomuqin,
			uquser:    "empty_user",
			recording: true,
			talker: &Talker{
				strID: "empty_talker",
			},
		}
	}

	stopServerRecord(cl)
}
