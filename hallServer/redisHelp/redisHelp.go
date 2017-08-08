package RdsHelp

import (
	"mj/hallServer/db"
	. "mj/common/cost"
	"fmt"
)
func AddRoomInfo(roomid int,info map[string]interface{}){
	db.RdsDB.HMSet(fmt.Sprintf(CreatorRoom, roomid), info )
}

