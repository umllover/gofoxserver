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

//onlineorder
//

// +gen *
type Onlineorder struct {
	OnLineID      int        `db:"OnLineID" json:"OnLineID"`           // 订单标识
	OperUserID    string     `db:"OperUserID" json:"OperUserID"`       // 操作用户
	ShareID       int        `db:"ShareID" json:"ShareID"`             // 服务标识
	UserID        int        `db:"UserID" json:"UserID"`               // 用户标识
	GameID        int        `db:"GameID" json:"GameID"`               // 游戏ID
	Accounts      string     `db:"Accounts" json:"Accounts"`           // 用户名
	OrderID       string     `db:"OrderID" json:"OrderID"`             // 订单号码
	CardTypeID    int        `db:"CardTypeID" json:"CardTypeID"`       // 卡类标识
	CardPrice     float64    `db:"CardPrice" json:"CardPrice"`         // 会员卡价格
	CardGold      int64      `db:"CardGold" json:"CardGold"`           // 卡金币
	CardTotal     int        `db:"CardTotal" json:"CardTotal"`         // 充卡数量
	OrderAmount   float64    `db:"OrderAmount" json:"OrderAmount"`     // 订单金额
	DiscountScale float64    `db:"DiscountScale" json:"DiscountScale"` // 折扣比例
	PayAmount     int        `db:"PayAmount" json:"PayAmount"`         // 实付金额
	OrderStatus   int8       `db:"OrderStatus" json:"OrderStatus"`     // 订单状态  0:未付款;1:已付款待处理;2:处理完成
	IPAddress     string     `db:"IPAddress" json:"IPAddress"`         // 订单地址
	ApplyDate     *time.Time `db:"ApplyDate" json:"ApplyDate"`         // 订单日期
	PhoneNum      string     `db:"PhoneNum" json:"PhoneNum"`           //
	GameName      string     `db:"GameName" json:"GameName"`           //
	NickName      string     `db:"NickName" json:"NickName"`           //
	GoodsNumber   int        `db:"GoodsNumber" json:"GoodsNumber"`     //
	GoodsID       int        `db:"GoodsID" json:"GoodsID"`             //
	GoodsName     string     `db:"GoodsName" json:"GoodsName"`         //
	OrderDate     *time.Time `db:"OrderDate" json:"OrderDate"`         //
	PayType       string     `db:"PayType" json:"PayType"`             //
}

type onlineorderOp struct{}

var OnlineorderOp = &onlineorderOp{}
var DefaultOnlineorder = &Onlineorder{}

// 按主键查询. 注:未找到记录的话将触发sql.ErrNoRows错误，返回nil, false
func (op *onlineorderOp) Get(OnLineID int) (*Onlineorder, bool) {
	obj := &Onlineorder{}
	sql := "select * from onlineorder where OnLineID=? "
	err := db.DB.Get(obj, sql,
		OnLineID,
	)

	if err != nil {
		log.Error("Get data error:%v", err.Error())
		return nil, false
	}
	return obj, true
}
func (op *onlineorderOp) SelectAll() ([]*Onlineorder, error) {
	objList := []*Onlineorder{}
	sql := "select * from onlineorder "
	err := db.DB.Select(&objList, sql)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return objList, nil
}

