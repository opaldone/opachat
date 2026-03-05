package serv

import (
	"encoding/json"
	"sort"
	"sync"
	"time"

	"github.com/pion/webrtc/v4"
	// "github.com/pion/webrtc/v3"
)

type Room struct {
	id          string
	perRoom     int
	talkers     map[string]*Talker
	removeMe    func(string)
	trackLocals map[string]*webrtc.TrackLocalStaticRTP
	lockRoom    sync.RWMutex
}

type TalkerDebType struct {
	Nik        string   `json:"nik"`
	Uquser     string   `json:"uquser"`
	StrID      string   `json:"strID"`
	Recording  bool     `json:"recording"`
	Crecording bool     `json:"crecording"`
	Screen     bool     `json:"screen"`
	Sound      bool     `json:"sound"`
	Video      bool     `json:"video"`
	Invis      bool     `json:"invis"`
	Virt       bool     `json:"virt"`
	Ices       []string `json:"ices"`
}

type RoomDebType struct {
	RoomID         string          `json:"room_id"`
	TalkersLen     int             `json:"talkers_len"`
	TrackLocalsLen int             `json:"trackLocals_len"`
	Talkers        []TalkerDebType `json:"talkers"`
}

type RoomsDebugType struct {
	Rooms []RoomDebType `json:"rooms,omitempty"`
}

// NewRoom creates new room
func NewRoom(uqroom string, perroom int, remFn func(string)) *Room {
	r := &Room{
		id:          uqroom,
		perRoom:     perroom,
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

			talker.wsc.sendMe(string(offerString), OFFER)
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
	lentalkers := len(r.talkers)
	r.lockRoom.RUnlock()

	if r.perRoom <= lentalkers {
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
	lentalkers := len(r.talkers)
	r.lockRoom.RUnlock()

	if lentalkers == 0 {
		r.removeMe(r.id)
	}
}

func (r *Room) getConnectedList(me string, onlyInvis bool) (res string) {
	r.lockRoom.RLock()
	defer r.lockRoom.RUnlock()

	lis := make(map[string]WConnected)
	for _, talker := range r.talkers {
		if talker.wsc.uquser == me {
			continue
		}
		if onlyInvis && !talker.wsc.invis {
			continue
		}
		if talker.wsc.virt {
			continue
		}

		lis[talker.strID] = WConnected{
			StrID:      talker.strID,
			Uquser:     talker.wsc.uquser,
			Nik:        talker.wsc.nik,
			Mic:        talker.sound,
			Cam:        talker.video,
			Recording:  talker.wsc.recording,
			Crecording: talker.wsc.crecording,
			ScreenOn:   talker.wsc.screen,
		}
	}

	str := ListConnected{List: lis}
	bont, _ := json.Marshal(str)
	res = string(bont)

	return
}

func (r *Room) notifTalkersHi(me *Client) {
	r.lockRoom.RLock()
	defer r.lockRoom.RUnlock()

	lis := make(map[string]WConnected)

	lis[me.talker.strID] = WConnected{
		StrID:     me.talker.strID,
		Uquser:    me.talker.wsc.uquser,
		Nik:       me.talker.wsc.nik,
		Mic:       me.talker.sound,
		Cam:       me.talker.video,
		Recording: me.talker.wsc.recording,
		ScreenOn:  me.talker.wsc.screen,
	}

	str := ListConnected{List: lis}
	bont, _ := json.Marshal(str)
	res := string(bont)

	for _, talker := range r.talkers {
		if talker.wsc.uquser == me.uquser {
			continue
		}

		talker.wsc.sendMe(res, TCON)
	}
}

func (r *Room) notifTalkersStop(me *Client) {
	if me.talker == nil {
		return
	}

	r.lockRoom.RLock()
	defer r.lockRoom.RUnlock()

	wc := WConnected{
		StrID: me.talker.strID,
	}

	bont, _ := json.Marshal(wc)
	res := string(bont)

	for _, talker := range r.talkers {
		if talker.wsc.uquser == me.uquser {
			continue
		}

		talker.wsc.sendMe(res, TALKERST)
	}
}

func (r *Room) getRecordingClient() *Client {
	r.lockRoom.RLock()
	defer r.lockRoom.RUnlock()

	for _, talker := range r.talkers {
		if talker.wsc.recording {
			return talker.wsc
		}
	}

	return nil
}

func (r *Room) notifStartedRecord(me *Client, msgtag string) {
	r.lockRoom.RLock()
	defer r.lockRoom.RUnlock()

	wc := WConnected{
		StrID: me.talker.strID,
	}

	bont, _ := json.Marshal(wc)
	res := string(bont)

	for _, talker := range r.talkers {
		if talker.wsc.uquser == me.uquser {
			continue
		}

		talker.wsc.sendMe(res, msgtag)
	}
}

func (r *Room) notifStoppedRecord(me *Client, msgtag string) {
	r.lockRoom.RLock()
	defer r.lockRoom.RUnlock()

	wc := WConnected{
		StrID:  me.talker.strID,
		Uquser: me.uquser,
	}

	bont, _ := json.Marshal(wc)
	res := string(bont)

	for _, talker := range r.talkers {
		talker.wsc.sendMe(res, msgtag)
	}
}

func (r *Room) notifTalkersChangedOpts(me *Client) {
	if me.talker == nil {
		return
	}

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

		talker.wsc.sendMe(res, AVCD)
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

		talker.wsc.sendMe(res, SCRECD)
	}
}

func (r *Room) chatMessage(me *Client, msg string) {
	r.lockRoom.RLock()
	defer r.lockRoom.RUnlock()

	wc := WConnected{
		StrID:   me.talker.strID,
		Uquser:  me.uquser,
		ChatMsg: msg,
	}

	bont, _ := json.Marshal(wc)
	res := string(bont)

	for _, talker := range r.talkers {
		talker.wsc.sendMe(res, CHAT)
	}
}

func (r *Room) getInfo() (ret RoomDebType) {
	r.lockRoom.RLock()
	defer r.lockRoom.RUnlock()

	ret.RoomID = r.id
	ret.TalkersLen = len(r.talkers)
	ret.TrackLocalsLen = len(r.trackLocals)

	for _, t := range r.talkers {
		ret.Talkers = append(ret.Talkers, t.getInfo())
	}

	sort.Slice(ret.Talkers, func(i, j int) bool {
		return ret.Talkers[i].Nik < ret.Talkers[j].Nik
	})

	return
}
