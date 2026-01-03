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
	uqroomIn := cl.uqroom

	return getRoom(uqroomIn)
}

func CheckKeRoom(uqroomIn string, keIn string) bool {
	roo := getRoom(uqroomIn)

	if roo == nil {
		return false
	}

	return roo.checkKe(keIn)
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

func whoConnectedRoom(uqroomIn string, me string, onlyInvis bool) (res string) {
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

func startRecord(cl *Client) {
	roo := getClientRoom(cl)

	if roo == nil {
		return
	}

	roo.startRecord(cl)
}

func stopRecord(cl *Client) {
	roo := getClientRoom(cl)

	if roo == nil {
		return
	}

	roo.stopRecord(cl)
}

func removeRecord(cl *Client) {
	roo := getClientRoom(cl)

	if roo == nil {
		return
	}

	roo.removeRecord()
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
