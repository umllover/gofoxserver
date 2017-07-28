package base

import (
	"mj/hallServer/db"

	"github.com/lovelly/leaf/log"
)

//This file is generate by scripts,don't edit it

//personal_table_fee_bak_7_28
//

// +gen
type PersonalTableFeeBak728 struct {
	KindID         int `db:"KindID" json:"KindID"`                 // 游戏标识
	ServerID       int `db:"ServerID" json:"ServerID"`             //
	DrawCountLimit int `db:"DrawCountLimit" json:"DrawCountLimit"` // 局数限制
	AATableFee     int `db:"AATableFee" json:"AATableFee"`         // 时间限制
	TableFee       int `db:"TableFee" json:"TableFee"`             // 创建费用
}

var DefaultPersonalTableFeeBak728 = PersonalTableFeeBak728{}

type personalTableFeeBak728Cache struct {
	objMap  map[int]map[int]map[int]*PersonalTableFeeBak728
	objList []*PersonalTableFeeBak728
}

var PersonalTableFeeBak728Cache = &personalTableFeeBak728Cache{}

func (c *personalTableFeeBak728Cache) LoadAll() {
	sql := "select * from personal_table_fee_bak_7_28"
	c.objList = make([]*PersonalTableFeeBak728, 0)
	err := db.BaseDB.Select(&c.objList, sql)
	if err != nil {
		log.Fatal(err.Error())
	}
	c.objMap = make(map[int]map[int]map[int]*PersonalTableFeeBak728)
	log.Debug("Load all personal_table_fee_bak_7_28 success %v", len(c.objList))
	for _, v := range c.objList {
		obj, ok := c.objMap[v.KindID]
		if !ok {
			obj = make(map[int]map[int]*PersonalTableFeeBak728)
			c.objMap[v.KindID] = obj
		}

		obj2, ok2 := obj[v.ServerID]
		if !ok2 {
			obj2 = make(map[int]*PersonalTableFeeBak728)
			obj[v.ServerID] = obj2
		}
		obj2[v.DrawCountLimit] = v

	}
}

func (c *personalTableFeeBak728Cache) All() []*PersonalTableFeeBak728 {
	return c.objList
}

func (c *personalTableFeeBak728Cache) Count() int {
	return len(c.objList)
}

func (c *personalTableFeeBak728Cache) Get(KindID int, ServerID int, DrawCountLimit int) (*PersonalTableFeeBak728, bool) {
	return c.GetKey3(KindID, ServerID, DrawCountLimit)
}

func (c *personalTableFeeBak728Cache) GetKey1(KindID int) (map[int]map[int]*PersonalTableFeeBak728, bool) {
	v, ok := c.objMap[KindID]
	return v, ok
}

func (c *personalTableFeeBak728Cache) GetKey2(KindID int, ServerID int) (map[int]*PersonalTableFeeBak728, bool) {
	v, ok := c.objMap[KindID]
	if !ok {
		return nil, false
	}
	v1, ok1 := v[ServerID]
	return v1, ok1
}
func (c *personalTableFeeBak728Cache) GetKey3(KindID int, ServerID int, DrawCountLimit int) (*PersonalTableFeeBak728, bool) {
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
