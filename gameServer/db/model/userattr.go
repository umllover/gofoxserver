package model

import (
	"errors"
	"fmt"
	"mj/gameServer/db"

	"github.com/jmoiron/sqlx"
	"github.com/lovelly/leaf/log"
)

//This file is generate by scripts,don't edit it

//userattr
//

// +gen *
type Userattr struct {
	UserID          int    `db:"UserID" json:"UserID"`                   //
	UnderWrite      string `db:"UnderWrite" json:"UnderWrite"`           // 个性签名
	FaceID          int8   `db:"FaceID" json:"FaceID"`                   // 头像标识
	CustomID        int    `db:"CustomID" json:"CustomID"`               // 自定标识
	UserMedal       int    `db:"UserMedal" json:"UserMedal"`             // 用户奖牌
	Experience      int    `db:"Experience" json:"Experience"`           // 经验数值
	LoveLiness      int    `db:"LoveLiness" json:"LoveLiness"`           // 用户魅力
	UserRight       int    `db:"UserRight" json:"UserRight"`             // 用户权限
	MasterRight     int    `db:"MasterRight" json:"MasterRight"`         // 管理权限
	MasterOrder     int8   `db:"MasterOrder" json:"MasterOrder"`         // 管理等级
	PlayTimeCount   int    `db:"PlayTimeCount" json:"PlayTimeCount"`     // 游戏时间
	OnLineTimeCount int    `db:"OnLineTimeCount" json:"OnLineTimeCount"` // 在线时间
	HeadImgUrl      string `db:"HeadImgUrl" json:"HeadImgUrl"`           // 头像
	Gender          int8   `db:"Gender" json:"Gender"`                   // 性别
	NickName        string `db:"NickName" json:"NickName"`               //
}

type userattrOp struct{}

var UserattrOp = &userattrOp{}
var DefaultUserattr = &Userattr{}

// 按主键查询. 注:未找到记录的话将触发sql.ErrNoRows错误，返回nil, false
func (op *userattrOp) Get(UserID int) (*Userattr, bool) {
	obj := &Userattr{}
	sql := "select * from userattr where UserID=? "
	err := db.DB.Get(obj, sql,
		UserID,
	)

	if err != nil {
		log.Error("Get data error:%v", err.Error())
		return nil, false
	}
	return obj, true
}
func (op *userattrOp) SelectAll() ([]*Userattr, error) {
	objList := []*Userattr{}
	sql := "select * from userattr "
	err := db.DB.Select(&objList, sql)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return objList, nil
}

func (op *userattrOp) QueryByMap(m map[string]interface{}) ([]*Userattr, error) {
	result := []*Userattr{}
	var params []interface{}

	sql := "select * from userattr where 1=1 "
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

func (op *userattrOp) QueryByMapQueryByMapComparison(m map[string]interface{}) ([]*Userattr, error) {
	result := []*Userattr{}
	var params []interface{}

	sql := "select * from userattr where 1=1 "
	for k, v := range m {
		sql += fmt.Sprintf(" and %s? ", k)
		params = append(params, v)
	}
	err := db.DB.Select(&result, sql, params...)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return result, nil
}

func (op *userattrOp) GetByMap(m map[string]interface{}) (*Userattr, error) {
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
func (i *Userattr) Insert() error {
    err := db.DBMap.Insert(i)
    if err != nil{
		log.Error("Insert sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 插入数据，自增长字段将被忽略
func (op *userattrOp) Insert(m *Userattr) (int64, error) {
	return op.InsertTx(db.DB, m)
}

// 插入数据，自增长字段将被忽略
func (op *userattrOp) InsertTx(ext sqlx.Ext, m *Userattr) (int64, error) {
	sql := "insert into userattr(UserID,UnderWrite,FaceID,CustomID,UserMedal,Experience,LoveLiness,UserRight,MasterRight,MasterOrder,PlayTimeCount,OnLineTimeCount,HeadImgUrl,Gender,NickName) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	result, err := ext.Exec(sql,
		m.UserID,
		m.UnderWrite,
		m.FaceID,
		m.CustomID,
		m.UserMedal,
		m.Experience,
		m.LoveLiness,
		m.UserRight,
		m.MasterRight,
		m.MasterOrder,
		m.PlayTimeCount,
		m.OnLineTimeCount,
		m.HeadImgUrl,
		m.Gender,
		m.NickName,
	)
	if err != nil {
		log.Error("InsertTx sql error:%v, data:%v", err.Error(), m)
		return -1, err
	}
	affected, _ := result.LastInsertId()
	return affected, nil
}

/*
func (i *Userattr) Update()  error {
    _,err := db.DBMap.Update(i)
    if err != nil{
		log.Error("update sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *userattrOp) Update(m *Userattr) error {
	return op.UpdateTx(db.DB, m)
}

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *userattrOp) UpdateTx(ext sqlx.Ext, m *Userattr) error {
	sql := `update userattr set UnderWrite=?,FaceID=?,CustomID=?,UserMedal=?,Experience=?,LoveLiness=?,UserRight=?,MasterRight=?,MasterOrder=?,PlayTimeCount=?,OnLineTimeCount=?,HeadImgUrl=?,Gender=?,NickName=? where UserID=?`
	_, err := ext.Exec(sql,
		m.UnderWrite,
		m.FaceID,
		m.CustomID,
		m.UserMedal,
		m.Experience,
		m.LoveLiness,
		m.UserRight,
		m.MasterRight,
		m.MasterOrder,
		m.PlayTimeCount,
		m.OnLineTimeCount,
		m.HeadImgUrl,
		m.Gender,
		m.NickName,
		m.UserID,
	)

	if err != nil {
		log.Error("update sql error:%v, data:%v", err.Error(), m)
		return err
	}

	return nil
}

// 用主键做条件，更新map里包含的字段名
func (op *userattrOp) UpdateWithMap(UserID int, m map[string]interface{}) error {
	return op.UpdateWithMapTx(db.DB, UserID, m)
}

// 用主键做条件，更新map里包含的字段名
func (op *userattrOp) UpdateWithMapTx(ext sqlx.Ext, UserID int, m map[string]interface{}) error {

	sql := `update userattr set %s where 1=1 and UserID=? ;`

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
func (i *Userattr) Delete() error{
    _,err := db.DBMap.Delete(i)
	log.Error("Delete sql error:%v", err.Error())
    return err
}
*/
// 根据主键删除相关记录
func (op *userattrOp) Delete(UserID int) error {
	return op.DeleteTx(db.DB, UserID)
}

// 根据主键删除相关记录,Tx
func (op *userattrOp) DeleteTx(ext sqlx.Ext, UserID int) error {
	sql := `delete from userattr where 1=1
        and UserID=?
        `
	_, err := ext.Exec(sql,
		UserID,
	)
	return err
}

// 返回符合查询条件的记录数
func (op *userattrOp) CountByMap(m map[string]interface{}) (int64, error) {

	var params []interface{}
	sql := `select count(*) from userattr where 1=1 `
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

func (op *userattrOp) DeleteByMap(m map[string]interface{}) (int64, error) {
	return op.DeleteByMapTx(db.DB, m)
}

func (op *userattrOp) DeleteByMapTx(ext sqlx.Ext, m map[string]interface{}) (int64, error) {
	var params []interface{}
	sql := "delete from userattr where 1=1 "
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
