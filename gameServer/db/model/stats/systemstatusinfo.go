package stats

import (
	"errors"
	"fmt"
	"mj/gameServer/db"

	"github.com/jmoiron/sqlx"
	"github.com/lovelly/leaf/log"
)

//This file is generate by scripts,don't edit it

//systemstatusinfo
//

// +gen *
type Systemstatusinfo struct {
	StatusName        string `db:"StatusName" json:"StatusName"`               // 状态名字
	StatusValue       int    `db:"StatusValue" json:"StatusValue"`             // 状态数值
	StatusString      string `db:"StatusString" json:"StatusString"`           // 状态字符
	StatusTip         string `db:"StatusTip" json:"StatusTip"`                 // 状态显示名称
	StatusDescription string `db:"StatusDescription" json:"StatusDescription"` // 字符的描述
	SortID            int    `db:"SortID" json:"SortID"`                       //
}

type systemstatusinfoOp struct{}

var SystemstatusinfoOp = &systemstatusinfoOp{}
var DefaultSystemstatusinfo = &Systemstatusinfo{}

// 按主键查询. 注:未找到记录的话将触发sql.ErrNoRows错误，返回nil, false
func (op *systemstatusinfoOp) Get(StatusName string) (*Systemstatusinfo, bool) {
	obj := &Systemstatusinfo{}
	sql := "select * from systemstatusinfo where StatusName=? "
	err := db.StatsDB.Get(obj, sql,
		StatusName,
	)

	if err != nil {
		log.Error("Get data error:%v", err.Error())
		return nil, false
	}
	return obj, true
}
func (op *systemstatusinfoOp) SelectAll() ([]*Systemstatusinfo, error) {
	objList := []*Systemstatusinfo{}
	sql := "select * from systemstatusinfo "
	err := db.StatsDB.Select(&objList, sql)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return objList, nil
}

func (op *systemstatusinfoOp) QueryByMap(m map[string]interface{}) ([]*Systemstatusinfo, error) {
	result := []*Systemstatusinfo{}
	var params []interface{}

	sql := "select * from systemstatusinfo where 1=1 "
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

func (op *systemstatusinfoOp) GetByMap(m map[string]interface{}) (*Systemstatusinfo, error) {
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
func (i *Systemstatusinfo) Insert() error {
    err := db.StatsDBMap.Insert(i)
    if err != nil{
		log.Error("Insert sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 插入数据，自增长字段将被忽略
func (op *systemstatusinfoOp) Insert(m *Systemstatusinfo) (int64, error) {
	return op.InsertTx(db.StatsDB, m)
}

// 插入数据，自增长字段将被忽略
func (op *systemstatusinfoOp) InsertTx(ext sqlx.Ext, m *Systemstatusinfo) (int64, error) {
	sql := "insert into systemstatusinfo(StatusName,StatusValue,StatusString,StatusTip,StatusDescription,SortID) values(?,?,?,?,?,?)"
	result, err := ext.Exec(sql,
		m.StatusName,
		m.StatusValue,
		m.StatusString,
		m.StatusTip,
		m.StatusDescription,
		m.SortID,
	)
	if err != nil {
		log.Error("InsertTx sql error:%v, data:%v", err.Error(), m)
		return -1, err
	}
	affected, _ := result.LastInsertId()
	return affected, nil
}

//存在就更新， 不存在就插入
func (op *systemstatusinfoOp) InsertUpdate(obj *Systemstatusinfo, m map[string]interface{}) ( error) {
    sql := "insert into systemstatusinfo(StatusName,StatusValue,StatusString,StatusTip,StatusDescription,SortID) values(?,?,?,?,?,?) ON DUPLICATE KEY UPDATE "
    var params = []interface{}{ obj.StatusName,
        obj.StatusValue,
        obj.StatusString,
        obj.StatusTip,
        obj.StatusDescription,
        obj.SortID,
        }
    var set_sql string
    for k, v := range m{
		if set_sql != "" {
			set_sql += ","
		}
        set_sql += fmt.Sprintf(" %s=? ", k)
        params = append(params, v)
    }

    _, err := db.StatsDB.Exec(sql + set_sql, params...)
    return err
}


/*
func (i *Systemstatusinfo) Update()  error {
    _,err := db.StatsDBMap.Update(i)
    if err != nil{
		log.Error("update sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *systemstatusinfoOp) Update(m *Systemstatusinfo) error {
	return op.UpdateTx(db.StatsDB, m)
}

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *systemstatusinfoOp) UpdateTx(ext sqlx.Ext, m *Systemstatusinfo) error {
	sql := `update systemstatusinfo set StatusValue=?,StatusString=?,StatusTip=?,StatusDescription=?,SortID=? where StatusName=?`
	_, err := ext.Exec(sql,
		m.StatusValue,
		m.StatusString,
		m.StatusTip,
		m.StatusDescription,
		m.SortID,
		m.StatusName,
	)

	if err != nil {
		log.Error("update sql error:%v, data:%v", err.Error(), m)
		return err
	}

	return nil
}

// 用主键做条件，更新map里包含的字段名
func (op *systemstatusinfoOp) UpdateWithMap(StatusName string, m map[string]interface{}) error {
	return op.UpdateWithMapTx(db.StatsDB, StatusName, m)
}

// 用主键做条件，更新map里包含的字段名
func (op *systemstatusinfoOp) UpdateWithMapTx(ext sqlx.Ext, StatusName string, m map[string]interface{}) error {

	sql := `update systemstatusinfo set %s where 1=1 and StatusName=? ;`

	var params []interface{}
	var set_sql string
	for k, v := range m {
		if set_sql != "" {
			set_sql += ","
		}
		set_sql += fmt.Sprintf(" %s=? ", k)
		params = append(params, v)
	}
	params = append(params, StatusName)
	_, err := ext.Exec(fmt.Sprintf(sql, set_sql), params...)
	return err
}

/*
func (i *Systemstatusinfo) Delete() error{
    _,err := db.StatsDBMap.Delete(i)
	log.Error("Delete sql error:%v", err.Error())
    return err
}
*/
// 根据主键删除相关记录
func (op *systemstatusinfoOp) Delete(StatusName string) error {
	return op.DeleteTx(db.StatsDB, StatusName)
}

// 根据主键删除相关记录,Tx
func (op *systemstatusinfoOp) DeleteTx(ext sqlx.Ext, StatusName string) error {
	sql := `delete from systemstatusinfo where 1=1
        and StatusName=?
        `
	_, err := ext.Exec(sql,
		StatusName,
	)
	return err
}

// 返回符合查询条件的记录数
func (op *systemstatusinfoOp) CountByMap(m map[string]interface{}) (int64, error) {

	var params []interface{}
	sql := `select count(*) from systemstatusinfo where 1=1 `
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

func (op *systemstatusinfoOp) DeleteByMap(m map[string]interface{}) (int64, error) {
	return op.DeleteByMapTx(db.StatsDB, m)
}

func (op *systemstatusinfoOp) DeleteByMapTx(ext sqlx.Ext, m map[string]interface{}) (int64, error) {
	var params []interface{}
	sql := "delete from systemstatusinfo where 1=1 "
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
