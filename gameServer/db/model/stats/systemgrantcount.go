package stats

import (
	"fmt"
	"mj/gameServer/db"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lovelly/leaf/log"
)

//This file is generate by scripts,don't edit it

//systemgrantcount
//

// +gen *
type Systemgrantcount struct {
	DateID          int        `db:"DateID" json:"DateID"`                   //
	RegisterIP      string     `db:"RegisterIP" json:"RegisterIP"`           // 注册地址
	RegisterMachine string     `db:"RegisterMachine" json:"RegisterMachine"` // 注册机器
	GrantScore      int64      `db:"GrantScore" json:"GrantScore"`           // 赠送金币
	GrantCount      int64      `db:"GrantCount" json:"GrantCount"`           // 赠送次数
	CollectDate     *time.Time `db:"CollectDate" json:"CollectDate"`         // 收集时间
}

type systemgrantcountOp struct{}

var SystemgrantcountOp = &systemgrantcountOp{}
var DefaultSystemgrantcount = &Systemgrantcount{}

// 按主键查询. 注:未找到记录的话将触发sql.ErrNoRows错误，返回nil, false
func (op *systemgrantcountOp) Get(DateID int, RegisterIP string) (*Systemgrantcount, bool) {
	obj := &Systemgrantcount{}
	sql := "select * from systemgrantcount where DateID=? and RegisterIP=? "
	err := db.StatsDB.Get(obj, sql,
		DateID,
		RegisterIP,
	)

	if err != nil {
		log.Error("Get data error:%v", err.Error())
		return nil, false
	}
	return obj, true
}
func (op *systemgrantcountOp) SelectAll() ([]*Systemgrantcount, error) {
	objList := []*Systemgrantcount{}
	sql := "select * from systemgrantcount "
	err := db.StatsDB.Select(&objList, sql)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return objList, nil
}

func (op *systemgrantcountOp) QueryByMap(m map[string]interface{}) ([]*Systemgrantcount, error) {
	result := []*Systemgrantcount{}
	var params []interface{}

	sql := "select * from systemgrantcount where 1=1 "
	for k, v := range m {
		sql += fmt.Sprintf(" and %s=? ", k)
		params = append(params, v)
	}
	err := db.StatsDB.Select(&result, sql, params...)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return result, nil
}

func (op *systemgrantcountOp) QueryByMapQueryByMapComparison(m map[string]interface{}) ([]*Systemgrantcount, error) {
	result := []*Systemgrantcount{}
	var params []interface{}

	sql := "select * from systemgrantcount where 1=1 "
	for k, v := range m {
		sql += fmt.Sprintf(" and %s? ", k)
		params = append(params, v)
	}
	err := db.StatsDB.Select(&result, sql, params...)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return result, nil
}

func (op *systemgrantcountOp) GetByMap(m map[string]interface{}) (*Systemgrantcount, error) {
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
func (i *Systemgrantcount) Insert() error {
    err := db.StatsDBMap.Insert(i)
    if err != nil{
		log.Error("Insert sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 插入数据，自增长字段将被忽略
func (op *systemgrantcountOp) Insert(m *Systemgrantcount) (int64, error) {
	return op.InsertTx(db.StatsDB, m)
}

// 插入数据，自增长字段将被忽略
func (op *systemgrantcountOp) InsertTx(ext sqlx.Ext, m *Systemgrantcount) (int64, error) {
	sql := "insert into systemgrantcount(DateID,RegisterIP,RegisterMachine,GrantScore,GrantCount,CollectDate) values(?,?,?,?,?,?)"
	result, err := ext.Exec(sql,
		m.DateID,
		m.RegisterIP,
		m.RegisterMachine,
		m.GrantScore,
		m.GrantCount,
		m.CollectDate,
	)
	if err != nil {
		log.Error("InsertTx sql error:%v, data:%v", err.Error(), m)
		return -1, err
	}
	affected, _ := result.LastInsertId()
	return affected, nil
}

/*
func (i *Systemgrantcount) Update()  error {
    _,err := db.StatsDBMap.Update(i)
    if err != nil{
		log.Error("update sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *systemgrantcountOp) Update(m *Systemgrantcount) error {
	return op.UpdateTx(db.StatsDB, m)
}

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *systemgrantcountOp) UpdateTx(ext sqlx.Ext, m *Systemgrantcount) error {
	sql := `update systemgrantcount set RegisterMachine=?,GrantScore=?,GrantCount=?,CollectDate=? where DateID=? and RegisterIP=?`
	_, err := ext.Exec(sql,
		m.RegisterMachine,
		m.GrantScore,
		m.GrantCount,
		m.CollectDate,
		m.DateID,
		m.RegisterIP,
	)

	if err != nil {
		log.Error("update sql error:%v, data:%v", err.Error(), m)
		return err
	}

	return nil
}

// 用主键做条件，更新map里包含的字段名
func (op *systemgrantcountOp) UpdateWithMap(DateID int, RegisterIP string, m map[string]interface{}) error {
	return op.UpdateWithMapTx(db.StatsDB, DateID, RegisterIP, m)
}

// 用主键做条件，更新map里包含的字段名
func (op *systemgrantcountOp) UpdateWithMapTx(ext sqlx.Ext, DateID int, RegisterIP string, m map[string]interface{}) error {

	sql := `update systemgrantcount set %s where 1=1 and DateID=? and RegisterIP=? ;`

	var params []interface{}
	var set_sql string
	for k, v := range m {
		set_sql += fmt.Sprintf(" %s=? ", k)
		params = append(params, v)
	}
	params = append(params, DateID, RegisterIP)
	_, err := ext.Exec(fmt.Sprintf(sql, set_sql), params...)
	return err
}

/*
func (i *Systemgrantcount) Delete() error{
    _,err := db.StatsDBMap.Delete(i)
	log.Error("Delete sql error:%v", err.Error())
    return err
}
*/
// 根据主键删除相关记录
func (op *systemgrantcountOp) Delete(DateID int, RegisterIP string) error {
	return op.DeleteTx(db.StatsDB, DateID, RegisterIP)
}

// 根据主键删除相关记录,Tx
func (op *systemgrantcountOp) DeleteTx(ext sqlx.Ext, DateID int, RegisterIP string) error {
	sql := `delete from systemgrantcount where 1=1
        and DateID=?
        and RegisterIP=?
        `
	_, err := ext.Exec(sql,
		DateID,
		RegisterIP,
	)
	return err
}

// 返回符合查询条件的记录数
func (op *systemgrantcountOp) CountByMap(m map[string]interface{}) (int64, error) {

	var params []interface{}
	sql := `select count(*) from systemgrantcount where 1=1 `
	for k, v := range m {
		sql += fmt.Sprintf(" and  %s=? ", k)
		params = append(params, v)
	}
	count := int64(-1)
	err := db.StatsDB.Get(&count, sql, params...)
	if err != nil {
		log.Error("CountByMap  error:%v data :%v", err.Error(), m)
		return 0, err
	}
	return count, nil
}

func (op *systemgrantcountOp) DeleteByMap(m map[string]interface{}) (int64, error) {
	return op.DeleteByMapTx(db.StatsDB, m)
}

func (op *systemgrantcountOp) DeleteByMapTx(ext sqlx.Ext, m map[string]interface{}) (int64, error) {
	var params []interface{}
	sql := "delete from systemgrantcount where 1=1 "
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
