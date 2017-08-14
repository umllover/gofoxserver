package account

import (
	"fmt"
	"mj/gameServer/db"

	"github.com/jmoiron/sqlx"
	"github.com/lovelly/leaf/log"
)

//This file is generate by scripts,don't edit it

//server_list
//

// +gen *
type ServerList struct {
	SvrId   int    `db:"svr_id" json:"svr_id"`     // 节点id
	SvrType int    `db:"svr_type" json:"svr_type"` // 服务器类型 1是大厅服
	Host    string `db:"host" json:"host"`         // ip
	Port    int    `db:"port" json:"port"`         // 端口
	Status  int    `db:"status" json:"status"`     // 状态 1是正常状态 2是维护
}

type serverListOp struct{}

var ServerListOp = &serverListOp{}
var DefaultServerList = &ServerList{}

// 按主键查询. 注:未找到记录的话将触发sql.ErrNoRows错误，返回nil, false
func (op *serverListOp) Get(svr_id int, svr_type int) (*ServerList, bool) {
	obj := &ServerList{}
	sql := "select * from server_list where svr_id=? and svr_type=? "
	err := db.AccountDB.Get(obj, sql,
		svr_id,
		svr_type,
	)

	if err != nil {
		log.Error("Get data error:%v", err.Error())
		return nil, false
	}
	return obj, true
}
func (op *serverListOp) SelectAll() ([]*ServerList, error) {
	objList := []*ServerList{}
	sql := "select * from server_list "
	err := db.AccountDB.Select(&objList, sql)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return objList, nil
}

func (op *serverListOp) QueryByMap(m map[string]interface{}) ([]*ServerList, error) {
	result := []*ServerList{}
	var params []interface{}

	sql := "select * from server_list where 1=1 "
	for k, v := range m {
		sql += fmt.Sprintf(" and %s=? ", k)
		params = append(params, v)
	}
	err := db.AccountDB.Select(&result, sql, params...)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return result, nil
}

func (op *serverListOp) GetByMap(m map[string]interface{}) (*ServerList, error) {
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
func (i *ServerList) Insert() error {
    err := db.AccountDBMap.Insert(i)
    if err != nil{
		log.Error("Insert sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 插入数据，自增长字段将被忽略
func (op *serverListOp) Insert(m *ServerList) (int64, error) {
	return op.InsertTx(db.AccountDB, m)
}

// 插入数据，自增长字段将被忽略
func (op *serverListOp) InsertTx(ext sqlx.Ext, m *ServerList) (int64, error) {
	sql := "insert into server_list(svr_id,svr_type,host,port,status) values(?,?,?,?,?)"
	result, err := ext.Exec(sql,
		m.SvrId,
		m.SvrType,
		m.Host,
		m.Port,
		m.Status,
	)
	if err != nil {
		log.Error("InsertTx sql error:%v, data:%v", err.Error(), m)
		return -1, err
	}
	affected, _ := result.LastInsertId()
	return affected, nil
}

//存在就更新， 不存在就插入
func (op *serverListOp) InsertUpdate(obj *ServerList, m map[string]interface{}) error {
	sql := "insert into server_list(svr_id,svr_type,host,port,status) values(?,?,?,?,?) ON DUPLICATE KEY UPDATE "
	var params = []interface{}{obj.SvrId,
		obj.SvrType,
		obj.Host,
		obj.Port,
		obj.Status,
	}
	var set_sql string
	for k, v := range m {
		if set_sql != "" {
			set_sql += ","
		}
		set_sql += fmt.Sprintf(" %s=? ", k)
		params = append(params, v)
	}

	_, err := db.AccountDB.Exec(sql+set_sql, params...)
	return err
}

/*
func (i *ServerList) Update()  error {
    _,err := db.AccountDBMap.Update(i)
    if err != nil{
		log.Error("update sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *serverListOp) Update(m *ServerList) error {
	return op.UpdateTx(db.AccountDB, m)
}

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *serverListOp) UpdateTx(ext sqlx.Ext, m *ServerList) error {
	sql := `update server_list set host=?,port=?,status=? where svr_id=? and svr_type=?`
	_, err := ext.Exec(sql,
		m.Host,
		m.Port,
		m.Status,
		m.SvrId,
		m.SvrType,
	)

	if err != nil {
		log.Error("update sql error:%v, data:%v", err.Error(), m)
		return err
	}

	return nil
}

// 用主键做条件，更新map里包含的字段名
func (op *serverListOp) UpdateWithMap(svr_id int, svr_type int, m map[string]interface{}) error {
	return op.UpdateWithMapTx(db.AccountDB, svr_id, svr_type, m)
}

// 用主键做条件，更新map里包含的字段名
func (op *serverListOp) UpdateWithMapTx(ext sqlx.Ext, svr_id int, svr_type int, m map[string]interface{}) error {

	sql := `update server_list set %s where 1=1 and svr_id=? and svr_type=? ;`

	var params []interface{}
	var set_sql string
	for k, v := range m {
		if set_sql != "" {
			set_sql += ","
		}
		set_sql += fmt.Sprintf(" %s=? ", k)
		params = append(params, v)
	}
	params = append(params, svr_id, svr_type)
	_, err := ext.Exec(fmt.Sprintf(sql, set_sql), params...)
	return err
}

/*
func (i *ServerList) Delete() error{
    _,err := db.AccountDBMap.Delete(i)
	log.Error("Delete sql error:%v", err.Error())
    return err
}
*/
// 根据主键删除相关记录
func (op *serverListOp) Delete(svr_id int, svr_type int) error {
	return op.DeleteTx(db.AccountDB, svr_id, svr_type)
}

// 根据主键删除相关记录,Tx
func (op *serverListOp) DeleteTx(ext sqlx.Ext, svr_id int, svr_type int) error {
	sql := `delete from server_list where 1=1
        and svr_id=?
        and svr_type=?
        `
	_, err := ext.Exec(sql,
		svr_id,
		svr_type,
	)
	return err
}

// 返回符合查询条件的记录数
func (op *serverListOp) CountByMap(m map[string]interface{}) (int64, error) {

	var params []interface{}
	sql := `select count(*) from server_list where 1=1 `
	for k, v := range m {
		sql += fmt.Sprintf(" and  %s=? ", k)
		params = append(params, v)
	}
	count := int64(-1)
	err := db.AccountDB.Get(&count, sql, params...)
	if err != nil {
		log.Error("CountByMap  error:%v data :%v", err.Error(), m)
		return 0, err
	}
	return count, nil
}

func (op *serverListOp) DeleteByMap(m map[string]interface{}) (int64, error) {
	return op.DeleteByMapTx(db.AccountDB, m)
}

func (op *serverListOp) DeleteByMapTx(ext sqlx.Ext, m map[string]interface{}) (int64, error) {
	var params []interface{}
	sql := "delete from server_list where 1=1 "
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
