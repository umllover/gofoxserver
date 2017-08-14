package model

import (
	"fmt"
	"mj/hallServer/db"

	"github.com/jmoiron/sqlx"
	"github.com/lovelly/leaf/log"
)

//This file is generate by scripts,don't edit it

//inc_userid
//

// +gen *
type IncUserid struct {
	NodeId int   `db:"node_id" json:"node_id"` //
	IncId  int64 `db:"inc_id" json:"inc_id"`   //
}

type incUseridOp struct{}

var IncUseridOp = &incUseridOp{}
var DefaultIncUserid = &IncUserid{}

// 按主键查询. 注:未找到记录的话将触发sql.ErrNoRows错误，返回nil, false
func (op *incUseridOp) Get(node_id int) (*IncUserid, bool) {
	obj := &IncUserid{}
	sql := "select * from inc_userid where node_id=? "
	err := db.DB.Get(obj, sql,
		node_id,
	)

	if err != nil {
		log.Error("Get data error:%v", err.Error())
		return nil, false
	}
	return obj, true
}
func (op *incUseridOp) SelectAll() ([]*IncUserid, error) {
	objList := []*IncUserid{}
	sql := "select * from inc_userid "
	err := db.DB.Select(&objList, sql)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return objList, nil
}

func (op *incUseridOp) QueryByMap(m map[string]interface{}) ([]*IncUserid, error) {
	result := []*IncUserid{}
	var params []interface{}

	sql := "select * from inc_userid where 1=1 "
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

func (op *incUseridOp) GetByMap(m map[string]interface{}) (*IncUserid, error) {
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
func (i *IncUserid) Insert() error {
    err := db.DBMap.Insert(i)
    if err != nil{
		log.Error("Insert sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 插入数据，自增长字段将被忽略
func (op *incUseridOp) Insert(m *IncUserid) (int64, error) {
	return op.InsertTx(db.DB, m)
}

// 插入数据，自增长字段将被忽略
func (op *incUseridOp) InsertTx(ext sqlx.Ext, m *IncUserid) (int64, error) {
	sql := "insert into inc_userid(node_id,inc_id) values(?,?)"
	result, err := ext.Exec(sql,
		m.NodeId,
		m.IncId,
	)
	if err != nil {
		log.Error("InsertTx sql error:%v, data:%v", err.Error(), m)
		return -1, err
	}
	affected, _ := result.LastInsertId()
	return affected, nil
}

//存在就更新， 不存在就插入
func (op *incUseridOp) InsertUpdate(obj *IncUserid, m map[string]interface{}) error {
	sql := "insert into inc_userid(node_id,inc_id) values(?,?) ON DUPLICATE KEY UPDATE "
	var params = []interface{}{obj.NodeId,
		obj.IncId,
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
func (i *IncUserid) Update()  error {
    _,err := db.DBMap.Update(i)
    if err != nil{
		log.Error("update sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *incUseridOp) Update(m *IncUserid) error {
	return op.UpdateTx(db.DB, m)
}

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *incUseridOp) UpdateTx(ext sqlx.Ext, m *IncUserid) error {
	sql := `update inc_userid set inc_id=? where node_id=?`
	_, err := ext.Exec(sql,
		m.IncId,
		m.NodeId,
	)

	if err != nil {
		log.Error("update sql error:%v, data:%v", err.Error(), m)
		return err
	}

	return nil
}

// 用主键做条件，更新map里包含的字段名
func (op *incUseridOp) UpdateWithMap(node_id int, m map[string]interface{}) error {
	return op.UpdateWithMapTx(db.DB, node_id, m)
}

// 用主键做条件，更新map里包含的字段名
func (op *incUseridOp) UpdateWithMapTx(ext sqlx.Ext, node_id int, m map[string]interface{}) error {

	sql := `update inc_userid set %s where 1=1 and node_id=? ;`

	var params []interface{}
	var set_sql string
	for k, v := range m {
		if set_sql != "" {
			set_sql += ","
		}
		set_sql += fmt.Sprintf(" %s=? ", k)
		params = append(params, v)
	}
	params = append(params, node_id)
	_, err := ext.Exec(fmt.Sprintf(sql, set_sql), params...)
	return err
}

/*
func (i *IncUserid) Delete() error{
    _,err := db.DBMap.Delete(i)
	log.Error("Delete sql error:%v", err.Error())
    return err
}
*/
// 根据主键删除相关记录
func (op *incUseridOp) Delete(node_id int) error {
	return op.DeleteTx(db.DB, node_id)
}

// 根据主键删除相关记录,Tx
func (op *incUseridOp) DeleteTx(ext sqlx.Ext, node_id int) error {
	sql := `delete from inc_userid where 1=1
        and node_id=?
        `
	_, err := ext.Exec(sql,
		node_id,
	)
	return err
}

// 返回符合查询条件的记录数
func (op *incUseridOp) CountByMap(m map[string]interface{}) (int64, error) {

	var params []interface{}
	sql := `select count(*) from inc_userid where 1=1 `
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

func (op *incUseridOp) DeleteByMap(m map[string]interface{}) (int64, error) {
	return op.DeleteByMapTx(db.DB, m)
}

func (op *incUseridOp) DeleteByMapTx(ext sqlx.Ext, m map[string]interface{}) (int64, error) {
	var params []interface{}
	sql := "delete from inc_userid where 1=1 "
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
