package stats

import (
	"errors"
	"fmt"
	"mj/gameServer/db"

	"github.com/jmoiron/sqlx"
	"github.com/lovelly/leaf/log"
)

//This file is generate by scripts,don't edit it

//activity
//

// +gen *
type Activity struct {
	ActivityName  string `db:"activity_name" json:"activity_name"`   // 活动名
	ActivityType  int    `db:"activity_type" json:"activity_type"`   // 活动类别
	ActivityBegin None   `db:"activity_begin" json:"activity_begin"` // 活动开始时间
	ActivityEnd   None   `db:"activity_end" json:"activity_end"`     // 活动开始时间
}

type activityOp struct{}

var ActivityOp = &activityOp{}
var DefaultActivity = &Activity{}

// 按主键查询. 注:未找到记录的话将触发sql.ErrNoRows错误，返回nil, false
func (op *activityOp) Get(activity_name string) (*Activity, bool) {
	obj := &Activity{}
	sql := "select * from activity where activity_name=? "
	err := db.StatsDB.Get(obj, sql,
		activity_name,
	)

	if err != nil {
		log.Error("Get data error:%v", err.Error())
		return nil, false
	}
	return obj, true
}
func (op *activityOp) SelectAll() ([]*Activity, error) {
	objList := []*Activity{}
	sql := "select * from activity "
	err := db.StatsDB.Select(&objList, sql)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return objList, nil
}

func (op *activityOp) QueryByMap(m map[string]interface{}) ([]*Activity, error) {
	result := []*Activity{}
	var params []interface{}

	sql := "select * from activity where 1=1 "
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

func (op *activityOp) GetByMap(m map[string]interface{}) (*Activity, error) {
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
func (i *Activity) Insert() error {
    err := db.StatsDBMap.Insert(i)
    if err != nil{
		log.Error("Insert sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 插入数据，自增长字段将被忽略
func (op *activityOp) Insert(m *Activity) (int64, error) {
	return op.InsertTx(db.StatsDB, m)
}

// 插入数据，自增长字段将被忽略
func (op *activityOp) InsertTx(ext sqlx.Ext, m *Activity) (int64, error) {
	sql := "insert into activity(activity_name,activity_type,activity_begin,activity_end) values(?,?,?,?)"
	result, err := ext.Exec(sql,
		m.ActivityName,
		m.ActivityType,
		m.ActivityBegin,
		m.ActivityEnd,
	)
	if err != nil {
		log.Error("InsertTx sql error:%v, data:%v", err.Error(), m)
		return -1, err
	}
	affected, _ := result.LastInsertId()
	return affected, nil
}

/*
func (i *Activity) Update()  error {
    _,err := db.StatsDBMap.Update(i)
    if err != nil{
		log.Error("update sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *activityOp) Update(m *Activity) error {
	return op.UpdateTx(db.StatsDB, m)
}

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *activityOp) UpdateTx(ext sqlx.Ext, m *Activity) error {
	sql := `update activity set activity_type=?,activity_begin=?,activity_end=? where activity_name=?`
	_, err := ext.Exec(sql,
		m.ActivityType,
		m.ActivityBegin,
		m.ActivityEnd,
		m.ActivityName,
	)

	if err != nil {
		log.Error("update sql error:%v, data:%v", err.Error(), m)
		return err
	}

	return nil
}

// 用主键做条件，更新map里包含的字段名
func (op *activityOp) UpdateWithMap(activity_name string, m map[string]interface{}) error {
	return op.UpdateWithMapTx(db.StatsDB, activity_name, m)
}

// 用主键做条件，更新map里包含的字段名
func (op *activityOp) UpdateWithMapTx(ext sqlx.Ext, activity_name string, m map[string]interface{}) error {

	sql := `update activity set %s where 1=1 and activity_name=? ;`

	var params []interface{}
	var set_sql string
	for k, v := range m {
		if set_sql != "" {
			set_sql += ","
		}
		set_sql += fmt.Sprintf(" %s=? ", k)
		params = append(params, v)
	}
	params = append(params, activity_name)
	_, err := ext.Exec(fmt.Sprintf(sql, set_sql), params...)
	return err
}

/*
func (i *Activity) Delete() error{
    _,err := db.StatsDBMap.Delete(i)
	log.Error("Delete sql error:%v", err.Error())
    return err
}
*/
// 根据主键删除相关记录
func (op *activityOp) Delete(activity_name string) error {
	return op.DeleteTx(db.StatsDB, activity_name)
}

// 根据主键删除相关记录,Tx
func (op *activityOp) DeleteTx(ext sqlx.Ext, activity_name string) error {
	sql := `delete from activity where 1=1
        and activity_name=?
        `
	_, err := ext.Exec(sql,
		activity_name,
	)
	return err
}

// 返回符合查询条件的记录数
func (op *activityOp) CountByMap(m map[string]interface{}) (int64, error) {

	var params []interface{}
	sql := `select count(*) from activity where 1=1 `
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

func (op *activityOp) DeleteByMap(m map[string]interface{}) (int64, error) {
	return op.DeleteByMapTx(db.StatsDB, m)
}

func (op *activityOp) DeleteByMapTx(ext sqlx.Ext, m map[string]interface{}) (int64, error) {
	var params []interface{}
	sql := "delete from activity where 1=1 "
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
