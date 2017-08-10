package account

import (
	"errors"
	"fmt"
	"mj/hallServer/db"

	"github.com/jmoiron/sqlx"
	"github.com/lovelly/leaf/log"
)

//This file is generate by scripts,don't edit it

//serverversion
//

// +gen *
type Serverversion struct {
	ServerId      string `db:"server_id" json:"server_id"`           //
	VersionName   string `db:"version_name" json:"version_name"`     //
	ServerVersion string `db:"server_version" json:"server_version"` //
	UpdateText    string `db:"update_text" json:"update_text"`       //
	DownloadPath  string `db:"download_path" json:"download_path"`   //
}

type serverversionOp struct{}

var ServerversionOp = &serverversionOp{}
var DefaultServerversion = &Serverversion{}

// 按主键查询. 注:未找到记录的话将触发sql.ErrNoRows错误，返回nil, false
func (op *serverversionOp) Get(server_id string) (*Serverversion, bool) {
	obj := &Serverversion{}
	sql := "select * from serverversion where server_id=? "
	err := db.AccountDB.Get(obj, sql,
		server_id,
	)

	if err != nil {
		log.Error("Get data error:%v", err.Error())
		return nil, false
	}
	return obj, true
}
func (op *serverversionOp) SelectAll() ([]*Serverversion, error) {
	objList := []*Serverversion{}
	sql := "select * from serverversion "
	err := db.AccountDB.Select(&objList, sql)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return objList, nil
}

func (op *serverversionOp) QueryByMap(m map[string]interface{}) ([]*Serverversion, error) {
	result := []*Serverversion{}
	var params []interface{}

	sql := "select * from serverversion where 1=1 "
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

func (op *serverversionOp) GetByMap(m map[string]interface{}) (*Serverversion, error) {
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
func (i *Serverversion) Insert() error {
    err := db.AccountDBMap.Insert(i)
    if err != nil{
		log.Error("Insert sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 插入数据，自增长字段将被忽略
func (op *serverversionOp) Insert(m *Serverversion) (int64, error) {
	return op.InsertTx(db.AccountDB, m)
}

// 插入数据，自增长字段将被忽略
func (op *serverversionOp) InsertTx(ext sqlx.Ext, m *Serverversion) (int64, error) {
	sql := "insert into serverversion(server_id,version_name,server_version,update_text,download_path) values(?,?,?,?,?)"
	result, err := ext.Exec(sql,
		m.ServerId,
		m.VersionName,
		m.ServerVersion,
		m.UpdateText,
		m.DownloadPath,
	)
	if err != nil {
		log.Error("InsertTx sql error:%v, data:%v", err.Error(), m)
		return -1, err
	}
	affected, _ := result.LastInsertId()
	return affected, nil
}

//存在就更新， 不存在就插入
func (op *serverversionOp) InsertUpdate(obj *Serverversion, m map[string]interface{}) error {
	sql := "insert into serverversion(server_id,version_name,server_version,update_text,download_path) values(?,?,?,?,?) ON DUPLICATE KEY UPDATE "
	var params = []interface{}{obj.ServerId,
		obj.VersionName,
		obj.ServerVersion,
		obj.UpdateText,
		obj.DownloadPath,
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
func (i *Serverversion) Update()  error {
    _,err := db.AccountDBMap.Update(i)
    if err != nil{
		log.Error("update sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *serverversionOp) Update(m *Serverversion) error {
	return op.UpdateTx(db.AccountDB, m)
}

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *serverversionOp) UpdateTx(ext sqlx.Ext, m *Serverversion) error {
	sql := `update serverversion set version_name=?,server_version=?,update_text=?,download_path=? where server_id=?`
	_, err := ext.Exec(sql,
		m.VersionName,
		m.ServerVersion,
		m.UpdateText,
		m.DownloadPath,
		m.ServerId,
	)

	if err != nil {
		log.Error("update sql error:%v, data:%v", err.Error(), m)
		return err
	}

	return nil
}

// 用主键做条件，更新map里包含的字段名
func (op *serverversionOp) UpdateWithMap(server_id string, m map[string]interface{}) error {
	return op.UpdateWithMapTx(db.AccountDB, server_id, m)
}

// 用主键做条件，更新map里包含的字段名
func (op *serverversionOp) UpdateWithMapTx(ext sqlx.Ext, server_id string, m map[string]interface{}) error {

	sql := `update serverversion set %s where 1=1 and server_id=? ;`

	var params []interface{}
	var set_sql string
	for k, v := range m {
		if set_sql != "" {
			set_sql += ","
		}
		set_sql += fmt.Sprintf(" %s=? ", k)
		params = append(params, v)
	}
	params = append(params, server_id)
	_, err := ext.Exec(fmt.Sprintf(sql, set_sql), params...)
	return err
}

/*
func (i *Serverversion) Delete() error{
    _,err := db.AccountDBMap.Delete(i)
	log.Error("Delete sql error:%v", err.Error())
    return err
}
*/
// 根据主键删除相关记录
func (op *serverversionOp) Delete(server_id string) error {
	return op.DeleteTx(db.AccountDB, server_id)
}

// 根据主键删除相关记录,Tx
func (op *serverversionOp) DeleteTx(ext sqlx.Ext, server_id string) error {
	sql := `delete from serverversion where 1=1
        and server_id=?
        `
	_, err := ext.Exec(sql,
		server_id,
	)
	return err
}

// 返回符合查询条件的记录数
func (op *serverversionOp) CountByMap(m map[string]interface{}) (int64, error) {

	var params []interface{}
	sql := `select count(*) from serverversion where 1=1 `
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

func (op *serverversionOp) DeleteByMap(m map[string]interface{}) (int64, error) {
	return op.DeleteByMapTx(db.AccountDB, m)
}

func (op *serverversionOp) DeleteByMapTx(ext sqlx.Ext, m map[string]interface{}) (int64, error) {
	var params []interface{}
	sql := "delete from serverversion where 1=1 "
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
