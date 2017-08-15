package account

import (
	"fmt"
	"mj/hallServer/db"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lovelly/leaf/log"
)

//This file is generate by scripts,don't edit it

//agentinfo
//

// +gen *
type Agentinfo struct {
	AgentId          int64      `db:"agent_id" json:"agent_id"`                   //
	Level            int        `db:"level" json:"level"`                         //
	Phone            string     `db:"phone" json:"phone"`                         //
	RegisterDate     *time.Time `db:"register_date" json:"register_date"`         //
	Balance          int        `db:"balance" json:"balance"`                     //
	DiscipleNum      int        `db:"disciple_num" json:"disciple_num"`           //
	Commission       float64    `db:"commission" json:"commission"`               //
	RecentCommission float64    `db:"recent_commission" json:"recent_commission"` //
	AgentNum         string     `db:"agent_num" json:"agent_num"`                 //
	FormatAgentNum   string     `db:"format_agent_num" json:"format_agent_num"`   //
	ParAgentNum      string     `db:"par_agent_num" json:"par_agent_num"`         //
	HeadImgUrl       string     `db:"head_img_url" json:"head_img_url"`           //
	Account          string     `db:"account" json:"account"`                     // 用户id
}

type agentinfoOp struct{}

var AgentinfoOp = &agentinfoOp{}
var DefaultAgentinfo = &Agentinfo{}

// 按主键查询. 注:未找到记录的话将触发sql.ErrNoRows错误，返回nil, false
func (op *agentinfoOp) Get(agent_id int64) (*Agentinfo, bool) {
	obj := &Agentinfo{}
	sql := "select * from agentinfo where agent_id=? "
	err := db.AccountDB.Get(obj, sql,
		agent_id,
	)

	if err != nil {
		log.Error("Get data error:%v", err.Error())
		return nil, false
	}
	return obj, true
}
func (op *agentinfoOp) SelectAll() ([]*Agentinfo, error) {
	objList := []*Agentinfo{}
	sql := "select * from agentinfo "
	err := db.AccountDB.Select(&objList, sql)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return objList, nil
}

func (op *agentinfoOp) QueryByMap(m map[string]interface{}) ([]*Agentinfo, error) {
	result := []*Agentinfo{}
	var params []interface{}

	sql := "select * from agentinfo where 1=1 "
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

func (op *agentinfoOp) GetByMap(m map[string]interface{}) (*Agentinfo, error) {
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
func (i *Agentinfo) Insert() error {
    err := db.AccountDBMap.Insert(i)
    if err != nil{
		log.Error("Insert sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 插入数据，自增长字段将被忽略
func (op *agentinfoOp) Insert(m *Agentinfo) (int64, error) {
	return op.InsertTx(db.AccountDB, m)
}

// 插入数据，自增长字段将被忽略
func (op *agentinfoOp) InsertTx(ext sqlx.Ext, m *Agentinfo) (int64, error) {
	sql := "insert into agentinfo(agent_id,level,phone,register_date,balance,disciple_num,commission,recent_commission,agent_num,format_agent_num,par_agent_num,head_img_url,account) values(?,?,?,?,?,?,?,?,?,?,?,?,?)"
	result, err := ext.Exec(sql,
		m.AgentId,
		m.Level,
		m.Phone,
		m.RegisterDate,
		m.Balance,
		m.DiscipleNum,
		m.Commission,
		m.RecentCommission,
		m.AgentNum,
		m.FormatAgentNum,
		m.ParAgentNum,
		m.HeadImgUrl,
		m.Account,
	)
	if err != nil {
		log.Error("InsertTx sql error:%v, data:%v", err.Error(), m)
		return -1, err
	}
	affected, _ := result.LastInsertId()
	return affected, nil
}

//存在就更新， 不存在就插入
func (op *agentinfoOp) InsertUpdate(obj *Agentinfo, m map[string]interface{}) error {
	sql := "insert into agentinfo(agent_id,level,phone,register_date,balance,disciple_num,commission,recent_commission,agent_num,format_agent_num,par_agent_num,head_img_url,account) values(?,?,?,?,?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE "
	var params = []interface{}{obj.AgentId,
		obj.Level,
		obj.Phone,
		obj.RegisterDate,
		obj.Balance,
		obj.DiscipleNum,
		obj.Commission,
		obj.RecentCommission,
		obj.AgentNum,
		obj.FormatAgentNum,
		obj.ParAgentNum,
		obj.HeadImgUrl,
		obj.Account,
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
func (i *Agentinfo) Update()  error {
    _,err := db.AccountDBMap.Update(i)
    if err != nil{
		log.Error("update sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *agentinfoOp) Update(m *Agentinfo) error {
	return op.UpdateTx(db.AccountDB, m)
}

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *agentinfoOp) UpdateTx(ext sqlx.Ext, m *Agentinfo) error {
	sql := `update agentinfo set level=?,phone=?,register_date=?,balance=?,disciple_num=?,commission=?,recent_commission=?,agent_num=?,format_agent_num=?,par_agent_num=?,head_img_url=?,account=? where agent_id=?`
	_, err := ext.Exec(sql,
		m.Level,
		m.Phone,
		m.RegisterDate,
		m.Balance,
		m.DiscipleNum,
		m.Commission,
		m.RecentCommission,
		m.AgentNum,
		m.FormatAgentNum,
		m.ParAgentNum,
		m.HeadImgUrl,
		m.Account,
		m.AgentId,
	)

	if err != nil {
		log.Error("update sql error:%v, data:%v", err.Error(), m)
		return err
	}

	return nil
}

// 用主键做条件，更新map里包含的字段名
func (op *agentinfoOp) UpdateWithMap(agent_id int64, m map[string]interface{}) error {
	return op.UpdateWithMapTx(db.AccountDB, agent_id, m)
}

// 用主键做条件，更新map里包含的字段名
func (op *agentinfoOp) UpdateWithMapTx(ext sqlx.Ext, agent_id int64, m map[string]interface{}) error {

	sql := `update agentinfo set %s where 1=1 and agent_id=? ;`

	var params []interface{}
	var set_sql string
	for k, v := range m {
		if set_sql != "" {
			set_sql += ","
		}
		set_sql += fmt.Sprintf(" %s=? ", k)
		params = append(params, v)
	}
	params = append(params, agent_id)
	_, err := ext.Exec(fmt.Sprintf(sql, set_sql), params...)
	return err
}

/*
func (i *Agentinfo) Delete() error{
    _,err := db.AccountDBMap.Delete(i)
	log.Error("Delete sql error:%v", err.Error())
    return err
}
*/
// 根据主键删除相关记录
func (op *agentinfoOp) Delete(agent_id int64) error {
	return op.DeleteTx(db.AccountDB, agent_id)
}

// 根据主键删除相关记录,Tx
func (op *agentinfoOp) DeleteTx(ext sqlx.Ext, agent_id int64) error {
	sql := `delete from agentinfo where 1=1
        and agent_id=?
        `
	_, err := ext.Exec(sql,
		agent_id,
	)
	return err
}

// 返回符合查询条件的记录数
func (op *agentinfoOp) CountByMap(m map[string]interface{}) (int64, error) {

	var params []interface{}
	sql := `select count(*) from agentinfo where 1=1 `
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

func (op *agentinfoOp) DeleteByMap(m map[string]interface{}) (int64, error) {
	return op.DeleteByMapTx(db.AccountDB, m)
}

func (op *agentinfoOp) DeleteByMapTx(ext sqlx.Ext, m map[string]interface{}) (int64, error) {
	var params []interface{}
	sql := "delete from agentinfo where 1=1 "
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
