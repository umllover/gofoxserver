package model

import (
	"errors"
	"fmt"
	"mj/hallServer/db"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lovelly/leaf/log"
)

//This file is generate by scripts,don't edit it

//gamescoreinfo
//

// +gen *
type Gamescoreinfo struct {
	UserID           int64      `db:"UserID" json:"UserID"`                     // 用户 ID
	Score            int64      `db:"Score" json:"Score"`                       // 用户积分（货币）
	Revenue          int64      `db:"Revenue" json:"Revenue"`                   // 游戏税收
	InsureScore      int64      `db:"InsureScore" json:"InsureScore"`           // 银行金币
	WinCount         int        `db:"WinCount" json:"WinCount"`                 // 胜局数目
	LostCount        int        `db:"LostCount" json:"LostCount"`               // 输局数目
	DrawCount        int        `db:"DrawCount" json:"DrawCount"`               // 和局数目
	FleeCount        int        `db:"FleeCount" json:"FleeCount"`               // 逃局数目
	AllLogonTimes    int        `db:"AllLogonTimes" json:"AllLogonTimes"`       // 总登陆次数
	PlayTimeCount    int        `db:"PlayTimeCount" json:"PlayTimeCount"`       // 游戏时间
	OnLineTimeCount  int        `db:"OnLineTimeCount" json:"OnLineTimeCount"`   // 在线时间
	LastLogonIP      string     `db:"LastLogonIP" json:"LastLogonIP"`           // 上次登陆 IP
	LastLogonDate    *time.Time `db:"LastLogonDate" json:"LastLogonDate"`       // 上次登陆时间
	LastLogonMachine string     `db:"LastLogonMachine" json:"LastLogonMachine"` // 登录机器
	RegisterIP       string     `db:"RegisterIP" json:"RegisterIP"`             // 注册 IP
	RegisterMachine  string     `db:"RegisterMachine" json:"RegisterMachine"`   // 注册机器
}

type gamescoreinfoOp struct{}

var GamescoreinfoOp = &gamescoreinfoOp{}
var DefaultGamescoreinfo = &Gamescoreinfo{}

// 按主键查询. 注:未找到记录的话将触发sql.ErrNoRows错误，返回nil, false
func (op *gamescoreinfoOp) Get(UserID int64) (*Gamescoreinfo, bool) {
	obj := &Gamescoreinfo{}
	sql := "select * from gamescoreinfo where UserID=? "
	err := db.DB.Get(obj, sql,
		UserID,
	)

	if err != nil {
		log.Error("Get data error:%v", err.Error())
		return nil, false
	}
	return obj, true
}
func (op *gamescoreinfoOp) SelectAll() ([]*Gamescoreinfo, error) {
	objList := []*Gamescoreinfo{}
	sql := "select * from gamescoreinfo "
	err := db.DB.Select(&objList, sql)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return objList, nil
}

