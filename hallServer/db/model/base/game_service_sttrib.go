package base

import (
	"mj/hallServer/db"

	"github.com/lovelly/leaf/log"
)

//This file is generate by scripts,don't edit it

//game_service_sttrib
//

// +gen
type GameServiceSttrib struct {
	KindID         int    `db:"KindID" json:"KindID"`                 // 名称号码
	ChairCount     int    `db:"ChairCount" json:"ChairCount"`         // 椅子数目
	SupporType     int    `db:"SupporType" json:"SupporType"`         // 支持类型
	GameName       string `db:"GameName" json:"GameName"`             // 游戏名字
	AndroidUser    int    `db:"AndroidUser" json:"AndroidUser"`       // 机器标志
	DynamicJoin    int    `db:"DynamicJoin" json:"DynamicJoin"`       // 动态加入
	OffLineTrustee int    `db:"OffLineTrustee" json:"OffLineTrustee"` // 断线代打
	AndroidActive  int    `db:"AndroidActive" json:"AndroidActive"`   // 主动陪打
	ServerVersion  int    `db:"ServerVersion" json:"ServerVersion"`   // 游戏版本
	ClientVersion  int    `db:"ClientVersion" json:"ClientVersion"`   //
}

var DefaultGameServiceSttrib = GameServiceSttrib{}

type gameServiceSttribCache struct {
	objMap  map[int]*GameServiceSttrib
	objList []*GameServiceSttrib
}

var GameServiceSttribCache = &gameServiceSttribCache{}

func (c *gameServiceSttribCache) LoadAll() {
	sql := "select * from game_service_sttrib"
	c.objList = make([]*GameServiceSttrib, 0)
	err := db.BaseDB.Select(&c.objList, sql)
	if err != nil {
		log.Fatal(err.Error())
	}
	c.objMap = make(map[int]*GameServiceSttrib)
	log.Debug("Load all game_service_sttrib success %v", len(c.objList))
	for _, v := range c.objList {
		c.objMap[v.KindID] = v
	}
}

func (c *gameServiceSttribCache) All() []*GameServiceSttrib {
	return c.objList
}

func (c *gameServiceSttribCache) Count() int {
	return len(c.objList)
}

func (c *gameServiceSttribCache) Get(KindID int) (*GameServiceSttrib, bool) {
	key := KindID
	v, ok := c.objMap[key]
	return v, ok
}

// 仅限运营后台实时刷新服务器数据用
func (c *gameServiceSttribCache) Update(v *GameServiceSttrib) {
	key := v.KindID
	c.objMap[key] = v
}
