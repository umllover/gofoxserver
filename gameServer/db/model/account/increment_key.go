package account

import (
	"fmt"
	"mj/gameServer/db"

	"github.com/jmoiron/sqlx"
	"github.com/lovelly/leaf/log"
)

//This file is generate by scripts,don't edit it

//increment_key
//

// +gen *
type IncrementKey struct {
	IncrementName  string `db:"increment_name" json:"increment_name"`   //
	IncrementValue int64  `db:"increment_value" json:"increment_value"` //
}

type incrementKeyOp struct{}

var IncrementKeyOp = &incrementKeyOp{}
var DefaultIncrementKey = &IncrementKey{}

// 按主键查询. 注:未找到记录的话将触发sql.ErrNoRows错误，返回nil, false
func (op *incrementKeyOp) Get(increment_name string) (*IncrementKey, bool) {
	obj := &IncrementKey{}
	sql := "select * from increment_key where increment_name=? "
	err := db.AccountDB.Get(obj, sql,
		increment_name,
	)

	if err != nil {
		log.Error("Get data error:%v", err.Error())
		return nil, false
	}
	return obj, true
}
func (op *incrementKeyOp) SelectAll() ([]*IncrementKey, error) {
	objList := []*IncrementKey{}
	sql := "select * from increment_key "
	err := db.AccountDB.Select(&objList, sql)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return objList, nil
}

func (op *incrementKeyOp) QueryByMap(m map[string]interface{}) ([]*IncrementKey, error) {
	result := []*IncrementKey{}
	var params []interface{}

	sql := "select * from increment_key where 1=1 "
	for k, v := range m {
		sql += fmt.Sprintf(" and %s=? ", k)
		params = append(params, v)
	}
	err := db.AccountDB.Select(&result, sql, params...)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return result, nil
}

func (op *incrementKeyOp) GetByMap(m map[string]interface{}) (*IncrementKey, error) {
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
func (i *IncrementKey) Insert() error {
    err := db.AccountDBMap.Insert(i)
    if err != nil{
		log.Error("Insert sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 插入数据，自增长字段将被忽略
func (op *incrementKeyOp) Insert(m *IncrementKey) (int64, error) {
	return op.InsertTx(db.AccountDB, m)
}

// 插入数据，自增长字段将被忽略
func (op *incrementKeyOp) InsertTx(ext sqlx.Ext, m *IncrementKey) (int64, error) {
	sql := "insert into increment_key(increment_name,increment_value) values(?,?)"
	result, err := ext.Exec(sql,
		m.IncrementName,
		m.IncrementValue,
	)
	if err != nil {
		log.Error("InsertTx sql error:%v, data:%v", err.Error(), m)
		return -1, err
	}
	affected, _ := result.LastInsertId()
	return affected, nil
}

//存在就更新， 不存在就插入
func (op *incrementKeyOp) InsertUpdate(obj *IncrementKey, m map[string]interface{}) error {
	sql := "insert into increment_key(increment_name,increment_value) values(?,?) ON DUPLICATE KEY UPDATE "
	var params = []interface{}{obj.IncrementName,
		obj.IncrementValue,
	}
	var set_sql string
	for k, v := range m {
		if set_sql != "" {
			set_sql += ","
		}
		set_sql += fmt.Sprintf(" %s=? ", k)
		params = append(params, v)
	}

	_, err := db.AccountDB.Exec(sql+set_sql, params...)
	return err
}

/*
func (i *IncrementKey) Update()  error {
    _,err := db.AccountDBMap.Update(i)
    if err != nil{
		log.Error("update sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *incrementKeyOp) Update(m *IncrementKey) error {
	return op.UpdateTx(db.AccountDB, m)
}

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *incrementKeyOp) UpdateTx(ext sqlx.Ext, m *IncrementKey) error {
	sql := `update increment_key set increment_value=? where increment_name=?`
	_, err := ext.Exec(sql,
		m.IncrementValue,
		m.IncrementName,
	)

	if err != nil {
		log.Error("update sql error:%v, data:%v", err.Error(), m)
		return err
	}

	return nil
}

// 用主键做条件，更新map里包含的字段名
func (op *incrementKeyOp) UpdateWithMap(increment_name string, m map[string]interface{}) error {
	return op.UpdateWithMapTx(db.AccountDB, increment_name, m)
}

// 用主键做条件，更新map里包含的字段名
func (op *incrementKeyOp) UpdateWithMapTx(ext sqlx.Ext, increment_name string, m map[string]interface{}) error {

	sql := `update increment_key set %s where 1=1 and increment_name=? ;`

	var params []interface{}
	var set_sql string
	for k, v := range m {
		if set_sql != "" {
			set_sql += ","
		}
		set_sql += fmt.Sprintf(" %s=? ", k)
		params = append(params, v)
	}
	params = append(params, increment_name)
	_, err := ext.Exec(fmt.Sprintf(sql, set_sql), params...)
	return err
}

/*
func (i *IncrementKey) Delete() error{
    _,err := db.AccountDBMap.Delete(i)
	log.Error("Delete sql error:%v", err.Error())
    return err
}
*/
// 根据主键删除相关记录
func (op *incrementKeyOp) Delete(increment_name string) error {
	return op.DeleteTx(db.AccountDB, increment_name)
}

// 根据主键删除相关记录,Tx
func (op *incrementKeyOp) DeleteTx(ext sqlx.Ext, increment_name string) error {
	sql := `delete from increment_key where 1=1
        and increment_name=?
        `
	_, err := ext.Exec(sql,
		increment_name,
	)
	return err
}

// 返回符合查询条件的记录数
func (op *incrementKeyOp) CountByMap(m map[string]interface{}) (int64, error) {

	var params []interface{}
	sql := `select count(*) from increment_key where 1=1 `
	for k, v := range m {
		sql += fmt.Sprintf(" and  %s=? ", k)
		params = append(params, v)
	}
	count := int64(-1)
	err := db.AccountDB.Get(&count, sql, params...)
	if err != nil {
		log.Error("CountByMap  error:%v data :%v", err.Error(), m)
		return 0, err
	}
	return count, nil
}

func (op *incrementKeyOp) DeleteByMap(m map[string]interface{}) (int64, error) {
	return op.DeleteByMapTx(db.AccountDB, m)
}

func (op *incrementKeyOp) DeleteByMapTx(ext sqlx.Ext, m map[string]interface{}) (int64, error) {
	var params []interface{}
	sql := "delete from increment_key where 1=1 "
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
