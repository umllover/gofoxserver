package model

import (
	"errors"
	"fmt"
	"mj/hallServer/db"

	"github.com/jmoiron/sqlx"
	"github.com/lovelly/leaf/log"
)

//This file is generate by scripts,don't edit it

//usertoken
//

// +gen *
type Usertoken struct {
	UserID   int64 `db:"UserID" json:"UserID"`     //
	Currency int   `db:"Currency" json:"Currency"` // 游戏豆
	RoomCard int   `db:"RoomCard" json:"RoomCard"` // 房卡数
}

type usertokenOp struct{}

var UsertokenOp = &usertokenOp{}
var DefaultUsertoken = &Usertoken{}

// 按主键查询. 注:未找到记录的话将触发sql.ErrNoRows错误，返回nil, false
func (op *usertokenOp) Get(UserID int64) (*Usertoken, bool) {
	obj := &Usertoken{}
	sql := "select * from usertoken where UserID=? "
	err := db.DB.Get(obj, sql,
		UserID,
	)

	if err != nil {
		log.Error("Get data error:%v", err.Error())
		return nil, false
	}
	return obj, true
}
func (op *usertokenOp) SelectAll() ([]*Usertoken, error) {
	objList := []*Usertoken{}
	sql := "select * from usertoken "
	err := db.DB.Select(&objList, sql)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return objList, nil
}

func (op *usertokenOp) QueryByMap(m map[string]interface{}) ([]*Usertoken, error) {
	result := []*Usertoken{}
	var params []interface{}

	sql := "select * from usertoken where 1=1 "
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

func (op *usertokenOp) GetByMap(m map[string]interface{}) (*Usertoken, error) {
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
func (i *Usertoken) Insert() error {
    err := db.DBMap.Insert(i)
    if err != nil{
		log.Error("Insert sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 插入数据，自增长字段将被忽略
func (op *usertokenOp) Insert(m *Usertoken) (int64, error) {
	return op.InsertTx(db.DB, m)
}

// 插入数据，自增长字段将被忽略
func (op *usertokenOp) InsertTx(ext sqlx.Ext, m *Usertoken) (int64, error) {
	sql := "insert into usertoken(UserID,Currency,RoomCard) values(?,?,?)"
	result, err := ext.Exec(sql,
		m.UserID,
		m.Currency,
		m.RoomCard,
	)
	if err != nil {
		log.Error("InsertTx sql error:%v, data:%v", err.Error(), m)
		return -1, err
	}
	affected, _ := result.LastInsertId()
	return affected, nil
}

//存在就更新， 不存在就插入
func (op *usertokenOp) InsertUpdate(obj *Usertoken, m map[string]interface{}) error {
	sql := "insert into usertoken(UserID,Currency,RoomCard) values(?,?,?) ON DUPLICATE KEY UPDATE "
	var params = []interface{}{obj.UserID,
		obj.Currency,
		obj.RoomCard,
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
func (i *Usertoken) Update()  error {
    _,err := db.DBMap.Update(i)
    if err != nil{
		log.Error("update sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *usertokenOp) Update(m *Usertoken) error {
	return op.UpdateTx(db.DB, m)
}

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *usertokenOp) UpdateTx(ext sqlx.Ext, m *Usertoken) error {
	sql := `update usertoken set Currency=?,RoomCard=? where UserID=?`
	_, err := ext.Exec(sql,
		m.Currency,
		m.RoomCard,
		m.UserID,
	)

	if err != nil {
		log.Error("update sql error:%v, data:%v", err.Error(), m)
		return err
	}

	return nil
}

// 用主键做条件，更新map里包含的字段名
func (op *usertokenOp) UpdateWithMap(UserID int64, m map[string]interface{}) error {
	return op.UpdateWithMapTx(db.DB, UserID, m)
}

// 用主键做条件，更新map里包含的字段名
func (op *usertokenOp) UpdateWithMapTx(ext sqlx.Ext, UserID int64, m map[string]interface{}) error {

	sql := `update usertoken set %s where 1=1 and UserID=? ;`

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
func (i *Usertoken) Delete() error{
    _,err := db.DBMap.Delete(i)
	log.Error("Delete sql error:%v", err.Error())
    return err
}
*/
// 根据主键删除相关记录
func (op *usertokenOp) Delete(UserID int64) error {
	return op.DeleteTx(db.DB, UserID)
}

// 根据主键删除相关记录,Tx
func (op *usertokenOp) DeleteTx(ext sqlx.Ext, UserID int64) error {
	sql := `delete from usertoken where 1=1
        and UserID=?
        `
	_, err := ext.Exec(sql,
		UserID,
	)
	return err
}

// 返回符合查询条件的记录数
func (op *usertokenOp) CountByMap(m map[string]interface{}) (int64, error) {

	var params []interface{}
	sql := `select count(*) from usertoken where 1=1 `
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

func (op *usertokenOp) DeleteByMap(m map[string]interface{}) (int64, error) {
	return op.DeleteByMapTx(db.DB, m)
}

func (op *usertokenOp) DeleteByMapTx(ext sqlx.Ext, m map[string]interface{}) (int64, error) {
	var params []interface{}
	sql := "delete from usertoken where 1=1 "
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
