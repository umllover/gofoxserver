package model

import (
	"fmt"
	"mj/gameServer/db"

	"github.com/jmoiron/sqlx"
	"github.com/lovelly/leaf/log"
)

//This file is generate by scripts,don't edit it

//record_outcard_ddz
//

// +gen *
type RecordOutcardDdz struct {
	RecordID   int64  `db:"RecordID" json:"RecordID"`     //
	CreateTime int    `db:"CreateTime" json:"CreateTime"` // 创建时间
	CardData   string `db:"CardData" json:"CardData"`     // 牌数据，数组转成字符串
}

type recordOutcardDdzOp struct{}

var RecordOutcardDdzOp = &recordOutcardDdzOp{}
var DefaultRecordOutcardDdz = &RecordOutcardDdz{}

// 按主键查询. 注:未找到记录的话将触发sql.ErrNoRows错误，返回nil, false
func (op *recordOutcardDdzOp) Get(RecordID int64) (*RecordOutcardDdz, bool) {
	obj := &RecordOutcardDdz{}
	sql := "select * from record_outcard_ddz where RecordID=? "
	err := db.DB.Get(obj, sql,
		RecordID,
	)

	if err != nil {
		log.Error("Get data error:%v", err.Error())
		return nil, false
	}
	return obj, true
}
func (op *recordOutcardDdzOp) SelectAll() ([]*RecordOutcardDdz, error) {
	objList := []*RecordOutcardDdz{}
	sql := "select * from record_outcard_ddz "
	err := db.DB.Select(&objList, sql)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return objList, nil
}

func (op *recordOutcardDdzOp) QueryByMap(m map[string]interface{}) ([]*RecordOutcardDdz, error) {
	result := []*RecordOutcardDdz{}
	var params []interface{}

	sql := "select * from record_outcard_ddz where 1=1 "
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

func (op *recordOutcardDdzOp) GetByMap(m map[string]interface{}) (*RecordOutcardDdz, error) {
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
func (i *RecordOutcardDdz) Insert() error {
    err := db.DBMap.Insert(i)
    if err != nil{
		log.Error("Insert sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 插入数据，自增长字段将被忽略
func (op *recordOutcardDdzOp) Insert(m *RecordOutcardDdz) (int64, error) {
	return op.InsertTx(db.DB, m)
}

// 插入数据，自增长字段将被忽略
func (op *recordOutcardDdzOp) InsertTx(ext sqlx.Ext, m *RecordOutcardDdz) (int64, error) {
	sql := "insert into record_outcard_ddz(CreateTime,CardData) values(?,?)"
	result, err := ext.Exec(sql,
		m.CreateTime,
		m.CardData,
	)
	if err != nil {
		log.Error("InsertTx sql error:%v, data:%v", err.Error(), m)
		return -1, err
	}
	affected, _ := result.LastInsertId()
	return affected, nil
}

//存在就更新， 不存在就插入
func (op *recordOutcardDdzOp) InsertUpdate(obj *RecordOutcardDdz, m map[string]interface{}) error {
	sql := "insert into record_outcard_ddz(CreateTime,CardData) values(?,?) ON DUPLICATE KEY UPDATE "
	var params = []interface{}{obj.CreateTime,
		obj.CardData,
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
func (i *RecordOutcardDdz) Update()  error {
    _,err := db.DBMap.Update(i)
    if err != nil{
		log.Error("update sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *recordOutcardDdzOp) Update(m *RecordOutcardDdz) error {
	return op.UpdateTx(db.DB, m)
}

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *recordOutcardDdzOp) UpdateTx(ext sqlx.Ext, m *RecordOutcardDdz) error {
	sql := `update record_outcard_ddz set CreateTime=?,CardData=? where RecordID=?`
	_, err := ext.Exec(sql,
		m.CreateTime,
		m.CardData,
		m.RecordID,
	)

	if err != nil {
		log.Error("update sql error:%v, data:%v", err.Error(), m)
		return err
	}

	return nil
}

// 用主键做条件，更新map里包含的字段名
func (op *recordOutcardDdzOp) UpdateWithMap(RecordID int64, m map[string]interface{}) error {
	return op.UpdateWithMapTx(db.DB, RecordID, m)
}

// 用主键做条件，更新map里包含的字段名
func (op *recordOutcardDdzOp) UpdateWithMapTx(ext sqlx.Ext, RecordID int64, m map[string]interface{}) error {

	sql := `update record_outcard_ddz set %s where 1=1 and RecordID=? ;`

	var params []interface{}
	var set_sql string
	for k, v := range m {
		if set_sql != "" {
			set_sql += ","
		}
		set_sql += fmt.Sprintf(" %s=? ", k)
		params = append(params, v)
	}
	params = append(params, RecordID)
	_, err := ext.Exec(fmt.Sprintf(sql, set_sql), params...)
	return err
}

/*
func (i *RecordOutcardDdz) Delete() error{
    _,err := db.DBMap.Delete(i)
	log.Error("Delete sql error:%v", err.Error())
    return err
}
*/
// 根据主键删除相关记录
func (op *recordOutcardDdzOp) Delete(RecordID int64) error {
	return op.DeleteTx(db.DB, RecordID)
}

// 根据主键删除相关记录,Tx
func (op *recordOutcardDdzOp) DeleteTx(ext sqlx.Ext, RecordID int64) error {
	sql := `delete from record_outcard_ddz where 1=1
        and RecordID=?
        `
	_, err := ext.Exec(sql,
		RecordID,
	)
	return err
}

// 返回符合查询条件的记录数
func (op *recordOutcardDdzOp) CountByMap(m map[string]interface{}) (int64, error) {

	var params []interface{}
	sql := `select count(*) from record_outcard_ddz where 1=1 `
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

func (op *recordOutcardDdzOp) DeleteByMap(m map[string]interface{}) (int64, error) {
	return op.DeleteByMapTx(db.DB, m)
}

func (op *recordOutcardDdzOp) DeleteByMapTx(ext sqlx.Ext, m map[string]interface{}) (int64, error) {
	var params []interface{}
	sql := "delete from record_outcard_ddz where 1=1 "
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
