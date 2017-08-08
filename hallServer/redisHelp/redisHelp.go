package RdsHelp

import (
	"fmt"
	. "mj/common/cost"
	"mj/hallServer/db"
	//"github.com/mitchellh/mapstructure"
)

func AddRoomInfo(roomid int, info map[string]interface{}) {
	db.RdsDB.HMSet(fmt.Sprintf(CreatorRoom, roomid), info)

}

func LoadRoomInfo(uid int64) {

}
