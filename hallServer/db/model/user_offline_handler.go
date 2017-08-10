package model

import (
	"errors"
	"fmt"
	"mj/hallServer/db"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lovelly/leaf/log"
)

//This file is generate by scripts,don't edit it

//user_offline_handler
//

// +gen *
type UserOfflineHandler struct {
	Id         int        `db:"id" json:"id"`                   //
	UserId     int64      `db:"user_id" json:"user_id"`         //
	HType      string     `db:"h_type" json:"h_type"`           //
	Context    string     `db:"context" json:"context"`         //
	ExpiryTime *time.Time `db:"expiry_time" json:"expiry_time"` //
}

type userOfflineHandlerOp struct{}

var UserOfflineHandlerOp = &userOfflineHandlerOp{}
var DefaultUserOfflineHandler = &UserOfflineHandler{}

// 按主键查询. 注:未找到记录的话将触发sql.ErrNoRows错误，返回nil, false
func (op *userOfflineHandlerOp) Get(id int) (*UserOfflineHandler, bool) {
	obj := &UserOfflineHandler{}
	sql := "select * from user_offline_handler where id=? "
	err := db.DB.Get(obj, sql,
		id,
	)

	if err != nil {
		log.Error("Get data error:%v", err.Error())
		return nil, false
	}
	return obj, true
}
func (op *userOfflineHandlerOp) SelectAll() ([]*UserOfflineHandler, error) {
	objList := []*UserOfflineHandler{}
	sql := "select * from user_offline_handler "
	err := db.DB.Select(&objList, sql)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return objList, nil
}

func (op *userOfflineHandlerOp) QueryByMap(m map[string]interface{}) ([]*UserOfflineHandler, error) {
	result := []*UserOfflineHandler{}
	var params []interface{}

	sql := "select * from user_offline_handler where 1=1 "
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

func (op *userOfflineHandlerOp) GetByMap(m map[string]interface{}) (*UserOfflineHandler, error) {
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
func (i *UserOfflineHandler) Insert() error {
    err := db.DBMap.Insert(i)
    if err != nil{
		log.Error("Insert sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 插入数据，自增长字段将被忽略
func (op *userOfflineHandlerOp) Insert(m *UserOfflineHandler) (int64, error) {
	return op.InsertTx(db.DB, m)
}

// 插入数据，自增长字段将被忽略
func (op *userOfflineHandlerOp) InsertTx(ext sqlx.Ext, m *UserOfflineHandler) (int64, error) {
	sql := "insert into user_offline_handler(user_id,h_type,context,expiry_time) values(?,?,?,?)"
	result, err := ext.Exec(sql,
		m.UserId,
		m.HType,
		m.Context,
		m.ExpiryTime,
	)
	if err != nil {
		log.Error("InsertTx sql error:%v, data:%v", err.Error(), m)
		return -1, err
	}
	affected, _ := result.LastInsertId()
	return affected, nil
}

//存在就更新， 不存在就插入
func (op *userOfflineHandlerOp) InsertUpdate(obj *UserOfflineHandler, m map[string]interface{}) error {
	sql := "insert into user_offline_handler(user_id,h_type,context,expiry_time) values(?,?,?,?) ON DUPLICATE KEY UPDATE "
	var params = []interface{}{obj.UserId,
		obj.HType,
		obj.Context,
		obj.ExpiryTime,
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
func (i *UserOfflineHandler) Update()  error {
    _,err := db.DBMap.Update(i)
    if err != nil{
		log.Error("update sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *userOfflineHandlerOp) Update(m *UserOfflineHandler) error {
	return op.UpdateTx(db.DB, m)
}

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *userOfflineHandlerOp) UpdateTx(ext sqlx.Ext, m *UserOfflineHandler) error {
	sql := `update user_offline_handler set user_id=?,h_type=?,context=?,expiry_time=? where id=?`
	_, err := ext.Exec(sql,
		m.UserId,
		m.HType,
		m.Context,
		m.ExpiryTime,
		m.Id,
	)

	if err != nil {
		log.Error("update sql error:%v, data:%v", err.Error(), m)
		return err
	}

	return nil
}

// 用主键做条件，更新map里包含的字段名
func (op *userOfflineHandlerOp) UpdateWithMap(id int, m map[string]interface{}) error {
	return op.UpdateWithMapTx(db.DB, id, m)
}

// 用主键做条件，更新map里包含的字段名
func (op *userOfflineHandlerOp) UpdateWithMapTx(ext sqlx.Ext, id int, m map[string]interface{}) error {

	sql := `update user_offline_handler set %s where 1=1 and id=? ;`

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
func (i *UserOfflineHandler) Delete() error{
    _,err := db.DBMap.Delete(i)
	log.Error("Delete sql error:%v", err.Error())
    return err
}
*/
// 根据主键删除相关记录
func (op *userOfflineHandlerOp) Delete(id int) error {
	return op.DeleteTx(db.DB, id)
}

// 根据主键删除相关记录,Tx
func (op *userOfflineHandlerOp) DeleteTx(ext sqlx.Ext, id int) error {
	sql := `delete from user_offline_handler where 1=1
        and id=?
        `
	_, err := ext.Exec(sql,
		id,
	)
	return err
}

// 返回符合查询条件的记录数
func (op *userOfflineHandlerOp) CountByMap(m map[string]interface{}) (int64, error) {

	var params []interface{}
	sql := `select count(*) from user_offline_handler where 1=1 `
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

func (op *userOfflineHandlerOp) DeleteByMap(m map[string]interface{}) (int64, error) {
	return op.DeleteByMapTx(db.DB, m)
}

func (op *userOfflineHandlerOp) DeleteByMapTx(ext sqlx.Ext, m map[string]interface{}) (int64, error) {
	var params []interface{}
	sql := "delete from user_offline_handler where 1=1 "
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
