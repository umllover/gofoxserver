package stats

import (
	"errors"
	"fmt"
	"mj/hallServer/db"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lovelly/leaf/log"
)

//This file is generate by scripts,don't edit it

//consum_log
//

// +gen *
type ConsumLog struct {
	RecodeId   int        `db:"recode_id" json:"recode_id"`     //
	UserId     int64      `db:"user_id" json:"user_id"`         // 用户索引
	ConsumType int        `db:"consum_type" json:"consum_type"` // 消费类型 0钻石 1开房 3道具
	ConsumNum  int        `db:"consum_num" json:"consum_num"`   // 消费数量
	ConsumTime *time.Time `db:"consum_time" json:"consum_time"` // 消费时间
}

type consumLogOp struct{}

var ConsumLogOp = &consumLogOp{}
var DefaultConsumLog = &ConsumLog{}

// 按主键查询. 注:未找到记录的话将触发sql.ErrNoRows错误，返回nil, false
func (op *consumLogOp) Get(recode_id int) (*ConsumLog, bool) {
	obj := &ConsumLog{}
	sql := "select * from consum_log where recode_id=? "
	err := db.StatsDB.Get(obj, sql,
		recode_id,
	)

	if err != nil {
		log.Error("Get data error:%v", err.Error())
		return nil, false
	}
	return obj, true
}
func (op *consumLogOp) SelectAll() ([]*ConsumLog, error) {
	objList := []*ConsumLog{}
	sql := "select * from consum_log "
	err := db.StatsDB.Select(&objList, sql)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return objList, nil
}

func (op *consumLogOp) QueryByMap(m map[string]interface{}) ([]*ConsumLog, error) {
	result := []*ConsumLog{}
	var params []interface{}

	sql := "select * from consum_log where 1=1 "
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

func (op *consumLogOp) GetByMap(m map[string]interface{}) (*ConsumLog, error) {
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
func (i *ConsumLog) Insert() error {
    err := db.StatsDBMap.Insert(i)
    if err != nil{
		log.Error("Insert sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 插入数据，自增长字段将被忽略
func (op *consumLogOp) Insert(m *ConsumLog) (int64, error) {
	return op.InsertTx(db.StatsDB, m)
}

// 插入数据，自增长字段将被忽略
func (op *consumLogOp) InsertTx(ext sqlx.Ext, m *ConsumLog) (int64, error) {
	sql := "insert into consum_log(recode_id,user_id,consum_type,consum_num,consum_time) values(?,?,?,?,?)"
	result, err := ext.Exec(sql,
		m.RecodeId,
		m.UserId,
		m.ConsumType,
		m.ConsumNum,
		m.ConsumTime,
	)
	if err != nil {
		log.Error("InsertTx sql error:%v, data:%v", err.Error(), m)
		return -1, err
	}
	affected, _ := result.LastInsertId()
	return affected, nil
}

//存在就更新， 不存在就插入
func (op *consumLogOp) InsertUpdate(obj *ConsumLog, m map[string]interface{}) error {
	sql := "insert into consum_log(recode_id,user_id,consum_type,consum_num,consum_time) values(?,?,?,?,?) ON DUPLICATE KEY UPDATE "
	var params = []interface{}{obj.RecodeId,
		obj.UserId,
		obj.ConsumType,
		obj.ConsumNum,
		obj.ConsumTime,
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
func (i *ConsumLog) Update()  error {
    _,err := db.StatsDBMap.Update(i)
    if err != nil{
		log.Error("update sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *consumLogOp) Update(m *ConsumLog) error {
	return op.UpdateTx(db.StatsDB, m)
}

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *consumLogOp) UpdateTx(ext sqlx.Ext, m *ConsumLog) error {
	sql := `update consum_log set user_id=?,consum_type=?,consum_num=?,consum_time=? where recode_id=?`
	_, err := ext.Exec(sql,
		m.UserId,
		m.ConsumType,
		m.ConsumNum,
		m.ConsumTime,
		m.RecodeId,
	)

	if err != nil {
		log.Error("update sql error:%v, data:%v", err.Error(), m)
		return err
	}

	return nil
}

// 用主键做条件，更新map里包含的字段名
func (op *consumLogOp) UpdateWithMap(recode_id int, m map[string]interface{}) error {
	return op.UpdateWithMapTx(db.StatsDB, recode_id, m)
}

// 用主键做条件，更新map里包含的字段名
func (op *consumLogOp) UpdateWithMapTx(ext sqlx.Ext, recode_id int, m map[string]interface{}) error {

	sql := `update consum_log set %s where 1=1 and recode_id=? ;`

	var params []interface{}
	var set_sql string
	for k, v := range m {
		if set_sql != "" {
			set_sql += ","
		}
		set_sql += fmt.Sprintf(" %s=? ", k)
		params = append(params, v)
	}
	params = append(params, recode_id)
	_, err := ext.Exec(fmt.Sprintf(sql, set_sql), params...)
	return err
}

/*
func (i *ConsumLog) Delete() error{
    _,err := db.StatsDBMap.Delete(i)
	log.Error("Delete sql error:%v", err.Error())
    return err
}
*/
// 根据主键删除相关记录
func (op *consumLogOp) Delete(recode_id int) error {
	return op.DeleteTx(db.StatsDB, recode_id)
}

// 根据主键删除相关记录,Tx
func (op *consumLogOp) DeleteTx(ext sqlx.Ext, recode_id int) error {
	sql := `delete from consum_log where 1=1
        and recode_id=?
        `
	_, err := ext.Exec(sql,
		recode_id,
	)
	return err
}

// 返回符合查询条件的记录数
func (op *consumLogOp) CountByMap(m map[string]interface{}) (int64, error) {

	var params []interface{}
	sql := `select count(*) from consum_log where 1=1 `
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

func (op *consumLogOp) DeleteByMap(m map[string]interface{}) (int64, error) {
	return op.DeleteByMapTx(db.StatsDB, m)
}

func (op *consumLogOp) DeleteByMapTx(ext sqlx.Ext, m map[string]interface{}) (int64, error) {
	var params []interface{}
	sql := "delete from consum_log where 1=1 "
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
