package model

import (
	"errors"
	"fmt"
	"mj/hallServer/db"

	"github.com/jmoiron/sqlx"
	"github.com/lovelly/leaf/log"
)

//This file is generate by scripts,don't edit it

//user_day_times
//

// +gen *
type UserDayTimes struct {
	UserId int   `db:"user_id" json:"user_id"` //
	KeyId  int   `db:"key_id" json:"key_id"`   //
	V      int64 `db:"v" json:"v"`             //
}

type userDayTimesOp struct{}

var UserDayTimesOp = &userDayTimesOp{}
var DefaultUserDayTimes = &UserDayTimes{}

// 按主键查询. 注:未找到记录的话将触发sql.ErrNoRows错误，返回nil, false
func (op *userDayTimesOp) Get(user_id int, key_id int) (*UserDayTimes, bool) {
	obj := &UserDayTimes{}
	sql := "select * from user_day_times where user_id=? and key_id=? "
	err := db.DB.Get(obj, sql,
		user_id,
		key_id,
	)

	if err != nil {
		log.Error("Get data error:%v", err.Error())
		return nil, false
	}
	return obj, true
}
func (op *userDayTimesOp) SelectAll() ([]*UserDayTimes, error) {
	objList := []*UserDayTimes{}
	sql := "select * from user_day_times "
	err := db.DB.Select(&objList, sql)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return objList, nil
}

func (op *userDayTimesOp) QueryByMap(m map[string]interface{}) ([]*UserDayTimes, error) {
	result := []*UserDayTimes{}
	var params []interface{}

	sql := "select * from user_day_times where 1=1 "
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

func (op *userDayTimesOp) QueryByMapQueryByMapComparison(m map[string]interface{}) ([]*UserDayTimes, error) {
	result := []*UserDayTimes{}
	var params []interface{}

	sql := "select * from user_day_times where 1=1 "
	for k, v := range m {
		sql += fmt.Sprintf(" and %s? ", k)
		params = append(params, v)
	}
	err := db.DB.Select(&result, sql, params...)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return result, nil
}

func (op *userDayTimesOp) GetByMap(m map[string]interface{}) (*UserDayTimes, error) {
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
func (i *UserDayTimes) Insert() error {
    err := db.DBMap.Insert(i)
    if err != nil{
		log.Error("Insert sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 插入数据，自增长字段将被忽略
func (op *userDayTimesOp) Insert(m *UserDayTimes) (int64, error) {
	return op.InsertTx(db.DB, m)
}

// 插入数据，自增长字段将被忽略
func (op *userDayTimesOp) InsertTx(ext sqlx.Ext, m *UserDayTimes) (int64, error) {
	sql := "insert into user_day_times(user_id,key_id,v) values(?,?,?)"
	result, err := ext.Exec(sql,
		m.UserId,
		m.KeyId,
		m.V,
	)
	if err != nil {
		log.Error("InsertTx sql error:%v, data:%v", err.Error(), m)
		return -1, err
	}
	affected, _ := result.LastInsertId()
	return affected, nil
}

/*
func (i *UserDayTimes) Update()  error {
    _,err := db.DBMap.Update(i)
    if err != nil{
		log.Error("update sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *userDayTimesOp) Update(m *UserDayTimes) error {
	return op.UpdateTx(db.DB, m)
}

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *userDayTimesOp) UpdateTx(ext sqlx.Ext, m *UserDayTimes) error {
	sql := `update user_day_times set v=? where user_id=? and key_id=?`
	_, err := ext.Exec(sql,
		m.V,
		m.UserId,
		m.KeyId,
	)

	if err != nil {
		log.Error("update sql error:%v, data:%v", err.Error(), m)
		return err
	}

	return nil
}

// 用主键做条件，更新map里包含的字段名
func (op *userDayTimesOp) UpdateWithMap(user_id int, key_id int, m map[string]interface{}) error {
	return op.UpdateWithMapTx(db.DB, user_id, key_id, m)
}

// 用主键做条件，更新map里包含的字段名
func (op *userDayTimesOp) UpdateWithMapTx(ext sqlx.Ext, user_id int, key_id int, m map[string]interface{}) error {

	sql := `update user_day_times set %s where 1=1 and user_id=? and key_id=? ;`

	var params []interface{}
	var set_sql string
	for k, v := range m {
		if set_sql != "" {
			set_sql += ","
		}
		set_sql += fmt.Sprintf(" %s=? ", k)
		params = append(params, v)
	}
	params = append(params, user_id, key_id)
	_, err := ext.Exec(fmt.Sprintf(sql, set_sql), params...)
	return err
}

/*
func (i *UserDayTimes) Delete() error{
    _,err := db.DBMap.Delete(i)
	log.Error("Delete sql error:%v", err.Error())
    return err
}
*/
// 根据主键删除相关记录
func (op *userDayTimesOp) Delete(user_id int, key_id int) error {
	return op.DeleteTx(db.DB, user_id, key_id)
}

// 根据主键删除相关记录,Tx
func (op *userDayTimesOp) DeleteTx(ext sqlx.Ext, user_id int, key_id int) error {
	sql := `delete from user_day_times where 1=1
        and user_id=?
        and key_id=?
        `
	_, err := ext.Exec(sql,
		user_id,
		key_id,
	)
	return err
}

// 返回符合查询条件的记录数
func (op *userDayTimesOp) CountByMap(m map[string]interface{}) (int64, error) {

	var params []interface{}
	sql := `select count(*) from user_day_times where 1=1 `
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

func (op *userDayTimesOp) DeleteByMap(m map[string]interface{}) (int64, error) {
	return op.DeleteByMapTx(db.DB, m)
}

func (op *userDayTimesOp) DeleteByMapTx(ext sqlx.Ext, m map[string]interface{}) (int64, error) {
	var params []interface{}
	sql := "delete from user_day_times where 1=1 "
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
