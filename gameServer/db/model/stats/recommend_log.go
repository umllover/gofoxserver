package stats

import (
	"errors"
	"fmt"
	"mj/gameServer/db"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lovelly/leaf/log"
)

//This file is generate by scripts,don't edit it

//recommend_log
//

// +gen *
type RecommendLog struct {
	SubElectUid int64      `db:"sub_elect_uid" json:"sub_elect_uid"` // 被推举人人id
	ElectUid    int64      `db:"elect_uid" json:"elect_uid"`         // 推举人id
	ElectTime   *time.Time `db:"elect_time" json:"elect_time"`       // 领取时间
}

type recommendLogOp struct{}

var RecommendLogOp = &recommendLogOp{}
var DefaultRecommendLog = &RecommendLog{}

// 按主键查询. 注:未找到记录的话将触发sql.ErrNoRows错误，返回nil, false
func (op *recommendLogOp) Get(sub_elect_uid int64) (*RecommendLog, bool) {
	obj := &RecommendLog{}
	sql := "select * from recommend_log where sub_elect_uid=? "
	err := db.StatsDB.Get(obj, sql,
		sub_elect_uid,
	)

	if err != nil {
		log.Error("Get data error:%v", err.Error())
		return nil, false
	}
	return obj, true
}
func (op *recommendLogOp) SelectAll() ([]*RecommendLog, error) {
	objList := []*RecommendLog{}
	sql := "select * from recommend_log "
	err := db.StatsDB.Select(&objList, sql)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return objList, nil
}

func (op *recommendLogOp) QueryByMap(m map[string]interface{}) ([]*RecommendLog, error) {
	result := []*RecommendLog{}
	var params []interface{}

	sql := "select * from recommend_log where 1=1 "
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

func (op *recommendLogOp) GetByMap(m map[string]interface{}) (*RecommendLog, error) {
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
func (i *RecommendLog) Insert() error {
    err := db.StatsDBMap.Insert(i)
    if err != nil{
		log.Error("Insert sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 插入数据，自增长字段将被忽略
func (op *recommendLogOp) Insert(m *RecommendLog) (int64, error) {
	return op.InsertTx(db.StatsDB, m)
}

// 插入数据，自增长字段将被忽略
func (op *recommendLogOp) InsertTx(ext sqlx.Ext, m *RecommendLog) (int64, error) {
	sql := "insert into recommend_log(sub_elect_uid,elect_uid,elect_time) values(?,?,?)"
	result, err := ext.Exec(sql,
		m.SubElectUid,
		m.ElectUid,
		m.ElectTime,
	)
	if err != nil {
		log.Error("InsertTx sql error:%v, data:%v", err.Error(), m)
		return -1, err
	}
	affected, _ := result.LastInsertId()
	return affected, nil
}

//存在就更新， 不存在就插入
func (op *recommendLogOp) InsertUpdate(obj *RecommendLog, m map[string]interface{}) error {
	sql := "insert into recommend_log(sub_elect_uid,elect_uid,elect_time) values(?,?,?) ON DUPLICATE KEY UPDATE "
	var params = []interface{}{obj.SubElectUid,
		obj.ElectUid,
		obj.ElectTime,
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
func (i *RecommendLog) Update()  error {
    _,err := db.StatsDBMap.Update(i)
    if err != nil{
		log.Error("update sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *recommendLogOp) Update(m *RecommendLog) error {
	return op.UpdateTx(db.StatsDB, m)
}

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *recommendLogOp) UpdateTx(ext sqlx.Ext, m *RecommendLog) error {
	sql := `update recommend_log set elect_uid=?,elect_time=? where sub_elect_uid=?`
	_, err := ext.Exec(sql,
		m.ElectUid,
		m.ElectTime,
		m.SubElectUid,
	)

	if err != nil {
		log.Error("update sql error:%v, data:%v", err.Error(), m)
		return err
	}

	return nil
}

// 用主键做条件，更新map里包含的字段名
func (op *recommendLogOp) UpdateWithMap(sub_elect_uid int64, m map[string]interface{}) error {
	return op.UpdateWithMapTx(db.StatsDB, sub_elect_uid, m)
}

// 用主键做条件，更新map里包含的字段名
func (op *recommendLogOp) UpdateWithMapTx(ext sqlx.Ext, sub_elect_uid int64, m map[string]interface{}) error {

	sql := `update recommend_log set %s where 1=1 and sub_elect_uid=? ;`

	var params []interface{}
	var set_sql string
	for k, v := range m {
		if set_sql != "" {
			set_sql += ","
		}
		set_sql += fmt.Sprintf(" %s=? ", k)
		params = append(params, v)
	}
	params = append(params, sub_elect_uid)
	_, err := ext.Exec(fmt.Sprintf(sql, set_sql), params...)
	return err
}

/*
func (i *RecommendLog) Delete() error{
    _,err := db.StatsDBMap.Delete(i)
	log.Error("Delete sql error:%v", err.Error())
    return err
}
*/
// 根据主键删除相关记录
func (op *recommendLogOp) Delete(sub_elect_uid int64) error {
	return op.DeleteTx(db.StatsDB, sub_elect_uid)
}

// 根据主键删除相关记录,Tx
func (op *recommendLogOp) DeleteTx(ext sqlx.Ext, sub_elect_uid int64) error {
	sql := `delete from recommend_log where 1=1
        and sub_elect_uid=?
        `
	_, err := ext.Exec(sql,
		sub_elect_uid,
	)
	return err
}

// 返回符合查询条件的记录数
func (op *recommendLogOp) CountByMap(m map[string]interface{}) (int64, error) {

	var params []interface{}
	sql := `select count(*) from recommend_log where 1=1 `
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

func (op *recommendLogOp) DeleteByMap(m map[string]interface{}) (int64, error) {
	return op.DeleteByMapTx(db.StatsDB, m)
}

func (op *recommendLogOp) DeleteByMapTx(ext sqlx.Ext, m map[string]interface{}) (int64, error) {
	var params []interface{}
	sql := "delete from recommend_log where 1=1 "
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
