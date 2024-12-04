package serv

import (
	"sync"

	"opachat/tools"
)

var rooms map[string]*Room

var lockRooms sync.RWMutex

func init() {
	rooms = make(map[string]*Room)
}

func CheckKeRoom(uqroom_in string, ke_in string) bool {
	lockRooms.RLock()
	roo, exists := rooms[uqroom_in]
	lockRooms.RUnlock()

	if !exists {
		return false
	}

	return roo.checkKe(ke_in)
}

func createRoom(uqroom string, perroom int) {
	lockRooms.RLock()
	_, exists := rooms[uqroom]
	lockRooms.RUnlock()

	if exists {
		return
	}

	r := NewRoom(uqroom, perroom, removeRoom)

	lockRooms.Lock()
	rooms[uqroom] = r
	lockRooms.Unlock()
}

func joinRoom(c *Client, av *AVConfig) *Talker {
	lockRooms.RLock()
	r, exists := rooms[c.uqroom]
	lockRooms.RUnlock()

	if !exists {
		return nil
	}

	talker := r.addTalker(c, av)

	return talker
}

func removeRoom(uqroom string) {
	lockRooms.RLock()
	_, exists := rooms[uqroom]
	lockRooms.RUnlock()

	if !exists {
		return
	}

	lockRooms.Lock()
	delete(rooms, uqroom)
	lockRooms.Unlock()
}

func whoConnectedRoom(uqroom_in string, me string) (res string) {
	lockRooms.RLock()
	roo, exists := rooms[uqroom_in]
	lockRooms.RUnlock()

	if !exists {
		return
	}

	res = roo.getConnectedList(me)

	return
}

func talkerChangedOpts(me *Client) {
	uqroom_in := me.uqroom

	lockRooms.RLock()
	roo, exists := rooms[uqroom_in]
	lockRooms.RUnlock()

	if !exists {
		return
	}

	roo.notifTalkersChangedOpts(me)
}

func talkerChangedScreen(me *Client, sv *AVConfig) {
	uqroom_in := me.uqroom

	lockRooms.RLock()
	roo, exists := rooms[uqroom_in]
	lockRooms.RUnlock()

	if !exists {
		return
	}

	roo.notifTalkersChangedScreen(me, sv)
}

func startRecord(cl *Client) {
	uqroom_in := cl.uqroom

	lockRooms.RLock()
	roo, exists := rooms[uqroom_in]
	lockRooms.RUnlock()

	if !exists {
		return
	}

	roo.startRecord(cl)
}

func stopRecord(cl *Client) {
	uqroom_in := cl.uqroom

	lockRooms.RLock()
	roo, exists := rooms[uqroom_in]
	lockRooms.RUnlock()

	if !exists {
		return
	}

	roo.stopRecord(cl)
}

func removeRecord(cl *Client) {
	uqroom_in := cl.uqroom

	lockRooms.RLock()
	roo, exists := rooms[uqroom_in]
	lockRooms.RUnlock()

	if !exists {
		return
	}

	roo.removeRecord()
}

func GetOsaFromRoom(uqroom_in string) *OutSaver {
	lockRooms.RLock()
	roo, exists := rooms[uqroom_in]
	lockRooms.RUnlock()

	if !exists {
		return nil
	}

	return roo.getOsa()
}

func GetShowRooms() []tools.RoomDebugType {
	list := []tools.RoomDebugType{}

	lockRooms.RLock()
	for _, r := range rooms {
		rt := tools.RoomDebugType{}
		rt.Room = r.getInfo()

		osa := r.getOsa()
		if osa != nil {
			rt.Osa = tools.DebugJ(osa, false, "", "")
		}

		tlist := []string{}
		for _, t := range r.talkers {
			tlist = append(tlist, t.getInfo())
		}

		rt.Talkers = tlist

		list = append(list, rt)
	}
	lockRooms.RUnlock()

	return list
}
