package account

import (
	"errors"
	"fmt"
	"mj/gameServer/db"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lovelly/leaf/log"
)

//This file is generate by scripts,don't edit it

//onlineorder
//

// +gen *
type Onlineorder struct {
	OnLineId      int        `db:"onLine_id" json:"onLine_id"`           // 订单标识
	UserId        int64      `db:"user_id" json:"user_id"`               // 用户标识
	OrderId       int64      `db:"order_id" json:"order_id"`             // 订单号码
	PayAmount     int        `db:"pay_amount" json:"pay_amount"`         // 实付金额
	OrderStatus   int8       `db:"order_status" json:"order_status"`     // 订单状态  0:未付款;1:已付款待处理;2:处理完成
	IpAddress     string     `db:"ip_address" json:"ip_address"`         // 订单地址
	ApplyDate     *time.Time `db:"apply_date" json:"apply_date"`         // 订单日期
	GoodsId       int        `db:"goods_id" json:"goods_id"`             //
	PayType       string     `db:"pay_type" json:"pay_type"`             // 支付类型
	AgentNum      string     `db:"agent_num" json:"agent_num"`           //
	PrepayId      string     `db:"prepay_id" json:"prepay_id"`           //
	AgentId       int        `db:"agent_id" json:"agent_id"`             //
	TransactionId string     `db:"transaction_id" json:"transaction_id"` //
	IsSettle      int8       `db:"is_settle" json:"is_settle"`           //
	Quantity      int        `db:"quantity" json:"quantity"`             //
}

type onlineorderOp struct{}

var OnlineorderOp = &onlineorderOp{}
var DefaultOnlineorder = &Onlineorder{}

// 按主键查询. 注:未找到记录的话将触发sql.ErrNoRows错误，返回nil, false
func (op *onlineorderOp) Get(onLine_id int) (*Onlineorder, bool) {
	obj := &Onlineorder{}
	sql := "select * from onlineorder where onLine_id=? "
	err := db.AccountDB.Get(obj, sql,
		onLine_id,
	)

	if err != nil {
		log.Error("Get data error:%v", err.Error())
		return nil, false
	}
	return obj, true
}
func (op *onlineorderOp) SelectAll() ([]*Onlineorder, error) {
	objList := []*Onlineorder{}
	sql := "select * from onlineorder "
	err := db.AccountDB.Select(&objList, sql)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return objList, nil
}

func (op *onlineorderOp) QueryByMap(m map[string]interface{}) ([]*Onlineorder, error) {
	result := []*Onlineorder{}
	var params []interface{}

	sql := "select * from onlineorder where 1=1 "
	for k, v := range m {
		sql += fmt.Sprintf(" and %s=? ", k)
		params = append(params, v)
	}
	err := db.AccountDB.Select(&result, sql, params...)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return result, nil
}

func (op *onlineorderOp) GetByMap(m map[string]interface{}) (*Onlineorder, error) {
	lst, err := op.QueryByMap(m)
	if err != nil {
		return nil, err
	}
	if len(lst) > 0 {
		return lst[0], nil
	}
	return nil, errors.New("no row in result")
}

