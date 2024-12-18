package serv

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"syscall"
	"time"

	"opachat/tools"

	"github.com/pion/webrtc/v4"
)

type OutSaver struct {
	Pxv  int `json:"pxv,omitempty"`
	Pff  int `json:"pff,omitempty"`
	Pgoo int `json:"pgoo,omitempty"`
}

// Room is a room
type Room struct {
	id          string
	perRoom     int
	keSaver     string
	talkers     map[string]*Talker
	removeMe    func(string)
	trackLocals map[string]*webrtc.TrackLocalStaticRTP
	lockRoom    sync.RWMutex
}

// NewRoom creates new room
func NewRoom(uqroom string, perroom int, remFn func(string)) *Room {
	r := &Room{
		id:          uqroom,
		perRoom:     perroom,
		keSaver:     "",
		removeMe:    remFn,
		trackLocals: map[string]*webrtc.TrackLocalStaticRTP{},
	}

	return r
}

// Add to list of tracks and fire renegotation for all PeerConnections
func (r *Room) addTrack(t *webrtc.TrackRemote) *webrtc.TrackLocalStaticRTP {
	r.lockRoom.Lock()
	defer func() {
		r.lockRoom.Unlock()
		r.signalPeerConnections()
	}()

	// Create a new TrackLocal with the same codec as our incoming
	trackLocal, err := webrtc.NewTrackLocalStaticRTP(t.Codec().RTPCodecCapability, t.ID(), t.StreamID())
	if err != nil {
		panic(err)
	}

	r.trackLocals[t.ID()] = trackLocal

	return trackLocal
}

// Remove from list of tracks and fire renegotation for all PeerConnections
func (r *Room) removeTrack(t *webrtc.TrackLocalStaticRTP) {
	r.lockRoom.Lock()
	defer func() {
		r.lockRoom.Unlock()
		r.signalPeerConnections()
	}()

	delete(r.trackLocals, t.ID())
}

// signalPeerConnections updates each PeerConnection so that it is getting all the expected media tracks
func (r *Room) signalPeerConnections() {
	r.lockRoom.Lock()
	defer func() {
		r.lockRoom.Unlock()
		// dispatchKeyFrame()
	}()

	attemptSync := func() (tryAgain bool) {
		for _, talker := range r.talkers {
			if talker.pc.ConnectionState() == webrtc.PeerConnectionStateClosed {
				return true
			}

			// map of sender we already are seanding, so we don't double send
			existingSenders := map[string]bool{}

			for _, sender := range talker.pc.GetSenders() {
				if sender.Track() == nil {
					continue
				}

				existingSenders[sender.Track().ID()] = true

				// If we have a RTPSender that doesn't map to a existing track remove and signal
				if _, ok := r.trackLocals[sender.Track().ID()]; !ok {
					if err := talker.pc.RemoveTrack(sender); err != nil {
						return true
					}
				}
			}

			// Don't receive videos we are sending, make sure we don't have loopback
			for _, receiver := range talker.pc.GetReceivers() {
				if receiver.Track() == nil {
					continue
				}

				existingSenders[receiver.Track().ID()] = true
			}

			// Add all track we aren't sending yet to the PeerConnection
			for trackID := range r.trackLocals {
				if _, ok := existingSenders[trackID]; !ok {
					if _, err := talker.pc.AddTrack(r.trackLocals[trackID]); err != nil {
						return true
					}
				}
			}

			offer, err := talker.pc.CreateOffer(nil)
			if err != nil {
				return true
			}

			talker.lockO.Lock()
			err = talker.pc.SetLocalDescription(offer)
			talker.lockO.Unlock()

			if err != nil {
				return true
			}

			offerString, err := json.Marshal(offer)
			if err != nil {
				return true
			}

			talker.wsc.sendMeOffer(string(offerString))
		}

		return
	}

	for syncAttempt := 0; ; syncAttempt++ {
		if syncAttempt == 25 {
			// Release the lock and attempt a sync in 3 seconds. We might be blocking a RemoveTrack or AddTrack
			go func() {
				time.Sleep(time.Second * 3)
				r.signalPeerConnections()
			}()
			return
		}

		if !attemptSync() {
			break
		}
	}
}

