package model

import (
	"errors"
	"fmt"
	"mj/gameServer/db"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lovelly/leaf/log"
)

//This file is generate by scripts,don't edit it

//accountsmember
//

// +gen *
type Accountsmember struct {
	UserID         int64      `db:"UserID" json:"UserID"`                 // 用户标识
	MemberOrder    int8       `db:"MemberOrder" json:"MemberOrder"`       // 会员标识
	UserRight      int        `db:"UserRight" json:"UserRight"`           // 用户权限
	MemberOverDate *time.Time `db:"MemberOverDate" json:"MemberOverDate"` // 会员期限
}

type accountsmemberOp struct{}

var AccountsmemberOp = &accountsmemberOp{}
var DefaultAccountsmember = &Accountsmember{}

// 按主键查询. 注:未找到记录的话将触发sql.ErrNoRows错误，返回nil, false
func (op *accountsmemberOp) Get(UserID int64) (*Accountsmember, bool) {
	obj := &Accountsmember{}
	sql := "select * from accountsmember where UserID=? "
	err := db.DB.Get(obj, sql,
		UserID,
	)

	if err != nil {
		log.Error("Get data error:%v", err.Error())
		return nil, false
	}
	return obj, true
}
func (op *accountsmemberOp) SelectAll() ([]*Accountsmember, error) {
	objList := []*Accountsmember{}
	sql := "select * from accountsmember "
	err := db.DB.Select(&objList, sql)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return objList, nil
}

func (op *accountsmemberOp) QueryByMap(m map[string]interface{}) ([]*Accountsmember, error) {
	result := []*Accountsmember{}
	var params []interface{}

	sql := "select * from accountsmember where 1=1 "
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

func (op *accountsmemberOp) GetByMap(m map[string]interface{}) (*Accountsmember, error) {
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
func (i *Accountsmember) Insert() error {
    err := db.DBMap.Insert(i)
    if err != nil{
		log.Error("Insert sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 插入数据，自增长字段将被忽略
func (op *accountsmemberOp) Insert(m *Accountsmember) (int64, error) {
	return op.InsertTx(db.DB, m)
}

// 插入数据，自增长字段将被忽略
func (op *accountsmemberOp) InsertTx(ext sqlx.Ext, m *Accountsmember) (int64, error) {
	sql := "insert into accountsmember(UserID,MemberOrder,UserRight,MemberOverDate) values(?,?,?,?)"
	result, err := ext.Exec(sql,
		m.UserID,
		m.MemberOrder,
		m.UserRight,
		m.MemberOverDate,
	)
	if err != nil {
		log.Error("InsertTx sql error:%v, data:%v", err.Error(), m)
		return -1, err
	}
	affected, _ := result.LastInsertId()
	return affected, nil
}

/*
func (i *Accountsmember) Update()  error {
    _,err := db.DBMap.Update(i)
    if err != nil{
		log.Error("update sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *accountsmemberOp) Update(m *Accountsmember) error {
	return op.UpdateTx(db.DB, m)
}

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *accountsmemberOp) UpdateTx(ext sqlx.Ext, m *Accountsmember) error {
	sql := `update accountsmember set MemberOrder=?,UserRight=?,MemberOverDate=? where UserID=?`
	_, err := ext.Exec(sql,
		m.MemberOrder,
		m.UserRight,
		m.MemberOverDate,
		m.UserID,
	)

	if err != nil {
		log.Error("update sql error:%v, data:%v", err.Error(), m)
		return err
	}

	return nil
}

// 用主键做条件，更新map里包含的字段名
func (op *accountsmemberOp) UpdateWithMap(UserID int64, m map[string]interface{}) error {
	return op.UpdateWithMapTx(db.DB, UserID, m)
}

// 用主键做条件，更新map里包含的字段名
func (op *accountsmemberOp) UpdateWithMapTx(ext sqlx.Ext, UserID int64, m map[string]interface{}) error {

	sql := `update accountsmember set %s where 1=1 and UserID=? ;`

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
func (i *Accountsmember) Delete() error{
    _,err := db.DBMap.Delete(i)
	log.Error("Delete sql error:%v", err.Error())
    return err
}
*/
// 根据主键删除相关记录
func (op *accountsmemberOp) Delete(UserID int64) error {
	return op.DeleteTx(db.DB, UserID)
}

// 根据主键删除相关记录,Tx
func (op *accountsmemberOp) DeleteTx(ext sqlx.Ext, UserID int64) error {
	sql := `delete from accountsmember where 1=1
        and UserID=?
        `
	_, err := ext.Exec(sql,
		UserID,
	)
	return err
}

// 返回符合查询条件的记录数
func (op *accountsmemberOp) CountByMap(m map[string]interface{}) (int64, error) {

	var params []interface{}
	sql := `select count(*) from accountsmember where 1=1 `
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

func (op *accountsmemberOp) DeleteByMap(m map[string]interface{}) (int64, error) {
	return op.DeleteByMapTx(db.DB, m)
}

func (op *accountsmemberOp) DeleteByMapTx(ext sqlx.Ext, m map[string]interface{}) (int64, error) {
	var params []interface{}
	sql := "delete from accountsmember where 1=1 "
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
