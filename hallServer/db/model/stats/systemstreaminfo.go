package stats

import (
	"fmt"
	"mj/hallServer/db"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lovelly/leaf/log"
)

//This file is generate by scripts,don't edit it

//systemstreaminfo
//

// +gen *
type Systemstreaminfo struct {
	DateID              int        `db:"DateID" json:"DateID"`                           // 日期标识
	WebLogonSuccess     int        `db:"WebLogonSuccess" json:"WebLogonSuccess"`         // 登录成功
	WebRegisterSuccess  int        `db:"WebRegisterSuccess" json:"WebRegisterSuccess"`   // 注册成功
	GameLogonSuccess    int        `db:"GameLogonSuccess" json:"GameLogonSuccess"`       // 登录成功
	GameRegisterSuccess int        `db:"GameRegisterSuccess" json:"GameRegisterSuccess"` // 注册成功
	CollectDate         *time.Time `db:"CollectDate" json:"CollectDate"`                 // 录入时间
}

type systemstreaminfoOp struct{}

var SystemstreaminfoOp = &systemstreaminfoOp{}
var DefaultSystemstreaminfo = &Systemstreaminfo{}

// 按主键查询. 注:未找到记录的话将触发sql.ErrNoRows错误，返回nil, false
func (op *systemstreaminfoOp) Get(DateID int) (*Systemstreaminfo, bool) {
	obj := &Systemstreaminfo{}
	sql := "select * from systemstreaminfo where DateID=? "
	err := db.StatsDB.Get(obj, sql,
		DateID,
	)

	if err != nil {
		log.Error("Get data error:%v", err.Error())
		return nil, false
	}
	return obj, true
}
func (op *systemstreaminfoOp) SelectAll() ([]*Systemstreaminfo, error) {
	objList := []*Systemstreaminfo{}
	sql := "select * from systemstreaminfo "
	err := db.StatsDB.Select(&objList, sql)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return objList, nil
}

func (op *systemstreaminfoOp) QueryByMap(m map[string]interface{}) ([]*Systemstreaminfo, error) {
	result := []*Systemstreaminfo{}
	var params []interface{}

	sql := "select * from systemstreaminfo where 1=1 "
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

func (op *systemstreaminfoOp) GetByMap(m map[string]interface{}) (*Systemstreaminfo, error) {
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
func (i *Systemstreaminfo) Insert() error {
    err := db.StatsDBMap.Insert(i)
    if err != nil{
		log.Error("Insert sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 插入数据，自增长字段将被忽略
func (op *systemstreaminfoOp) Insert(m *Systemstreaminfo) (int64, error) {
	return op.InsertTx(db.StatsDB, m)
}

// 插入数据，自增长字段将被忽略
func (op *systemstreaminfoOp) InsertTx(ext sqlx.Ext, m *Systemstreaminfo) (int64, error) {
	sql := "insert into systemstreaminfo(DateID,WebLogonSuccess,WebRegisterSuccess,GameLogonSuccess,GameRegisterSuccess,CollectDate) values(?,?,?,?,?,?)"
	result, err := ext.Exec(sql,
		m.DateID,
		m.WebLogonSuccess,
		m.WebRegisterSuccess,
		m.GameLogonSuccess,
		m.GameRegisterSuccess,
		m.CollectDate,
	)
	if err != nil {
		log.Error("InsertTx sql error:%v, data:%v", err.Error(), m)
		return -1, err
	}
	affected, _ := result.LastInsertId()
	return affected, nil
}

//存在就更新， 不存在就插入
func (op *systemstreaminfoOp) InsertUpdate(obj *Systemstreaminfo, m map[string]interface{}) error {
	sql := "insert into systemstreaminfo(DateID,WebLogonSuccess,WebRegisterSuccess,GameLogonSuccess,GameRegisterSuccess,CollectDate) values(?,?,?,?,?,?) ON DUPLICATE KEY UPDATE "
	var params = []interface{}{obj.DateID,
		obj.WebLogonSuccess,
		obj.WebRegisterSuccess,
		obj.GameLogonSuccess,
		obj.GameRegisterSuccess,
		obj.CollectDate,
	}
	var set_sql string
	for k, v := range m {
		if set_sql != "" {
			set_sql += ","
		}
		set_sql += fmt.Sprintf(" %s=? ", k)
		params = append(params, v)
	}

	_, err := db.StatsDB.Exec(sql+set_sql, params...)
	return err
}

/*
func (i *Systemstreaminfo) Update()  error {
    _,err := db.StatsDBMap.Update(i)
    if err != nil{
		log.Error("update sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *systemstreaminfoOp) Update(m *Systemstreaminfo) error {
	return op.UpdateTx(db.StatsDB, m)
}

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *systemstreaminfoOp) UpdateTx(ext sqlx.Ext, m *Systemstreaminfo) error {
	sql := `update systemstreaminfo set WebLogonSuccess=?,WebRegisterSuccess=?,GameLogonSuccess=?,GameRegisterSuccess=?,CollectDate=? where DateID=?`
	_, err := ext.Exec(sql,
		m.WebLogonSuccess,
		m.WebRegisterSuccess,
		m.GameLogonSuccess,
		m.GameRegisterSuccess,
		m.CollectDate,
		m.DateID,
	)

	if err != nil {
		log.Error("update sql error:%v, data:%v", err.Error(), m)
		return err
	}

	return nil
}

// 用主键做条件，更新map里包含的字段名
func (op *systemstreaminfoOp) UpdateWithMap(DateID int, m map[string]interface{}) error {
	return op.UpdateWithMapTx(db.StatsDB, DateID, m)
}

// 用主键做条件，更新map里包含的字段名
func (op *systemstreaminfoOp) UpdateWithMapTx(ext sqlx.Ext, DateID int, m map[string]interface{}) error {

	sql := `update systemstreaminfo set %s where 1=1 and DateID=? ;`

	var params []interface{}
	var set_sql string
	for k, v := range m {
		if set_sql != "" {
			set_sql += ","
		}
		set_sql += fmt.Sprintf(" %s=? ", k)
		params = append(params, v)
	}
	params = append(params, DateID)
	_, err := ext.Exec(fmt.Sprintf(sql, set_sql), params...)
	return err
}

/*
func (i *Systemstreaminfo) Delete() error{
    _,err := db.StatsDBMap.Delete(i)
	log.Error("Delete sql error:%v", err.Error())
    return err
}
*/
// 根据主键删除相关记录
func (op *systemstreaminfoOp) Delete(DateID int) error {
	return op.DeleteTx(db.StatsDB, DateID)
}

// 根据主键删除相关记录,Tx
func (op *systemstreaminfoOp) DeleteTx(ext sqlx.Ext, DateID int) error {
	sql := `delete from systemstreaminfo where 1=1
        and DateID=?
        `
	_, err := ext.Exec(sql,
		DateID,
	)
	return err
}

// 返回符合查询条件的记录数
func (op *systemstreaminfoOp) CountByMap(m map[string]interface{}) (int64, error) {

	var params []interface{}
	sql := `select count(*) from systemstreaminfo where 1=1 `
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

func (op *systemstreaminfoOp) DeleteByMap(m map[string]interface{}) (int64, error) {
	return op.DeleteByMapTx(db.StatsDB, m)
}

func (op *systemstreaminfoOp) DeleteByMapTx(ext sqlx.Ext, m map[string]interface{}) (int64, error) {
	var params []interface{}
	sql := "delete from systemstreaminfo where 1=1 "
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
