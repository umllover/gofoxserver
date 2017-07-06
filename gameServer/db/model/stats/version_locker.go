package stats

import(
    "mj/gameServer/db"
    "github.com/lovelly/leaf/log"
    "github.com/jmoiron/sqlx"
    "fmt"
    "strings"
)

//This file is generate by scripts,don't edit it

//version_locker
//

// +gen *
type VersionLocker struct {
    Id int `db:"id" json:"id"` // 
    }

type versionLockerOp struct{}

var VersionLockerOp = &versionLockerOp{}
var DefaultVersionLocker = &VersionLocker{}
// 按主键查询. 注:未找到记录的话将触发sql.ErrNoRows错误，返回nil, false
func (op *versionLockerOp) Get(id int) (*VersionLocker, bool) {
    obj := &VersionLocker{}
    sql := "select * from version_locker where id=? "
    err := db.StatsDB.Get(obj, sql, 
        id,
        )
    
    if err != nil{
        log.Error("Get data error:%v", err.Error())
        return nil,false
    }
    return obj, true
} 
func(op *versionLockerOp) SelectAll() ([]*VersionLocker, error) {
	objList := []*VersionLocker{}
	sql := "select * from version_locker "
	err := db.StatsDB.Select(&objList, sql)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return objList, nil
}

func(op *versionLockerOp) QueryByMap(m map[string]interface{}) ([]*VersionLocker, error) {
	result := []*VersionLocker{}
    var params []interface{}

	sql := "select * from version_locker where 1=1 "
    for k, v := range m{
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


func(op *versionLockerOp) GetByMap(m map[string]interface{}) (*VersionLocker, error) {
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
func (i *VersionLocker) Insert() error {
    err := db.StatsDBMap.Insert(i)
    if err != nil{
		log.Error("Insert sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 插入数据，自增长字段将被忽略
func (op *versionLockerOp) Insert(m *VersionLocker) (int64, error) {
    return op.InsertTx(db.StatsDB, m)
}

// 插入数据，自增长字段将被忽略
func (op *versionLockerOp) InsertTx(ext sqlx.Ext, m *VersionLocker) (int64, error) {
    sql := "insert into version_locker(id) values(?)"
    result, err := ext.Exec(sql,
    m.Id,
        )
    if err != nil{
        log.Error("InsertTx sql error:%v, data:%v", err.Error(),m)
        return -1, err
    }
    affected, _ := result.LastInsertId()
        return affected, nil
    }

/*
func (i *VersionLocker) Update()  error {
    _,err := db.StatsDBMap.Update(i)
    if err != nil{
		log.Error("update sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *versionLockerOp) Update(m *VersionLocker) (error) {
    return op.UpdateTx(db.StatsDB, m)
}

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *versionLockerOp) UpdateTx(ext sqlx.Ext, m *VersionLocker) (error) {
    sql := `update version_locker set  where id=?`
    _, err := ext.Exec(sql,
    m.Id,
        )

    if err != nil{
		log.Error("update sql error:%v, data:%v", err.Error(),m)
        return err
    }

    return nil
}

// 用主键做条件，更新map里包含的字段名
func (op *versionLockerOp) UpdateWithMap(id int, m map[string]interface{}) (error) {
    return op.UpdateWithMapTx(db.StatsDB, id, m)
}

// 用主键做条件，更新map里包含的字段名
func (op *versionLockerOp) UpdateWithMapTx(ext sqlx.Ext, id int, m map[string]interface{}) (error) {

    sql := `update version_locker set %s where 1=1 and id=? ;`

    var params []interface{}
    var set_sql string
    for k, v := range m{
		if set_sql != "" {
			set_sql += ","
		}
        set_sql += fmt.Sprintf(" %s=? ", k)
        params = append(params, v)
    }
	params = append(params, id)
    _, err := ext.Exec(fmt.Sprintf(sql, set_sql), params...)
    return err
}

/*
func (i *VersionLocker) Delete() error{
    _,err := db.StatsDBMap.Delete(i)
	log.Error("Delete sql error:%v", err.Error())
    return err
}
*/
// 根据主键删除相关记录
func (op *versionLockerOp) Delete(id int) error{
    return op.DeleteTx(db.StatsDB, id)
}

// 根据主键删除相关记录,Tx
func (op *versionLockerOp) DeleteTx(ext sqlx.Ext, id int) error{
    sql := `delete from version_locker where 1=1
        and id=?
        `
    _, err := ext.Exec(sql, 
        id,
        )
    return err
}

// 返回符合查询条件的记录数
func (op *versionLockerOp) CountByMap(m map[string]interface{}) (int64, error) {

    var params []interface{}
    sql := `select count(*) from version_locker where 1=1 `
    for k, v := range m{
        sql += fmt.Sprintf(" and  %s=? ",k)
        params = append(params, v)
    }
    count := int64(-1)
    err := db.StatsDB.Get(&count, sql, params...)
    if err != nil {
        log.Error("CountByMap  error:%v data :%v", err.Error(), m)
		return 0,err
    }
    return count, nil
}

func (op *versionLockerOp) DeleteByMap(m map[string]interface{})(int64, error){
	return op.DeleteByMapTx(db.StatsDB, m)
}

func (op *versionLockerOp) DeleteByMapTx(ext sqlx.Ext, m map[string]interface{}) (int64, error){
	var params []interface{}
	sql := "delete from version_locker where 1=1 "
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

