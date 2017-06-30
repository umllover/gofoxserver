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

//token_record
//

// +gen *
type TokenRecord struct {
	RoomId      int        `db:"room_id" json:"room_id"`           //
	UserId      int        `db:"user_id" json:"user_id"`           //
	TokenType   int        `db:"tokenType" json:"tokenType"`       //
	Amount      int        `db:"amount" json:"amount"`             //
	Status      int        `db:"status" json:"status"`             //
	CreatorTime *time.Time `db:"creator_time" json:"creator_time"` //
	KindID      int        `db:"KindID" json:"KindID"`             //
	ServerId    int        `db:"ServerId" json:"ServerId"`         //
}

type tokenRecordOp struct{}

var TokenRecordOp = &tokenRecordOp{}
var DefaultTokenRecord = &TokenRecord{}

// 按主键查询. 注:未找到记录的话将触发sql.ErrNoRows错误，返回nil, false
func (op *tokenRecordOp) Get(room_id int, user_id int) (*TokenRecord, bool) {
	obj := &TokenRecord{}
	sql := "select * from token_record where room_id=? and user_id=? "
	err := db.DB.Get(obj, sql,
		room_id,
		user_id,
	)

	if err != nil {
		log.Error("Get data error:%v", err.Error())
		return nil, false
	}
	return obj, true
}
func (op *tokenRecordOp) SelectAll() ([]*TokenRecord, error) {
	objList := []*TokenRecord{}
	sql := "select * from token_record "
	err := db.DB.Select(&objList, sql)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return objList, nil
}

func (op *tokenRecordOp) QueryByMap(m map[string]interface{}) ([]*TokenRecord, error) {
	result := []*TokenRecord{}
	var params []interface{}

	sql := "select * from token_record where 1=1 "
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

func (op *tokenRecordOp) QueryByMapQueryByMapComparison(m map[string]interface{}) ([]*TokenRecord, error) {
	result := []*TokenRecord{}
	var params []interface{}

	sql := "select * from token_record where 1=1 "
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

func (op *tokenRecordOp) GetByMap(m map[string]interface{}) (*TokenRecord, error) {
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
func (i *TokenRecord) Insert() error {
    err := db.DBMap.Insert(i)
    if err != nil{
		log.Error("Insert sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 插入数据，自增长字段将被忽略
func (op *tokenRecordOp) Insert(m *TokenRecord) (int64, error) {
	return op.InsertTx(db.DB, m)
}

// 插入数据，自增长字段将被忽略
func (op *tokenRecordOp) InsertTx(ext sqlx.Ext, m *TokenRecord) (int64, error) {
	sql := "insert into token_record(room_id,user_id,tokenType,amount,status,creator_time,KindID,ServerId) values(?,?,?,?,?,?,?,?)"
	result, err := ext.Exec(sql,
		m.RoomId,
		m.UserId,
		m.TokenType,
		m.Amount,
		m.Status,
		m.CreatorTime,
		m.KindID,
		m.ServerId,
	)
	if err != nil {
		log.Error("InsertTx sql error:%v, data:%v", err.Error(), m)
		return -1, err
	}
	affected, _ := result.LastInsertId()
	return affected, nil
}

/*
func (i *TokenRecord) Update()  error {
    _,err := db.DBMap.Update(i)
    if err != nil{
		log.Error("update sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *tokenRecordOp) Update(m *TokenRecord) error {
	return op.UpdateTx(db.DB, m)
}

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *tokenRecordOp) UpdateTx(ext sqlx.Ext, m *TokenRecord) error {
	sql := `update token_record set tokenType=?,amount=?,status=?,creator_time=?,KindID=?,ServerId=? where room_id=? and user_id=?`
	_, err := ext.Exec(sql,
		m.TokenType,
		m.Amount,
		m.Status,
		m.CreatorTime,
		m.KindID,
		m.ServerId,
		m.RoomId,
		m.UserId,
	)

	if err != nil {
		log.Error("update sql error:%v, data:%v", err.Error(), m)
		return err
	}

	return nil
}

// 用主键做条件，更新map里包含的字段名
func (op *tokenRecordOp) UpdateWithMap(room_id int, user_id int, m map[string]interface{}) error {
	return op.UpdateWithMapTx(db.DB, room_id, user_id, m)
}

// 用主键做条件，更新map里包含的字段名
func (op *tokenRecordOp) UpdateWithMapTx(ext sqlx.Ext, room_id int, user_id int, m map[string]interface{}) error {

	sql := `update token_record set %s where 1=1 and room_id=? and user_id=? ;`

	var params []interface{}
	var set_sql string
	for k, v := range m {
		if set_sql != "" {
			set_sql += ","
		}
		set_sql += fmt.Sprintf(" %s=? ", k)
		params = append(params, v)
	}
	params = append(params, room_id, user_id)
	_, err := ext.Exec(fmt.Sprintf(sql, set_sql), params...)
	return err
}

/*
func (i *TokenRecord) Delete() error{
    _,err := db.DBMap.Delete(i)
	log.Error("Delete sql error:%v", err.Error())
    return err
}
*/
// 根据主键删除相关记录
func (op *tokenRecordOp) Delete(room_id int, user_id int) error {
	return op.DeleteTx(db.DB, room_id, user_id)
}

// 根据主键删除相关记录,Tx
func (op *tokenRecordOp) DeleteTx(ext sqlx.Ext, room_id int, user_id int) error {
	sql := `delete from token_record where 1=1
        and room_id=?
        and user_id=?
        `
	_, err := ext.Exec(sql,
		room_id,
		user_id,
	)
	return err
}

// 返回符合查询条件的记录数
func (op *tokenRecordOp) CountByMap(m map[string]interface{}) (int64, error) {

	var params []interface{}
	sql := `select count(*) from token_record where 1=1 `
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

func (op *tokenRecordOp) DeleteByMap(m map[string]interface{}) (int64, error) {
	return op.DeleteByMapTx(db.DB, m)
}

func (op *tokenRecordOp) DeleteByMapTx(ext sqlx.Ext, m map[string]interface{}) (int64, error) {
	var params []interface{}
	sql := "delete from token_record where 1=1 "
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
