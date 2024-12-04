package serv

import (
	"os/exec"
	"strconv"

	"opachat/tools"
)

func startRec(r *Room) {
	rid, rke, js_file := r.setKeSaver()

	e := tools.Env(true)

	cmd := exec.Command("./scr/s_s",
		rid,
		rke,
		js_file,
		e.Saver.UrlVirt,
		e.Saver.IHw,
		e.Saver.ScrRes,
		e.Saver.LogLevel,
		strconv.Itoa(e.Saver.Timeout),
	)

	err := cmd.Start()
	if err != nil {
		tools.Danger("startRec", err)
	}
}

func stopRec(r *Room) {
	osa := r.getOsa()

	if osa == nil {
		return
	}

	cmd := exec.Command("./scr/k_s",
		strconv.Itoa(osa.Pgoo),
		strconv.Itoa(osa.Pff),
		strconv.Itoa(osa.Pxv),
	)

	err := cmd.Start()
	if err != nil {
		tools.Danger("stopRec", err)
	}
}
