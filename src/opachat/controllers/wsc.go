package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"opachat/serv"
	"opachat/tools"
	"strconv"

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
	w.Header().Set("Content-Type", "application/json")

	uq := ps.ByName("uq")

	if uq != tools.GetKeyCSRF() {
		fmt.Printf("\ncsrf\t\t%s\n",
			tools.GetKeyCSRF(),
		)

		return
	}

	some := serv.GetShowRooms()

	output, _ := json.MarshalIndent(some, "", "\t")

	w.Write(output)
}
