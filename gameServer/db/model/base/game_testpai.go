package base

import (
	"mj/gameServer/db"

	"github.com/lovelly/leaf/log"
)

//This file is generate by scripts,don't edit it

//game_testpai
//

// +gen
type GameTestpai struct {
	KindID        int    `db:"KindID" json:"KindID"`               // 游戏类型
	ServerID      int    `db:"ServerID" json:"ServerID"`           // 服务器ID
	CardsName     string `db:"CardsName" json:"CardsName"`         // 卡牌类型名称
	BankerCard    string `db:"BankerCard" json:"BankerCard"`       // 庄家卡牌
	AllPlayerCard string `db:"AllPlayerCard" json:"AllPlayerCard"` // 所有玩家卡牌
	IsAcivate     int8   `db:"IsAcivate" json:"IsAcivate"`         // 是否激活牌型
}

var DefaultGameTestpai = GameTestpai{}

type gameTestpaiCache struct {
	objMap  map[int]map[int]*GameTestpai
	objList []*GameTestpai
}

var GameTestpaiCache = &gameTestpaiCache{}

func (c *gameTestpaiCache) LoadAll() {
	sql := "select * from game_testpai"
	c.objList = make([]*GameTestpai, 0)
	err := db.BaseDB.Select(&c.objList, sql)
	if err != nil {
		log.Fatal(err.Error())
	}
	c.objMap = make(map[int]map[int]*GameTestpai)
	log.Debug("Load all game_testpai success %v", len(c.objList))
	for _, v := range c.objList {
		obj, ok := c.objMap[v.KindID]
		if !ok {
			obj = make(map[int]*GameTestpai)
			c.objMap[v.KindID] = obj
		}
		obj[v.ServerID] = v

	}
}

func (c *gameTestpaiCache) All() []*GameTestpai {
	return c.objList
}

func (c *gameTestpaiCache) Count() int {
	return len(c.objList)
}

func (c *gameTestpaiCache) Get(KindID int, ServerID int) (*GameTestpai, bool) {
	return c.GetKey2(KindID, ServerID)
}

func (c *gameTestpaiCache) GetKey1(KindID int) (map[int]*GameTestpai, bool) {
	v, ok := c.objMap[KindID]
	return v, ok
}

func (c *gameTestpaiCache) GetKey2(KindID int, ServerID int) (*GameTestpai, bool) {
	v, ok := c.objMap[KindID]
	if !ok {
		return nil, false
	}
	v1, ok1 := v[ServerID]
	return v1, ok1
}
