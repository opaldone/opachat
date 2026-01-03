package controllers

import (
	"fmt"
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
	ke := qv.Get("ke")

	perroom, err := strconv.Atoi(ps.ByName("perroom"))
	if err != nil {
		tools.Danger("perroom convert", err)
	}

	serv.ServeWs(roomuq, useruq, perroom, nik, ke, hub, w, r)
}

func Di(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	uq := ps.ByName("uq")

	if uq != "shpa" && uq != tools.GetKeyCSRF() {
		fmt.Printf("\ncsrf\t\t%s\n",
			tools.GetKeyCSRF(),
		)

		return
	}

	deb := serv.GetShowRooms()

	GenerateHTMLEmp(w, r, deb, "stru/dix")
}
