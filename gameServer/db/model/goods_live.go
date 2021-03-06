package model

import (
	"fmt"
	"mj/gameServer/db"

	"github.com/jmoiron/sqlx"
	"github.com/lovelly/leaf/log"
)

//This file is generate by scripts,don't edit it

//goods_live
//

// +gen *
type GoodsLive struct {
	Id         int   `db:"id" json:"id"`                   // 物品id
	LeftAmount int   `db:"left_amount" json:"left_amount"` // 剩余的数量
	TradeTime  int64 `db:"trade_time" json:"trade_time"`   // 交易次数
}

type goodsLiveOp struct{}

var GoodsLiveOp = &goodsLiveOp{}
var DefaultGoodsLive = &GoodsLive{}

// 按主键查询. 注:未找到记录的话将触发sql.ErrNoRows错误，返回nil, false
func (op *goodsLiveOp) Get(id int) (*GoodsLive, bool) {
	obj := &GoodsLive{}
	sql := "select * from goods_live where id=? "
	err := db.DB.Get(obj, sql,
		id,
	)

	if err != nil {
		log.Error("Get data error:%v", err.Error())
		return nil, false
	}
	return obj, true
}
func (op *goodsLiveOp) SelectAll() ([]*GoodsLive, error) {
	objList := []*GoodsLive{}
	sql := "select * from goods_live "
	err := db.DB.Select(&objList, sql)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return objList, nil
}

func (op *goodsLiveOp) QueryByMap(m map[string]interface{}) ([]*GoodsLive, error) {
	result := []*GoodsLive{}
	var params []interface{}

	sql := "select * from goods_live where 1=1 "
	for k, v := range m {
		sql += fmt.Sprintf(" and %s=? ", k)
		params = append(params, v)
	}
	err := db.DB.Select(&result, sql, params...)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return result, nil
}

func (op *goodsLiveOp) GetByMap(m map[string]interface{}) (*GoodsLive, error) {
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
func (i *GoodsLive) Insert() error {
    err := db.DBMap.Insert(i)
    if err != nil{
		log.Error("Insert sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 插入数据，自增长字段将被忽略
func (op *goodsLiveOp) Insert(m *GoodsLive) (int64, error) {
	return op.InsertTx(db.DB, m)
}

// 插入数据，自增长字段将被忽略
func (op *goodsLiveOp) InsertTx(ext sqlx.Ext, m *GoodsLive) (int64, error) {
	sql := "insert into goods_live(id,left_amount,trade_time) values(?,?,?)"
	result, err := ext.Exec(sql,
		m.Id,
		m.LeftAmount,
		m.TradeTime,
	)
	if err != nil {
		log.Error("InsertTx sql error:%v, data:%v", err.Error(), m)
		return -1, err
	}
	affected, _ := result.LastInsertId()
	return affected, nil
}

//存在就更新， 不存在就插入
func (op *goodsLiveOp) InsertUpdate(obj *GoodsLive, m map[string]interface{}) error {
	sql := "insert into goods_live(id,left_amount,trade_time) values(?,?,?) ON DUPLICATE KEY UPDATE "
	var params = []interface{}{obj.Id,
		obj.LeftAmount,
		obj.TradeTime,
	}
	var set_sql string
	for k, v := range m {
		if set_sql != "" {
			set_sql += ","
		}
		set_sql += fmt.Sprintf(" %s=? ", k)
		params = append(params, v)
	}

	_, err := db.DB.Exec(sql+set_sql, params...)
	return err
}

/*
func (i *GoodsLive) Update()  error {
    _,err := db.DBMap.Update(i)
    if err != nil{
		log.Error("update sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *goodsLiveOp) Update(m *GoodsLive) error {
	return op.UpdateTx(db.DB, m)
}

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *goodsLiveOp) UpdateTx(ext sqlx.Ext, m *GoodsLive) error {
	sql := `update goods_live set left_amount=?,trade_time=? where id=?`
	_, err := ext.Exec(sql,
		m.LeftAmount,
		m.TradeTime,
		m.Id,
	)

	if err != nil {
		log.Error("update sql error:%v, data:%v", err.Error(), m)
		return err
	}

	return nil
}

// 用主键做条件，更新map里包含的字段名
func (op *goodsLiveOp) UpdateWithMap(id int, m map[string]interface{}) error {
	return op.UpdateWithMapTx(db.DB, id, m)
}

// 用主键做条件，更新map里包含的字段名
func (op *goodsLiveOp) UpdateWithMapTx(ext sqlx.Ext, id int, m map[string]interface{}) error {

	sql := `update goods_live set %s where 1=1 and id=? ;`

	var params []interface{}
	var set_sql string
	for k, v := range m {
		if set_sql != "" {
			set_sql += ","
		}
		set_sql += fmt.Sprintf(" %s=? ", k)
		params = append(params, v)
	}
	params = append(params, id)
	_, err := ext.Exec(fmt.Sprintf(sql, set_sql), params...)
	return err
}

/*
func (i *GoodsLive) Delete() error{
    _,err := db.DBMap.Delete(i)
	log.Error("Delete sql error:%v", err.Error())
    return err
}
*/
// 根据主键删除相关记录
func (op *goodsLiveOp) Delete(id int) error {
	return op.DeleteTx(db.DB, id)
}

// 根据主键删除相关记录,Tx
func (op *goodsLiveOp) DeleteTx(ext sqlx.Ext, id int) error {
	sql := `delete from goods_live where 1=1
        and id=?
        `
	_, err := ext.Exec(sql,
		id,
	)
	return err
}

// 返回符合查询条件的记录数
func (op *goodsLiveOp) CountByMap(m map[string]interface{}) (int64, error) {

	var params []interface{}
	sql := `select count(*) from goods_live where 1=1 `
	for k, v := range m {
		sql += fmt.Sprintf(" and  %s=? ", k)
		params = append(params, v)
	}
	count := int64(-1)
	err := db.DB.Get(&count, sql, params...)
	if err != nil {
		log.Error("CountByMap  error:%v data :%v", err.Error(), m)
		return 0, err
	}
	return count, nil
}

func (op *goodsLiveOp) DeleteByMap(m map[string]interface{}) (int64, error) {
	return op.DeleteByMapTx(db.DB, m)
}

func (op *goodsLiveOp) DeleteByMapTx(ext sqlx.Ext, m map[string]interface{}) (int64, error) {
	var params []interface{}
	sql := "delete from goods_live where 1=1 "
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
