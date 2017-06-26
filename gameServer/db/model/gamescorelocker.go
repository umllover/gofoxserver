package model

import (
	"fmt"
	"mj/gameServer/db"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lovelly/leaf/log"
)

//This file is generate by scripts,don't edit it

//gamescorelocker
//

// +gen *
type Gamescorelocker struct {
	UserID       int        `db:"UserID" json:"UserID"`             // 用户索引
	KindID       int        `db:"KindID" json:"KindID"`             // 房间索引
	NodeID       int        `db:"NodeID" json:"NodeID"`             // 在哪个服务器上
	ServerID     int        `db:"ServerID" json:"ServerID"`         // 游戏标识
	Roomid       int        `db:"roomid" json:"roomid"`             // 进出索引
	EnterIP      string     `db:"EnterIP" json:"EnterIP"`           // 进入地址
	EnterMachine string     `db:"EnterMachine" json:"EnterMachine"` // 进入机器
	CollectDate  *time.Time `db:"CollectDate" json:"CollectDate"`   // 录入日期
}

type gamescorelockerOp struct{}

var GamescorelockerOp = &gamescorelockerOp{}
var DefaultGamescorelocker = &Gamescorelocker{}

// 按主键查询. 注:未找到记录的话将触发sql.ErrNoRows错误，返回nil, false
func (op *gamescorelockerOp) Get(UserID int) (*Gamescorelocker, bool) {
	obj := &Gamescorelocker{}
	sql := "select * from gamescorelocker where UserID=? "
	err := db.DB.Get(obj, sql,
		UserID,
	)

	if err != nil {
		log.Error("Get data error:%v", err.Error())
		return nil, false
	}
	return obj, true
}
func (op *gamescorelockerOp) SelectAll() ([]*Gamescorelocker, error) {
	objList := []*Gamescorelocker{}
	sql := "select * from gamescorelocker "
	err := db.DB.Select(&objList, sql)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return objList, nil
}

func (op *gamescorelockerOp) QueryByMap(m map[string]interface{}) ([]*Gamescorelocker, error) {
	result := []*Gamescorelocker{}
	var params []interface{}

	sql := "select * from gamescorelocker where 1=1 "
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

func (op *gamescorelockerOp) QueryByMapQueryByMapComparison(m map[string]interface{}) ([]*Gamescorelocker, error) {
	result := []*Gamescorelocker{}
	var params []interface{}

	sql := "select * from gamescorelocker where 1=1 "
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

func (op *gamescorelockerOp) GetByMap(m map[string]interface{}) (*Gamescorelocker, error) {
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
func (i *Gamescorelocker) Insert() error {
    err := db.DBMap.Insert(i)
    if err != nil{
		log.Error("Insert sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 插入数据，自增长字段将被忽略
func (op *gamescorelockerOp) Insert(m *Gamescorelocker) (int64, error) {
	return op.InsertTx(db.DB, m)
}

// 插入数据，自增长字段将被忽略
func (op *gamescorelockerOp) InsertTx(ext sqlx.Ext, m *Gamescorelocker) (int64, error) {
	sql := "insert into gamescorelocker(UserID,KindID,NodeID,ServerID,roomid,EnterIP,EnterMachine,CollectDate) values(?,?,?,?,?,?,?,?)"
	result, err := ext.Exec(sql,
		m.UserID,
		m.KindID,
		m.NodeID,
		m.ServerID,
		m.Roomid,
		m.EnterIP,
		m.EnterMachine,
		m.CollectDate,
	)
	if err != nil {
		log.Error("InsertTx sql error:%v, data:%v", err.Error(), m)
		return -1, err
	}
	affected, _ := result.LastInsertId()
	return affected, nil
}

/*
func (i *Gamescorelocker) Update()  error {
    _,err := db.DBMap.Update(i)
    if err != nil{
		log.Error("update sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *gamescorelockerOp) Update(m *Gamescorelocker) error {
	return op.UpdateTx(db.DB, m)
}

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *gamescorelockerOp) UpdateTx(ext sqlx.Ext, m *Gamescorelocker) error {
	sql := `update gamescorelocker set KindID=?,NodeID=?,ServerID=?,roomid=?,EnterIP=?,EnterMachine=?,CollectDate=? where UserID=?`
	_, err := ext.Exec(sql,
		m.KindID,
		m.NodeID,
		m.ServerID,
		m.Roomid,
		m.EnterIP,
		m.EnterMachine,
		m.CollectDate,
		m.UserID,
	)

	if err != nil {
		log.Error("update sql error:%v, data:%v", err.Error(), m)
		return err
	}

	return nil
}

// 用主键做条件，更新map里包含的字段名
func (op *gamescorelockerOp) UpdateWithMap(UserID int, m map[string]interface{}) error {
	return op.UpdateWithMapTx(db.DB, UserID, m)
}

// 用主键做条件，更新map里包含的字段名
func (op *gamescorelockerOp) UpdateWithMapTx(ext sqlx.Ext, UserID int, m map[string]interface{}) error {

	sql := `update gamescorelocker set %s where 1=1 and UserID=? ;`

	var params []interface{}
	var set_sql string
	for k, v := range m {
		if set_sql != "" {
			set_sql += ","
		}
		set_sql += fmt.Sprintf(" %s=? ", k)
		params = append(params, v)
	}
	params = append(params, UserID)
	_, err := ext.Exec(fmt.Sprintf(sql, set_sql), params...)
	return err
}

/*
func (i *Gamescorelocker) Delete() error{
    _,err := db.DBMap.Delete(i)
	log.Error("Delete sql error:%v", err.Error())
    return err
}
*/
// 根据主键删除相关记录
func (op *gamescorelockerOp) Delete(UserID int) error {
	return op.DeleteTx(db.DB, UserID)
}

// 根据主键删除相关记录,Tx
func (op *gamescorelockerOp) DeleteTx(ext sqlx.Ext, UserID int) error {
	sql := `delete from gamescorelocker where 1=1
        and UserID=?
        `
	_, err := ext.Exec(sql,
		UserID,
	)
	return err
}

// 返回符合查询条件的记录数
func (op *gamescorelockerOp) CountByMap(m map[string]interface{}) (int64, error) {

	var params []interface{}
	sql := `select count(*) from gamescorelocker where 1=1 `
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

func (op *gamescorelockerOp) DeleteByMap(m map[string]interface{}) (int64, error) {
	return op.DeleteByMapTx(db.DB, m)
}

func (op *gamescorelockerOp) DeleteByMapTx(ext sqlx.Ext, m map[string]interface{}) (int64, error) {
	var params []interface{}
	sql := "delete from gamescorelocker where 1=1 "
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
