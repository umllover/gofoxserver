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

//mall_buy_log
//

// +gen *
type MallBuyLog struct {
	GoodsId           int        `db:"goods_id" json:"goods_id"`                       //
	Rmb               int        `db:"rmb" json:"rmb"`                                 //
	Diamond           int        `db:"diamond" json:"diamond"`                         //
	Name              string     `db:"name" json:"name"`                               // 商品名称
	LeftCnt           int        `db:"left_cnt" json:"left_cnt"`                       // 剩余数量
	SpecialOffer      int        `db:"special_offer" json:"special_offer"`             // 特价
	GivePresent       int        `db:"give_present" json:"give_present"`               // 赠送
	SpecialOfferBegin *time.Time `db:"special_offer_begin" json:"special_offer_begin"` // 特价开始时间
	SpecialOfferEnd   *time.Time `db:"special_offer_end" json:"special_offer_end"`     // 特价结束时间
	GoodsType         string     `db:"goods_type" json:"goods_type"`                   // 类别
	BuyTime           *time.Time `db:"buy_time" json:"buy_time"`                       // 购买时间
}

type mallBuyLogOp struct{}

var MallBuyLogOp = &mallBuyLogOp{}
var DefaultMallBuyLog = &MallBuyLog{}

// 按主键查询. 注:未找到记录的话将触发sql.ErrNoRows错误，返回nil, false
func (op *mallBuyLogOp) Get(goods_id int) (*MallBuyLog, bool) {
	obj := &MallBuyLog{}
	sql := "select * from mall_buy_log where goods_id=? "
	err := db.StatsDB.Get(obj, sql,
		goods_id,
	)

	if err != nil {
		log.Error("Get data error:%v", err.Error())
		return nil, false
	}
	return obj, true
}
func (op *mallBuyLogOp) SelectAll() ([]*MallBuyLog, error) {
	objList := []*MallBuyLog{}
	sql := "select * from mall_buy_log "
	err := db.StatsDB.Select(&objList, sql)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return objList, nil
}

func (op *mallBuyLogOp) QueryByMap(m map[string]interface{}) ([]*MallBuyLog, error) {
	result := []*MallBuyLog{}
	var params []interface{}

	sql := "select * from mall_buy_log where 1=1 "
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

func (op *mallBuyLogOp) GetByMap(m map[string]interface{}) (*MallBuyLog, error) {
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
func (i *MallBuyLog) Insert() error {
    err := db.StatsDBMap.Insert(i)
    if err != nil{
		log.Error("Insert sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 插入数据，自增长字段将被忽略
func (op *mallBuyLogOp) Insert(m *MallBuyLog) (int64, error) {
	return op.InsertTx(db.StatsDB, m)
}

// 插入数据，自增长字段将被忽略
func (op *mallBuyLogOp) InsertTx(ext sqlx.Ext, m *MallBuyLog) (int64, error) {
	sql := "insert into mall_buy_log(goods_id,rmb,diamond,name,left_cnt,special_offer,give_present,special_offer_begin,special_offer_end,goods_type,buy_time) values(?,?,?,?,?,?,?,?,?,?,?)"
	result, err := ext.Exec(sql,
		m.GoodsId,
		m.Rmb,
		m.Diamond,
		m.Name,
		m.LeftCnt,
		m.SpecialOffer,
		m.GivePresent,
		m.SpecialOfferBegin,
		m.SpecialOfferEnd,
		m.GoodsType,
		m.BuyTime,
	)
	if err != nil {
		log.Error("InsertTx sql error:%v, data:%v", err.Error(), m)
		return -1, err
	}
	affected, _ := result.LastInsertId()
	return affected, nil
}

//存在就更新， 不存在就插入
func (op *mallBuyLogOp) InsertUpdate(obj *MallBuyLog, m map[string]interface{}) error {
	sql := "insert into mall_buy_log(goods_id,rmb,diamond,name,left_cnt,special_offer,give_present,special_offer_begin,special_offer_end,goods_type,buy_time) values(?,?,?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE "
	var params = []interface{}{obj.GoodsId,
		obj.Rmb,
		obj.Diamond,
		obj.Name,
		obj.LeftCnt,
		obj.SpecialOffer,
		obj.GivePresent,
		obj.SpecialOfferBegin,
		obj.SpecialOfferEnd,
		obj.GoodsType,
		obj.BuyTime,
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
func (i *MallBuyLog) Update()  error {
    _,err := db.StatsDBMap.Update(i)
    if err != nil{
		log.Error("update sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *mallBuyLogOp) Update(m *MallBuyLog) error {
	return op.UpdateTx(db.StatsDB, m)
}

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *mallBuyLogOp) UpdateTx(ext sqlx.Ext, m *MallBuyLog) error {
	sql := `update mall_buy_log set rmb=?,diamond=?,name=?,left_cnt=?,special_offer=?,give_present=?,special_offer_begin=?,special_offer_end=?,goods_type=?,buy_time=? where goods_id=?`
	_, err := ext.Exec(sql,
		m.Rmb,
		m.Diamond,
		m.Name,
		m.LeftCnt,
		m.SpecialOffer,
		m.GivePresent,
		m.SpecialOfferBegin,
		m.SpecialOfferEnd,
		m.GoodsType,
		m.BuyTime,
		m.GoodsId,
	)

	if err != nil {
		log.Error("update sql error:%v, data:%v", err.Error(), m)
		return err
	}

	return nil
}

// 用主键做条件，更新map里包含的字段名
func (op *mallBuyLogOp) UpdateWithMap(goods_id int, m map[string]interface{}) error {
	return op.UpdateWithMapTx(db.StatsDB, goods_id, m)
}

// 用主键做条件，更新map里包含的字段名
func (op *mallBuyLogOp) UpdateWithMapTx(ext sqlx.Ext, goods_id int, m map[string]interface{}) error {

	sql := `update mall_buy_log set %s where 1=1 and goods_id=? ;`

	var params []interface{}
	var set_sql string
	for k, v := range m {
		if set_sql != "" {
			set_sql += ","
		}
		set_sql += fmt.Sprintf(" %s=? ", k)
		params = append(params, v)
	}
	params = append(params, goods_id)
	_, err := ext.Exec(fmt.Sprintf(sql, set_sql), params...)
	return err
}

/*
func (i *MallBuyLog) Delete() error{
    _,err := db.StatsDBMap.Delete(i)
	log.Error("Delete sql error:%v", err.Error())
    return err
}
*/
// 根据主键删除相关记录
func (op *mallBuyLogOp) Delete(goods_id int) error {
	return op.DeleteTx(db.StatsDB, goods_id)
}

// 根据主键删除相关记录,Tx
func (op *mallBuyLogOp) DeleteTx(ext sqlx.Ext, goods_id int) error {
	sql := `delete from mall_buy_log where 1=1
        and goods_id=?
        `
	_, err := ext.Exec(sql,
		goods_id,
	)
	return err
}

// 返回符合查询条件的记录数
func (op *mallBuyLogOp) CountByMap(m map[string]interface{}) (int64, error) {

	var params []interface{}
	sql := `select count(*) from mall_buy_log where 1=1 `
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

func (op *mallBuyLogOp) DeleteByMap(m map[string]interface{}) (int64, error) {
	return op.DeleteByMapTx(db.StatsDB, m)
}

func (op *mallBuyLogOp) DeleteByMapTx(ext sqlx.Ext, m map[string]interface{}) (int64, error) {
	var params []interface{}
	sql := "delete from mall_buy_log where 1=1 "
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
