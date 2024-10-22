package serv

import (
	"encoding/json"
	"fmt"
	"opachat/tools"

	"sync"

	"github.com/pion/interceptor"
	"github.com/pion/interceptor/pkg/intervalpli"

	"github.com/pion/webrtc/v3"
)

type Talker struct {
	wsc   *Client
	room  *Room
	pc    *webrtc.PeerConnection
	sound bool
	video bool
	strID string
	lockO sync.RWMutex
}

// NewTalker creates a new talker
func NewTalker(c_in *Client, room_in *Room, av *AVConfig) *Talker {
	newTalker := &Talker{
		wsc:  c_in,
		room: room_in,
	}

	newTalker.sound = av.Sound
	newTalker.video = av.Video

	newTalker.connect()

	return newTalker
}

func (t *Talker) getPeerConnectionConfig() (peerConnectionConfig webrtc.Configuration) {
	// turn server is here
	// /mnt/terik/a_my/test/gol/turns/tser

	urls_out, username_out, credential_out := tools.GetIceList()

	peerConnectionConfig = webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				// URLs: []string{"stun:stun.l.google.com:19302"},

				// URLs: []string{"stun:192.168.0.104:3478"},

				URLs:       urls_out,
				Username:   username_out,
				Credential: credential_out,
			},
		},
	}

	return
}

func (t *Talker) myOnTrack(jsTrack *webrtc.TrackRemote, _ *webrtc.RTPReceiver) {
	t.strID = jsTrack.StreamID()

	vid := true
	if jsTrack.Kind() == webrtc.RTPCodecTypeAudio {
		vid = false
	}

	trackLocal := t.room.addTrack(jsTrack)
	defer t.room.removeTrack(trackLocal)

	buf := make([]byte, 1500)
	for {
		i, _, err := jsTrack.Read(buf)

		if err != nil {
			return
		}

		ch := false
		t.lockO.RLock()
		if vid {
			ch = t.video
		} else {
			ch = t.sound
		}
		t.lockO.RUnlock()

		if !ch {
			continue
		}

		if _, err = trackLocal.Write(buf[:i]); err != nil {
			return
		}
	}
}

func (t *Talker) iceCandidate(i *webrtc.ICECandidate) {
	if i == nil {
		return
	}

	candidateString, err := json.Marshal(i.ToJSON())

	if err != nil {
		tools.Danger("iceCandidate Marshal", err)
		return
	}

	t.wsc.sendMeCandidate(string(candidateString))
}

func (t *Talker) connectionStateChange(p webrtc.PeerConnectionState) {
	switch p {
	case webrtc.PeerConnectionStateFailed:
		if err := t.pc.Close(); err != nil {
			tools.Danger("connectionStateChange close", err)
		}
	case webrtc.PeerConnectionStateClosed:
		t.room.signalPeerConnections()
	default:
	}
}

func (t *Talker) setAnswer(cont string) {
	answer := webrtc.SessionDescription{}

	if err := json.Unmarshal([]byte(cont), &answer); err != nil {
		tools.Danger("setAnswer Unmarshal", err)
		return
	}

	t.lockO.Lock()
	err := t.pc.SetRemoteDescription(answer)
	t.lockO.Unlock()

	if err != nil {
		tools.Danger("setAnswer SetRemoteDescription", err)
		return
	}
}

func (t *Talker) setCandidate(cont string) {
	candidate := webrtc.ICECandidateInit{}

	if err := json.Unmarshal([]byte(cont), &candidate); err != nil {
		tools.Danger("setCandidate Unmarshal", err)
		return
	}

	if err := t.pc.AddICECandidate(candidate); err != nil {
		tools.Danger("setCandidate AddICECandidate", err)
		return
	}
}

func (t *Talker) connect() {
	var err error

	m := &webrtc.MediaEngine{}
	err = m.RegisterDefaultCodecs()

	if err != nil {
		tools.Danger("RegisterDefaultCodecs", err)
	}

	i := &interceptor.Registry{}
	err = webrtc.RegisterDefaultInterceptors(m, i)

	if err != nil {
		tools.Danger("RegisterDefaultInterceptors", err)
	}

	intervalPliFactory, err := intervalpli.NewReceiverInterceptor()

	if err != nil {
		tools.Danger("NewReceiverInterceptor", err)
	}

	i.Add(intervalPliFactory)

	peer_conf := t.getPeerConnectionConfig()

	t.pc, err = webrtc.
		NewAPI(webrtc.WithMediaEngine(m), webrtc.WithInterceptorRegistry(i)).
		NewPeerConnection(peer_conf)

	if err != nil {
		tools.Danger("New peer connection", err)
	}

	for _, typ := range []webrtc.RTPCodecType{webrtc.RTPCodecTypeVideo, webrtc.RTPCodecTypeAudio} {
		if _, err := t.pc.AddTransceiverFromKind(typ, webrtc.RTPTransceiverInit{
			Direction: webrtc.RTPTransceiverDirectionRecvonly,
		}); err != nil {
			tools.Danger("AddTransceiverFromKind", err)
			return
		}
	}

	t.pc.OnICECandidate(t.iceCandidate)
	t.pc.OnConnectionStateChange(t.connectionStateChange)
	t.pc.OnTrack(t.myOnTrack)
}

func (t *Talker) changeOpts(av *AVConfig) {
	t.lockO.Lock()
	t.sound = av.Sound
	t.video = av.Video
	t.lockO.Unlock()
}

func (t *Talker) stopTalker() {
	t.pc.Close()
	t.room.removeTalker(t.wsc.uquser)
}

func (t *Talker) getInfo() (ret string) {
	t.lockO.RLock()
	defer t.lockO.RUnlock()

	ret = fmt.Sprintf(
		"nik=%s "+
			"uquser=%s "+
			"strID=%s "+
			"rec=%t "+
			"scr=%t "+
			"sound=%t "+
			"video=%t "+
			"ke=%s",
		t.wsc.nik, t.wsc.uquser, t.strID,
		t.wsc.recording, t.wsc.screen, t.sound, t.video, t.wsc.ke,
	)

	return
}
