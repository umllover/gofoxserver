package model

import (
	"errors"
	"fmt"
	"mj/hallServer/db"

	"github.com/jmoiron/sqlx"
	"github.com/lovelly/leaf/log"
)

//This file is generate by scripts,don't edit it

//user_spread
//

// +gen *
type UserSpread struct {
	UserId    int `db:"user_id" json:"user_id"`       //
	SpreadUid int `db:"spread_uid" json:"spread_uid"` // 被我领取的推广人id
	Status    int `db:"status" json:"status"`         // 是否领取了这个人的奖励
}

type userSpreadOp struct{}

var UserSpreadOp = &userSpreadOp{}
var DefaultUserSpread = &UserSpread{}

// 按主键查询. 注:未找到记录的话将触发sql.ErrNoRows错误，返回nil, false
func (op *userSpreadOp) Get(user_id int, spread_uid int) (*UserSpread, bool) {
	obj := &UserSpread{}
	sql := "select * from user_spread where user_id=? and spread_uid=? "
	err := db.DB.Get(obj, sql,
		user_id,
		spread_uid,
	)

	if err != nil {
		log.Error("Get data error:%v", err.Error())
		return nil, false
	}
	return obj, true
}
func (op *userSpreadOp) SelectAll() ([]*UserSpread, error) {
	objList := []*UserSpread{}
	sql := "select * from user_spread "
	err := db.DB.Select(&objList, sql)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return objList, nil
}

func (op *userSpreadOp) QueryByMap(m map[string]interface{}) ([]*UserSpread, error) {
	result := []*UserSpread{}
	var params []interface{}

	sql := "select * from user_spread where 1=1 "
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

func (op *userSpreadOp) GetByMap(m map[string]interface{}) (*UserSpread, error) {
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
func (i *UserSpread) Insert() error {
    err := db.DBMap.Insert(i)
    if err != nil{
		log.Error("Insert sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 插入数据，自增长字段将被忽略
func (op *userSpreadOp) Insert(m *UserSpread) (int64, error) {
	return op.InsertTx(db.DB, m)
}

// 插入数据，自增长字段将被忽略
func (op *userSpreadOp) InsertTx(ext sqlx.Ext, m *UserSpread) (int64, error) {
	sql := "insert into user_spread(user_id,spread_uid,status) values(?,?,?)"
	result, err := ext.Exec(sql,
		m.UserId,
		m.SpreadUid,
		m.Status,
	)
	if err != nil {
		log.Error("InsertTx sql error:%v, data:%v", err.Error(), m)
		return -1, err
	}
	affected, _ := result.LastInsertId()
	return affected, nil
}

/*
func (i *UserSpread) Update()  error {
    _,err := db.DBMap.Update(i)
    if err != nil{
		log.Error("update sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *userSpreadOp) Update(m *UserSpread) error {
	return op.UpdateTx(db.DB, m)
}

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *userSpreadOp) UpdateTx(ext sqlx.Ext, m *UserSpread) error {
	sql := `update user_spread set status=? where user_id=? and spread_uid=?`
	_, err := ext.Exec(sql,
		m.Status,
		m.UserId,
		m.SpreadUid,
	)

	if err != nil {
		log.Error("update sql error:%v, data:%v", err.Error(), m)
		return err
	}

	return nil
}

// 用主键做条件，更新map里包含的字段名
func (op *userSpreadOp) UpdateWithMap(user_id int, spread_uid int, m map[string]interface{}) error {
	return op.UpdateWithMapTx(db.DB, user_id, spread_uid, m)
}

// 用主键做条件，更新map里包含的字段名
func (op *userSpreadOp) UpdateWithMapTx(ext sqlx.Ext, user_id int, spread_uid int, m map[string]interface{}) error {

	sql := `update user_spread set %s where 1=1 and user_id=? and spread_uid=? ;`

	var params []interface{}
	var set_sql string
	for k, v := range m {
		if set_sql != "" {
			set_sql += ","
		}
		set_sql += fmt.Sprintf(" %s=? ", k)
		params = append(params, v)
	}
	params = append(params, user_id, spread_uid)
	_, err := ext.Exec(fmt.Sprintf(sql, set_sql), params...)
	return err
}

/*
func (i *UserSpread) Delete() error{
    _,err := db.DBMap.Delete(i)
	log.Error("Delete sql error:%v", err.Error())
    return err
}
*/
// 根据主键删除相关记录
func (op *userSpreadOp) Delete(user_id int, spread_uid int) error {
	return op.DeleteTx(db.DB, user_id, spread_uid)
}

// 根据主键删除相关记录,Tx
func (op *userSpreadOp) DeleteTx(ext sqlx.Ext, user_id int, spread_uid int) error {
	sql := `delete from user_spread where 1=1
        and user_id=?
        and spread_uid=?
        `
	_, err := ext.Exec(sql,
		user_id,
		spread_uid,
	)
	return err
}

// 返回符合查询条件的记录数
func (op *userSpreadOp) CountByMap(m map[string]interface{}) (int64, error) {

	var params []interface{}
	sql := `select count(*) from user_spread where 1=1 `
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

func (op *userSpreadOp) DeleteByMap(m map[string]interface{}) (int64, error) {
	return op.DeleteByMapTx(db.DB, m)
}

func (op *userSpreadOp) DeleteByMapTx(ext sqlx.Ext, m map[string]interface{}) (int64, error) {
	var params []interface{}
	sql := "delete from user_spread where 1=1 "
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
