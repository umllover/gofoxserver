package stats

import (
	"fmt"
	"mj/gameServer/db"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lovelly/leaf/log"
)

//This file is generate by scripts,don't edit it

//recharge_log
//

// +gen *
type RechargeLog struct {
	OnLineID     int        `db:"OnLineID" json:"OnLineID"`         // 订单标识
	PayAmount    int        `db:"PayAmount" json:"PayAmount"`       // 实付金额
	UserID       int64      `db:"UserID" json:"UserID"`             // 用户标识
	PayType      string     `db:"PayType" json:"PayType"`           // 支付类型
	GoodsID      int        `db:"GoodsID" json:"GoodsID"`           // 物品id
	RechangeTime *time.Time `db:"RechangeTime" json:"RechangeTime"` // 冲值时间
}

type rechargeLogOp struct{}

var RechargeLogOp = &rechargeLogOp{}
var DefaultRechargeLog = &RechargeLog{}

// 按主键查询. 注:未找到记录的话将触发sql.ErrNoRows错误，返回nil, false
func (op *rechargeLogOp) Get(OnLineID int) (*RechargeLog, bool) {
	obj := &RechargeLog{}
	sql := "select * from recharge_log where OnLineID=? "
	err := db.StatsDB.Get(obj, sql,
		OnLineID,
	)

	if err != nil {
		log.Error("Get data error:%v", err.Error())
		return nil, false
	}
	return obj, true
}
func (op *rechargeLogOp) SelectAll() ([]*RechargeLog, error) {
	objList := []*RechargeLog{}
	sql := "select * from recharge_log "
	err := db.StatsDB.Select(&objList, sql)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return objList, nil
}

func (op *rechargeLogOp) QueryByMap(m map[string]interface{}) ([]*RechargeLog, error) {
	result := []*RechargeLog{}
	var params []interface{}

	sql := "select * from recharge_log where 1=1 "
	for k, v := range m {
		sql += fmt.Sprintf(" and %s=? ", k)
		params = append(params, v)
	}
	err := db.StatsDB.Select(&result, sql, params...)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return result, nil
}

func (op *rechargeLogOp) GetByMap(m map[string]interface{}) (*RechargeLog, error) {
	lst, err := op.QueryByMap(m)
	if err != nil {
		return nil, err
	}
	if len(lst) > 0 {
		return lst[0], nil
	}
	return nil, nil
}

/*
func (i *RechargeLog) Insert() error {
    err := db.StatsDBMap.Insert(i)
    if err != nil{
		log.Error("Insert sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 插入数据，自增长字段将被忽略
func (op *rechargeLogOp) Insert(m *RechargeLog) (int64, error) {
	return op.InsertTx(db.StatsDB, m)
}

// 插入数据，自增长字段将被忽略
func (op *rechargeLogOp) InsertTx(ext sqlx.Ext, m *RechargeLog) (int64, error) {
	sql := "insert into recharge_log(OnLineID,PayAmount,UserID,PayType,GoodsID,RechangeTime) values(?,?,?,?,?,?)"
	result, err := ext.Exec(sql,
		m.OnLineID,
		m.PayAmount,
		m.UserID,
		m.PayType,
		m.GoodsID,
		m.RechangeTime,
	)
	if err != nil {
		log.Error("InsertTx sql error:%v, data:%v", err.Error(), m)
		return -1, err
	}
	affected, _ := result.LastInsertId()
	return affected, nil
}

//存在就更新， 不存在就插入
func (op *rechargeLogOp) InsertUpdate(obj *RechargeLog, m map[string]interface{}) error {
	sql := "insert into recharge_log(OnLineID,PayAmount,UserID,PayType,GoodsID,RechangeTime) values(?,?,?,?,?,?) ON DUPLICATE KEY UPDATE "
	var params = []interface{}{obj.OnLineID,
		obj.PayAmount,
		obj.UserID,
		obj.PayType,
		obj.GoodsID,
		obj.RechangeTime,
	}
	var set_sql string
	for k, v := range m {
		if set_sql != "" {
			set_sql += ","
		}
		set_sql += fmt.Sprintf(" %s=? ", k)
		params = append(params, v)
	}

	_, err := db.StatsDB.Exec(sql+set_sql, params...)
	return err
}

/*
func (i *RechargeLog) Update()  error {
    _,err := db.StatsDBMap.Update(i)
    if err != nil{
		log.Error("update sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *rechargeLogOp) Update(m *RechargeLog) error {
	return op.UpdateTx(db.StatsDB, m)
}

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *rechargeLogOp) UpdateTx(ext sqlx.Ext, m *RechargeLog) error {
	sql := `update recharge_log set PayAmount=?,UserID=?,PayType=?,GoodsID=?,RechangeTime=? where OnLineID=?`
	_, err := ext.Exec(sql,
		m.PayAmount,
		m.UserID,
		m.PayType,
		m.GoodsID,
		m.RechangeTime,
		m.OnLineID,
	)

	if err != nil {
		log.Error("update sql error:%v, data:%v", err.Error(), m)
		return err
	}

	return nil
}

// 用主键做条件，更新map里包含的字段名
func (op *rechargeLogOp) UpdateWithMap(OnLineID int, m map[string]interface{}) error {
	return op.UpdateWithMapTx(db.StatsDB, OnLineID, m)
}

// 用主键做条件，更新map里包含的字段名
func (op *rechargeLogOp) UpdateWithMapTx(ext sqlx.Ext, OnLineID int, m map[string]interface{}) error {

	sql := `update recharge_log set %s where 1=1 and OnLineID=? ;`

	var params []interface{}
	var set_sql string
	for k, v := range m {
		if set_sql != "" {
			set_sql += ","
		}
		set_sql += fmt.Sprintf(" %s=? ", k)
		params = append(params, v)
	}
	params = append(params, OnLineID)
	_, err := ext.Exec(fmt.Sprintf(sql, set_sql), params...)
	return err
}

/*
func (i *RechargeLog) Delete() error{
    _,err := db.StatsDBMap.Delete(i)
	log.Error("Delete sql error:%v", err.Error())
    return err
}
*/
// 根据主键删除相关记录
func (op *rechargeLogOp) Delete(OnLineID int) error {
	return op.DeleteTx(db.StatsDB, OnLineID)
}

// 根据主键删除相关记录,Tx
func (op *rechargeLogOp) DeleteTx(ext sqlx.Ext, OnLineID int) error {
	sql := `delete from recharge_log where 1=1
        and OnLineID=?
        `
	_, err := ext.Exec(sql,
		OnLineID,
	)
	return err
}

// 返回符合查询条件的记录数
func (op *rechargeLogOp) CountByMap(m map[string]interface{}) (int64, error) {

	var params []interface{}
	sql := `select count(*) from recharge_log where 1=1 `
	for k, v := range m {
		sql += fmt.Sprintf(" and  %s=? ", k)
		params = append(params, v)
	}
	count := int64(-1)
	err := db.StatsDB.Get(&count, sql, params...)
	if err != nil {
		log.Error("CountByMap  error:%v data :%v", err.Error(), m)
		return 0, err
	}
	return count, nil
}

func (op *rechargeLogOp) DeleteByMap(m map[string]interface{}) (int64, error) {
	return op.DeleteByMapTx(db.StatsDB, m)
}

func (op *rechargeLogOp) DeleteByMapTx(ext sqlx.Ext, m map[string]interface{}) (int64, error) {
	var params []interface{}
	sql := "delete from recharge_log where 1=1 "
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
