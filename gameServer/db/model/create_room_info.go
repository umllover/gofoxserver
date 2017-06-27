package model

import (
	"errors"
	"fmt"
	"mj/gameServer/db"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lovelly/leaf/log"
)

//This file is generate by scripts,don't edit it

//create_room_info
//

// +gen *
type CreateRoomInfo struct {
	UserId       int        `db:"user_id" json:"user_id"`               // 用户索引
	KindId       int        `db:"kind_id" json:"kind_id"`               // 房间索引
	ServiceId    int        `db:"service_id" json:"service_id"`         // 游戏标识
	CreateTime   *time.Time `db:"create_time" json:"create_time"`       // 录入日期
	NodeId       int        `db:"node_id" json:"node_id"`               // 在哪个服务器上
	RoomId       int        `db:"room_id" json:"room_id"`               // 房间id
	Num          int        `db:"num" json:"num"`                       // 局数
	Status       int        `db:"status" json:"status"`                 //
	MaxPlayerCnt int        `db:"max_player_cnt" json:"max_player_cnt"` // 最多几个玩家进入
	PayType      int        `db:"pay_type" json:"pay_type"`             // 支付方式 1是全服 2是AA
	OtherInfo    string     `db:"other_info" json:"other_info"`         // 其他配置 json格式
}

type createRoomInfoOp struct{}

var CreateRoomInfoOp = &createRoomInfoOp{}
var DefaultCreateRoomInfo = &CreateRoomInfo{}

// 按主键查询. 注:未找到记录的话将触发sql.ErrNoRows错误，返回nil, false
func (op *createRoomInfoOp) Get(room_id int) (*CreateRoomInfo, bool) {
	obj := &CreateRoomInfo{}
	sql := "select * from create_room_info where room_id=? "
	err := db.DB.Get(obj, sql,
		room_id,
	)

	if err != nil {
		log.Error("Get data error:%v", err.Error())
		return nil, false
	}
	return obj, true
}
func (op *createRoomInfoOp) SelectAll() ([]*CreateRoomInfo, error) {
	objList := []*CreateRoomInfo{}
	sql := "select * from create_room_info "
	err := db.DB.Select(&objList, sql)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return objList, nil
}

func (op *createRoomInfoOp) QueryByMap(m map[string]interface{}) ([]*CreateRoomInfo, error) {
	result := []*CreateRoomInfo{}
	var params []interface{}

	sql := "select * from create_room_info where 1=1 "
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

func (op *createRoomInfoOp) QueryByMapQueryByMapComparison(m map[string]interface{}) ([]*CreateRoomInfo, error) {
	result := []*CreateRoomInfo{}
	var params []interface{}

	sql := "select * from create_room_info where 1=1 "
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

func (op *createRoomInfoOp) GetByMap(m map[string]interface{}) (*CreateRoomInfo, error) {
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
func (i *CreateRoomInfo) Insert() error {
    err := db.DBMap.Insert(i)
    if err != nil{
		log.Error("Insert sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 插入数据，自增长字段将被忽略
func (op *createRoomInfoOp) Insert(m *CreateRoomInfo) (int64, error) {
	return op.InsertTx(db.DB, m)
}

// 插入数据，自增长字段将被忽略
func (op *createRoomInfoOp) InsertTx(ext sqlx.Ext, m *CreateRoomInfo) (int64, error) {
	sql := "insert into create_room_info(user_id,kind_id,service_id,create_time,node_id,room_id,num,status,max_player_cnt,pay_type,other_info) values(?,?,?,?,?,?,?,?,?,?,?)"
	result, err := ext.Exec(sql,
		m.UserId,
		m.KindId,
		m.ServiceId,
		m.CreateTime,
		m.NodeId,
		m.RoomId,
		m.Num,
		m.Status,
		m.MaxPlayerCnt,
		m.PayType,
		m.OtherInfo,
	)
	if err != nil {
		log.Error("InsertTx sql error:%v, data:%v", err.Error(), m)
		return -1, err
	}
	affected, _ := result.LastInsertId()
	return affected, nil
}

/*
func (i *CreateRoomInfo) Update()  error {
    _,err := db.DBMap.Update(i)
    if err != nil{
		log.Error("update sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *createRoomInfoOp) Update(m *CreateRoomInfo) error {
	return op.UpdateTx(db.DB, m)
}

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *createRoomInfoOp) UpdateTx(ext sqlx.Ext, m *CreateRoomInfo) error {
	sql := `update create_room_info set user_id=?,kind_id=?,service_id=?,create_time=?,node_id=?,num=?,status=?,max_player_cnt=?,pay_type=?,other_info=? where room_id=?`
	_, err := ext.Exec(sql,
		m.UserId,
		m.KindId,
		m.ServiceId,
		m.CreateTime,
		m.NodeId,
		m.Num,
		m.Status,
		m.MaxPlayerCnt,
		m.PayType,
		m.OtherInfo,
		m.RoomId,
	)

	if err != nil {
		log.Error("update sql error:%v, data:%v", err.Error(), m)
		return err
	}

	return nil
}

// 用主键做条件，更新map里包含的字段名
func (op *createRoomInfoOp) UpdateWithMap(room_id int, m map[string]interface{}) error {
	return op.UpdateWithMapTx(db.DB, room_id, m)
}

// 用主键做条件，更新map里包含的字段名
func (op *createRoomInfoOp) UpdateWithMapTx(ext sqlx.Ext, room_id int, m map[string]interface{}) error {

	sql := `update create_room_info set %s where 1=1 and room_id=? ;`

	var params []interface{}
	var set_sql string
	for k, v := range m {
		if set_sql != "" {
			set_sql += ","
		}
		set_sql += fmt.Sprintf(" %s=? ", k)
		params = append(params, v)
	}
	params = append(params, room_id)
	_, err := ext.Exec(fmt.Sprintf(sql, set_sql), params...)
	return err
}

/*
func (i *CreateRoomInfo) Delete() error{
    _,err := db.DBMap.Delete(i)
	log.Error("Delete sql error:%v", err.Error())
    return err
}
*/
// 根据主键删除相关记录
func (op *createRoomInfoOp) Delete(room_id int) error {
	return op.DeleteTx(db.DB, room_id)
}

// 根据主键删除相关记录,Tx
func (op *createRoomInfoOp) DeleteTx(ext sqlx.Ext, room_id int) error {
	sql := `delete from create_room_info where 1=1
        and room_id=?
        `
	_, err := ext.Exec(sql,
		room_id,
	)
	return err
}

// 返回符合查询条件的记录数
func (op *createRoomInfoOp) CountByMap(m map[string]interface{}) (int64, error) {

	var params []interface{}
	sql := `select count(*) from create_room_info where 1=1 `
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

func (op *createRoomInfoOp) DeleteByMap(m map[string]interface{}) (int64, error) {
	return op.DeleteByMapTx(db.DB, m)
}

func (op *createRoomInfoOp) DeleteByMapTx(ext sqlx.Ext, m map[string]interface{}) (int64, error) {
	var params []interface{}
	sql := "delete from create_room_info where 1=1 "
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
