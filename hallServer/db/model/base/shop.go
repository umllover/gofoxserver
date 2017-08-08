package base

import (
	"mj/hallServer/db"

	"github.com/lovelly/leaf/log"
)

//This file is generate by scripts,don't edit it

//shop
//

// +gen
type Shop struct {
	Id    int `db:"id" json:"id"`       // 商品
	Name  int `db:"name" json:"name"`   // 商品名字
	Price int `db:"price" json:"price"` // 价格
}

var DefaultShop = Shop{}

type shopCache struct {
	objMap  map[int]*Shop
	objList []*Shop
}

var ShopCache = &shopCache{}

func (c *shopCache) LoadAll() {
	sql := "select * from shop"
	c.objList = make([]*Shop, 0)
	err := db.BaseDB.Select(&c.objList, sql)
	if err != nil {
		log.Fatal(err.Error())
	}
	c.objMap = make(map[int]*Shop)
	log.Debug("Load all shop success %v", len(c.objList))
	for _, v := range c.objList {
		c.objMap[v.Id] = v
	}
}

func (c *shopCache) All() []*Shop {
	return c.objList
}

func (c *shopCache) Count() int {
	return len(c.objList)
}

func (c *shopCache) Get(id int) (*Shop, bool) {
	return c.GetKey1(id)
}

func (c *shopCache) GetKey1(id int) (*Shop, bool) {
	v, ok := c.objMap[id]
	return v, ok
}
