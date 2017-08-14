package stats

import (
	"fmt"
	"mj/hallServer/db"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lovelly/leaf/log"
)

//This file is generate by scripts,don't edit it

//draw_award_log
//

// +gen *
type DrawAwardLog struct {
	Id          int        `db:"id" json:"id"`                   // 活动id。 和程序保持一致
	DrawId      int        `db:"draw_id" json:"draw_id"`         // 领取奖励的key
	Description string     `db:"description" json:"description"` // 活动描述
	DrawCount   int64      `db:"draw_count" json:"draw_count"`   // 活动可以领取的次数
	DrawType    int        `db:"draw_type" json:"draw_type"`     // 领取类型，1是永久，2是每日领取，3是每周领取
	Amount      int        `db:"amount" json:"amount"`           // 奖励数量
	ItemType    int        `db:"item_type" json:"item_type"`     // 领取的物品类型， 1是钻石，
	DrawTime    *time.Time `db:"draw_time" json:"draw_time"`     // 领取奖励的时间
}

type drawAwardLogOp struct{}

var DrawAwardLogOp = &drawAwardLogOp{}
var DefaultDrawAwardLog = &DrawAwardLog{}

// 按主键查询. 注:未找到记录的话将触发sql.ErrNoRows错误，返回nil, false
func (op *drawAwardLogOp) Get(id int) (*DrawAwardLog, bool) {
	obj := &DrawAwardLog{}
	sql := "select * from draw_award_log where id=? "
	err := db.StatsDB.Get(obj, sql,
		id,
	)

	if err != nil {
		log.Error("Get data error:%v", err.Error())
		return nil, false
	}
	return obj, true
}
func (op *drawAwardLogOp) SelectAll() ([]*DrawAwardLog, error) {
	objList := []*DrawAwardLog{}
	sql := "select * from draw_award_log "
	err := db.StatsDB.Select(&objList, sql)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return objList, nil
}

func (op *drawAwardLogOp) QueryByMap(m map[string]interface{}) ([]*DrawAwardLog, error) {
	result := []*DrawAwardLog{}
	var params []interface{}

	sql := "select * from draw_award_log where 1=1 "
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

func (op *drawAwardLogOp) GetByMap(m map[string]interface{}) (*DrawAwardLog, error) {
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
func (i *DrawAwardLog) Insert() error {
    err := db.StatsDBMap.Insert(i)
    if err != nil{
		log.Error("Insert sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 插入数据，自增长字段将被忽略
func (op *drawAwardLogOp) Insert(m *DrawAwardLog) (int64, error) {
	return op.InsertTx(db.StatsDB, m)
}

// 插入数据，自增长字段将被忽略
func (op *drawAwardLogOp) InsertTx(ext sqlx.Ext, m *DrawAwardLog) (int64, error) {
	sql := "insert into draw_award_log(id,draw_id,description,draw_count,draw_type,amount,item_type,draw_time) values(?,?,?,?,?,?,?,?)"
	result, err := ext.Exec(sql,
		m.Id,
		m.DrawId,
		m.Description,
		m.DrawCount,
		m.DrawType,
		m.Amount,
		m.ItemType,
		m.DrawTime,
	)
	if err != nil {
		log.Error("InsertTx sql error:%v, data:%v", err.Error(), m)
		return -1, err
	}
	affected, _ := result.LastInsertId()
	return affected, nil
}

//存在就更新， 不存在就插入
func (op *drawAwardLogOp) InsertUpdate(obj *DrawAwardLog, m map[string]interface{}) error {
	sql := "insert into draw_award_log(id,draw_id,description,draw_count,draw_type,amount,item_type,draw_time) values(?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE "
	var params = []interface{}{obj.Id,
		obj.DrawId,
		obj.Description,
		obj.DrawCount,
		obj.DrawType,
		obj.Amount,
		obj.ItemType,
		obj.DrawTime,
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
func (i *DrawAwardLog) Update()  error {
    _,err := db.StatsDBMap.Update(i)
    if err != nil{
		log.Error("update sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *drawAwardLogOp) Update(m *DrawAwardLog) error {
	return op.UpdateTx(db.StatsDB, m)
}

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *drawAwardLogOp) UpdateTx(ext sqlx.Ext, m *DrawAwardLog) error {
	sql := `update draw_award_log set draw_id=?,description=?,draw_count=?,draw_type=?,amount=?,item_type=?,draw_time=? where id=?`
	_, err := ext.Exec(sql,
		m.DrawId,
		m.Description,
		m.DrawCount,
		m.DrawType,
		m.Amount,
		m.ItemType,
		m.DrawTime,
		m.Id,
	)

	if err != nil {
		log.Error("update sql error:%v, data:%v", err.Error(), m)
		return err
	}

	return nil
}

// 用主键做条件，更新map里包含的字段名
func (op *drawAwardLogOp) UpdateWithMap(id int, m map[string]interface{}) error {
	return op.UpdateWithMapTx(db.StatsDB, id, m)
}

// 用主键做条件，更新map里包含的字段名
func (op *drawAwardLogOp) UpdateWithMapTx(ext sqlx.Ext, id int, m map[string]interface{}) error {

	sql := `update draw_award_log set %s where 1=1 and id=? ;`

	var params []interface{}
	var set_sql string
	for k, v := range m {
		if set_sql != "" {
			set_sql += ","
		}
		set_sql += fmt.Sprintf(" %s=? ", k)
		params = append(params, v)
	}
	params = append(params, id)
	_, err := ext.Exec(fmt.Sprintf(sql, set_sql), params...)
	return err
}

/*
func (i *DrawAwardLog) Delete() error{
    _,err := db.StatsDBMap.Delete(i)
	log.Error("Delete sql error:%v", err.Error())
    return err
}
*/
// 根据主键删除相关记录
func (op *drawAwardLogOp) Delete(id int) error {
	return op.DeleteTx(db.StatsDB, id)
}

// 根据主键删除相关记录,Tx
func (op *drawAwardLogOp) DeleteTx(ext sqlx.Ext, id int) error {
	sql := `delete from draw_award_log where 1=1
        and id=?
        `
	_, err := ext.Exec(sql,
		id,
	)
	return err
}

// 返回符合查询条件的记录数
func (op *drawAwardLogOp) CountByMap(m map[string]interface{}) (int64, error) {

	var params []interface{}
	sql := `select count(*) from draw_award_log where 1=1 `
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

func (op *drawAwardLogOp) DeleteByMap(m map[string]interface{}) (int64, error) {
	return op.DeleteByMapTx(db.StatsDB, m)
}

func (op *drawAwardLogOp) DeleteByMapTx(ext sqlx.Ext, m map[string]interface{}) (int64, error) {
	var params []interface{}
	sql := "delete from draw_award_log where 1=1 "
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
