package serv

import (
	"sync"
)

var rooms map[string]*Room

var lockRooms sync.RWMutex

func init() {
	rooms = make(map[string]*Room)
}

func getRoom(uqroom_in string) *Room {
	lockRooms.RLock()
	roo, exists := rooms[uqroom_in]
	lockRooms.RUnlock()

	if !exists {
		return nil
	}

	return roo
}

func getClientRoom(cl *Client) *Room {
	uqroom_in := cl.uqroom

	return getRoom(uqroom_in)
}

func CheckKeRoom(uqroom_in string, ke_in string) bool {
	roo := getRoom(uqroom_in)

	if roo == nil {
		return false
	}

	return roo.checkKe(ke_in)
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

func whoConnectedRoom(uqroom_in string, me string) (res string) {
	roo := getRoom(uqroom_in)

	if roo == nil {
		return
	}

	res = roo.getConnectedList(me)

	return
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
