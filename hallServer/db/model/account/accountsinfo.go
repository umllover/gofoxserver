package account

import (
	"errors"
	"fmt"
	"mj/hallServer/db"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lovelly/leaf/log"
)

//This file is generate by scripts,don't edit it

//accountsinfo
//

// +gen *
type Accountsinfo struct {
	UserID           int64      `db:"UserID" json:"UserID"`                     // 用户标识
	ProtectID        int        `db:"ProtectID" json:"ProtectID"`               // 密保标识
	SpreaderID       int        `db:"SpreaderID" json:"SpreaderID"`             // 推广员标识
	Accounts         string     `db:"Accounts" json:"Accounts"`                 // 用户帐号
	NickName         string     `db:"NickName" json:"NickName"`                 // 用户昵称
	PassPortID       string     `db:"PassPortID" json:"PassPortID"`             // 身份证号
	Compellation     string     `db:"Compellation" json:"Compellation"`         // 真实名字
	LogonPass        string     `db:"LogonPass" json:"LogonPass"`               // 登录密码
	IsAndroid        int8       `db:"IsAndroid" json:"IsAndroid"`               //
	InsurePass       string     `db:"InsurePass" json:"InsurePass"`             // 安全密码
	MasterOrder      int8       `db:"MasterOrder" json:"MasterOrder"`           // 管理等级
	Gender           int8       `db:"Gender" json:"Gender"`                     // 用户性别
	Nullity          int8       `db:"Nullity" json:"Nullity"`                   // 禁止服务
	NullityOverDate  *time.Time `db:"NullityOverDate" json:"NullityOverDate"`   // 禁止时间
	StunDown         int8       `db:"StunDown" json:"StunDown"`                 // 关闭标志
	MoorMachine      int8       `db:"MoorMachine" json:"MoorMachine"`           // 固定机器
	WebLogonTimes    int        `db:"WebLogonTimes" json:"WebLogonTimes"`       // 登录次数
	GameLogonTimes   int        `db:"GameLogonTimes" json:"GameLogonTimes"`     // 登录次数
	LastLogonIP      string     `db:"LastLogonIP" json:"LastLogonIP"`           // 登录地址
	LastLogonDate    *time.Time `db:"LastLogonDate" json:"LastLogonDate"`       // 登录时间
	LastLogonMobile  string     `db:"LastLogonMobile" json:"LastLogonMobile"`   // 登录手机
	LastLogonMachine string     `db:"LastLogonMachine" json:"LastLogonMachine"` // 登录机器
	RegisterIP       string     `db:"RegisterIP" json:"RegisterIP"`             // 注册地址
	RegisterDate     *time.Time `db:"RegisterDate" json:"RegisterDate"`         // 注册时间
	RegisterMobile   string     `db:"RegisterMobile" json:"RegisterMobile"`     // 注册手机
	RegisterMachine  string     `db:"RegisterMachine" json:"RegisterMachine"`   // 注册机器
	QQID             string     `db:"QQID" json:"QQID"`                         // QQ对应ID
	WXID             string     `db:"WXID" json:"WXID"`                         // 微信对应ID
	AgentID          int        `db:"AgentID" json:"AgentID"`                   //
	AgentNumber      string     `db:"AgentNumber" json:"AgentNumber"`           //
	HeadImgUrl       string     `db:"HeadImgUrl" json:"HeadImgUrl"`             //
	UnionID          string     `db:"UnionID" json:"UnionID"`                   //
	QQ               string     `db:"QQ" json:"QQ"`                             // QQ 号码
	EMail            string     `db:"EMail" json:"EMail"`                       //
	DwellingPlace    string     `db:"DwellingPlace" json:"DwellingPlace"`       // 详细住址
	PostalCode       string     `db:"PostalCode" json:"PostalCode"`             // 邮政编码
	Birthday         *time.Time `db:"Birthday" json:"Birthday"`                 // 生日
	OpenID           string     `db:"OpenID" json:"OpenID"`                     //
	SessionKey       string     `db:"SessionKey" json:"SessionKey"`             //
}

type accountsinfoOp struct{}

var AccountsinfoOp = &accountsinfoOp{}
var DefaultAccountsinfo = &Accountsinfo{}

// 按主键查询. 注:未找到记录的话将触发sql.ErrNoRows错误，返回nil, false
func (op *accountsinfoOp) Get(UserID int64) (*Accountsinfo, bool) {
	obj := &Accountsinfo{}
	sql := "select * from accountsinfo where UserID=? "
	err := db.AccountDB.Get(obj, sql,
		UserID,
	)

	if err != nil {
		log.Error("Get data error:%v", err.Error())
		return nil, false
	}
	return obj, true
}
func (op *accountsinfoOp) SelectAll() ([]*Accountsinfo, error) {
	objList := []*Accountsinfo{}
	sql := "select * from accountsinfo "
	err := db.AccountDB.Select(&objList, sql)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return objList, nil
}

