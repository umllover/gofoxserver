package base

import (
	"mj/hallServer/db"
	"time"

	"github.com/lovelly/leaf/log"
)

//This file is generate by scripts,don't edit it

//goods
//

// +gen
type Goods struct {
	GoodsId           int        `db:"goods_id" json:"goods_id"`                       //
	Rmb               int        `db:"rmb" json:"rmb"`                                 //
	Diamond           int        `db:"diamond" json:"diamond"`                         //
	Name              string     `db:"name" json:"name"`                               // 商品名称
	LeftCnt           int        `db:"left_cnt" json:"left_cnt"`                       // 剩余数量
	SpecialOffer      int        `db:"special_offer" json:"special_offer"`             // 特价
	GivePresent       int        `db:"give_present" json:"give_present"`               // 赠送
	SpecialOfferBegin *time.Time `db:"special_offer_begin" json:"special_offer_begin"` // 特价开始时间
	SpecialOfferEnd   *time.Time `db:"special_offer_end" json:"special_offer_end"`     // 特价结束时间
	Type              string     `db:"type" json:"type"`                               // 类别
}

var DefaultGoods = Goods{}

type goodsCache struct {
	objMap  map[int]*Goods
	objList []*Goods
}

var GoodsCache = &goodsCache{}

func (c *goodsCache) LoadAll() {
	sql := "select * from goods"
	c.objList = make([]*Goods, 0)
	err := db.BaseDB.Select(&c.objList, sql)
	if err != nil {
		log.Fatal(err.Error())
	}
	c.objMap = make(map[int]*Goods)
	log.Debug("Load all goods success %v", len(c.objList))
	for _, v := range c.objList {
		c.objMap[v.GoodsId] = v
	}
}

func (c *goodsCache) All() []*Goods {
	return c.objList
}

func (c *goodsCache) Count() int {
	return len(c.objList)
}

func (c *goodsCache) Get(goods_id int) (*Goods, bool) {
	return c.GetKey1(goods_id)
}

func (c *goodsCache) GetKey1(goods_id int) (*Goods, bool) {
	v, ok := c.objMap[goods_id]
	return v, ok
}