func (r *Room) canPutTalker(uquser string) bool {
	r.lockRoom.Lock()
	if r.talkers == nil {
		r.talkers = make(map[string]*Talker)
	}
	r.lockRoom.Unlock()

	r.lockRoom.RLock()
	len_talkers := len(r.talkers)
	r.lockRoom.RUnlock()

	if r.perRoom <= len_talkers {
		return false
	}

	r.lockRoom.RLock()
	_, exists := r.talkers[uquser]
	r.lockRoom.RUnlock()

	return !exists
}

func (r *Room) addTalker(c *Client, av *AVConfig) *Talker {
	if !r.canPutTalker(c.uquser) {
		return nil
	}

	newTalker := NewTalker(c, r, av)

	r.lockRoom.Lock()
	r.talkers[newTalker.wsc.uquser] = newTalker
	r.lockRoom.Unlock()

	r.signalPeerConnections()

	return newTalker
}

func (r *Room) getConnectedList(me string) (res string) {
	r.lockRoom.RLock()
	defer r.lockRoom.RUnlock()

	lis := make(map[string]WConnected)
	for _, talker := range r.talkers {
		if talker.wsc.uquser == me {
			continue
		}
		if len(talker.wsc.ke) > 0 {
			continue
		}

		lis[talker.strID] = WConnected{
			Uquser:    talker.wsc.uquser,
			Nik:       talker.wsc.nik,
			Mic:       talker.sound,
			Cam:       talker.video,
			Recording: talker.wsc.recording,
			ScreenOn:  talker.wsc.screen,
		}
	}

	str := ListConnected{List: lis}
	bont, _ := json.Marshal(str)
	res = string(bont)

	return
}

func (r *Room) notifTalkersStartedRecord(me *Client) {
	r.lockRoom.RLock()
	defer r.lockRoom.RUnlock()

	me.recording = true

	wc := WConnected{
		StrID: me.talker.strID,
	}

	bont, _ := json.Marshal(wc)
	res := string(bont)

	for _, talker := range r.talkers {
		if talker.wsc.uquser == me.uquser {
			continue
		}

		talker.wsc.sendMeStartedRecord(res)
	}
}

func (r *Room) notifTalkersStoppedRecord(me *Client) {
	r.lockRoom.RLock()
	defer r.lockRoom.RUnlock()

	me.recording = false

	wc := WConnected{
		StrID:  me.talker.strID,
		Uquser: me.uquser,
		Vili:   r.keSaver,
	}

	bont, _ := json.Marshal(wc)
	res := string(bont)

	for _, talker := range r.talkers {
		talker.wsc.sendMeStoppedRecord(res)
	}
}

func (r *Room) notifTalkerAnotherRecord(c *Client) {
	r.lockRoom.RLock()
	defer r.lockRoom.RUnlock()

	for _, talker := range r.talkers {
		if !talker.wsc.recording {
			continue
		}

		wc := WConnected{
			StrID:     talker.strID,
			Uquser:    talker.wsc.uquser,
			Nik:       talker.wsc.nik,
			Recording: talker.wsc.recording,
		}

		bont, _ := json.Marshal(wc)
		res := string(bont)

		c.sendMeAnotherRecord(res)

		return
	}
}

func (r *Room) notifTalkersChangedOpts(me *Client) {
	r.lockRoom.RLock()
	defer r.lockRoom.RUnlock()

	wc := WConnected{
		StrID:  me.talker.strID,
		Uquser: me.talker.wsc.uquser,
		Mic:    me.talker.sound,
		Cam:    me.talker.video,
	}

	bont, _ := json.Marshal(wc)
	res := string(bont)

	for _, talker := range r.talkers {
		if talker.wsc.uquser == me.uquser {
			continue
		}

		talker.wsc.sendMeAvcChanged(res)
	}
}

func (r *Room) notifTalkersChangedScreen(me *Client, sv *AVConfig) {
	r.lockRoom.RLock()
	defer r.lockRoom.RUnlock()

	wc := WConnected{
		StrID:    me.talker.strID,
		Uquser:   me.uquser,
		ScreenOn: sv.ScreenOn,
		Mic:      sv.Sound,
		Cam:      sv.Video,
	}

	bont, _ := json.Marshal(wc)
	res := string(bont)

	for _, talker := range r.talkers {
		if talker.wsc.uquser == me.uquser {
			continue
		}

		talker.wsc.sendMeScreenChanged(res)
	}
}

