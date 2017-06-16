package base

import (
	"mj/hallServer/db"

	"github.com/lovelly/leaf/log"
)

//This file is generate by scripts,don't edit it

//global_var
//

// +gen
type GlobalVar struct {
	K      string `db:"K" json:"K"`           //
	V      string `db:"V" json:"V"`           //
	Remark string `db:"Remark" json:"Remark"` //
}

var DefaultGlobalVar = GlobalVar{}

type globalVarCache struct {
	objMap  map[string]*GlobalVar
	objList []*GlobalVar
}

var GlobalVarCache = &globalVarCache{}

func (c *globalVarCache) LoadAll() {
	sql := "select * from global_var"
	c.objList = make([]*GlobalVar, 0)
	err := db.BaseDB.Select(&c.objList, sql)
	if err != nil {
		log.Fatal(err.Error())
	}
	c.objMap = make(map[string]*GlobalVar)
	log.Debug("Load all global_var success %v", len(c.objList))
	for _, v := range c.objList {
		c.objMap[v.K] = v
	}
}

func (c *globalVarCache) All() []*GlobalVar {
	return c.objList
}

func (c *globalVarCache) Count() int {
	return len(c.objList)
}

func (c *globalVarCache) Get(K string) (*GlobalVar, bool) {
	return c.GetKey1(K)
}

func (c *globalVarCache) GetKey1(K string) (*GlobalVar, bool) {
	v, ok := c.objMap[K]
	return v, ok
}
