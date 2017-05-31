package base

import (
	"mj/gameServer/db"

	"github.com/lovelly/leaf/log"
)

//This file is generate by scripts,don't edit it

//personal_table_fee
//

// +gen
type PersonalTableFee struct {
	ServerID       int `db:"ServerID" json:"ServerID"`             //
	KindID         int `db:"KindID" json:"KindID"`                 // 游戏标识
	DrawCountLimit int `db:"DrawCountLimit" json:"DrawCountLimit"` // 局数限制
	DrawTimeLimit  int `db:"DrawTimeLimit" json:"DrawTimeLimit"`   // 时间限制
	TableFee       int `db:"TableFee" json:"TableFee"`             // 创建费用
	IniScore       int `db:"IniScore" json:"IniScore"`             // 初始分数
}

var DefaultPersonalTableFee = PersonalTableFee{}

type personalTableFeeCache struct {
	objMap  map[int]map[int]map[int]map[int]*PersonalTableFee
	objList []*PersonalTableFee
}

var PersonalTableFeeCache = &personalTableFeeCache{}

func (c *personalTableFeeCache) LoadAll() {
	sql := "select * from personal_table_fee"
	c.objList = make([]*PersonalTableFee, 0)
	err := db.BaseDB.Select(&c.objList, sql)
	if err != nil {
		log.Fatal(err.Error())
	}
	c.objMap = make(map[int]map[int]map[int]map[int]*PersonalTableFee)
	log.Debug("Load all personal_table_fee success %v", len(c.objList))
	for _, v := range c.objList {
		obj, ok := c.objMap[v.ServerID]
		if !ok {
			obj = make(map[int]map[int]map[int]*PersonalTableFee)
			c.objMap[v.ServerID] = obj
		}

		obj2, ok2 := obj[v.KindID]
		if !ok2 {
			obj2 = make(map[int]map[int]*PersonalTableFee)
			obj[v.KindID] = obj2
		}

		obj3, ok3 := obj2[v.DrawCountLimit]
		if !ok3 {
			obj3 = make(map[int]*PersonalTableFee)
			obj2[v.DrawCountLimit] = obj3
		}
		obj3[v.DrawTimeLimit] = v
	}
}

func (c *personalTableFeeCache) All() []*PersonalTableFee {
	return c.objList
}

func (c *personalTableFeeCache) Count() int {
	return len(c.objList)
}

func (c *personalTableFeeCache) Get(ServerID int, KindID int, DrawCountLimit int, DrawTimeLimit int) (*PersonalTableFee, bool) {
	return c.GetKey4(ServerID, KindID, DrawCountLimit, DrawTimeLimit)
}

func (c *personalTableFeeCache) GetKey1(ServerID int) (map[int]map[int]map[int]*PersonalTableFee, bool) {
	v, ok := c.objMap[ServerID]
	return v, ok
}

func (c *personalTableFeeCache) GetKey2(ServerID int, KindID int) (map[int]map[int]*PersonalTableFee, bool) {
	v, ok := c.objMap[ServerID]
	if !ok {
		return nil, false
	}
	v1, ok1 := v[KindID]
	return v1, ok1
}
func (c *personalTableFeeCache) GetKey3(ServerID int, KindID int, DrawCountLimit int) (map[int]*PersonalTableFee, bool) {
	v, ok := c.objMap[ServerID]
	if !ok {
		return nil, false
	}

	v1, ok1 := v[KindID]
	if !ok1 {
		return nil, false
	}

	v2, ok2 := v1[DrawCountLimit]
	return v2, ok2
}
func (c *personalTableFeeCache) GetKey4(ServerID int, KindID int, DrawCountLimit int, DrawTimeLimit int) (*PersonalTableFee, bool) {
	v, ok := c.objMap[ServerID]
	if !ok {
		return nil, false
	}

	v1, ok1 := v[KindID]
	if !ok1 {
		return nil, false
	}

	v2, ok2 := v1[DrawCountLimit]
	if !ok2 {
		return nil, false
	}

	v3, ok3 := v2[DrawTimeLimit]
	return v3, ok3
}
