package base

import (
	"mj/gameServer/db"

	"github.com/lovelly/leaf/log"
)

//This file is generate by scripts,don't edit it

//server_list
//

// +gen
type ServerList struct {
	SvrId   int    `db:"svr_id" json:"svr_id"`     // 节点id
	SvrType int    `db:"svr_type" json:"svr_type"` // 服务器类型 1是大厅服
	Host    string `db:"host" json:"host"`         // ip
	Port    int    `db:"port" json:"port"`         // 端口
	Status  int    `db:"status" json:"status"`     // 状态 1是正常状态 2是维护
}

var DefaultServerList = ServerList{}

type serverListCache struct {
	objMap  map[int]map[int]*ServerList
	objList []*ServerList
}

var ServerListCache = &serverListCache{}

func (c *serverListCache) LoadAll() {
	sql := "select * from server_list"
	c.objList = make([]*ServerList, 0)
	err := db.BaseDB.Select(&c.objList, sql)
	if err != nil {
		log.Fatal(err.Error())
	}
	c.objMap = make(map[int]map[int]*ServerList)
	log.Debug("Load all server_list success %v", len(c.objList))
	for _, v := range c.objList {
		obj, ok := c.objMap[v.SvrId]
		if !ok {
			obj = make(map[int]*ServerList)
			c.objMap[v.SvrId] = obj
		}
		obj[v.SvrType] = v

	}
}

func (c *serverListCache) All() []*ServerList {
	return c.objList
}

func (c *serverListCache) Count() int {
	return len(c.objList)
}

func (c *serverListCache) Get(svr_id int, svr_type int) (*ServerList, bool) {
	return c.GetKey2(svr_id, svr_type)
}

func (c *serverListCache) GetKey1(svr_id int) (map[int]*ServerList, bool) {
	v, ok := c.objMap[svr_id]
	return v, ok
}

func (c *serverListCache) GetKey2(svr_id int, svr_type int) (*ServerList, bool) {
	v, ok := c.objMap[svr_id]
	if !ok {
		return nil, false
	}
	v1, ok1 := v[svr_type]
	return v1, ok1
}
