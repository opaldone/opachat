package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"opachat/serv"
	"opachat/tools"

	"github.com/julienschmidt/httprouter"
)

func setupResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers",
		"access-control-allow-origin, x-requested-with",
	)
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

func Lir(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	setupResponse(w)

	if r.Method == "OPTIONS" {
		return
	}

	roomuq := ps.ByName("roomuq")

	str := serv.WhoConnectedRoom(roomuq, "", false)

	ret := struct {
		Res string `json:"res"`
	}{
		Res: str,
	}

	output, _ := json.Marshal(ret)

	w.Write(output)
}
