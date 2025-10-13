package serv

import (
	"os/exec"
	"strconv"

	"opachat/tools"
)

func startRec(r *Room) {
	e := tools.Env(true)

	if e.Recorder == nil {
		tools.Log("Recorder", "config of Recorder is not set")
		return
	}

	rid, rke, jsfile := r.setKeRecorder()

	cmd := exec.Command("./scr/s_s",
		rid,
		rke,
		jsfile,
		e.Recorder.URLVirt,
		e.Recorder.SoundLib,
		e.Recorder.IHw,
		e.Recorder.ScrRes,
		e.Recorder.LogLevel,
		strconv.Itoa(e.Recorder.Timeout),
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
