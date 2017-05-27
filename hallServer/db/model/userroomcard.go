package model

import (
	"fmt"
	"gate/game_error"
	"mj/hallServer/db"

	"github.com/jmoiron/sqlx"
	"github.com/lovelly/leaf/log"
)

//This file is generate by scripts,don't edit it

//userroomcard
//

// +gen *
type Userroomcard struct {
	UserID   int `db:"UserID" json:"UserID"`     // 用户标识
	RoomCard int `db:"RoomCard" json:"RoomCard"` // 房卡数
}

type userroomcardOp struct{}

var UserroomcardOp = &userroomcardOp{}
var DefaultUserroomcard = &Userroomcard{}

// 按主键查询. 注:未找到记录的话将触发sql.ErrNoRows错误，返回nil, false
func (op *userroomcardOp) Get(UserID int) (*Userroomcard, bool) {
	obj := &Userroomcard{}
	sql := "select * from userroomcard where UserID=? "
	err := db.DB.Get(obj, sql,
		UserID,
	)

	if err != nil {
		log.Error(err.Error())
		return nil, false
	}
	return obj, true
}
func (op *userroomcardOp) SelectAll() ([]*Userroomcard, error) {
	objList := []*Userroomcard{}
	sql := "select * from userroomcard "
	err := db.DB.Select(&objList, sql)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return objList, nil
}

func (op *userroomcardOp) QueryByMap(m map[string]interface{}) ([]*Userroomcard, error) {
	result := []*Userroomcard{}
	var params []interface{}

	sql := "select * from userroomcard where 1=1 "
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

func (op *userroomcardOp) QueryByMapQueryByMapComparison(m map[string]interface{}) ([]*Userroomcard, error) {
	result := []*Userroomcard{}
	var params []interface{}

	sql := "select * from userroomcard where 1=1 "
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

func (op *userroomcardOp) GetByMap(m map[string]interface{}) (*Userroomcard, error) {
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
func (i *Userroomcard) Insert() {
    err := db.DBMap.Insert(i)
    if err != nil{
        game_error.RaiseError(err)
    }
}
*/

// 插入数据，自增长字段将被忽略
func (op *userroomcardOp) Insert(m *Userroomcard) (int64, error) {
	return op.InsertTx(db.DB, m)
}

// 插入数据，自增长字段将被忽略
func (op *userroomcardOp) InsertTx(ext sqlx.Ext, m *Userroomcard) (int64, error) {
	sql := "insert into userroomcard(UserID,RoomCard) values(?,?)"
	result, err := ext.Exec(sql,
		m.UserID,
		m.RoomCard,
	)
	if err != nil {
		game_error.RaiseError(err)
		return -1, err
	}
	affected, _ := result.LastInsertId()
	return affected, nil
}

/*
func (i *Userroomcard) Update() {
    _,err := db.DBMap.Update(i)
    if err != nil{
        game_error.RaiseError(err)
    }
}
*/

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *userroomcardOp) Update(m *Userroomcard) error {
	return op.UpdateTx(db.DB, m)
}

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *userroomcardOp) UpdateTx(ext sqlx.Ext, m *Userroomcard) error {
	sql := `update userroomcard set RoomCard=? where UserID=?`
	_, err := ext.Exec(sql,
		m.RoomCard,
		m.UserID,
	)

	if err != nil {
		game_error.RaiseError(err)
		return err
	}

	return nil
}

// 用主键做条件，更新map里包含的字段名
func (op *userroomcardOp) UpdateWithMap(UserID int, m map[string]interface{}) error {
	return op.UpdateWithMapTx(db.DB, UserID, m)
}

// 用主键做条件，更新map里包含的字段名
func (op *userroomcardOp) UpdateWithMapTx(ext sqlx.Ext, UserID int, m map[string]interface{}) error {

	sql := `update userroomcard set %s where 1=1 and UserID=? ;`

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
func (i *Userroomcard) Delete(){
    _,err := db.DBMap.Delete(i)
    if err != nil{
        game_error.RaiseError(err)
    }
}
*/
// 根据主键删除相关记录
func (op *userroomcardOp) Delete(UserID int) error {
	return op.DeleteTx(db.DB, UserID)
}

// 根据主键删除相关记录,Tx
func (op *userroomcardOp) DeleteTx(ext sqlx.Ext, UserID int) error {
	sql := `delete from userroomcard where 1=1
        and UserID=?
        `
	_, err := ext.Exec(sql,
		UserID,
	)
	return err
}

// 返回符合查询条件的记录数
func (op *userroomcardOp) CountByMap(m map[string]interface{}) int64 {

	var params []interface{}
	sql := `select count(*) from userroomcard where 1=1 `
	for k, v := range m {
		sql += fmt.Sprintf(" and  %s=? ", k)
		params = append(params, v)
	}
	count := int64(-1)
	err := db.DB.Get(&count, sql, params...)
	if err != nil {
		game_error.RaiseError(err)
	}
	return count
}

func (op *userroomcardOp) DeleteByMap(m map[string]interface{}) (int64, error) {
	return op.DeleteByMapTx(db.DB, m)
}

func (op *userroomcardOp) DeleteByMapTx(ext sqlx.Ext, m map[string]interface{}) (int64, error) {
	var params []interface{}
	sql := "delete from userroomcard where 1=1 "
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
