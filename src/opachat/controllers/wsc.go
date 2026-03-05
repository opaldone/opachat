package controllers

import (
	"net/http"
	"strconv"

	"opachat/serv"
	"opachat/tools"

	"github.com/julienschmidt/httprouter"
)

var hub *serv.Hub

func init() {
	hub = serv.NewHub()
	go hub.Run()
}

// Ws handler to create client
func Ws(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	qv := r.URL.Query()
	roomuq := ps.ByName("roomuq")
	useruq := ps.ByName("useruq")
	nik := qv.Get("nik")

	perroom, err := strconv.Atoi(ps.ByName("perroom"))
	if err != nil {
		tools.Danger("perroom convert", err)
	}

	serv.ServeWs(roomuq, useruq, perroom, nik, hub, w, r)
}

func WsErec(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	roomuq := ps.ByName("roomuq")
	serv.ServeWsErec(roomuq, w, r)
}