func (op *accountsinfoOp) QueryByMap(m map[string]interface{}) ([]*Accountsinfo, error) {
	result := []*Accountsinfo{}
	var params []interface{}

	sql := "select * from accountsinfo where 1=1 "
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

func (op *accountsinfoOp) GetByMap(m map[string]interface{}) (*Accountsinfo, error) {
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
func (i *Accountsinfo) Insert() error {
    err := db.AccountDBMap.Insert(i)
    if err != nil{
		log.Error("Insert sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 插入数据，自增长字段将被忽略
func (op *accountsinfoOp) Insert(m *Accountsinfo) (int64, error) {
	return op.InsertTx(db.AccountDB, m)
}

// 插入数据，自增长字段将被忽略
func (op *accountsinfoOp) InsertTx(ext sqlx.Ext, m *Accountsinfo) (int64, error) {
	sql := "insert into accountsinfo(UserID,ProtectID,SpreaderID,Accounts,NickName,PassPortID,Compellation,LogonPass,IsAndroid,InsurePass,MasterOrder,Gender,Nullity,NullityOverDate,StunDown,MoorMachine,WebLogonTimes,GameLogonTimes,LastLogonIP,LastLogonDate,LastLogonMobile,LastLogonMachine,RegisterIP,RegisterDate,RegisterMobile,RegisterMachine,QQID,WXID,AgentID,AgentNumber,HeadImgUrl,UnionID,QQ,EMail,DwellingPlace,PostalCode,Birthday,OpenID,SessionKey) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	result, err := ext.Exec(sql,
		m.UserID,
		m.ProtectID,
		m.SpreaderID,
		m.Accounts,
		m.NickName,
		m.PassPortID,
		m.Compellation,
		m.LogonPass,
		m.IsAndroid,
		m.InsurePass,
		m.MasterOrder,
		m.Gender,
		m.Nullity,
		m.NullityOverDate,
		m.StunDown,
		m.MoorMachine,
		m.WebLogonTimes,
		m.GameLogonTimes,
		m.LastLogonIP,
		m.LastLogonDate,
		m.LastLogonMobile,
		m.LastLogonMachine,
		m.RegisterIP,
		m.RegisterDate,
		m.RegisterMobile,
		m.RegisterMachine,
		m.QQID,
		m.WXID,
		m.AgentID,
		m.AgentNumber,
		m.HeadImgUrl,
		m.UnionID,
		m.QQ,
		m.EMail,
		m.DwellingPlace,
		m.PostalCode,
		m.Birthday,
		m.OpenID,
		m.SessionKey,
	)
	if err != nil {
		log.Error("InsertTx sql error:%v, data:%v", err.Error(), m)
		return -1, err
	}
	affected, _ := result.LastInsertId()
	return affected, nil
}

//存在就更新， 不存在就插入
func (op *accountsinfoOp) InsertUpdate(obj *Accountsinfo, m map[string]interface{}) error {
	sql := "insert into accountsinfo(UserID,ProtectID,SpreaderID,Accounts,NickName,PassPortID,Compellation,LogonPass,IsAndroid,InsurePass,MasterOrder,Gender,Nullity,NullityOverDate,StunDown,MoorMachine,WebLogonTimes,GameLogonTimes,LastLogonIP,LastLogonDate,LastLogonMobile,LastLogonMachine,RegisterIP,RegisterDate,RegisterMobile,RegisterMachine,QQID,WXID,AgentID,AgentNumber,HeadImgUrl,UnionID,QQ,EMail,DwellingPlace,PostalCode,Birthday,OpenID,SessionKey) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE "
	var params = []interface{}{obj.UserID,
		obj.ProtectID,
		obj.SpreaderID,
		obj.Accounts,
		obj.NickName,
		obj.PassPortID,
		obj.Compellation,
		obj.LogonPass,
		obj.IsAndroid,
		obj.InsurePass,
		obj.MasterOrder,
		obj.Gender,
		obj.Nullity,
		obj.NullityOverDate,
		obj.StunDown,
		obj.MoorMachine,
		obj.WebLogonTimes,
		obj.GameLogonTimes,
		obj.LastLogonIP,
		obj.LastLogonDate,
		obj.LastLogonMobile,
		obj.LastLogonMachine,
		obj.RegisterIP,
		obj.RegisterDate,
		obj.RegisterMobile,
		obj.RegisterMachine,
		obj.QQID,
		obj.WXID,
		obj.AgentID,
		obj.AgentNumber,
		obj.HeadImgUrl,
		obj.UnionID,
		obj.QQ,
		obj.EMail,
		obj.DwellingPlace,
		obj.PostalCode,
		obj.Birthday,
		obj.OpenID,
		obj.SessionKey,
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
func (i *Accountsinfo) Update()  error {
    _,err := db.AccountDBMap.Update(i)
    if err != nil{
		log.Error("update sql error:%v, data:%v", err.Error(),i)
        return err
    }
}
*/

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *accountsinfoOp) Update(m *Accountsinfo) error {
	return op.UpdateTx(db.AccountDB, m)
}

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *accountsinfoOp) UpdateTx(ext sqlx.Ext, m *Accountsinfo) error {
	sql := `update accountsinfo set ProtectID=?,SpreaderID=?,Accounts=?,NickName=?,PassPortID=?,Compellation=?,LogonPass=?,IsAndroid=?,InsurePass=?,MasterOrder=?,Gender=?,Nullity=?,NullityOverDate=?,StunDown=?,MoorMachine=?,WebLogonTimes=?,GameLogonTimes=?,LastLogonIP=?,LastLogonDate=?,LastLogonMobile=?,LastLogonMachine=?,RegisterIP=?,RegisterDate=?,RegisterMobile=?,RegisterMachine=?,QQID=?,WXID=?,AgentID=?,AgentNumber=?,HeadImgUrl=?,UnionID=?,QQ=?,EMail=?,DwellingPlace=?,PostalCode=?,Birthday=?,OpenID=?,SessionKey=? where UserID=?`
	_, err := ext.Exec(sql,
		m.ProtectID,
		m.SpreaderID,
		m.Accounts,
		m.NickName,
		m.PassPortID,
		m.Compellation,
		m.LogonPass,
		m.IsAndroid,
		m.InsurePass,
		m.MasterOrder,
		m.Gender,
		m.Nullity,
		m.NullityOverDate,
		m.StunDown,
		m.MoorMachine,
		m.WebLogonTimes,
		m.GameLogonTimes,
		m.LastLogonIP,
		m.LastLogonDate,
		m.LastLogonMobile,
		m.LastLogonMachine,
		m.RegisterIP,
		m.RegisterDate,
		m.RegisterMobile,
		m.RegisterMachine,
		m.QQID,
		m.WXID,
		m.AgentID,
		m.AgentNumber,
		m.HeadImgUrl,
		m.UnionID,
		m.QQ,
		m.EMail,
		m.DwellingPlace,
		m.PostalCode,
		m.Birthday,
		m.OpenID,
		m.SessionKey,
		m.UserID,
	)

	if err != nil {
		log.Error("update sql error:%v, data:%v", err.Error(), m)
		return err
	}

	return nil
}

// 用主键做条件，更新map里包含的字段名
func (op *accountsinfoOp) UpdateWithMap(UserID int64, m map[string]interface{}) error {
	return op.UpdateWithMapTx(db.AccountDB, UserID, m)
}

// 用主键做条件，更新map里包含的字段名
func (op *accountsinfoOp) UpdateWithMapTx(ext sqlx.Ext, UserID int64, m map[string]interface{}) error {

	sql := `update accountsinfo set %s where 1=1 and UserID=? ;`

	var params []interface{}
	var set_sql string
	for k, v := range m {
		if set_sql != "" {
			set_sql += ","
		}
		set_sql += fmt.Sprintf(" %s=? ", k)
		params = append(params, v)
	}
	params = append(params, UserID)
	_, err := ext.Exec(fmt.Sprintf(sql, set_sql), params...)
	return err
}

/*
func (i *Accountsinfo) Delete() error{
    _,err := db.AccountDBMap.Delete(i)
	log.Error("Delete sql error:%v", err.Error())
    return err
}
*/
// 根据主键删除相关记录
func (op *accountsinfoOp) Delete(UserID int64) error {
	return op.DeleteTx(db.AccountDB, UserID)
}

// 根据主键删除相关记录,Tx
func (op *accountsinfoOp) DeleteTx(ext sqlx.Ext, UserID int64) error {
	sql := `delete from accountsinfo where 1=1
        and UserID=?
        `
	_, err := ext.Exec(sql,
		UserID,
	)
	return err
}

// 返回符合查询条件的记录数
func (op *accountsinfoOp) CountByMap(m map[string]interface{}) (int64, error) {

	var params []interface{}
	sql := `select count(*) from accountsinfo where 1=1 `
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

func (op *accountsinfoOp) DeleteByMap(m map[string]interface{}) (int64, error) {
	return op.DeleteByMapTx(db.AccountDB, m)
}

func (op *accountsinfoOp) DeleteByMapTx(ext sqlx.Ext, m map[string]interface{}) (int64, error) {
	var params []interface{}
	sql := "delete from accountsinfo where 1=1 "
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
