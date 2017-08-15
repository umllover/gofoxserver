package model

import (
	"fmt"
	"mj/gameServer/db"

	"github.com/jmoiron/sqlx"
	"github.com/lovelly/leaf/log"
)

//This file is generate by scripts,don't edit it

//room_record
//

// +gen *
type RoomRecord struct {
	RecordId  int64  `db:"record_id" json:"record_id"`   // 视频id
	StartInfo string `db:"start_info" json:"start_info"` // 开始信息
	PlayInfo  string `db:"play_info" json:"play_info"`   // 玩的数据
	EndInfo   string `db:"end_info" json:"end_info"`     // 结束数据
}

type roomRecordOp struct{}

var RoomRecordOp = &roomRecordOp{}
var DefaultRoomRecord = &RoomRecord{}

// 按主键查询. 注:未找到记录的话将触发sql.ErrNoRows错误，返回nil, false
func (op *roomRecordOp) Get(record_id int64) (*RoomRecord, bool) {
	obj := &RoomRecord{}
	sql := "select * from room_record where record_id=? "
	err := db.DB.Get(obj, sql,
		record_id,
	)

	if err != nil {
		log.Error("Get data error:%v", err.Error())
		return nil, false
	}
	return obj, true
}
func (op *roomRecordOp) SelectAll() ([]*RoomRecord, error) {
	objList := []*RoomRecord{}
	sql := "select * from room_record "
	err := db.DB.Select(&objList, sql)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return objList, nil
}

func (op *roomRecordOp) QueryByMap(m map[string]interface{}) ([]*RoomRecord, error) {
	result := []*RoomRecord{}
	var params []interface{}

	sql := "select * from room_record where 1=1 "
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

func (op *roomRecordOp) GetByMap(m map[string]interface{}) (*RoomRecord, error) {
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
func (i *RoomRecord) Insert() error {
    err := db.DBMap.Insert(i)
    if err != nil{
		log.Error("Insert sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 插入数据，自增长字段将被忽略
func (op *roomRecordOp) Insert(m *RoomRecord) (int64, error) {
	return op.InsertTx(db.DB, m)
}

// 插入数据，自增长字段将被忽略
func (op *roomRecordOp) InsertTx(ext sqlx.Ext, m *RoomRecord) (int64, error) {
	sql := "insert into room_record(start_info,play_info,end_info) values(?,?,?)"
	result, err := ext.Exec(sql,
		m.StartInfo,
		m.PlayInfo,
		m.EndInfo,
	)
	if err != nil {
		log.Error("InsertTx sql error:%v, data:%v", err.Error(), m)
		return -1, err
	}
	affected, _ := result.LastInsertId()
	return affected, nil
}

//存在就更新， 不存在就插入
func (op *roomRecordOp) InsertUpdate(obj *RoomRecord, m map[string]interface{}) error {
	sql := "insert into room_record(start_info,play_info,end_info) values(?,?,?) ON DUPLICATE KEY UPDATE "
	var params = []interface{}{obj.StartInfo,
		obj.PlayInfo,
		obj.EndInfo,
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
func (i *RoomRecord) Update()  error {
    _,err := db.DBMap.Update(i)
    if err != nil{
		log.Error("update sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *roomRecordOp) Update(m *RoomRecord) error {
	return op.UpdateTx(db.DB, m)
}

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *roomRecordOp) UpdateTx(ext sqlx.Ext, m *RoomRecord) error {
	sql := `update room_record set start_info=?,play_info=?,end_info=? where record_id=?`
	_, err := ext.Exec(sql,
		m.StartInfo,
		m.PlayInfo,
		m.EndInfo,
		m.RecordId,
	)

	if err != nil {
		log.Error("update sql error:%v, data:%v", err.Error(), m)
		return err
	}

	return nil
}

// 用主键做条件，更新map里包含的字段名
func (op *roomRecordOp) UpdateWithMap(record_id int64, m map[string]interface{}) error {
	return op.UpdateWithMapTx(db.DB, record_id, m)
}

// 用主键做条件，更新map里包含的字段名
func (op *roomRecordOp) UpdateWithMapTx(ext sqlx.Ext, record_id int64, m map[string]interface{}) error {

	sql := `update room_record set %s where 1=1 and record_id=? ;`

	var params []interface{}
	var set_sql string
	for k, v := range m {
		if set_sql != "" {
			set_sql += ","
		}
		set_sql += fmt.Sprintf(" %s=? ", k)
		params = append(params, v)
	}
	params = append(params, record_id)
	_, err := ext.Exec(fmt.Sprintf(sql, set_sql), params...)
	return err
}

/*
func (i *RoomRecord) Delete() error{
    _,err := db.DBMap.Delete(i)
	log.Error("Delete sql error:%v", err.Error())
    return err
}
*/
// 根据主键删除相关记录
func (op *roomRecordOp) Delete(record_id int64) error {
	return op.DeleteTx(db.DB, record_id)
}

// 根据主键删除相关记录,Tx
func (op *roomRecordOp) DeleteTx(ext sqlx.Ext, record_id int64) error {
	sql := `delete from room_record where 1=1
        and record_id=?
        `
	_, err := ext.Exec(sql,
		record_id,
	)
	return err
}

// 返回符合查询条件的记录数
func (op *roomRecordOp) CountByMap(m map[string]interface{}) (int64, error) {

	var params []interface{}
	sql := `select count(*) from room_record where 1=1 `
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

func (op *roomRecordOp) DeleteByMap(m map[string]interface{}) (int64, error) {
	return op.DeleteByMapTx(db.DB, m)
}

func (op *roomRecordOp) DeleteByMapTx(ext sqlx.Ext, m map[string]interface{}) (int64, error) {
	var params []interface{}
	sql := "delete from room_record where 1=1 "
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