func (op *gamescoreinfoOp) QueryByMap(m map[string]interface{}) ([]*Gamescoreinfo, error) {
	result := []*Gamescoreinfo{}
	var params []interface{}

	sql := "select * from gamescoreinfo where 1=1 "
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

func (op *gamescoreinfoOp) GetByMap(m map[string]interface{}) (*Gamescoreinfo, error) {
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
func (i *Gamescoreinfo) Insert() error {
    err := db.DBMap.Insert(i)
    if err != nil{
		log.Error("Insert sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 插入数据，自增长字段将被忽略
func (op *gamescoreinfoOp) Insert(m *Gamescoreinfo) (int64, error) {
	return op.InsertTx(db.DB, m)
}

// 插入数据，自增长字段将被忽略
func (op *gamescoreinfoOp) InsertTx(ext sqlx.Ext, m *Gamescoreinfo) (int64, error) {
	sql := "insert into gamescoreinfo(UserID,Score,Revenue,InsureScore,WinCount,LostCount,DrawCount,FleeCount,AllLogonTimes,PlayTimeCount,OnLineTimeCount,LastLogonIP,LastLogonDate,LastLogonMachine,RegisterIP,RegisterMachine) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	result, err := ext.Exec(sql,
		m.UserID,
		m.Score,
		m.Revenue,
		m.InsureScore,
		m.WinCount,
		m.LostCount,
		m.DrawCount,
		m.FleeCount,
		m.AllLogonTimes,
		m.PlayTimeCount,
		m.OnLineTimeCount,
		m.LastLogonIP,
		m.LastLogonDate,
		m.LastLogonMachine,
		m.RegisterIP,
		m.RegisterMachine,
	)
	if err != nil {
		log.Error("InsertTx sql error:%v, data:%v", err.Error(), m)
		return -1, err
	}
	affected, _ := result.LastInsertId()
	return affected, nil
}

//存在就更新， 不存在就插入
func (op *gamescoreinfoOp) InsertUpdate(obj *Gamescoreinfo, m map[string]interface{}) error {
	sql := "insert into gamescoreinfo(UserID,Score,Revenue,InsureScore,WinCount,LostCount,DrawCount,FleeCount,AllLogonTimes,PlayTimeCount,OnLineTimeCount,LastLogonIP,LastLogonDate,LastLogonMachine,RegisterIP,RegisterMachine) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE "
	var params = []interface{}{obj.UserID,
		obj.Score,
		obj.Revenue,
		obj.InsureScore,
		obj.WinCount,
		obj.LostCount,
		obj.DrawCount,
		obj.FleeCount,
		obj.AllLogonTimes,
		obj.PlayTimeCount,
		obj.OnLineTimeCount,
		obj.LastLogonIP,
		obj.LastLogonDate,
		obj.LastLogonMachine,
		obj.RegisterIP,
		obj.RegisterMachine,
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
func (i *Gamescoreinfo) Update()  error {
    _,err := db.DBMap.Update(i)
    if err != nil{
		log.Error("update sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *gamescoreinfoOp) Update(m *Gamescoreinfo) error {
	return op.UpdateTx(db.DB, m)
}

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *gamescoreinfoOp) UpdateTx(ext sqlx.Ext, m *Gamescoreinfo) error {
	sql := `update gamescoreinfo set Score=?,Revenue=?,InsureScore=?,WinCount=?,LostCount=?,DrawCount=?,FleeCount=?,AllLogonTimes=?,PlayTimeCount=?,OnLineTimeCount=?,LastLogonIP=?,LastLogonDate=?,LastLogonMachine=?,RegisterIP=?,RegisterMachine=? where UserID=?`
	_, err := ext.Exec(sql,
		m.Score,
		m.Revenue,
		m.InsureScore,
		m.WinCount,
		m.LostCount,
		m.DrawCount,
		m.FleeCount,
		m.AllLogonTimes,
		m.PlayTimeCount,
		m.OnLineTimeCount,
		m.LastLogonIP,
		m.LastLogonDate,
		m.LastLogonMachine,
		m.RegisterIP,
		m.RegisterMachine,
		m.UserID,
	)

	if err != nil {
		log.Error("update sql error:%v, data:%v", err.Error(), m)
		return err
	}

	return nil
}

// 用主键做条件，更新map里包含的字段名
func (op *gamescoreinfoOp) UpdateWithMap(UserID int64, m map[string]interface{}) error {
	return op.UpdateWithMapTx(db.DB, UserID, m)
}

// 用主键做条件，更新map里包含的字段名
func (op *gamescoreinfoOp) UpdateWithMapTx(ext sqlx.Ext, UserID int64, m map[string]interface{}) error {

	sql := `update gamescoreinfo set %s where 1=1 and UserID=? ;`

	var params []interface{}
	var set_sql string
	for k, v := range m {
		if set_sql != "" {
			set_sql += ","
		}
		set_sql += fmt.Sprintf(" %s=? ", k)
		params = append(params, v)
	}
	params = append(params, UserID)
	_, err := ext.Exec(fmt.Sprintf(sql, set_sql), params...)
	return err
}

/*
func (i *Gamescoreinfo) Delete() error{
    _,err := db.DBMap.Delete(i)
	log.Error("Delete sql error:%v", err.Error())
    return err
}
*/
// 根据主键删除相关记录
func (op *gamescoreinfoOp) Delete(UserID int64) error {
	return op.DeleteTx(db.DB, UserID)
}

// 根据主键删除相关记录,Tx
func (op *gamescoreinfoOp) DeleteTx(ext sqlx.Ext, UserID int64) error {
	sql := `delete from gamescoreinfo where 1=1
        and UserID=?
        `
	_, err := ext.Exec(sql,
		UserID,
	)
	return err
}

// 返回符合查询条件的记录数
func (op *gamescoreinfoOp) CountByMap(m map[string]interface{}) (int64, error) {

	var params []interface{}
	sql := `select count(*) from gamescoreinfo where 1=1 `
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

func (op *gamescoreinfoOp) DeleteByMap(m map[string]interface{}) (int64, error) {
	return op.DeleteByMapTx(db.DB, m)
}

func (op *gamescoreinfoOp) DeleteByMapTx(ext sqlx.Ext, m map[string]interface{}) (int64, error) {
	var params []interface{}
	sql := "delete from gamescoreinfo where 1=1 "
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