/*
func (i *Onlineorder) Insert() error {
    err := db.AccountDBMap.Insert(i)
    if err != nil{
		log.Error("Insert sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 插入数据，自增长字段将被忽略
func (op *onlineorderOp) Insert(m *Onlineorder) (int64, error) {
	return op.InsertTx(db.AccountDB, m)
}

// 插入数据，自增长字段将被忽略
func (op *onlineorderOp) InsertTx(ext sqlx.Ext, m *Onlineorder) (int64, error) {
	sql := "insert into onlineorder(user_id,order_id,pay_amount,order_status,ip_address,apply_date,goods_id,pay_type,agent_num,prepay_id,agent_id,transaction_id,is_settle,quantity) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	result, err := ext.Exec(sql,
		m.UserId,
		m.OrderId,
		m.PayAmount,
		m.OrderStatus,
		m.IpAddress,
		m.ApplyDate,
		m.GoodsId,
		m.PayType,
		m.AgentNum,
		m.PrepayId,
		m.AgentId,
		m.TransactionId,
		m.IsSettle,
		m.Quantity,
	)
	if err != nil {
		log.Error("InsertTx sql error:%v, data:%v", err.Error(), m)
		return -1, err
	}
	affected, _ := result.LastInsertId()
	return affected, nil
}

//存在就更新， 不存在就插入
func (op *onlineorderOp) InsertUpdate(obj *Onlineorder, m map[string]interface{}) error {
	sql := "insert into onlineorder(user_id,order_id,pay_amount,order_status,ip_address,apply_date,goods_id,pay_type,agent_num,prepay_id,agent_id,transaction_id,is_settle,quantity) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE "
	var params = []interface{}{obj.UserId,
		obj.OrderId,
		obj.PayAmount,
		obj.OrderStatus,
		obj.IpAddress,
		obj.ApplyDate,
		obj.GoodsId,
		obj.PayType,
		obj.AgentNum,
		obj.PrepayId,
		obj.AgentId,
		obj.TransactionId,
		obj.IsSettle,
		obj.Quantity,
	}
	var set_sql string
	for k, v := range m {
		if set_sql != "" {
			set_sql += ","
		}
		set_sql += fmt.Sprintf(" %s=? ", k)
		params = append(params, v)
	}

	_, err := db.AccountDB.Exec(sql+set_sql, params...)
	return err
}

/*
func (i *Onlineorder) Update()  error {
    _,err := db.AccountDBMap.Update(i)
    if err != nil{
		log.Error("update sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *onlineorderOp) Update(m *Onlineorder) error {
	return op.UpdateTx(db.AccountDB, m)
}

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *onlineorderOp) UpdateTx(ext sqlx.Ext, m *Onlineorder) error {
	sql := `update onlineorder set user_id=?,order_id=?,pay_amount=?,order_status=?,ip_address=?,apply_date=?,goods_id=?,pay_type=?,agent_num=?,prepay_id=?,agent_id=?,transaction_id=?,is_settle=?,quantity=? where onLine_id=?`
	_, err := ext.Exec(sql,
		m.UserId,
		m.OrderId,
		m.PayAmount,
		m.OrderStatus,
		m.IpAddress,
		m.ApplyDate,
		m.GoodsId,
		m.PayType,
		m.AgentNum,
		m.PrepayId,
		m.AgentId,
		m.TransactionId,
		m.IsSettle,
		m.Quantity,
		m.OnLineId,
	)

	if err != nil {
		log.Error("update sql error:%v, data:%v", err.Error(), m)
		return err
	}

	return nil
}

// 用主键做条件，更新map里包含的字段名
func (op *onlineorderOp) UpdateWithMap(onLine_id int, m map[string]interface{}) error {
	return op.UpdateWithMapTx(db.AccountDB, onLine_id, m)
}

// 用主键做条件，更新map里包含的字段名
func (op *onlineorderOp) UpdateWithMapTx(ext sqlx.Ext, onLine_id int, m map[string]interface{}) error {

	sql := `update onlineorder set %s where 1=1 and onLine_id=? ;`

	var params []interface{}
	var set_sql string
	for k, v := range m {
		if set_sql != "" {
			set_sql += ","
		}
		set_sql += fmt.Sprintf(" %s=? ", k)
		params = append(params, v)
	}
	params = append(params, onLine_id)
	_, err := ext.Exec(fmt.Sprintf(sql, set_sql), params...)
	return err
}

/*
func (i *Onlineorder) Delete() error{
    _,err := db.AccountDBMap.Delete(i)
	log.Error("Delete sql error:%v", err.Error())
    return err
}
*/
// 根据主键删除相关记录
func (op *onlineorderOp) Delete(onLine_id int) error {
	return op.DeleteTx(db.AccountDB, onLine_id)
}

// 根据主键删除相关记录,Tx
func (op *onlineorderOp) DeleteTx(ext sqlx.Ext, onLine_id int) error {
	sql := `delete from onlineorder where 1=1
        and onLine_id=?
        `
	_, err := ext.Exec(sql,
		onLine_id,
	)
	return err
}

// 返回符合查询条件的记录数
func (op *onlineorderOp) CountByMap(m map[string]interface{}) (int64, error) {

	var params []interface{}
	sql := `select count(*) from onlineorder where 1=1 `
	for k, v := range m {
		sql += fmt.Sprintf(" and  %s=? ", k)
		params = append(params, v)
	}
	count := int64(-1)
	err := db.AccountDB.Get(&count, sql, params...)
	if err != nil {
		log.Error("CountByMap  error:%v data :%v", err.Error(), m)
		return 0, err
	}
	return count, nil
}

func (op *onlineorderOp) DeleteByMap(m map[string]interface{}) (int64, error) {
	return op.DeleteByMapTx(db.AccountDB, m)
}

func (op *onlineorderOp) DeleteByMapTx(ext sqlx.Ext, m map[string]interface{}) (int64, error) {
	var params []interface{}
	sql := "delete from onlineorder where 1=1 "
	for k, v := range m {
		sql += fmt.Sprintf(" and %s=? ", k)
		params = append(params, v)
	}
	result, err := ext.Exec(sql, params...)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}