func (r *Room) removeTalker(idTalker string) {
	r.lockRoom.RLock()
	_, exists := r.talkers[idTalker]
	r.lockRoom.RUnlock()

	if !exists {
		return
	}

	r.lockRoom.Lock()
	delete(r.talkers, idTalker)
	r.lockRoom.Unlock()

	r.lockRoom.RLock()
	len_talkers := len(r.talkers)
	r.lockRoom.RUnlock()

	if len_talkers == 0 {
		r.removeMe(r.id)
	}
}

func (r *Room) checkKe(ke_in string) bool {
	r.lockRoom.RLock()
	defer r.lockRoom.RUnlock()

	if len(r.keSaver) == 0 {
		return false
	}

	if r.keSaver != ke_in {
		return false
	}

	return true
}

func (r *Room) getWriter() *Client {
	r.lockRoom.RLock()
	defer r.lockRoom.RUnlock()

	for _, talker := range r.talkers {
		if talker.wsc.recording {
			return talker.wsc
		}
	}

	return nil
}

func (r *Room) monitorRecording() {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	_running := func(pid int) bool {
		proc, _ := os.FindProcess(pid)

		err := proc.Signal(syscall.Signal(0))

		return err == nil
	}

	osa_in := r.getOsa()
	pids := []int{}

	for range ticker.C {
		if osa_in != nil && len(pids) == 0 {
			pids = append(pids, osa_in.Pxv)
			pids = append(pids, osa_in.Pff)
			pids = append(pids, osa_in.Pgoo)
		}

		if osa_in == nil {
			osa_in = r.getOsa()
		}

		if len(pids) == 0 {
			continue
		}

		for _, pid := range pids {
			if _running(pid) {
				continue
			}

			wr_cl := r.getWriter()
			if wr_cl != nil {
				r.stopRecord(wr_cl)
				return
			}
		}
	}
}

func (r *Room) startRecord(c *Client) {
	r.lockRoom.RLock()
	empty_ke := len(r.keSaver) == 0
	r.lockRoom.RUnlock()

	if empty_ke {
		startRec(r)

		go r.monitorRecording()

		r.notifTalkersStartedRecord(c)
		return
	}

	r.notifTalkerAnotherRecord(c)
}

func (r *Room) stopRecord(c *Client) {
	r.lockRoom.RLock()
	empty_ke := len(r.keSaver) == 0
	r.lockRoom.RUnlock()

	if empty_ke {
		return
	}

	stopRec(r)

	r.notifTalkersStoppedRecord(c)
}

func (r *Room) getPathOsa() (string, string) {
	r.lockRoom.RLock()
	rid := r.id
	r.lockRoom.RUnlock()

	js_file := fmt.Sprintf("./prcs/pr_%s.json", rid)

	return rid, js_file
}

func (r *Room) setKeRecorder() (string, string, string) {
	ke_new := tools.CreateUUID()

	r.lockRoom.Lock()
	r.keSaver = ke_new
	r.lockRoom.Unlock()

	rid, js_file := r.getPathOsa()

	return rid, ke_new, js_file
}

func (r *Room) clearKeSaver() {
	r.lockRoom.Lock()
	r.keSaver = ""
	r.lockRoom.Unlock()
}

func (r *Room) getOsa() *OutSaver {
	_, js_file := r.getPathOsa()

	os_json_str, err := os.Open(js_file)
	if err != nil {
		return nil
	}

	decoder := json.NewDecoder(os_json_str)
	ous := &OutSaver{}
	err = decoder.Decode(ous)
	if err != nil {
		tools.Danger(fmt.Sprintf("Cannot parse %s", js_file), err)
		return nil
	}

	return ous
}

func (r *Room) removeRecord() {
	_, js_file := r.getPathOsa()

	if len(js_file) == 0 {
		return
	}

	err := os.Remove(js_file)
	if err != nil {
		tools.Danger("Removing js file", err)
	}

	r.clearKeSaver()
}

func (r *Room) getInfo() (ret string) {
	r.lockRoom.RLock()
	defer r.lockRoom.RUnlock()

	ret = fmt.Sprintf("room_id=%s "+
		"talkers_len=%d "+
		"trackLocals_len=%d "+
		"keSaver=%s",
		r.id, len(r.talkers), len(r.trackLocals), r.keSaver,
	)

	return
}
