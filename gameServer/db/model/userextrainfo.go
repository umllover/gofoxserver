package model

import(
    "mj/gameServer/db"
    "github.com/lovelly/leaf/log"
    "github.com/jmoiron/sqlx"
    "fmt"
    "strings"
)

//This file is generate by scripts,don't edit it

//userextrainfo
//

// +gen *
type Userextrainfo struct {
    UserId int64 `db:"UserId" json:"UserId"` // 
    MbPayTotal int `db:"MbPayTotal" json:"MbPayTotal"` // 手机充值总额
    MbVipLevel int `db:"MbVipLevel" json:"MbVipLevel"` // 手机VIP等级
    PayMbVipUpgrade int `db:"PayMbVipUpgrade" json:"PayMbVipUpgrade"` // 手机VIP升级，所需充值数（vip最高级时该值为0）
    MbTicket int `db:"MbTicket" json:"MbTicket"` // 手机兑换券数量
    }

type userextrainfoOp struct{}

var UserextrainfoOp = &userextrainfoOp{}
var DefaultUserextrainfo = &Userextrainfo{}
// 按主键查询. 注:未找到记录的话将触发sql.ErrNoRows错误，返回nil, false
func (op *userextrainfoOp) Get(UserId int64) (*Userextrainfo, bool) {
    obj := &Userextrainfo{}
    sql := "select * from userextrainfo where UserId=? "
    err := db.DB.Get(obj, sql, 
        UserId,
        )
    
    if err != nil{
        log.Error("Get data error:%v", err.Error())
        return nil,false
    }
    return obj, true
} 
func(op *userextrainfoOp) SelectAll() ([]*Userextrainfo, error) {
	objList := []*Userextrainfo{}
	sql := "select * from userextrainfo "
	err := db.DB.Select(&objList, sql)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return objList, nil
}

func(op *userextrainfoOp) QueryByMap(m map[string]interface{}) ([]*Userextrainfo, error) {
	result := []*Userextrainfo{}
    var params []interface{}

	sql := "select * from userextrainfo where 1=1 "
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


func(op *userextrainfoOp) GetByMap(m map[string]interface{}) (*Userextrainfo, error) {
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
func (i *Userextrainfo) Insert() error {
    err := db.DBMap.Insert(i)
    if err != nil{
		log.Error("Insert sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 插入数据，自增长字段将被忽略
func (op *userextrainfoOp) Insert(m *Userextrainfo) (int64, error) {
    return op.InsertTx(db.DB, m)
}

// 插入数据，自增长字段将被忽略
func (op *userextrainfoOp) InsertTx(ext sqlx.Ext, m *Userextrainfo) (int64, error) {
    sql := "insert into userextrainfo(UserId,MbPayTotal,MbVipLevel,PayMbVipUpgrade,MbTicket) values(?,?,?,?,?)"
    result, err := ext.Exec(sql,
    m.UserId,
        m.MbPayTotal,
        m.MbVipLevel,
        m.PayMbVipUpgrade,
        m.MbTicket,
        )
    if err != nil{
        log.Error("InsertTx sql error:%v, data:%v", err.Error(),m)
        return -1, err
    }
    affected, _ := result.LastInsertId()
        return affected, nil
    }

/*
func (i *Userextrainfo) Update()  error {
    _,err := db.DBMap.Update(i)
    if err != nil{
		log.Error("update sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *userextrainfoOp) Update(m *Userextrainfo) (error) {
    return op.UpdateTx(db.DB, m)
}

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *userextrainfoOp) UpdateTx(ext sqlx.Ext, m *Userextrainfo) (error) {
    sql := `update userextrainfo set MbPayTotal=?,MbVipLevel=?,PayMbVipUpgrade=?,MbTicket=? where UserId=?`
    _, err := ext.Exec(sql,
    m.MbPayTotal,
        m.MbVipLevel,
        m.PayMbVipUpgrade,
        m.MbTicket,
        m.UserId,
        )

    if err != nil{
		log.Error("update sql error:%v, data:%v", err.Error(),m)
        return err
    }

    return nil
}

// 用主键做条件，更新map里包含的字段名
func (op *userextrainfoOp) UpdateWithMap(UserId int64, m map[string]interface{}) (error) {
    return op.UpdateWithMapTx(db.DB, UserId, m)
}

// 用主键做条件，更新map里包含的字段名
func (op *userextrainfoOp) UpdateWithMapTx(ext sqlx.Ext, UserId int64, m map[string]interface{}) (error) {

    sql := `update userextrainfo set %s where 1=1 and UserId=? ;`

    var params []interface{}
    var set_sql string
    for k, v := range m{
		if set_sql != "" {
			set_sql += ","
		}
        set_sql += fmt.Sprintf(" %s=? ", k)
        params = append(params, v)
    }
	params = append(params, UserId)
    _, err := ext.Exec(fmt.Sprintf(sql, set_sql), params...)
    return err
}

/*
func (i *Userextrainfo) Delete() error{
    _,err := db.DBMap.Delete(i)
	log.Error("Delete sql error:%v", err.Error())
    return err
}
*/
// 根据主键删除相关记录
func (op *userextrainfoOp) Delete(UserId int64) error{
    return op.DeleteTx(db.DB, UserId)
}

// 根据主键删除相关记录,Tx
func (op *userextrainfoOp) DeleteTx(ext sqlx.Ext, UserId int64) error{
    sql := `delete from userextrainfo where 1=1
        and UserId=?
        `
    _, err := ext.Exec(sql, 
        UserId,
        )
    return err
}

// 返回符合查询条件的记录数
func (op *userextrainfoOp) CountByMap(m map[string]interface{}) (int64, error) {

    var params []interface{}
    sql := `select count(*) from userextrainfo where 1=1 `
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

func (op *userextrainfoOp) DeleteByMap(m map[string]interface{})(int64, error){
	return op.DeleteByMapTx(db.DB, m)
}

func (op *userextrainfoOp) DeleteByMapTx(ext sqlx.Ext, m map[string]interface{}) (int64, error){
	var params []interface{}
	sql := "delete from userextrainfo where 1=1 "
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