func (op *onlineorderOp) QueryByMap(m map[string]interface{}) ([]*Onlineorder, error) {
	result := []*Onlineorder{}
	var params []interface{}

	sql := "select * from onlineorder where 1=1 "
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

func (op *onlineorderOp) QueryByMapQueryByMapComparison(m map[string]interface{}) ([]*Onlineorder, error) {
	result := []*Onlineorder{}
	var params []interface{}

	sql := "select * from onlineorder where 1=1 "
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

func (op *onlineorderOp) GetByMap(m map[string]interface{}) (*Onlineorder, error) {
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
func (i *Onlineorder) Insert() error {
    err := db.DBMap.Insert(i)
    if err != nil{
		log.Error("Insert sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 插入数据，自增长字段将被忽略
func (op *onlineorderOp) Insert(m *Onlineorder) (int64, error) {
	return op.InsertTx(db.DB, m)
}

// 插入数据，自增长字段将被忽略
func (op *onlineorderOp) InsertTx(ext sqlx.Ext, m *Onlineorder) (int64, error) {
	sql := "insert into onlineorder(OnLineID,OperUserID,ShareID,UserID,GameID,Accounts,OrderID,CardTypeID,CardPrice,CardGold,CardTotal,OrderAmount,DiscountScale,PayAmount,OrderStatus,IPAddress,ApplyDate,PhoneNum,GameName,NickName,GoodsNumber,GoodsID,GoodsName,OrderDate,PayType) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	result, err := ext.Exec(sql,
		m.OnLineID,
		m.OperUserID,
		m.ShareID,
		m.UserID,
		m.GameID,
		m.Accounts,
		m.OrderID,
		m.CardTypeID,
		m.CardPrice,
		m.CardGold,
		m.CardTotal,
		m.OrderAmount,
		m.DiscountScale,
		m.PayAmount,
		m.OrderStatus,
		m.IPAddress,
		m.ApplyDate,
		m.PhoneNum,
		m.GameName,
		m.NickName,
		m.GoodsNumber,
		m.GoodsID,
		m.GoodsName,
		m.OrderDate,
		m.PayType,
	)
	if err != nil {
		log.Error("InsertTx sql error:%v, data:%v", err.Error(), m)
		return -1, err
	}
	affected, _ := result.LastInsertId()
	return affected, nil
}

/*
func (i *Onlineorder) Update()  error {
    _,err := db.DBMap.Update(i)
    if err != nil{
		log.Error("update sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *onlineorderOp) Update(m *Onlineorder) error {
	return op.UpdateTx(db.DB, m)
}

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *onlineorderOp) UpdateTx(ext sqlx.Ext, m *Onlineorder) error {
	sql := `update onlineorder set OperUserID=?,ShareID=?,UserID=?,GameID=?,Accounts=?,OrderID=?,CardTypeID=?,CardPrice=?,CardGold=?,CardTotal=?,OrderAmount=?,DiscountScale=?,PayAmount=?,OrderStatus=?,IPAddress=?,ApplyDate=?,PhoneNum=?,GameName=?,NickName=?,GoodsNumber=?,GoodsID=?,GoodsName=?,OrderDate=?,PayType=? where OnLineID=?`
	_, err := ext.Exec(sql,
		m.OperUserID,
		m.ShareID,
		m.UserID,
		m.GameID,
		m.Accounts,
		m.OrderID,
		m.CardTypeID,
		m.CardPrice,
		m.CardGold,
		m.CardTotal,
		m.OrderAmount,
		m.DiscountScale,
		m.PayAmount,
		m.OrderStatus,
		m.IPAddress,
		m.ApplyDate,
		m.PhoneNum,
		m.GameName,
		m.NickName,
		m.GoodsNumber,
		m.GoodsID,
		m.GoodsName,
		m.OrderDate,
		m.PayType,
		m.OnLineID,
	)

	if err != nil {
		log.Error("update sql error:%v, data:%v", err.Error(), m)
		return err
	}

	return nil
}

// 用主键做条件，更新map里包含的字段名
func (op *onlineorderOp) UpdateWithMap(OnLineID int, m map[string]interface{}) error {
	return op.UpdateWithMapTx(db.DB, OnLineID, m)
}

// 用主键做条件，更新map里包含的字段名
func (op *onlineorderOp) UpdateWithMapTx(ext sqlx.Ext, OnLineID int, m map[string]interface{}) error {

	sql := `update onlineorder set %s where 1=1 and OnLineID=? ;`

	var params []interface{}
	var set_sql string
	for k, v := range m {
		if set_sql != "" {
			set_sql += ","
		}
		set_sql += fmt.Sprintf(" %s=? ", k)
		params = append(params, v)
	}
	params = append(params, OnLineID)
	_, err := ext.Exec(fmt.Sprintf(sql, set_sql), params...)
	return err
}

/*
func (i *Onlineorder) Delete() error{
    _,err := db.DBMap.Delete(i)
	log.Error("Delete sql error:%v", err.Error())
    return err
}
*/
// 根据主键删除相关记录
func (op *onlineorderOp) Delete(OnLineID int) error {
	return op.DeleteTx(db.DB, OnLineID)
}

// 根据主键删除相关记录,Tx
func (op *onlineorderOp) DeleteTx(ext sqlx.Ext, OnLineID int) error {
	sql := `delete from onlineorder where 1=1
        and OnLineID=?
        `
	_, err := ext.Exec(sql,
		OnLineID,
	)
	return err
}

// 返回符合查询条件的记录数
func (op *onlineorderOp) CountByMap(m map[string]interface{}) (int64, error) {

	var params []interface{}
	sql := `select count(*) from onlineorder where 1=1 `
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

func (op *onlineorderOp) DeleteByMap(m map[string]interface{}) (int64, error) {
	return op.DeleteByMapTx(db.DB, m)
}

func (op *onlineorderOp) DeleteByMapTx(ext sqlx.Ext, m map[string]interface{}) (int64, error) {
	var params []interface{}
	sql := "delete from onlineorder where 1=1 "
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
