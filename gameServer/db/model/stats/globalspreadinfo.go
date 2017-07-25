package stats

import (
	"errors"
	"fmt"
	"mj/gameServer/db"

	"github.com/jmoiron/sqlx"
	"github.com/lovelly/leaf/log"
)

//This file is generate by scripts,don't edit it

//globalspreadinfo
//

// +gen *
type Globalspreadinfo struct {
	ID                 int     `db:"ID" json:"ID"`                                 //
	RegisterGrantScore int     `db:"RegisterGrantScore" json:"RegisterGrantScore"` // 注册时赠送金币数目
	PlayTimeCount      int     `db:"PlayTimeCount" json:"PlayTimeCount"`           // 在线时长（单位：秒）
	PlayTimeGrantScore int     `db:"PlayTimeGrantScore" json:"PlayTimeGrantScore"` // 根据在线时长赠送金币数目
	FillGrantRate      float64 `db:"FillGrantRate" json:"FillGrantRate"`           // 充值赠送比率
	BalanceRate        float64 `db:"BalanceRate" json:"BalanceRate"`               // 结算赠送比率
	MinBalanceScore    int     `db:"MinBalanceScore" json:"MinBalanceScore"`       // 结算最小值
}

type globalspreadinfoOp struct{}

var GlobalspreadinfoOp = &globalspreadinfoOp{}
var DefaultGlobalspreadinfo = &Globalspreadinfo{}

// 按主键查询. 注:未找到记录的话将触发sql.ErrNoRows错误，返回nil, false
func (op *globalspreadinfoOp) Get(ID int) (*Globalspreadinfo, bool) {
	obj := &Globalspreadinfo{}
	sql := "select * from globalspreadinfo where ID=? "
	err := db.StatsDB.Get(obj, sql,
		ID,
	)

	if err != nil {
		log.Error("Get data error:%v", err.Error())
		return nil, false
	}
	return obj, true
}
func (op *globalspreadinfoOp) SelectAll() ([]*Globalspreadinfo, error) {
	objList := []*Globalspreadinfo{}
	sql := "select * from globalspreadinfo "
	err := db.StatsDB.Select(&objList, sql)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return objList, nil
}

func (op *globalspreadinfoOp) QueryByMap(m map[string]interface{}) ([]*Globalspreadinfo, error) {
	result := []*Globalspreadinfo{}
	var params []interface{}

	sql := "select * from globalspreadinfo where 1=1 "
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

func (op *globalspreadinfoOp) GetByMap(m map[string]interface{}) (*Globalspreadinfo, error) {
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
func (i *Globalspreadinfo) Insert() error {
    err := db.StatsDBMap.Insert(i)
    if err != nil{
		log.Error("Insert sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 插入数据，自增长字段将被忽略
func (op *globalspreadinfoOp) Insert(m *Globalspreadinfo) (int64, error) {
	return op.InsertTx(db.StatsDB, m)
}

// 插入数据，自增长字段将被忽略
func (op *globalspreadinfoOp) InsertTx(ext sqlx.Ext, m *Globalspreadinfo) (int64, error) {
	sql := "insert into globalspreadinfo(ID,RegisterGrantScore,PlayTimeCount,PlayTimeGrantScore,FillGrantRate,BalanceRate,MinBalanceScore) values(?,?,?,?,?,?,?)"
	result, err := ext.Exec(sql,
		m.ID,
		m.RegisterGrantScore,
		m.PlayTimeCount,
		m.PlayTimeGrantScore,
		m.FillGrantRate,
		m.BalanceRate,
		m.MinBalanceScore,
	)
	if err != nil {
		log.Error("InsertTx sql error:%v, data:%v", err.Error(), m)
		return -1, err
	}
	affected, _ := result.LastInsertId()
	return affected, nil
}

//存在就更新， 不存在就插入
func (op *globalspreadinfoOp) InsertUpdate(obj *Globalspreadinfo, m map[string]interface{}) ( error) {
    sql := "insert into globalspreadinfo(ID,RegisterGrantScore,PlayTimeCount,PlayTimeGrantScore,FillGrantRate,BalanceRate,MinBalanceScore) values(?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE "
    var params = []interface{}{ obj.ID,
        obj.RegisterGrantScore,
        obj.PlayTimeCount,
        obj.PlayTimeGrantScore,
        obj.FillGrantRate,
        obj.BalanceRate,
        obj.MinBalanceScore,
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
func (i *Globalspreadinfo) Update()  error {
    _,err := db.StatsDBMap.Update(i)
    if err != nil{
		log.Error("update sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *globalspreadinfoOp) Update(m *Globalspreadinfo) error {
	return op.UpdateTx(db.StatsDB, m)
}

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *globalspreadinfoOp) UpdateTx(ext sqlx.Ext, m *Globalspreadinfo) error {
	sql := `update globalspreadinfo set RegisterGrantScore=?,PlayTimeCount=?,PlayTimeGrantScore=?,FillGrantRate=?,BalanceRate=?,MinBalanceScore=? where ID=?`
	_, err := ext.Exec(sql,
		m.RegisterGrantScore,
		m.PlayTimeCount,
		m.PlayTimeGrantScore,
		m.FillGrantRate,
		m.BalanceRate,
		m.MinBalanceScore,
		m.ID,
	)

	if err != nil {
		log.Error("update sql error:%v, data:%v", err.Error(), m)
		return err
	}

	return nil
}

// 用主键做条件，更新map里包含的字段名
func (op *globalspreadinfoOp) UpdateWithMap(ID int, m map[string]interface{}) error {
	return op.UpdateWithMapTx(db.StatsDB, ID, m)
}

// 用主键做条件，更新map里包含的字段名
func (op *globalspreadinfoOp) UpdateWithMapTx(ext sqlx.Ext, ID int, m map[string]interface{}) error {

	sql := `update globalspreadinfo set %s where 1=1 and ID=? ;`

	var params []interface{}
	var set_sql string
	for k, v := range m {
		if set_sql != "" {
			set_sql += ","
		}
		set_sql += fmt.Sprintf(" %s=? ", k)
		params = append(params, v)
	}
	params = append(params, ID)
	_, err := ext.Exec(fmt.Sprintf(sql, set_sql), params...)
	return err
}

/*
func (i *Globalspreadinfo) Delete() error{
    _,err := db.StatsDBMap.Delete(i)
	log.Error("Delete sql error:%v", err.Error())
    return err
}
*/
// 根据主键删除相关记录
func (op *globalspreadinfoOp) Delete(ID int) error {
	return op.DeleteTx(db.StatsDB, ID)
}

// 根据主键删除相关记录,Tx
func (op *globalspreadinfoOp) DeleteTx(ext sqlx.Ext, ID int) error {
	sql := `delete from globalspreadinfo where 1=1
        and ID=?
        `
	_, err := ext.Exec(sql,
		ID,
	)
	return err
}

// 返回符合查询条件的记录数
func (op *globalspreadinfoOp) CountByMap(m map[string]interface{}) (int64, error) {

	var params []interface{}
	sql := `select count(*) from globalspreadinfo where 1=1 `
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

func (op *globalspreadinfoOp) DeleteByMap(m map[string]interface{}) (int64, error) {
	return op.DeleteByMapTx(db.StatsDB, m)
}

func (op *globalspreadinfoOp) DeleteByMapTx(ext sqlx.Ext, m map[string]interface{}) (int64, error) {
	var params []interface{}
	sql := "delete from globalspreadinfo where 1=1 "
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
