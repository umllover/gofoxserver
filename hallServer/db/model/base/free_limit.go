package base

import (
	"mj/hallServer/db"
	"time"

	"github.com/lovelly/leaf/log"
)

//This file is generate by scripts,don't edit it

//free_limit
//

// +gen
type FreeLimit struct {
	FreeName  string     `db:"free_name" json:"free_name"`   //
	FreeType  string     `db:"free_type" json:"free_type"`   //
	FreeBegin *time.Time `db:"free_begin" json:"free_begin"` //
	FreeEnd   *time.Time `db:"free_end" json:"free_end"`     //
}

var DefaultFreeLimit = FreeLimit{}

type freeLimitCache struct {
	objMap  map[string]*FreeLimit
	objList []*FreeLimit
}

var FreeLimitCache = &freeLimitCache{}

func (c *freeLimitCache) LoadAll() {
	sql := "select * from free_limit"
	c.objList = make([]*FreeLimit, 0)
	err := db.BaseDB.Select(&c.objList, sql)
	if err != nil {
		log.Fatal(err.Error())
	}
	c.objMap = make(map[string]*FreeLimit)
	log.Debug("Load all free_limit success %v", len(c.objList))
	for _, v := range c.objList {
		c.objMap[v.FreeName] = v
	}
}

func (c *freeLimitCache) All() []*FreeLimit {
	return c.objList
}

func (c *freeLimitCache) Count() int {
	return len(c.objList)
}

func (c *freeLimitCache) Get(free_name string) (*FreeLimit, bool) {
	return c.GetKey1(free_name)
}

func (c *freeLimitCache) GetKey1(free_name string) (*FreeLimit, bool) {
	v, ok := c.objMap[free_name]
	return v, ok
}
