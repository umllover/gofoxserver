package model

import (
	"fmt"
	"mj/gameServer/db"

	"github.com/jmoiron/sqlx"
	"github.com/lovelly/leaf/log"
)

//This file is generate by scripts,don't edit it

//user_room_record
//

// +gen *
type UserRoomRecord struct {
	UserId   int64 `db:"user_id" json:"user_id"`     // 视频id
	RecordId int   `db:"record_id" json:"record_id"` // 记录id
}

type userRoomRecordOp struct{}

var UserRoomRecordOp = &userRoomRecordOp{}
var DefaultUserRoomRecord = &UserRoomRecord{}

// 按主键查询. 注:未找到记录的话将触发sql.ErrNoRows错误，返回nil, false
func (op *userRoomRecordOp) Get(user_id int64, record_id int) (*UserRoomRecord, bool) {
	obj := &UserRoomRecord{}
	sql := "select * from user_room_record where user_id=? and record_id=? "
	err := db.DB.Get(obj, sql,
		user_id,
		record_id,
	)

	if err != nil {
		log.Error("Get data error:%v", err.Error())
		return nil, false
	}
	return obj, true
}
func (op *userRoomRecordOp) SelectAll() ([]*UserRoomRecord, error) {
	objList := []*UserRoomRecord{}
	sql := "select * from user_room_record "
	err := db.DB.Select(&objList, sql)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return objList, nil
}

func (op *userRoomRecordOp) QueryByMap(m map[string]interface{}) ([]*UserRoomRecord, error) {
	result := []*UserRoomRecord{}
	var params []interface{}

	sql := "select * from user_room_record where 1=1 "
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

func (op *userRoomRecordOp) GetByMap(m map[string]interface{}) (*UserRoomRecord, error) {
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
func (i *UserRoomRecord) Insert() error {
    err := db.DBMap.Insert(i)
    if err != nil{
		log.Error("Insert sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 插入数据，自增长字段将被忽略
func (op *userRoomRecordOp) Insert(m *UserRoomRecord) (int64, error) {
	return op.InsertTx(db.DB, m)
}

// 插入数据，自增长字段将被忽略
func (op *userRoomRecordOp) InsertTx(ext sqlx.Ext, m *UserRoomRecord) (int64, error) {
	sql := "insert into user_room_record(user_id,record_id) values(?,?)"
	result, err := ext.Exec(sql,
		m.UserId,
		m.RecordId,
	)
	if err != nil {
		log.Error("InsertTx sql error:%v, data:%v", err.Error(), m)
		return -1, err
	}
	affected, _ := result.LastInsertId()
	return affected, nil
}

//存在就更新， 不存在就插入
func (op *userRoomRecordOp) InsertUpdate(obj *UserRoomRecord, m map[string]interface{}) error {
	sql := "insert into user_room_record(user_id,record_id) values(?,?) ON DUPLICATE KEY UPDATE "
	var params = []interface{}{obj.UserId,
		obj.RecordId,
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
func (i *UserRoomRecord) Update()  error {
    _,err := db.DBMap.Update(i)
    if err != nil{
		log.Error("update sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *userRoomRecordOp) Update(m *UserRoomRecord) error {
	return op.UpdateTx(db.DB, m)
}

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *userRoomRecordOp) UpdateTx(ext sqlx.Ext, m *UserRoomRecord) error {
	sql := `update user_room_record set  where user_id=? and record_id=?`
	_, err := ext.Exec(sql,
		m.UserId,
		m.RecordId,
	)

	if err != nil {
		log.Error("update sql error:%v, data:%v", err.Error(), m)
		return err
	}

	return nil
}

// 用主键做条件，更新map里包含的字段名
func (op *userRoomRecordOp) UpdateWithMap(user_id int64, record_id int, m map[string]interface{}) error {
	return op.UpdateWithMapTx(db.DB, user_id, record_id, m)
}

// 用主键做条件，更新map里包含的字段名
func (op *userRoomRecordOp) UpdateWithMapTx(ext sqlx.Ext, user_id int64, record_id int, m map[string]interface{}) error {

	sql := `update user_room_record set %s where 1=1 and user_id=? and record_id=? ;`

	var params []interface{}
	var set_sql string
	for k, v := range m {
		if set_sql != "" {
			set_sql += ","
		}
		set_sql += fmt.Sprintf(" %s=? ", k)
		params = append(params, v)
	}
	params = append(params, user_id, record_id)
	_, err := ext.Exec(fmt.Sprintf(sql, set_sql), params...)
	return err
}

/*
func (i *UserRoomRecord) Delete() error{
    _,err := db.DBMap.Delete(i)
	log.Error("Delete sql error:%v", err.Error())
    return err
}
*/
// 根据主键删除相关记录
func (op *userRoomRecordOp) Delete(user_id int64, record_id int) error {
	return op.DeleteTx(db.DB, user_id, record_id)
}

// 根据主键删除相关记录,Tx
func (op *userRoomRecordOp) DeleteTx(ext sqlx.Ext, user_id int64, record_id int) error {
	sql := `delete from user_room_record where 1=1
        and user_id=?
        and record_id=?
        `
	_, err := ext.Exec(sql,
		user_id,
		record_id,
	)
	return err
}

// 返回符合查询条件的记录数
func (op *userRoomRecordOp) CountByMap(m map[string]interface{}) (int64, error) {

	var params []interface{}
	sql := `select count(*) from user_room_record where 1=1 `
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

func (op *userRoomRecordOp) DeleteByMap(m map[string]interface{}) (int64, error) {
	return op.DeleteByMapTx(db.DB, m)
}

func (op *userRoomRecordOp) DeleteByMapTx(ext sqlx.Ext, m map[string]interface{}) (int64, error) {
	var params []interface{}
	sql := "delete from user_room_record where 1=1 "
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
