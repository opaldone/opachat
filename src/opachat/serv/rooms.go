package serv

import (
	"sync"
)

var rooms map[string]*Room

var lockRooms sync.RWMutex

func init() {
	rooms = make(map[string]*Room)
}

func getRoom(uqroomIn string) *Room {
	lockRooms.RLock()
	roo, exists := rooms[uqroomIn]
	lockRooms.RUnlock()

	if !exists {
		return nil
	}

	return roo
}

func getClientRoom(cl *Client) *Room {
	if cl == nil {
		return nil
	}

	uqroomIn := cl.uqroom

	return getRoom(uqroomIn)
}

func createRoom(uqroom string, perroom int) {
	roo := getRoom(uqroom)

	if roo != nil {
		return
	}

	r := NewRoom(uqroom, perroom, removeRoom)

	lockRooms.Lock()
	rooms[uqroom] = r
	lockRooms.Unlock()
}

func joinRoom(c *Client, av *AVConfig) *Talker {
	roo := getClientRoom(c)

	if roo == nil {
		return nil
	}

	talker := roo.addTalker(c, av)

	return talker
}

func removeRoom(uqroom string) {
	roo := getRoom(uqroom)

	if roo == nil {
		return
	}

	lockRooms.Lock()
	delete(rooms, uqroom)
	lockRooms.Unlock()
}

func WhoConnectedRoom(uqroomIn string, me string, onlyInvis bool) (res string) {
	roo := getRoom(uqroomIn)

	if roo == nil {
		return
	}

	res = roo.getConnectedList(me, onlyInvis)

	return
}

func talkerHi(me *Client) {
	roo := getClientRoom(me)

	if roo == nil {
		return
	}

	roo.notifTalkersHi(me)
}

func talkerStop(me *Client) {
	roo := getClientRoom(me)

	if roo == nil {
		return
	}

	roo.notifTalkersStop(me)
}

func talkerChangedOpts(me *Client) {
	roo := getClientRoom(me)

	if roo == nil {
		return
	}

	roo.notifTalkersChangedOpts(me)
}

func talkerChangedScreen(me *Client, sv *AVConfig) {
	roo := getClientRoom(me)

	if roo == nil {
		return
	}

	roo.notifTalkersChangedScreen(me, sv)
}

func startServerRecord(cl *Client) {
	roo := getClientRoom(cl)

	if roo == nil {
		return
	}

	roo.notifStartedRecord(cl, BREC)
	cl.setRecording(true)
}

func stopServerRecord(cl *Client) {
	roo := getClientRoom(cl)

	if roo == nil {
		return
	}

	roo.notifStoppedRecord(cl, EREC)
	cl.setRecording(false)
}

func startClientRecord(cl *Client) {
	roo := getClientRoom(cl)

	if roo == nil {
		return
	}

	roo.notifStartedRecord(cl, CLBREC)
	cl.setCrecording(true)
}

func stopClientRecord(cl *Client) {
	roo := getClientRoom(cl)

	if roo == nil {
		return
	}

	roo.notifStoppedRecord(cl, CLEREC)
	cl.setCrecording(false)
}

func chatMessage(cl *Client, msg string) {
	roo := getClientRoom(cl)

	if roo == nil {
		return
	}

	roo.chatMessage(cl, msg)
}

func GetShowRooms() (ret RoomsDebugType) {
	lockRooms.RLock()
	defer lockRooms.RUnlock()

	for _, r := range rooms {
		ret.Rooms = append(ret.Rooms, r.getInfo())
	}

	return
}
