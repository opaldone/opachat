package serv

import (
	"encoding/json"
	"fmt"
	"sync"

	"opachat/tools"

	"github.com/pion/interceptor"
	"github.com/pion/interceptor/pkg/intervalpli"
	"github.com/pion/rtp"

	"github.com/pion/webrtc/v4"
	// "github.com/pion/webrtc/v3"
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
func NewTalker(cin *Client, roomin *Room, av *AVConfig) *Talker {
	nt := &Talker{
		wsc:  cin,
		room: roomin,
	}

	if cin.invis {
		nt.strID = cin.uquser
	}

	nt.sound = av.Sound
	nt.video = av.Video

	nt.connect()

	return nt
}

func (t *Talker) getPeerConnectionConfig() (peerConnectionConfig webrtc.Configuration) {
	urlsout, usernameout, credentialout := tools.GetIceList()

	peerConnectionConfig = webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				// URLs: []string{"stun:stun.l.google.com:19302"},

				URLs:       urlsout,
				Username:   usernameout,
				Credential: credentialout,
			},
		},
	}

	return
}

func (t *Talker) myOnTrack(jsTrack *webrtc.TrackRemote, _ *webrtc.RTPReceiver) {
	var vid bool

	vid = true
	if jsTrack.Kind() == webrtc.RTPCodecTypeAudio {
		vid = false
	}

	t.strID = jsTrack.StreamID()

	trackLocal := t.room.addTrack(jsTrack)
	defer t.room.removeTrack(trackLocal)

	buf := make([]byte, 1500)
	rtpPkt := &rtp.Packet{}

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

		if err = rtpPkt.Unmarshal(buf[:i]); err != nil {
			return
		}

		rtpPkt.Extension = false
		rtpPkt.Extensions = nil

		if err = trackLocal.WriteRTP(rtpPkt); err != nil {
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

	t.wsc.sendMe(string(candidateString), CANDIDATE)
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

	med := &webrtc.MediaEngine{}
	err = med.RegisterDefaultCodecs()
	if err != nil {
		tools.Danger("RegisterDefaultCodecs", err)
	}

	i := &interceptor.Registry{}
	err = webrtc.RegisterDefaultInterceptors(med, i)
	if err != nil {
		tools.Danger("RegisterDefaultInterceptors", err)
	}

	intervalPliFactory, err := intervalpli.NewReceiverInterceptor()
	if err != nil {
		tools.Danger("NewReceiverInterceptor", err)
	}

	i.Add(intervalPliFactory)

	peerconf := t.getPeerConnectionConfig()

	t.pc, err = webrtc.
		NewAPI(webrtc.WithMediaEngine(med), webrtc.WithInterceptorRegistry(i)).
		NewPeerConnection(peerconf)
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

func (t *Talker) getInfo() (ret TalkerDebType) {
	t.lockO.RLock()
	defer t.lockO.RUnlock()

	ret.Nik = t.wsc.nik
	ret.Uquser = t.wsc.uquser
	ret.StrID = t.strID
	ret.Recording = t.wsc.recording
	ret.Crecording = t.wsc.crecording
	ret.Screen = t.wsc.screen
	ret.Sound = t.sound
	ret.Video = t.video
	ret.Invis = t.wsc.invis
	ret.Ke = t.wsc.ke

	for _, s := range t.pc.GetStats() {
		switch stat := s.(type) {
		case webrtc.ICECandidateStats:
			if stat.Type == webrtc.StatsTypeRemoteCandidate {
				ret.Ices = append(ret.Ices, fmt.Sprintf("%s %s:%d", stat.Type, stat.IP, stat.Port))
			}
		default:
		}
	}

	return
}
