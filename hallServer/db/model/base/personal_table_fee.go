package base

import (
	"mj/hallServer/db"

	"github.com/lovelly/leaf/log"
)

//This file is generate by scripts,don't edit it

//personal_table_fee
//

// +gen
type PersonalTableFee struct {
	KindID         int `db:"KindID" json:"KindID"`                 // 游戏标识
	ServerID       int `db:"ServerID" json:"ServerID"`             //
	DrawCountLimit int `db:"DrawCountLimit" json:"DrawCountLimit"` // 局数限制
	AATableFee     int `db:"AATableFee" json:"AATableFee"`         // 时间限制
	TableFee       int `db:"TableFee" json:"TableFee"`             // 创建费用
	IniScore       int `db:"IniScore" json:"IniScore"`             // 初始分数
}

var DefaultPersonalTableFee = PersonalTableFee{}

type personalTableFeeCache struct {
	objMap  map[int]map[int]map[int]*PersonalTableFee
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
	c.objMap = make(map[int]map[int]map[int]*PersonalTableFee)
	log.Debug("Load all personal_table_fee success %v", len(c.objList))
	for _, v := range c.objList {
		obj, ok := c.objMap[v.KindID]
		if !ok {
			obj = make(map[int]map[int]*PersonalTableFee)
			c.objMap[v.KindID] = obj
		}

		obj2, ok2 := obj[v.ServerID]
		if !ok2 {
			obj2 = make(map[int]*PersonalTableFee)
			obj[v.ServerID] = obj2
		}
		obj2[v.DrawCountLimit] = v

	}
}

func (c *personalTableFeeCache) All() []*PersonalTableFee {
	return c.objList
}

func (c *personalTableFeeCache) Count() int {
	return len(c.objList)
}

func (c *personalTableFeeCache) Get(KindID int, ServerID int, DrawCountLimit int) (*PersonalTableFee, bool) {
	return c.GetKey3(KindID, ServerID, DrawCountLimit)
}

func (c *personalTableFeeCache) GetKey1(KindID int) (map[int]map[int]*PersonalTableFee, bool) {
	v, ok := c.objMap[KindID]
	return v, ok
}

func (c *personalTableFeeCache) GetKey2(KindID int, ServerID int) (map[int]*PersonalTableFee, bool) {
	v, ok := c.objMap[KindID]
	if !ok {
		return nil, false
	}
	v1, ok1 := v[ServerID]
	return v1, ok1
}
func (c *personalTableFeeCache) GetKey3(KindID int, ServerID int, DrawCountLimit int) (*PersonalTableFee, bool) {
	v, ok := c.objMap[KindID]
	if !ok {
		return nil, false
	}

	v1, ok1 := v[ServerID]
	if !ok1 {
		return nil, false
	}

	v2, ok2 := v1[DrawCountLimit]
	return v2, ok2
}
