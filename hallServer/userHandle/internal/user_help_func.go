package internal

import (
	"errors"
	"fmt"
	. "mj/common/cost"
	"mj/hallServer/common"
	"mj/hallServer/db"
	"mj/hallServer/http_service"
	"mj/hallServer/user"
)

type Order struct {
	OnLineID    int    `db:"OnLineID"`    //订单标识
	PayAmount   int    `db:"PayAmount"`   // 价格
	UserID      int64  `db:"UserID"`      // uid
	PayType     string `db：dayType"`      //支付类型
	GoodsID     int    `db:"GoodsID"`     // 物品id
	OrderStatus int    `db:"OrderStatus"` //状态
}

func (m *UserModule) GetUser(args []interface{}) (interface{}, error) {
	u, ok := m.a.UserData().(*user.User)
	if !ok {
		return nil, errors.New("not foud user Data at GetUser")
	}
	return u, nil
}

func GetOrders(uid int64) []*Order {
	obj := make([]*Order, 0)
	err := db.AccountDB.Select(&obj, "select OnLineID,PayAmount, UserID, PayType, GoodsID, OrderStatus FROM `order` where UserID=? and OrderStatus = 1", uid)
	if err != nil {
		return nil
	}

	if len(obj) < 1 {
		return nil
	}
	return obj
}

func UpdateOrderStats(OnLineID int) bool {
	_, err := db.AccountDB.Exec(`UPDATE onlineorder SET OrderStatus=2 WHERE OnLineID = ?;`, OnLineID)
	if err != nil {
		return false
	}
	return true
}

func ReqGetMaskCode(phome string, maskCode int) {
	http_service.PostJSON("https://sms.yunpian.com/v2/sms/single_send.json", map[string]interface{}{
		"apikey": "fce482d259d86ca9b0490d400889a9b8",
		"mobile": phome,
		"text":   fmt.Sprintf(common.GetGlobalVar(MASK_CODE_TEXT), maskCode),
	})
}
