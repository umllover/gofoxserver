package model

import (
	"errors"
	"fmt"
	"mj/gameServer/db"

	"github.com/jmoiron/sqlx"
	"github.com/lovelly/leaf/log"
)

//This file is generate by scripts,don't edit it

//room_id
//

// +gen *
type RoomId struct {
	Id     int `db:"id" json:"id"`           //
	NodeId int `db:"node_id" json:"node_id"` //
}

type roomIdOp struct{}

var RoomIdOp = &roomIdOp{}
var DefaultRoomId = &RoomId{}

// 按主键查询. 注:未找到记录的话将触发sql.ErrNoRows错误，返回nil, false
func (op *roomIdOp) Get(id int) (*RoomId, bool) {
	obj := &RoomId{}
	sql := "select * from room_id where id=? "
	err := db.DB.Get(obj, sql,
		id,
	)

	if err != nil {
		log.Error("Get data error:%v", err.Error())
		return nil, false
	}
	return obj, true
}
func (op *roomIdOp) SelectAll() ([]*RoomId, error) {
	objList := []*RoomId{}
	sql := "select * from room_id "
	err := db.DB.Select(&objList, sql)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return objList, nil
}

func (op *roomIdOp) QueryByMap(m map[string]interface{}) ([]*RoomId, error) {
	result := []*RoomId{}
	var params []interface{}

	sql := "select * from room_id where 1=1 "
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

func (op *roomIdOp) QueryByMapQueryByMapComparison(m map[string]interface{}) ([]*RoomId, error) {
	result := []*RoomId{}
	var params []interface{}

	sql := "select * from room_id where 1=1 "
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

func (op *roomIdOp) GetByMap(m map[string]interface{}) (*RoomId, error) {
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
func (i *RoomId) Insert() error {
    err := db.DBMap.Insert(i)
    if err != nil{
		log.Error("Insert sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 插入数据，自增长字段将被忽略
func (op *roomIdOp) Insert(m *RoomId) (int64, error) {
	return op.InsertTx(db.DB, m)
}

// 插入数据，自增长字段将被忽略
func (op *roomIdOp) InsertTx(ext sqlx.Ext, m *RoomId) (int64, error) {
	sql := "insert into room_id(id,node_id) values(?,?)"
	result, err := ext.Exec(sql,
		m.Id,
		m.NodeId,
	)
	if err != nil {
		log.Error("InsertTx sql error:%v, data:%v", err.Error(), m)
		return -1, err
	}
	affected, _ := result.LastInsertId()
	return affected, nil
}

/*
func (i *RoomId) Update()  error {
    _,err := db.DBMap.Update(i)
    if err != nil{
		log.Error("update sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *roomIdOp) Update(m *RoomId) error {
	return op.UpdateTx(db.DB, m)
}

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *roomIdOp) UpdateTx(ext sqlx.Ext, m *RoomId) error {
	sql := `update room_id set node_id=? where id=?`
	_, err := ext.Exec(sql,
		m.NodeId,
		m.Id,
	)

	if err != nil {
		log.Error("update sql error:%v, data:%v", err.Error(), m)
		return err
	}

	return nil
}

// 用主键做条件，更新map里包含的字段名
func (op *roomIdOp) UpdateWithMap(id int, m map[string]interface{}) error {
	return op.UpdateWithMapTx(db.DB, id, m)
}

// 用主键做条件，更新map里包含的字段名
func (op *roomIdOp) UpdateWithMapTx(ext sqlx.Ext, id int, m map[string]interface{}) error {

	sql := `update room_id set %s where 1=1 and id=? ;`

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
func (i *RoomId) Delete() error{
    _,err := db.DBMap.Delete(i)
	log.Error("Delete sql error:%v", err.Error())
    return err
}
*/
// 根据主键删除相关记录
func (op *roomIdOp) Delete(id int) error {
	return op.DeleteTx(db.DB, id)
}

// 根据主键删除相关记录,Tx
func (op *roomIdOp) DeleteTx(ext sqlx.Ext, id int) error {
	sql := `delete from room_id where 1=1
        and id=?
        `
	_, err := ext.Exec(sql,
		id,
	)
	return err
}

// 返回符合查询条件的记录数
func (op *roomIdOp) CountByMap(m map[string]interface{}) (int64, error) {

	var params []interface{}
	sql := `select count(*) from room_id where 1=1 `
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

func (op *roomIdOp) DeleteByMap(m map[string]interface{}) (int64, error) {
	return op.DeleteByMapTx(db.DB, m)
}

func (op *roomIdOp) DeleteByMapTx(ext sqlx.Ext, m map[string]interface{}) (int64, error) {
	var params []interface{}
	sql := "delete from room_id where 1=1 "
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
