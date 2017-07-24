package stats

import (
	"errors"
	"fmt"
	"mj/hallServer/db"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lovelly/leaf/log"
)

//This file is generate by scripts,don't edit it

//create_room_log
//

// +gen *
type CreateRoomLog struct {
	RoomId         int        `db:"room_id" json:"room_id"`                 // 房间id
	UserId         int64      `db:"user_id" json:"user_id"`                 // 用户索引
	RoomName       string     `db:"room_name" json:"room_name"`             //
	KindId         int        `db:"kind_id" json:"kind_id"`                 // 房间索引
	NodeId         int        `db:"node_id" json:"node_id"`                 // 在哪个服务器上
	CreateTime     *time.Time `db:"create_time" json:"create_time"`         // 录入日期
	CreateOthers   int        `db:"create_others" json:"create_others"`     // 是否为他人开房 0否，1是
	PayType        int        `db:"pay_type" json:"pay_type"`               // 支付方式 1是全服 2是AA
	TimeoutNostart int        `db:"timeout_nostart" json:"timeout_nostart"` // 是否超时未开始游戏 0否  1是
	StartEnderror  int        `db:"start_endError" json:"start_endError"`   // 是否开始但非正常解散房间 0 否 1是
	NomalOpen      int        `db:"nomal_open" json:"nomal_open"`           // 是否正常开房 0否 1是
}

type createRoomLogOp struct{}

var CreateRoomLogOp = &createRoomLogOp{}
var DefaultCreateRoomLog = &CreateRoomLog{}

// 按主键查询. 注:未找到记录的话将触发sql.ErrNoRows错误，返回nil, false
func (op *createRoomLogOp) Get(room_id int) (*CreateRoomLog, bool) {
	obj := &CreateRoomLog{}
	sql := "select * from create_room_log where room_id=? "
	err := db.StatsDB.Get(obj, sql,
		room_id,
	)

	if err != nil {
		log.Error("Get data error:%v", err.Error())
		return nil, false
	}
	return obj, true
}
func (op *createRoomLogOp) SelectAll() ([]*CreateRoomLog, error) {
	objList := []*CreateRoomLog{}
	sql := "select * from create_room_log "
	err := db.StatsDB.Select(&objList, sql)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return objList, nil
}

func (op *createRoomLogOp) QueryByMap(m map[string]interface{}) ([]*CreateRoomLog, error) {
	result := []*CreateRoomLog{}
	var params []interface{}

	sql := "select * from create_room_log where 1=1 "
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

func (op *createRoomLogOp) GetByMap(m map[string]interface{}) (*CreateRoomLog, error) {
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
func (i *CreateRoomLog) Insert() error {
    err := db.StatsDBMap.Insert(i)
    if err != nil{
		log.Error("Insert sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 插入数据，自增长字段将被忽略
func (op *createRoomLogOp) Insert(m *CreateRoomLog) (int64, error) {
	return op.InsertTx(db.StatsDB, m)
}

// 插入数据，自增长字段将被忽略
func (op *createRoomLogOp) InsertTx(ext sqlx.Ext, m *CreateRoomLog) (int64, error) {
	sql := "insert into create_room_log(room_id,user_id,room_name,kind_id,node_id,create_time,create_others,pay_type,timeout_nostart,start_endError,nomal_open) values(?,?,?,?,?,?,?,?,?,?,?)"
	result, err := ext.Exec(sql,
		m.RoomId,
		m.UserId,
		m.RoomName,
		m.KindId,
		m.NodeId,
		m.CreateTime,
		m.CreateOthers,
		m.PayType,
		m.TimeoutNostart,
		m.StartEnderror,
		m.NomalOpen,
	)
	if err != nil {
		log.Error("InsertTx sql error:%v, data:%v", err.Error(), m)
		return -1, err
	}
	affected, _ := result.LastInsertId()
	return affected, nil
}

/*
func (i *CreateRoomLog) Update()  error {
    _,err := db.StatsDBMap.Update(i)
    if err != nil{
		log.Error("update sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *createRoomLogOp) Update(m *CreateRoomLog) error {
	return op.UpdateTx(db.StatsDB, m)
}

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *createRoomLogOp) UpdateTx(ext sqlx.Ext, m *CreateRoomLog) error {
	sql := `update create_room_log set user_id=?,room_name=?,kind_id=?,node_id=?,create_time=?,create_others=?,pay_type=?,timeout_nostart=?,start_endError=?,nomal_open=? where room_id=?`
	_, err := ext.Exec(sql,
		m.UserId,
		m.RoomName,
		m.KindId,
		m.NodeId,
		m.CreateTime,
		m.CreateOthers,
		m.PayType,
		m.TimeoutNostart,
		m.StartEnderror,
		m.NomalOpen,
		m.RoomId,
	)

	if err != nil {
		log.Error("update sql error:%v, data:%v", err.Error(), m)
		return err
	}

	return nil
}

// 用主键做条件，更新map里包含的字段名
func (op *createRoomLogOp) UpdateWithMap(room_id int, m map[string]interface{}) error {
	return op.UpdateWithMapTx(db.StatsDB, room_id, m)
}

// 用主键做条件，更新map里包含的字段名
func (op *createRoomLogOp) UpdateWithMapTx(ext sqlx.Ext, room_id int, m map[string]interface{}) error {

	sql := `update create_room_log set %s where 1=1 and room_id=? ;`

	var params []interface{}
	var set_sql string
	for k, v := range m {
		if set_sql != "" {
			set_sql += ","
		}
		set_sql += fmt.Sprintf(" %s=? ", k)
		params = append(params, v)
	}
	params = append(params, room_id)
	_, err := ext.Exec(fmt.Sprintf(sql, set_sql), params...)
	return err
}

/*
func (i *CreateRoomLog) Delete() error{
    _,err := db.StatsDBMap.Delete(i)
	log.Error("Delete sql error:%v", err.Error())
    return err
}
*/
// 根据主键删除相关记录
func (op *createRoomLogOp) Delete(room_id int) error {
	return op.DeleteTx(db.StatsDB, room_id)
}

// 根据主键删除相关记录,Tx
func (op *createRoomLogOp) DeleteTx(ext sqlx.Ext, room_id int) error {
	sql := `delete from create_room_log where 1=1
        and room_id=?
        `
	_, err := ext.Exec(sql,
		room_id,
	)
	return err
}

// 返回符合查询条件的记录数
func (op *createRoomLogOp) CountByMap(m map[string]interface{}) (int64, error) {

	var params []interface{}
	sql := `select count(*) from create_room_log where 1=1 `
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

func (op *createRoomLogOp) DeleteByMap(m map[string]interface{}) (int64, error) {
	return op.DeleteByMapTx(db.StatsDB, m)
}

func (op *createRoomLogOp) DeleteByMapTx(ext sqlx.Ext, m map[string]interface{}) (int64, error) {
	var params []interface{}
	sql := "delete from create_room_log where 1=1 "
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
