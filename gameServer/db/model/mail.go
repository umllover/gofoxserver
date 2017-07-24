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

//mail
//

// +gen *
type Mail struct {
	UserId      int        `db:"user_id" json:"user_id"`           //
	MailType    int        `db:"mail_type" json:"mail_type"`       //
	Context     string     `db:"context" json:"context"`           //
	CreatorTime *time.Time `db:"creator_time" json:"creator_time"` //
	Sender      string     `db:"sender" json:"sender"`             //
	Title       string     `db:"title" json:"title"`               //
}

type mailOp struct{}

var MailOp = &mailOp{}
var DefaultMail = &Mail{}

// 按主键查询. 注:未找到记录的话将触发sql.ErrNoRows错误，返回nil, false
func (op *mailOp) Get(user_id int) (*Mail, bool) {
	obj := &Mail{}
	sql := "select * from mail where user_id=? "
	err := db.DB.Get(obj, sql,
		user_id,
	)

	if err != nil {
		log.Error("Get data error:%v", err.Error())
		return nil, false
	}
	return obj, true
}
func (op *mailOp) SelectAll() ([]*Mail, error) {
	objList := []*Mail{}
	sql := "select * from mail "
	err := db.DB.Select(&objList, sql)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return objList, nil
}

func (op *mailOp) QueryByMap(m map[string]interface{}) ([]*Mail, error) {
	result := []*Mail{}
	var params []interface{}

	sql := "select * from mail where 1=1 "
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

func (op *mailOp) GetByMap(m map[string]interface{}) (*Mail, error) {
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
func (i *Mail) Insert() error {
    err := db.DBMap.Insert(i)
    if err != nil{
		log.Error("Insert sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 插入数据，自增长字段将被忽略
func (op *mailOp) Insert(m *Mail) (int64, error) {
	return op.InsertTx(db.DB, m)
}

// 插入数据，自增长字段将被忽略
func (op *mailOp) InsertTx(ext sqlx.Ext, m *Mail) (int64, error) {
	sql := "insert into mail(user_id,mail_type,context,creator_time,sender,title) values(?,?,?,?,?,?)"
	result, err := ext.Exec(sql,
		m.UserId,
		m.MailType,
		m.Context,
		m.CreatorTime,
		m.Sender,
		m.Title,
	)
	if err != nil {
		log.Error("InsertTx sql error:%v, data:%v", err.Error(), m)
		return -1, err
	}
	affected, _ := result.LastInsertId()
	return affected, nil
}

/*
func (i *Mail) Update()  error {
    _,err := db.DBMap.Update(i)
    if err != nil{
		log.Error("update sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *mailOp) Update(m *Mail) error {
	return op.UpdateTx(db.DB, m)
}

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *mailOp) UpdateTx(ext sqlx.Ext, m *Mail) error {
	sql := `update mail set mail_type=?,context=?,creator_time=?,sender=?,title=? where user_id=?`
	_, err := ext.Exec(sql,
		m.MailType,
		m.Context,
		m.CreatorTime,
		m.Sender,
		m.Title,
		m.UserId,
	)

	if err != nil {
		log.Error("update sql error:%v, data:%v", err.Error(), m)
		return err
	}

	return nil
}

// 用主键做条件，更新map里包含的字段名
func (op *mailOp) UpdateWithMap(user_id int, m map[string]interface{}) error {
	return op.UpdateWithMapTx(db.DB, user_id, m)
}

// 用主键做条件，更新map里包含的字段名
func (op *mailOp) UpdateWithMapTx(ext sqlx.Ext, user_id int, m map[string]interface{}) error {

	sql := `update mail set %s where 1=1 and user_id=? ;`

	var params []interface{}
	var set_sql string
	for k, v := range m {
		if set_sql != "" {
			set_sql += ","
		}
		set_sql += fmt.Sprintf(" %s=? ", k)
		params = append(params, v)
	}
	params = append(params, user_id)
	_, err := ext.Exec(fmt.Sprintf(sql, set_sql), params...)
	return err
}

/*
func (i *Mail) Delete() error{
    _,err := db.DBMap.Delete(i)
	log.Error("Delete sql error:%v", err.Error())
    return err
}
*/
// 根据主键删除相关记录
func (op *mailOp) Delete(user_id int) error {
	return op.DeleteTx(db.DB, user_id)
}

// 根据主键删除相关记录,Tx
func (op *mailOp) DeleteTx(ext sqlx.Ext, user_id int) error {
	sql := `delete from mail where 1=1
        and user_id=?
        `
	_, err := ext.Exec(sql,
		user_id,
	)
	return err
}

// 返回符合查询条件的记录数
func (op *mailOp) CountByMap(m map[string]interface{}) (int64, error) {

	var params []interface{}
	sql := `select count(*) from mail where 1=1 `
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

func (op *mailOp) DeleteByMap(m map[string]interface{}) (int64, error) {
	return op.DeleteByMapTx(db.DB, m)
}

func (op *mailOp) DeleteByMapTx(ext sqlx.Ext, m map[string]interface{}) (int64, error) {
	var params []interface{}
	sql := "delete from mail where 1=1 "
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
