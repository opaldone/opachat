package serv

import (
	"opachat/tools"
	"os/exec"
	"strconv"
)

func startRec(r *Room) {
	rid, rke, js_file := r.setKeSaver()

	e := tools.Env()
	cmd := exec.Command("./scr/s_s",
		rid,
		rke,
		js_file,
		e.Saver.UVirt,
		e.Saver.Loop,
		e.Saver.Screen,
		e.Saver.Loglevel,
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
