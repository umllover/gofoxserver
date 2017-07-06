package stats

import(
    "mj/gameServer/db"
    "github.com/lovelly/leaf/log"
    "github.com/jmoiron/sqlx"
    "fmt"
    "strings"
)

//This file is generate by scripts,don't edit it

//version_update_test
//

// +gen *
type VersionUpdateTest struct {
    Id int `db:"id" json:"id"` // 
    Test21 int `db:"test21" json:"test21"` // 
    Test22 int `db:"test22" json:"test22"` // 
    Test23 int `db:"test23" json:"test23"` // 
    Test31 int `db:"test31" json:"test31"` // 
    Test32 int `db:"test32" json:"test32"` // 
    Test33 int `db:"test33" json:"test33"` // 
    }

type versionUpdateTestOp struct{}

var VersionUpdateTestOp = &versionUpdateTestOp{}
var DefaultVersionUpdateTest = &VersionUpdateTest{}
// 按主键查询. 注:未找到记录的话将触发sql.ErrNoRows错误，返回nil, false
func (op *versionUpdateTestOp) Get(id int) (*VersionUpdateTest, bool) {
    obj := &VersionUpdateTest{}
    sql := "select * from version_update_test where id=? "
    err := db.StatsDB.Get(obj, sql, 
        id,
        )
    
    if err != nil{
        log.Error("Get data error:%v", err.Error())
        return nil,false
    }
    return obj, true
} 
func(op *versionUpdateTestOp) SelectAll() ([]*VersionUpdateTest, error) {
	objList := []*VersionUpdateTest{}
	sql := "select * from version_update_test "
	err := db.StatsDB.Select(&objList, sql)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return objList, nil
}

func(op *versionUpdateTestOp) QueryByMap(m map[string]interface{}) ([]*VersionUpdateTest, error) {
	result := []*VersionUpdateTest{}
    var params []interface{}

	sql := "select * from version_update_test where 1=1 "
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


func(op *versionUpdateTestOp) GetByMap(m map[string]interface{}) (*VersionUpdateTest, error) {
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
func (i *VersionUpdateTest) Insert() error {
    err := db.StatsDBMap.Insert(i)
    if err != nil{
		log.Error("Insert sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 插入数据，自增长字段将被忽略
func (op *versionUpdateTestOp) Insert(m *VersionUpdateTest) (int64, error) {
    return op.InsertTx(db.StatsDB, m)
}

// 插入数据，自增长字段将被忽略
func (op *versionUpdateTestOp) InsertTx(ext sqlx.Ext, m *VersionUpdateTest) (int64, error) {
    sql := "insert into version_update_test(id,test21,test22,test23,test31,test32,test33) values(?,?,?,?,?,?,?)"
    result, err := ext.Exec(sql,
    m.Id,
        m.Test21,
        m.Test22,
        m.Test23,
        m.Test31,
        m.Test32,
        m.Test33,
        )
    if err != nil{
        log.Error("InsertTx sql error:%v, data:%v", err.Error(),m)
        return -1, err
    }
    affected, _ := result.LastInsertId()
        return affected, nil
    }

/*
func (i *VersionUpdateTest) Update()  error {
    _,err := db.StatsDBMap.Update(i)
    if err != nil{
		log.Error("update sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *versionUpdateTestOp) Update(m *VersionUpdateTest) (error) {
    return op.UpdateTx(db.StatsDB, m)
}

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *versionUpdateTestOp) UpdateTx(ext sqlx.Ext, m *VersionUpdateTest) (error) {
    sql := `update version_update_test set test21=?,test22=?,test23=?,test31=?,test32=?,test33=? where id=?`
    _, err := ext.Exec(sql,
    m.Test21,
        m.Test22,
        m.Test23,
        m.Test31,
        m.Test32,
        m.Test33,
        m.Id,
        )

    if err != nil{
		log.Error("update sql error:%v, data:%v", err.Error(),m)
        return err
    }

    return nil
}

// 用主键做条件，更新map里包含的字段名
func (op *versionUpdateTestOp) UpdateWithMap(id int, m map[string]interface{}) (error) {
    return op.UpdateWithMapTx(db.StatsDB, id, m)
}

// 用主键做条件，更新map里包含的字段名
func (op *versionUpdateTestOp) UpdateWithMapTx(ext sqlx.Ext, id int, m map[string]interface{}) (error) {

    sql := `update version_update_test set %s where 1=1 and id=? ;`

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
func (i *VersionUpdateTest) Delete() error{
    _,err := db.StatsDBMap.Delete(i)
	log.Error("Delete sql error:%v", err.Error())
    return err
}
*/
// 根据主键删除相关记录
func (op *versionUpdateTestOp) Delete(id int) error{
    return op.DeleteTx(db.StatsDB, id)
}

// 根据主键删除相关记录,Tx
func (op *versionUpdateTestOp) DeleteTx(ext sqlx.Ext, id int) error{
    sql := `delete from version_update_test where 1=1
        and id=?
        `
    _, err := ext.Exec(sql, 
        id,
        )
    return err
}

// 返回符合查询条件的记录数
func (op *versionUpdateTestOp) CountByMap(m map[string]interface{}) (int64, error) {

    var params []interface{}
    sql := `select count(*) from version_update_test where 1=1 `
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

func (op *versionUpdateTestOp) DeleteByMap(m map[string]interface{})(int64, error){
	return op.DeleteByMapTx(db.StatsDB, m)
}

func (op *versionUpdateTestOp) DeleteByMapTx(ext sqlx.Ext, m map[string]interface{}) (int64, error){
	var params []interface{}
	sql := "delete from version_update_test where 1=1 "
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

