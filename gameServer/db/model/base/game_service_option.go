package base

import (
	"mj/gameServer/db"

	"github.com/lovelly/leaf/log"
)

//This file is generate by scripts,don't edit it

//game_service_option
//

// +gen
type GameServiceOption struct {
	KindID             int    `db:"KindID" json:"KindID"`                         // 名称号码
	ServerID           int    `db:"ServerID" json:"ServerID"`                     // 房间标识
	SortID             int    `db:"SortID" json:"SortID"`                         // 排列标识
	Source             int    `db:"Source" json:"Source"`                         // 单位积分
	MinEnterScore      int    `db:"MinEnterScore" json:"MinEnterScore"`           // 最低进入积分
	MaxEnterScore      int    `db:"MaxEnterScore" json:"MaxEnterScore"`           // 最高积分
	MinPlayer          int    `db:"MinPlayer" json:"MinPlayer"`                   // 最少几个人才能玩
	MaxPlayer          int    `db:"MaxPlayer" json:"MaxPlayer"`                   // 最多多少人一起玩
	GameType           int    `db:"GameType" json:"GameType"`                     // 游戏类型 1是开放类型， 2是比赛类型
	RoomName           string `db:"RoomName" json:"RoomName"`                     // 房间名称
	OffLineTrustee     int    `db:"OffLineTrustee" json:"OffLineTrustee"`         // 是否短线代打 0是不托管 1是托管
	IniScore           int    `db:"IniScore" json:"IniScore"`                     // 游戏开始玩家默认积分
	PlayTurnCount      int    `db:"PlayTurnCount" json:"PlayTurnCount"`           // 房间能够进行游戏的最大局数
	TimeAfterBeginTime int    `db:"TimeAfterBeginTime" json:"TimeAfterBeginTime"` // 游戏开始后多长时间后解散桌子
	TimeOffLineCount   int    `db:"TimeOffLineCount" json:"TimeOffLineCount"`     // 玩家掉线多长时间后解散桌子
	TimeNotBeginGame   int    `db:"TimeNotBeginGame" json:"TimeNotBeginGame"`     // 多长时间未开始游戏解散桌子	 单位秒
	DynamicJoin        int    `db:"DynamicJoin" json:"DynamicJoin"`               // 是够允许游戏开始后加入 1是允许
	OutCardTime        int    `db:"OutCardTime" json:"OutCardTime"`               // 多久没出牌自动出牌
	OperateCardTime    int    `db:"OperateCardTime" json:"OperateCardTime"`       // 操作最大时间
	TimeRoomTrustee    int    `db:"TimeRoomTrustee" json:"TimeRoomTrustee"`       // 房间进入托管时间
}

var DefaultGameServiceOption = GameServiceOption{}

type gameServiceOptionCache struct {
	objMap  map[int]map[int]*GameServiceOption
	objList []*GameServiceOption
}

var GameServiceOptionCache = &gameServiceOptionCache{}

func (c *gameServiceOptionCache) LoadAll() {
	sql := "select * from game_service_option"
	c.objList = make([]*GameServiceOption, 0)
	err := db.BaseDB.Select(&c.objList, sql)
	if err != nil {
		log.Fatal(err.Error())
	}
	c.objMap = make(map[int]map[int]*GameServiceOption)
	log.Debug("Load all game_service_option success %v", len(c.objList))
	for _, v := range c.objList {
		obj, ok := c.objMap[v.KindID]
		if !ok {
			obj = make(map[int]*GameServiceOption)
			c.objMap[v.KindID] = obj
		}
		obj[v.ServerID] = v

	}
}

func (c *gameServiceOptionCache) All() []*GameServiceOption {
	return c.objList
}

func (c *gameServiceOptionCache) Count() int {
	return len(c.objList)
}

func (c *gameServiceOptionCache) Get(KindID int, ServerID int) (*GameServiceOption, bool) {
	return c.GetKey2(KindID, ServerID)
}

func (c *gameServiceOptionCache) GetKey1(KindID int) (map[int]*GameServiceOption, bool) {
	v, ok := c.objMap[KindID]
	return v, ok
}

func (c *gameServiceOptionCache) GetKey2(KindID int, ServerID int) (*GameServiceOption, bool) {
	v, ok := c.objMap[KindID]
	if !ok {
		return nil, false
	}
	v1, ok1 := v[ServerID]
	return v1, ok1
}
