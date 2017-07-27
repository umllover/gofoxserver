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

//getin_room_log
//

// +gen *
type GetinRoomLog struct {
	GetinRecodeId int        `db:"getin_recode_id" json:"getin_recode_id"` // 加入游戏数据记录的Id
	RoomId        int        `db:"room_id" json:"room_id"`                 // 房间id
	UserId        int64      `db:"user_id" json:"user_id"`                 // 用户索引
	PublicGetin   int        `db:"public_getIn" json:"public_getIn"`       // 公房加入 0否 1是
	PrivateGetin  int        `db:"private_getIn" json:"private_getIn"`     // 私房加入 0否 1是
	TypeGetin     int        `db:"type_getIn" json:"type_getIn"`           // 加入房间类型 0列表加入 2输房号加入 3快速加入 4点击链接加入
	GetInTime     *time.Time `db:"getIn_time" json:"getIn_time"`           // 进入房间时间
}

type getinRoomLogOp struct{}

var GetinRoomLogOp = &getinRoomLogOp{}
var DefaultGetinRoomLog = &GetinRoomLog{}

// 按主键查询. 注:未找到记录的话将触发sql.ErrNoRows错误，返回nil, false
func (op *getinRoomLogOp) Get(getin_recode_id int) (*GetinRoomLog, bool) {
	obj := &GetinRoomLog{}
	sql := "select * from getin_room_log where getin_recode_id=? "
	err := db.StatsDB.Get(obj, sql,
		getin_recode_id,
	)

	if err != nil {
		log.Error("Get data error:%v", err.Error())
		return nil, false
	}
	return obj, true
}
func (op *getinRoomLogOp) SelectAll() ([]*GetinRoomLog, error) {
	objList := []*GetinRoomLog{}
	sql := "select * from getin_room_log "
	err := db.StatsDB.Select(&objList, sql)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return objList, nil
}

func (op *getinRoomLogOp) QueryByMap(m map[string]interface{}) ([]*GetinRoomLog, error) {
	result := []*GetinRoomLog{}
	var params []interface{}

	sql := "select * from getin_room_log where 1=1 "
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

func (op *getinRoomLogOp) GetByMap(m map[string]interface{}) (*GetinRoomLog, error) {
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
func (i *GetinRoomLog) Insert() error {
    err := db.StatsDBMap.Insert(i)
    if err != nil{
		log.Error("Insert sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 插入数据，自增长字段将被忽略
func (op *getinRoomLogOp) Insert(m *GetinRoomLog) (int64, error) {
	return op.InsertTx(db.StatsDB, m)
}

// 插入数据，自增长字段将被忽略
func (op *getinRoomLogOp) InsertTx(ext sqlx.Ext, m *GetinRoomLog) (int64, error) {
	sql := "insert into getin_room_log(room_id,user_id,public_getIn,private_getIn,type_getIn,getIn_time) values(?,?,?,?,?,?)"
	result, err := ext.Exec(sql,
		m.RoomId,
		m.UserId,
		m.PublicGetin,
		m.PrivateGetin,
		m.TypeGetin,
		m.GetInTime,
	)
	if err != nil {
		log.Error("InsertTx sql error:%v, data:%v", err.Error(), m)
		return -1, err
	}
	affected, _ := result.LastInsertId()
	return affected, nil
}

//存在就更新， 不存在就插入
func (op *getinRoomLogOp) InsertUpdate(obj *GetinRoomLog, m map[string]interface{}) error {
	sql := "insert into getin_room_log(room_id,user_id,public_getIn,private_getIn,type_getIn,getIn_time) values(?,?,?,?,?,?) ON DUPLICATE KEY UPDATE "
	var params = []interface{}{obj.RoomId,
		obj.UserId,
		obj.PublicGetin,
		obj.PrivateGetin,
		obj.TypeGetin,
		obj.GetInTime,
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
func (i *GetinRoomLog) Update()  error {
    _,err := db.StatsDBMap.Update(i)
    if err != nil{
		log.Error("update sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *getinRoomLogOp) Update(m *GetinRoomLog) error {
	return op.UpdateTx(db.StatsDB, m)
}

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *getinRoomLogOp) UpdateTx(ext sqlx.Ext, m *GetinRoomLog) error {
	sql := `update getin_room_log set room_id=?,user_id=?,public_getIn=?,private_getIn=?,type_getIn=?,getIn_time=? where getin_recode_id=?`
	_, err := ext.Exec(sql,
		m.RoomId,
		m.UserId,
		m.PublicGetin,
		m.PrivateGetin,
		m.TypeGetin,
		m.GetInTime,
		m.GetinRecodeId,
	)

	if err != nil {
		log.Error("update sql error:%v, data:%v", err.Error(), m)
		return err
	}

	return nil
}

// 用主键做条件，更新map里包含的字段名
func (op *getinRoomLogOp) UpdateWithMap(getin_recode_id int, m map[string]interface{}) error {
	return op.UpdateWithMapTx(db.StatsDB, getin_recode_id, m)
}

// 用主键做条件，更新map里包含的字段名
func (op *getinRoomLogOp) UpdateWithMapTx(ext sqlx.Ext, getin_recode_id int, m map[string]interface{}) error {

	sql := `update getin_room_log set %s where 1=1 and getin_recode_id=? ;`

	var params []interface{}
	var set_sql string
	for k, v := range m {
		if set_sql != "" {
			set_sql += ","
		}
		set_sql += fmt.Sprintf(" %s=? ", k)
		params = append(params, v)
	}
	params = append(params, getin_recode_id)
	_, err := ext.Exec(fmt.Sprintf(sql, set_sql), params...)
	return err
}

/*
func (i *GetinRoomLog) Delete() error{
    _,err := db.StatsDBMap.Delete(i)
	log.Error("Delete sql error:%v", err.Error())
    return err
}
*/
// 根据主键删除相关记录
func (op *getinRoomLogOp) Delete(getin_recode_id int) error {
	return op.DeleteTx(db.StatsDB, getin_recode_id)
}

// 根据主键删除相关记录,Tx
func (op *getinRoomLogOp) DeleteTx(ext sqlx.Ext, getin_recode_id int) error {
	sql := `delete from getin_room_log where 1=1
        and getin_recode_id=?
        `
	_, err := ext.Exec(sql,
		getin_recode_id,
	)
	return err
}

// 返回符合查询条件的记录数
func (op *getinRoomLogOp) CountByMap(m map[string]interface{}) (int64, error) {

	var params []interface{}
	sql := `select count(*) from getin_room_log where 1=1 `
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

func (op *getinRoomLogOp) DeleteByMap(m map[string]interface{}) (int64, error) {
	return op.DeleteByMapTx(db.StatsDB, m)
}

func (op *getinRoomLogOp) DeleteByMapTx(ext sqlx.Ext, m map[string]interface{}) (int64, error) {
	var params []interface{}
	sql := "delete from getin_room_log where 1=1 "
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
