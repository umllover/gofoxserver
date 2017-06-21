package model

import (
	"fmt"
	"mj/gameServer/db"

	"github.com/jmoiron/sqlx"
	"github.com/lovelly/leaf/log"
)

//This file is generate by scripts,don't edit it

//race_msg_info
//

// +gen *
type RaceMsgInfo struct {
	MsgID        int    `db:"MsgID" json:"MsgID"`               //
	StartTime    int    `db:"StartTime" json:"StartTime"`       //
	EndTime      int    `db:"EndTime" json:"EndTime"`           //
	IntervalTime int    `db:"IntervalTime" json:"IntervalTime"` //
	Context      string `db:"Context" json:"Context"`           //
	MsgType      int    `db:"MsgType" json:"MsgType"`           //
}

type raceMsgInfoOp struct{}

var RaceMsgInfoOp = &raceMsgInfoOp{}
var DefaultRaceMsgInfo = &RaceMsgInfo{}

// 按主键查询. 注:未找到记录的话将触发sql.ErrNoRows错误，返回nil, false
func (op *raceMsgInfoOp) Get(MsgID int) (*RaceMsgInfo, bool) {
	obj := &RaceMsgInfo{}
	sql := "select * from race_msg_info where MsgID=? "
	err := db.DB.Get(obj, sql,
		MsgID,
	)

	if err != nil {
		log.Error("Get data error:%v", err.Error())
		return nil, false
	}
	return obj, true
}
func (op *raceMsgInfoOp) SelectAll() ([]*RaceMsgInfo, error) {
	objList := []*RaceMsgInfo{}
	sql := "select * from race_msg_info "
	err := db.DB.Select(&objList, sql)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return objList, nil
}

func (op *raceMsgInfoOp) QueryByMap(m map[string]interface{}) ([]*RaceMsgInfo, error) {
	result := []*RaceMsgInfo{}
	var params []interface{}

	sql := "select * from race_msg_info where 1=1 "
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

func (op *raceMsgInfoOp) QueryByMapQueryByMapComparison(m map[string]interface{}) ([]*RaceMsgInfo, error) {
	result := []*RaceMsgInfo{}
	var params []interface{}

	sql := "select * from race_msg_info where 1=1 "
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

func (op *raceMsgInfoOp) GetByMap(m map[string]interface{}) (*RaceMsgInfo, error) {
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
func (i *RaceMsgInfo) Insert() error {
    err := db.DBMap.Insert(i)
    if err != nil{
		log.Error("Insert sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 插入数据，自增长字段将被忽略
func (op *raceMsgInfoOp) Insert(m *RaceMsgInfo) (int64, error) {
	return op.InsertTx(db.DB, m)
}

// 插入数据，自增长字段将被忽略
func (op *raceMsgInfoOp) InsertTx(ext sqlx.Ext, m *RaceMsgInfo) (int64, error) {
	sql := "insert into race_msg_info(StartTime,EndTime,IntervalTime,Context,MsgType) values(?,?,?,?,?)"
	result, err := ext.Exec(sql,
		m.StartTime,
		m.EndTime,
		m.IntervalTime,
		m.Context,
		m.MsgType,
	)
	if err != nil {
		log.Error("InsertTx sql error:%v, data:%v", err.Error(), m)
		return -1, err
	}
	affected, _ := result.LastInsertId()
	return affected, nil
}

/*
func (i *RaceMsgInfo) Update()  error {
    _,err := db.DBMap.Update(i)
    if err != nil{
		log.Error("update sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *raceMsgInfoOp) Update(m *RaceMsgInfo) error {
	return op.UpdateTx(db.DB, m)
}

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *raceMsgInfoOp) UpdateTx(ext sqlx.Ext, m *RaceMsgInfo) error {
	sql := `update race_msg_info set StartTime=?,EndTime=?,IntervalTime=?,Context=?,MsgType=? where MsgID=?`
	_, err := ext.Exec(sql,
		m.StartTime,
		m.EndTime,
		m.IntervalTime,
		m.Context,
		m.MsgType,
		m.MsgID,
	)

	if err != nil {
		log.Error("update sql error:%v, data:%v", err.Error(), m)
		return err
	}

	return nil
}

// 用主键做条件，更新map里包含的字段名
func (op *raceMsgInfoOp) UpdateWithMap(MsgID int, m map[string]interface{}) error {
	return op.UpdateWithMapTx(db.DB, MsgID, m)
}

// 用主键做条件，更新map里包含的字段名
func (op *raceMsgInfoOp) UpdateWithMapTx(ext sqlx.Ext, MsgID int, m map[string]interface{}) error {

	sql := `update race_msg_info set %s where 1=1 and MsgID=? ;`

	var params []interface{}
	var set_sql string
	for k, v := range m {
		if set_sql != "" {
			set_sql += ","
		}
		set_sql += fmt.Sprintf(" %s=? ", k)
		params = append(params, v)
	}
	params = append(params, MsgID)
	_, err := ext.Exec(fmt.Sprintf(sql, set_sql), params...)
	return err
}

/*
func (i *RaceMsgInfo) Delete() error{
    _,err := db.DBMap.Delete(i)
	log.Error("Delete sql error:%v", err.Error())
    return err
}
*/
// 根据主键删除相关记录
func (op *raceMsgInfoOp) Delete(MsgID int) error {
	return op.DeleteTx(db.DB, MsgID)
}

// 根据主键删除相关记录,Tx
func (op *raceMsgInfoOp) DeleteTx(ext sqlx.Ext, MsgID int) error {
	sql := `delete from race_msg_info where 1=1
        and MsgID=?
        `
	_, err := ext.Exec(sql,
		MsgID,
	)
	return err
}

// 返回符合查询条件的记录数
func (op *raceMsgInfoOp) CountByMap(m map[string]interface{}) (int64, error) {

	var params []interface{}
	sql := `select count(*) from race_msg_info where 1=1 `
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

func (op *raceMsgInfoOp) DeleteByMap(m map[string]interface{}) (int64, error) {
	return op.DeleteByMapTx(db.DB, m)
}

func (op *raceMsgInfoOp) DeleteByMapTx(ext sqlx.Ext, m map[string]interface{}) (int64, error) {
	var params []interface{}
	sql := "delete from race_msg_info where 1=1 "
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
