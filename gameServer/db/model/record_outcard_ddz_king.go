package model

import (
	"errors"
	"fmt"
	"mj/gameServer/db"

	"github.com/jmoiron/sqlx"
	"github.com/lovelly/leaf/log"
)

//This file is generate by scripts,don't edit it

//record_outcard_ddz_king
//

// +gen *
type RecordOutcardDdzKing struct {
	RecordID   int64  `db:"RecordID" json:"RecordID"`     // 记录ID八王表
	CreateTime int    `db:"CreateTime" json:"CreateTime"` // 创建时间
	CardData   string `db:"CardData" json:"CardData"`     // 牌数据，数组转成字符串
}

type recordOutcardDdzKingOp struct{}

var RecordOutcardDdzKingOp = &recordOutcardDdzKingOp{}
var DefaultRecordOutcardDdzKing = &RecordOutcardDdzKing{}

// 按主键查询. 注:未找到记录的话将触发sql.ErrNoRows错误，返回nil, false
func (op *recordOutcardDdzKingOp) Get(RecordID int64) (*RecordOutcardDdzKing, bool) {
	obj := &RecordOutcardDdzKing{}
	sql := "select * from record_outcard_ddz_king where RecordID=? "
	err := db.DB.Get(obj, sql,
		RecordID,
	)

	if err != nil {
		log.Error("Get data error:%v", err.Error())
		return nil, false
	}
	return obj, true
}
func (op *recordOutcardDdzKingOp) SelectAll() ([]*RecordOutcardDdzKing, error) {
	objList := []*RecordOutcardDdzKing{}
	sql := "select * from record_outcard_ddz_king "
	err := db.DB.Select(&objList, sql)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return objList, nil
}

func (op *recordOutcardDdzKingOp) QueryByMap(m map[string]interface{}) ([]*RecordOutcardDdzKing, error) {
	result := []*RecordOutcardDdzKing{}
	var params []interface{}

	sql := "select * from record_outcard_ddz_king where 1=1 "
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

func (op *recordOutcardDdzKingOp) GetByMap(m map[string]interface{}) (*RecordOutcardDdzKing, error) {
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
func (i *RecordOutcardDdzKing) Insert() error {
    err := db.DBMap.Insert(i)
    if err != nil{
		log.Error("Insert sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 插入数据，自增长字段将被忽略
func (op *recordOutcardDdzKingOp) Insert(m *RecordOutcardDdzKing) (int64, error) {
	return op.InsertTx(db.DB, m)
}

// 插入数据，自增长字段将被忽略
func (op *recordOutcardDdzKingOp) InsertTx(ext sqlx.Ext, m *RecordOutcardDdzKing) (int64, error) {
	sql := "insert into record_outcard_ddz_king(RecordID,CreateTime,CardData) values(?,?,?)"
	result, err := ext.Exec(sql,
		m.RecordID,
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
func (op *recordOutcardDdzKingOp) InsertUpdate(obj *RecordOutcardDdzKing, m map[string]interface{}) error {
	sql := "insert into record_outcard_ddz_king(RecordID,CreateTime,CardData) values(?,?,?) ON DUPLICATE KEY UPDATE "
	var params = []interface{}{obj.RecordID,
		obj.CreateTime,
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
func (i *RecordOutcardDdzKing) Update()  error {
    _,err := db.DBMap.Update(i)
    if err != nil{
		log.Error("update sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *recordOutcardDdzKingOp) Update(m *RecordOutcardDdzKing) error {
	return op.UpdateTx(db.DB, m)
}

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *recordOutcardDdzKingOp) UpdateTx(ext sqlx.Ext, m *RecordOutcardDdzKing) error {
	sql := `update record_outcard_ddz_king set CreateTime=?,CardData=? where RecordID=?`
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
func (op *recordOutcardDdzKingOp) UpdateWithMap(RecordID int64, m map[string]interface{}) error {
	return op.UpdateWithMapTx(db.DB, RecordID, m)
}

// 用主键做条件，更新map里包含的字段名
func (op *recordOutcardDdzKingOp) UpdateWithMapTx(ext sqlx.Ext, RecordID int64, m map[string]interface{}) error {

	sql := `update record_outcard_ddz_king set %s where 1=1 and RecordID=? ;`

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
func (i *RecordOutcardDdzKing) Delete() error{
    _,err := db.DBMap.Delete(i)
	log.Error("Delete sql error:%v", err.Error())
    return err
}
*/
// 根据主键删除相关记录
func (op *recordOutcardDdzKingOp) Delete(RecordID int64) error {
	return op.DeleteTx(db.DB, RecordID)
}

// 根据主键删除相关记录,Tx
func (op *recordOutcardDdzKingOp) DeleteTx(ext sqlx.Ext, RecordID int64) error {
	sql := `delete from record_outcard_ddz_king where 1=1
        and RecordID=?
        `
	_, err := ext.Exec(sql,
		RecordID,
	)
	return err
}

// 返回符合查询条件的记录数
func (op *recordOutcardDdzKingOp) CountByMap(m map[string]interface{}) (int64, error) {

	var params []interface{}
	sql := `select count(*) from record_outcard_ddz_king where 1=1 `
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

func (op *recordOutcardDdzKingOp) DeleteByMap(m map[string]interface{}) (int64, error) {
	return op.DeleteByMapTx(db.DB, m)
}

func (op *recordOutcardDdzKingOp) DeleteByMapTx(ext sqlx.Ext, m map[string]interface{}) (int64, error) {
	var params []interface{}
	sql := "delete from record_outcard_ddz_king where 1=1 "
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
