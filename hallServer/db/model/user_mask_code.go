package model

import (
	"errors"
	"fmt"
	"mj/hallServer/db"

	"github.com/jmoiron/sqlx"
	"github.com/lovelly/leaf/log"
)

//This file is generate by scripts,don't edit it

//user_mask_code
//

// +gen *
type UserMaskCode struct {
	UserId      int64  `db:"user_id" json:"user_id"`           //
	PhomeNumber string `db:"phome_number" json:"phome_number"` // 电话号码
	MaskCode    int    `db:"mask_code" json:"mask_code"`       // 验证按
	CreatorTime string `db:"creator_time" json:"creator_time"` //
}

type userMaskCodeOp struct{}

var UserMaskCodeOp = &userMaskCodeOp{}
var DefaultUserMaskCode = &UserMaskCode{}

// 按主键查询. 注:未找到记录的话将触发sql.ErrNoRows错误，返回nil, false
func (op *userMaskCodeOp) Get(user_id int64) (*UserMaskCode, bool) {
	obj := &UserMaskCode{}
	sql := "select * from user_mask_code where user_id=? "
	err := db.DB.Get(obj, sql,
		user_id,
	)

	if err != nil {
		log.Error("Get data error:%v", err.Error())
		return nil, false
	}
	return obj, true
}
func (op *userMaskCodeOp) SelectAll() ([]*UserMaskCode, error) {
	objList := []*UserMaskCode{}
	sql := "select * from user_mask_code "
	err := db.DB.Select(&objList, sql)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return objList, nil
}

func (op *userMaskCodeOp) QueryByMap(m map[string]interface{}) ([]*UserMaskCode, error) {
	result := []*UserMaskCode{}
	var params []interface{}

	sql := "select * from user_mask_code where 1=1 "
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

func (op *userMaskCodeOp) GetByMap(m map[string]interface{}) (*UserMaskCode, error) {
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
func (i *UserMaskCode) Insert() error {
    err := db.DBMap.Insert(i)
    if err != nil{
		log.Error("Insert sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 插入数据，自增长字段将被忽略
func (op *userMaskCodeOp) Insert(m *UserMaskCode) (int64, error) {
	return op.InsertTx(db.DB, m)
}

// 插入数据，自增长字段将被忽略
func (op *userMaskCodeOp) InsertTx(ext sqlx.Ext, m *UserMaskCode) (int64, error) {
	sql := "insert into user_mask_code(user_id,phome_number,mask_code,creator_time) values(?,?,?,?)"
	result, err := ext.Exec(sql,
		m.UserId,
		m.PhomeNumber,
		m.MaskCode,
		m.CreatorTime,
	)
	if err != nil {
		log.Error("InsertTx sql error:%v, data:%v", err.Error(), m)
		return -1, err
	}
	affected, _ := result.LastInsertId()
	return affected, nil
}

/*
func (i *UserMaskCode) Update()  error {
    _,err := db.DBMap.Update(i)
    if err != nil{
		log.Error("update sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *userMaskCodeOp) Update(m *UserMaskCode) error {
	return op.UpdateTx(db.DB, m)
}

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *userMaskCodeOp) UpdateTx(ext sqlx.Ext, m *UserMaskCode) error {
	sql := `update user_mask_code set phome_number=?,mask_code=?,creator_time=? where user_id=?`
	_, err := ext.Exec(sql,
		m.PhomeNumber,
		m.MaskCode,
		m.CreatorTime,
		m.UserId,
	)

	if err != nil {
		log.Error("update sql error:%v, data:%v", err.Error(), m)
		return err
	}

	return nil
}

// 用主键做条件，更新map里包含的字段名
func (op *userMaskCodeOp) UpdateWithMap(user_id int64, m map[string]interface{}) error {
	return op.UpdateWithMapTx(db.DB, user_id, m)
}

// 用主键做条件，更新map里包含的字段名
func (op *userMaskCodeOp) UpdateWithMapTx(ext sqlx.Ext, user_id int64, m map[string]interface{}) error {

	sql := `update user_mask_code set %s where 1=1 and user_id=? ;`

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
func (i *UserMaskCode) Delete() error{
    _,err := db.DBMap.Delete(i)
	log.Error("Delete sql error:%v", err.Error())
    return err
}
*/
// 根据主键删除相关记录
func (op *userMaskCodeOp) Delete(user_id int64) error {
	return op.DeleteTx(db.DB, user_id)
}

// 根据主键删除相关记录,Tx
func (op *userMaskCodeOp) DeleteTx(ext sqlx.Ext, user_id int64) error {
	sql := `delete from user_mask_code where 1=1
        and user_id=?
        `
	_, err := ext.Exec(sql,
		user_id,
	)
	return err
}

// 返回符合查询条件的记录数
func (op *userMaskCodeOp) CountByMap(m map[string]interface{}) (int64, error) {

	var params []interface{}
	sql := `select count(*) from user_mask_code where 1=1 `
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

func (op *userMaskCodeOp) DeleteByMap(m map[string]interface{}) (int64, error) {
	return op.DeleteByMapTx(db.DB, m)
}

func (op *userMaskCodeOp) DeleteByMapTx(ext sqlx.Ext, m map[string]interface{}) (int64, error) {
	var params []interface{}
	sql := "delete from user_mask_code where 1=1 "
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
