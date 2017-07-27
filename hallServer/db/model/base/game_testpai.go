package base

import (
	"mj/hallServer/db"

	"github.com/lovelly/leaf/log"
)

//This file is generate by scripts,don't edit it

//game_testpai
//

// +gen
type GameTestpai struct {
	Id        int    `db:"id" json:"id"`               //
	KindID    int    `db:"KindID" json:"KindID"`       // 游戏第一类型
	RoomID    int    `db:"RoomID" json:"RoomID"`       //
	ServerID  int    `db:"ServerID" json:"ServerID"`   // 游戏第二类型
	CardsName string `db:"CardsName" json:"CardsName"` // 卡牌类型名称
	ChairId   string `db:"chair_id" json:"chair_id"`   // 座位id
	Cards     string `db:"Cards" json:"Cards"`         // 所有玩家卡牌
	IsAcivate int8   `db:"IsAcivate" json:"IsAcivate"` // 是否激活牌型
}

var DefaultGameTestpai = GameTestpai{}

type gameTestpaiCache struct {
	objMap  map[int]*GameTestpai
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
	c.objMap = make(map[int]*GameTestpai)
	log.Debug("Load all game_testpai success %v", len(c.objList))
	for _, v := range c.objList {
		c.objMap[v.Id] = v
	}
}

func (c *gameTestpaiCache) All() []*GameTestpai {
	return c.objList
}

func (c *gameTestpaiCache) Count() int {
	return len(c.objList)
}

func (c *gameTestpaiCache) Get(id int) (*GameTestpai, bool) {
	return c.GetKey1(id)
}

func (c *gameTestpaiCache) GetKey1(id int) (*GameTestpai, bool) {
	v, ok := c.objMap[id]
	return v, ok
}
