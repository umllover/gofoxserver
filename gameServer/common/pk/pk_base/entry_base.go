package pk_base

import (
	. "mj/common/cost"
	"mj/common/msg"
	"mj/gameServer/common/pk"
	"mj/gameServer/common/room_base"
	"mj/gameServer/db/model/base"
	datalog "mj/gameServer/log"

	"github.com/lovelly/leaf/log"
)

//创建的配置文件
type NewPKCtlConfig struct {
	BaseMgr  room_base.BaseManager
	TimerMgr room_base.TimerManager
	UserMgr  room_base.UserManager
	DataMgr  pk.DataManager
	LogicMgr pk.LogicManager
}

//消息入口文件
type Entry_base struct {
	*room_base.Entry_base
	LogicMgr pk.LogicManager
}

func (r *Entry_base) GetDataMgr() pk.DataManager {
	return r.DataMgr.(pk.DataManager)
}

func NewPKBase(info *msg.L2G_CreatorRoom) *Entry_base {
	Temp, ok1 := base.GameServiceOptionCache.Get(info.KindId, info.ServiceId)
	log.Debug("new pk base %d %d", info.KindId, info.ServiceId)
	if !ok1 {
		log.Error("at NewPKBase not foud config .... ")
		return nil
	}

	pk := new(Entry_base)
	pk.Entry_base = room_base.NewEntryBase(info.KindId, info.ServiceId)
	pk.Temp = Temp
	return pk
}

func (r *Entry_base) Init(cfg *NewPKCtlConfig) {
	r.UserMgr = cfg.UserMgr
	r.DataMgr = cfg.DataMgr
	r.BaseManager = cfg.BaseMgr
	r.LogicMgr = cfg.LogicMgr
	r.TimerMgr = cfg.TimerMgr
	r.RoomRun(r.DataMgr.GetRoomId())
	r.GetDataMgr().OnCreateRoom()
	r.TimerMgr.StartCreatorTimer(func() {
		roomLogData := datalog.RoomLog{}
		logData := roomLogData.GetRoomLogRecode(r.DataMgr.GetRoomId(), r.Temp.KindID, r.Temp.ServerID)
		roomLogData.UpdateGameLogRecode(logData, 4)
		r.OnEventGameConclude(NO_START_GER_DISMISS)
	})
}

func (r *Entry_base) GetRoomId() int {
	return r.DataMgr.GetRoomId()
}

//计算税收 暂时未实现
func (room *Entry_base) CalculateRevenue(ChairId, lScore int) int {
	//效验参数
	UserCnt := room.UserMgr.GetMaxPlayerCnt()
	if ChairId >= UserCnt {
		return 0
	}
	return 0
}
