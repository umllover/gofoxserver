package base

import (
	"mj/hallServer/db"

	"github.com/lovelly/leaf/log"
)

//This file is generate by scripts,don't edit it

//token_record
//

// +gen
type TokenRecord struct {
	UserId    int `db:"user_id" json:"user_id"`     //
	TokenType int `db:"tokenType" json:"tokenType"` //
	Amount    int `db:"amount" json:"amount"`       //
	Status    int `db:"status" json:"status"`       //
}

var DefaultTokenRecord = TokenRecord{}

type tokenRecordCache struct {
	objMap  map[int]*TokenRecord
	objList []*TokenRecord
}

var TokenRecordCache = &tokenRecordCache{}

func (c *tokenRecordCache) LoadAll() {
	sql := "select * from token_record"
	c.objList = make([]*TokenRecord, 0)
	err := db.BaseDB.Select(&c.objList, sql)
	if err != nil {
		log.Fatal(err.Error())
	}
	c.objMap = make(map[int]*TokenRecord)
	log.Debug("Load all token_record success %v", len(c.objList))
	for _, v := range c.objList {
		c.objMap[v.UserId] = v
	}
}

func (c *tokenRecordCache) All() []*TokenRecord {
	return c.objList
}

func (c *tokenRecordCache) Count() int {
	return len(c.objList)
}

func (c *tokenRecordCache) Get(user_id int) (*TokenRecord, bool) {
	return c.GetKey1(user_id)
}

func (c *tokenRecordCache) GetKey1(user_id int) (*TokenRecord, bool) {
	v, ok := c.objMap[user_id]
	return v, ok
}
