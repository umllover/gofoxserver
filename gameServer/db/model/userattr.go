package model

import (
	"fmt"
	"mj/gameServer/db"

	"github.com/jmoiron/sqlx"
	"github.com/lovelly/leaf/log"
)

//This file is generate by scripts,don't edit it

//userattr
//

// +gen *
type Userattr struct {
	UserID   int `db:"UserID" json:"UserID"`     //
	Currency int `db:"Currency" json:"Currency"` // 游戏豆
	RoomCard int `db:"RoomCard" json:"RoomCard"` // 房卡数
}

type userattrOp struct{}

var UserattrOp = &userattrOp{}
var DefaultUserattr = &Userattr{}

// 按主键查询. 注:未找到记录的话将触发sql.ErrNoRows错误，返回nil, false
func (op *userattrOp) Get(UserID int) (*Userattr, bool) {
	obj := &Userattr{}
	sql := "select * from userattr where UserID=? "
	err := db.DB.Get(obj, sql,
		UserID,
	)

	if err != nil {
		log.Error("Get data error:%v", err.Error())
		return nil, false
	}
	return obj, true
}
func (op *userattrOp) SelectAll() ([]*Userattr, error) {
	objList := []*Userattr{}
	sql := "select * from userattr "
	err := db.DB.Select(&objList, sql)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return objList, nil
}

func (op *userattrOp) QueryByMap(m map[string]interface{}) ([]*Userattr, error) {
	result := []*Userattr{}
	var params []interface{}

	sql := "select * from userattr where 1=1 "
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

func (op *userattrOp) QueryByMapQueryByMapComparison(m map[string]interface{}) ([]*Userattr, error) {
	result := []*Userattr{}
	var params []interface{}

	sql := "select * from userattr where 1=1 "
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

func (op *userattrOp) GetByMap(m map[string]interface{}) (*Userattr, error) {
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
func (i *Userattr) Insert() error {
    err := db.DBMap.Insert(i)
    if err != nil{
		log.Error("Insert sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 插入数据，自增长字段将被忽略
func (op *userattrOp) Insert(m *Userattr) (int64, error) {
	return op.InsertTx(db.DB, m)
}

// 插入数据，自增长字段将被忽略
func (op *userattrOp) InsertTx(ext sqlx.Ext, m *Userattr) (int64, error) {
	sql := "insert into userattr(UserID,Currency,RoomCard) values(?,?,?)"
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

/*
func (i *Userattr) Update()  error {
    _,err := db.DBMap.Update(i)
    if err != nil{
		log.Error("update sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *userattrOp) Update(m *Userattr) error {
	return op.UpdateTx(db.DB, m)
}

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *userattrOp) UpdateTx(ext sqlx.Ext, m *Userattr) error {
	sql := `update userattr set Currency=?,RoomCard=? where UserID=?`
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
func (op *userattrOp) UpdateWithMap(UserID int, m map[string]interface{}) error {
	return op.UpdateWithMapTx(db.DB, UserID, m)
}

// 用主键做条件，更新map里包含的字段名
func (op *userattrOp) UpdateWithMapTx(ext sqlx.Ext, UserID int, m map[string]interface{}) error {

	sql := `update userattr set %s where 1=1 and UserID=? ;`

	var params []interface{}
	var set_sql string
	for k, v := range m {
		set_sql += fmt.Sprintf(" %s=? ", k)
		params = append(params, v)
	}
	params = append(params, UserID)
	_, err := ext.Exec(fmt.Sprintf(sql, set_sql), params...)
	return err
}

/*
func (i *Userattr) Delete() error{
    _,err := db.DBMap.Delete(i)
	log.Error("Delete sql error:%v", err.Error())
    return err
}
*/
// 根据主键删除相关记录
func (op *userattrOp) Delete(UserID int) error {
	return op.DeleteTx(db.DB, UserID)
}

// 根据主键删除相关记录,Tx
func (op *userattrOp) DeleteTx(ext sqlx.Ext, UserID int) error {
	sql := `delete from userattr where 1=1
        and UserID=?
        `
	_, err := ext.Exec(sql,
		UserID,
	)
	return err
}

// 返回符合查询条件的记录数
func (op *userattrOp) CountByMap(m map[string]interface{}) (int64, error) {

	var params []interface{}
	sql := `select count(*) from userattr where 1=1 `
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

func (op *userattrOp) DeleteByMap(m map[string]interface{}) (int64, error) {
	return op.DeleteByMapTx(db.DB, m)
}

func (op *userattrOp) DeleteByMapTx(ext sqlx.Ext, m map[string]interface{}) (int64, error) {
	var params []interface{}
	sql := "delete from userattr where 1=1 "
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
