package base

import (
	"mj/gameServer/db"

	"github.com/lovelly/leaf/log"
)

//This file is generate by scripts,don't edit it

//achievement_limit
//

// +gen
type AchievementLimit struct {
	Level  int    `db:"level" json:"level"`   //
	Min    int    `db:"min" json:"min"`       //
	Max    int    `db:"max" json:"max"`       //
	Remark string `db:"remark" json:"remark"` //
}

var DefaultAchievementLimit = AchievementLimit{}

type achievementLimitCache struct {
	objMap  map[int]*AchievementLimit
	objList []*AchievementLimit
}

var AchievementLimitCache = &achievementLimitCache{}

func (c *achievementLimitCache) LoadAll() {
	sql := "select * from achievement_limit"
	c.objList = make([]*AchievementLimit, 0)
	err := db.BaseDB.Select(&c.objList, sql)
	if err != nil {
		log.Fatal(err.Error())
	}
	c.objMap = make(map[int]*AchievementLimit)
	log.Debug("Load all achievement_limit success %v", len(c.objList))
	for _, v := range c.objList {
		c.objMap[v.Level] = v
	}
}

func (c *achievementLimitCache) All() []*AchievementLimit {
	return c.objList
}

func (c *achievementLimitCache) Count() int {
	return len(c.objList)
}

func (c *achievementLimitCache) Get(level int) (*AchievementLimit, bool) {
	return c.GetKey1(level)
}

func (c *achievementLimitCache) GetKey1(level int) (*AchievementLimit, bool) {
	v, ok := c.objMap[level]
	return v, ok
}
