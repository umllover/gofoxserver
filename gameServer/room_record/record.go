package room_record

import (
	"encoding/json"
	"mj/common/msg"
	"mj/gameServer/db/model"

	"time"

	"github.com/lovelly/leaf/log"
)

type StartRecord struct {
	Users     []*msg.G2C_UserEnter
	Info      map[string]interface{}
	BeginTime int64
}

type PlayingRecord struct {
	Cmd []map[string]interface{}
}

type EndRecord struct {
	TotalTimes int64
}

type record struct {
	Start   *StartRecord
	Playing *PlayingRecord
	End     *EndRecord
	KindID  int
}

var (
	records = make(map[int]*record) //key is room id
)

func HasRoomRecord(roomId int) bool {
	_, ok := records[roomId]
	return ok
}

func AddRoomRecordInfo(roomId int, info *msg.G2C_PersonalTableTip) {
	roominfo, ok := records[roomId]
	if !ok {
		roominfo = new(record)
		records[roomId] = roominfo
		roominfo.Start = new(StartRecord)
		roominfo.Playing = new(PlayingRecord)
		roominfo.End = new(EndRecord)
		roominfo.Start.Info = make(map[string]interface{})
		roominfo.KindID = info.KindID
	}

	roominfo.Start.Info["G2C_PersonalTableTip"] = info
}

func AddStartInfo(roomId int, Users []*msg.G2C_UserEnter) {
	roominfo, ok := records[roomId]
	if ok {
		roominfo.Start.BeginTime = time.Now().Unix()
		roominfo.Start.Users = Users
	}
}
func AddPlayCmd(roomId int, data map[string]interface{}) {
	roominfo, ok := records[roomId]
	if ok {
		roominfo.Playing.Cmd = append(roominfo.Playing.Cmd, data)
	}
}

func AddEndInfo(roomId int) {
	roominfo, ok := records[roomId]
	if ok {
		roominfo.End.TotalTimes = time.Now().Unix() - roominfo.Start.BeginTime
		saveRecor := &model.RoomRecord{}
		str, err := json.Marshal(roominfo.Start)
		if err != nil {
			log.Error("at AddEndInfo json.Marshal error:%s", err.Error())
			return
		}
		saveRecor.StartInfo = string(str)

		str, err = json.Marshal(roominfo.Playing)
		if err != nil {
			log.Error("at AddEndInfo json.Marshal 1 error:%s", err.Error())
			return
		}
		saveRecor.PlayInfo = string(str)

		str, err = json.Marshal(roominfo.End)
		if err != nil {
			log.Error("at AddEndInfo json.Marshal 2 error:%s", err.Error())
			return
		}
		saveRecor.EndInfo = string(str)
		id, ierr := model.RoomRecordOp.Insert(saveRecor)
		if ierr == nil {
			userRecord := &model.UserRoomRecord{}
			for _, v := range roominfo.Start.Users {
				userRecord.UserId = v.UserID
				userRecord.RecordId = int(id)
				userRecord.KindId = roominfo.KindID
				_, err := model.UserRoomRecordOp.Insert(userRecord)
				if err != nil {
					log.Error("at AddEndInfo inser error:%s", err.Error())
				}
			}
		} else {
			log.Error("at AddEndInfo inser 11 error:%s", err.Error())
		}

		delete(records, roomId)
	}
}
