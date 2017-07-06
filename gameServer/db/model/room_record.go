package model

import(
    "mj/gameServer/db"
    "github.com/lovelly/leaf/log"
    "github.com/jmoiron/sqlx"
    "fmt"
    "strings"
)

//This file is generate by scripts,don't edit it

//room_record
//

// +gen *
type RoomRecord struct {
    RoomId int `db:"room_id" json:"room_id"` // 
    KindId int `db:"kind_id" json:"kind_id"` // 
    UserId int `db:"user_id" json:"user_id"` // 创建房间的玩家id
    Status int `db:"status" json:"status"` // 游戏状态
    RoomName string `db:"room_name" json:"room_name"` // 房间名字
    JionUser string `db:"jion_user" json:"jion_user"` // 进入的玩家id
    }

type roomRecordOp struct{}

var RoomRecordOp = &roomRecordOp{}
var DefaultRoomRecord = &RoomRecord{}
// 按主键查询. 注:未找到记录的话将触发sql.ErrNoRows错误，返回nil, false
func (op *roomRecordOp) Get(room_id int) (*RoomRecord, bool) {
    obj := &RoomRecord{}
    sql := "select * from room_record where room_id=? "
    err := db.DB.Get(obj, sql, 
        room_id,
        )
    
    if err != nil{
        log.Error("Get data error:%v", err.Error())
        return nil,false
    }
    return obj, true
} 
func(op *roomRecordOp) SelectAll() ([]*RoomRecord, error) {
	objList := []*RoomRecord{}
	sql := "select * from room_record "
	err := db.DB.Select(&objList, sql)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return objList, nil
}

func(op *roomRecordOp) QueryByMap(m map[string]interface{}) ([]*RoomRecord, error) {
	result := []*RoomRecord{}
    var params []interface{}

	sql := "select * from room_record where 1=1 "
    for k, v := range m{
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


func(op *roomRecordOp) GetByMap(m map[string]interface{}) (*RoomRecord, error) {
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
func (i *RoomRecord) Insert() error {
    err := db.DBMap.Insert(i)
    if err != nil{
		log.Error("Insert sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 插入数据，自增长字段将被忽略
func (op *roomRecordOp) Insert(m *RoomRecord) (int64, error) {
    return op.InsertTx(db.DB, m)
}

// 插入数据，自增长字段将被忽略
func (op *roomRecordOp) InsertTx(ext sqlx.Ext, m *RoomRecord) (int64, error) {
    sql := "insert into room_record(room_id,kind_id,user_id,status,room_name,jion_user) values(?,?,?,?,?,?)"
    result, err := ext.Exec(sql,
    m.RoomId,
        m.KindId,
        m.UserId,
        m.Status,
        m.RoomName,
        m.JionUser,
        )
    if err != nil{
        log.Error("InsertTx sql error:%v, data:%v", err.Error(),m)
        return -1, err
    }
    affected, _ := result.LastInsertId()
        return affected, nil
    }

/*
func (i *RoomRecord) Update()  error {
    _,err := db.DBMap.Update(i)
    if err != nil{
		log.Error("update sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *roomRecordOp) Update(m *RoomRecord) (error) {
    return op.UpdateTx(db.DB, m)
}

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *roomRecordOp) UpdateTx(ext sqlx.Ext, m *RoomRecord) (error) {
    sql := `update room_record set kind_id=?,user_id=?,status=?,room_name=?,jion_user=? where room_id=?`
    _, err := ext.Exec(sql,
    m.KindId,
        m.UserId,
        m.Status,
        m.RoomName,
        m.JionUser,
        m.RoomId,
        )

    if err != nil{
		log.Error("update sql error:%v, data:%v", err.Error(),m)
        return err
    }

    return nil
}

// 用主键做条件，更新map里包含的字段名
func (op *roomRecordOp) UpdateWithMap(room_id int, m map[string]interface{}) (error) {
    return op.UpdateWithMapTx(db.DB, room_id, m)
}

// 用主键做条件，更新map里包含的字段名
func (op *roomRecordOp) UpdateWithMapTx(ext sqlx.Ext, room_id int, m map[string]interface{}) (error) {

    sql := `update room_record set %s where 1=1 and room_id=? ;`

    var params []interface{}
    var set_sql string
    for k, v := range m{
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
func (i *RoomRecord) Delete() error{
    _,err := db.DBMap.Delete(i)
	log.Error("Delete sql error:%v", err.Error())
    return err
}
*/
// 根据主键删除相关记录
func (op *roomRecordOp) Delete(room_id int) error{
    return op.DeleteTx(db.DB, room_id)
}

// 根据主键删除相关记录,Tx
func (op *roomRecordOp) DeleteTx(ext sqlx.Ext, room_id int) error{
    sql := `delete from room_record where 1=1
        and room_id=?
        `
    _, err := ext.Exec(sql, 
        room_id,
        )
    return err
}

// 返回符合查询条件的记录数
func (op *roomRecordOp) CountByMap(m map[string]interface{}) (int64, error) {

    var params []interface{}
    sql := `select count(*) from room_record where 1=1 `
    for k, v := range m{
        sql += fmt.Sprintf(" and  %s=? ",k)
        params = append(params, v)
    }
    count := int64(-1)
    err := db.DB.Get(&count, sql, params...)
    if err != nil {
        log.Error("CountByMap  error:%v data :%v", err.Error(), m)
		return 0,err
    }
    return count, nil
}

func (op *roomRecordOp) DeleteByMap(m map[string]interface{})(int64, error){
	return op.DeleteByMapTx(db.DB, m)
}

func (op *roomRecordOp) DeleteByMapTx(ext sqlx.Ext, m map[string]interface{}) (int64, error){
	var params []interface{}
	sql := "delete from room_record where 1=1 "
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

