package internal

import (
	"errors"
	"mj/hallServer/user"
	"mj/hallServer/db"
)

type Order struct {
	Orderid  int `db:"id"`        //
	PlayerId int `db:"player_id"` // 玩家id
	ServerId int `db:"server_id"` // 服务器id
	Price    int `db:"price"`     // 价格
	Uid      int `db:"user_id"`   // uid
	StoneId  int `db:"stone_id"`  // 物品id
	Status   int `db:"status"`    //状态
}

func (m *UserModule) GetUser(args []interface{}) (interface{}, error) {
	u, ok := m.a.UserData().(*user.User)
	if !ok {
		return nil, errors.New("not foud user Data at GetUser")
	}
	return u, nil
}

func GetOrder(orderid int) *Order {
	obj := make([]*Order, 0)
	err := db.AccountDB.Select(&obj, "select id,player_id, server_id, price, user_id, stone_id, status FROM `order` where id=?", orderid)
	if err != nil {
		return nil
	}

	if len(obj) < 1 {
		return nil
	}
	return obj[0]
}
