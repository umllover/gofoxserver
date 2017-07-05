package base

import (
	"mj/hallServer/db"

	"github.com/lovelly/leaf/log"
)

//This file is generate by scripts,don't edit it

//activity
//

// +gen
type Activity struct {
	Id          int    `db:"id" json:"id"`                   // 活动id。 和程序保持一致
	Description string `db:"description" json:"description"` // 活动描述
	DrawTimes   int64  `db:"draw_times" json:"draw_times"`   // 活动可以领取的次数
	DrawType    int    `db:"draw_type" json:"draw_type"`     // 领取类型，1是永久，2是每日领取，3是每周领取
	Amount      int    `db:"amount" json:"amount"`           // 奖励数量
	ItemType    int    `db:"item_type" json:"item_type"`     // 领取的物品类型， 1是钻石，
}

var DefaultActivity = Activity{}

type activityCache struct {
	objMap  map[int]*Activity
	objList []*Activity
}

var ActivityCache = &activityCache{}

func (c *activityCache) LoadAll() {
	sql := "select * from activity"
	c.objList = make([]*Activity, 0)
	err := db.BaseDB.Select(&c.objList, sql)
	if err != nil {
		log.Fatal(err.Error())
	}
	c.objMap = make(map[int]*Activity)
	log.Debug("Load all activity success %v", len(c.objList))
	for _, v := range c.objList {
		c.objMap[v.Id] = v
	}
}

func (c *activityCache) All() []*Activity {
	return c.objList
}

func (c *activityCache) Count() int {
	return len(c.objList)
}

func (c *activityCache) Get(id int) (*Activity, bool) {
	return c.GetKey1(id)
}

func (c *activityCache) GetKey1(id int) (*Activity, bool) {
	v, ok := c.objMap[id]
	return v, ok
}
