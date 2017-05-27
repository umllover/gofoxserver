package {{package_name}}

import(
    "mj/hallServer/db"
    "github.com/lovelly/leaf/log"
    "github.com/jmoiron/sqlx"
    "fmt"
    "strings"
    "ChessCard/server/utils"
)

//This file is generate by scripts,don't edit it

//{{table_name}}
//{{table_comment}}

{% if is_base_db -%}
// +gen 
{% else -%}
// +gen *
{% endif -%}

type {{struct_name}} struct {
    {% for column in column_list -%}
    {{column.field_name}} {{column.type}} `db:"{{column.name}}" json:"{{column.name}}"` // {{column.comment}}
    {% endfor -%}
}

{% if not is_base_db -%}
type {{op_struct_name}} struct{}

var {{op_name}} = &{{op_struct_name}}{}
{% endif -%}

{% if is_base_db -%}
var Default{{struct_name}} = {{struct_name}}{}
{% else -%}
var Default{{struct_name}} = &{{struct_name}}{}
{% endif -%}

{% if not is_base_db -%}
{% if primary_key -%}
// 按主键查询. 注:未找到记录的话将触发sql.ErrNoRows错误，返回nil, false
func (op *{{op_struct_name}}) Get({{primary_key_params}}) (*{{struct_name}}, bool) {
    obj := &{{struct_name}}{}
    sql := "{{get_by_pk_sql2}}"
    err := db.{{db_sel}}.Get(obj, sql, 
        {% for key in primary_key -%}
        {{key}},
        {% endfor -%}
        )
    
    if err != nil{
        log.Error(err.Error())
        return nil,false
    }
    return obj, true
} 
{% endif -%}

func(op *{{op_struct_name}}) SelectAll() ([]*{{struct_name}}, error) {
	objList := []*{{struct_name}}{}
	sql := "select * from {{table_name}} "
	err := db.{{db_sel}}.Select(&objList, sql)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return objList, nil
}

func(op *{{op_struct_name}}) QueryByMap(m map[string]interface{}) ([]*{{struct_name}}, error) {
	result := []*{{struct_name}}{}
    var params []interface{}

	sql := "select * from {{table_name}} where 1=1 "
    for k, v := range m{
        sql += fmt.Sprintf(" and %s=? ", k)
        params = append(params, v)
    }
	err := db.{{db_sel}}.Select(&result, sql, params...)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return result, nil
}

func(op *{{op_struct_name}}) QueryByMapQueryByMapComparison(m map[string]interface{}) ([]*{{struct_name}}, error) {
    result := []*{{struct_name}}{}
    var params []interface{}

    sql := "select * from {{table_name}} where 1=1 "
    for k, v := range m{
        sql += fmt.Sprintf(" and %s? ", k)
        params = append(params, v)
    }
    err := db.{{db_sel}}.Select(&result, sql, params...)
    if err != nil {
        log.Error(err.Error())
        return nil, err
    }
    return result, nil
}

func(op *{{op_struct_name}}) GetByMap(m map[string]interface{}) (*{{struct_name}}, error) {
    lst, err := op.QueryByMap(m)
    if err != nil {
        return nil, err
    }
    if len(lst) > 0 {
        return lst[0], nil
    }
    return nil, nil
}

{% if not is_view -%}
/*
func (i *{{struct_name}}) Insert() {
    err := db.{{db_map}}.Insert(i)
    if err != nil{
        game_error.RaiseError(err)
    }
}
*/

// 插入数据，自增长字段将被忽略
func (op *{{op_struct_name}}) Insert(m *{{struct_name}}) (int64, error) {
    return op.InsertTx(db.{{db_sel}}, m)
}

// 插入数据，自增长字段将被忽略
func (op *{{op_struct_name}}) InsertTx(ext sqlx.Ext, m *{{struct_name}}) (int64, error) {
    sql := "{{insert_sql}}"
    result, err := ext.Exec(sql,
    {% for column in column_list -%}
        {% if not column.auto_incr -%}
            m.{{column.field_name}},
        {% endif -%}
    {% endfor -%}
    )
    if err != nil{
        game_error.RaiseError(err)
        return -1, err
    }
    {%if is_auto_incr -%}
        id, _ := result.LastInsertId()
        return id, nil
    {% else -%}
        affected, _ := result.LastInsertId()
        return affected, nil
    {% endif -%}
}

/*
func (i *{{struct_name}}) Update() {
    _,err := db.{{db_map}}.Update(i)
    if err != nil{
        game_error.RaiseError(err)
    }
}
*/

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *{{op_struct_name}}) Update(m *{{struct_name}}) (error) {
    return op.UpdateTx(db.{{db_sel}}, m)
}

// 用主键(属性)做条件，更新除主键外的所有字段
func (op *{{op_struct_name}}) UpdateTx(ext sqlx.Ext, m *{{struct_name}}) (error) {
    sql := `{{update_sql}}`
    _, err := ext.Exec(sql,
    {% for column in column_list -%}
        {% if not column.is_pk -%}
            m.{{column.field_name}},
        {% endif -%}
    {% endfor -%}

        {% for field in primary_field -%}
            m.{{field}},
        {% endfor -%}
    )

    if err != nil{
        game_error.RaiseError(err)
        return err
    }

    return nil
}

// 用主键做条件，更新map里包含的字段名
func (op *{{op_struct_name}}) UpdateWithMap({{primary_key_params}}, m map[string]interface{}) (error) {
    return op.UpdateWithMapTx(db.{{db_sel}}, {{primary_key_param_names}}, m)
}

// 用主键做条件，更新map里包含的字段名
func (op *{{op_struct_name}}) UpdateWithMapTx(ext sqlx.Ext, {{primary_key_params}}, m map[string]interface{}) (error) {

    sql := `update {{table_name}} set %s where 1=1 {% for key in primary_key -%} and {{key}}=? {% endfor -%};`

    var params []interface{}
    var set_sql string
    for k, v := range m{
        set_sql += fmt.Sprintf(" %s=? ", k)
        params = append(params, v)
    }
	params = append(params, {{primary_key_param_names}})
    _, err := ext.Exec(fmt.Sprintf(sql, set_sql), params...)
    return err
}

/*
func (i *{{struct_name}}) Delete(){
    _,err := db.{{db_map}}.Delete(i)
    if err != nil{
        game_error.RaiseError(err)
    }
}
*/
// 根据主键删除相关记录
func (op *{{op_struct_name}}) Delete({{primary_key_params}}) error{
    return op.DeleteTx(db.{{db_sel}}, {{primary_key_param_names}})
}

// 根据主键删除相关记录,Tx
func (op *{{op_struct_name}}) DeleteTx(ext sqlx.Ext, {{primary_key_params}}) error{
    sql := `delete from {{table_name}} where 1=1
        {% for key in primary_key -%}
and {{key}}=?
        {% endfor -%}
`
    _, err := ext.Exec(sql, 
        {% for key in primary_key -%}
        {{key}},
        {% endfor -%}
        )
    return err
}

{% if has_player_id -%}
// 根据玩家id删除相关记录
func (op *{{op_struct_name}}) DeleteByPlayerId(player_id int) (int64, error){
    return op.DeleteByPlayerIdTx(db.{{db_sel}}, player_id)
}

func (op *{{op_struct_name}}) DeleteByPlayerIdTx(ext sqlx.Ext, player_id int) (int64, error){
    sql := "delete from {{table_name}} where player_id=?"
    result, err := ext.Exec(sql, player_id) 
    if err != nil {
        return -1, err
    }
    return result.RowsAffected()
}

{% endif -%} 

{% endif -%}

// 返回符合查询条件的记录数
func (op *{{op_struct_name}}) CountByMap(m map[string]interface{}) int64 {

    var params []interface{}
    sql := `select count(*) from {{table_name}} where 1=1 `
    for k, v := range m{
        sql += fmt.Sprintf(" and  %s=? ",k)
        params = append(params, v)
    }
    count := int64(-1)
    err := db.{{db_sel}}.Get(&count, sql, params...)
    if err != nil {
        game_error.RaiseError(err)
    }
    return count
}

func (op *{{op_struct_name}}) DeleteByMap(m map[string]interface{})(int64, error){
	return op.DeleteByMapTx(db.{{db_sel}}, m)
}

func (op *{{op_struct_name}}) DeleteByMapTx(ext sqlx.Ext, m map[string]interface{}) (int64, error){
	var params []interface{}
	sql := "delete from {{table_name}} where 1=1 "
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

{% if has_player_id -%}
// 返回指定玩家ID的记录数量
func (op *{{op_struct_name}}) CountByPlayerId(player_id int) (int64){
    return op.CountByMap(map[string]interface{}{
        "player_id" : player_id,
    })
}

func (op *{{op_struct_name}}) QueryByPlayerId(player_id int) ([]*{{struct_name}}){
    sql := "select * from {{table_name}} where player_id=?"
    result := []*{{struct_name}}{}
    err := db.{{db_sel}}.Select(&result, sql, player_id) 
    if err != nil{
        game_error.RaiseError(err)
    }
    return result
}

{% endif -%} 

{% else -%}

type {{cache_struct_name}} struct{
    {% if primary_key_length == 1 -%}
    objMap map[{{primary_field_type[0]}}]*{{struct_name}}
    {% elif primary_key_length == 0 -%}
	objMap map[string]*{{struct_name}}
    {% else-%}
    objMap map[string]*{{struct_name}}
    {% endif -%}
    objList []*{{struct_name}}
}

var {{cache_name}} = &{{cache_struct_name}}{}

func (c *{{cache_struct_name}}) LoadAll() {
    sql := "select * from {{table_name}}"
	c.objList = make([]*{{struct_name}}, 0)
    err := db.{{db_sel}}.Select(&c.objList,sql)
    if err != nil{
        log.Fatal(err.Error())
    }
    {% if primary_key_length == 1 -%}
        c.objMap = make(map[{{primary_field_type[0]}}]*{{struct_name}})
    {% else-%}
        c.objMap = make(map[string]*{{struct_name}})
    {% endif -%}

    log.Debug("Load all {{table_name}} success %v", len(c.objList))
    for _,v := range c.objList{
        {% if primary_key_length == 1 -%}
        c.objMap[v.{{primary_field[0]}}] = v
        {% else -%}
        var key string
        {% for k in primary_field -%}
        key += fmt.Sprintf("%v-",v.{{k}})
        {% endfor -%}
        c.objMap[key] = v
        {% endif -%}
    }
}

func (c *{{cache_struct_name}}) All() []*{{struct_name}}{
    return c.objList
}

func (c *{{cache_struct_name}}) Count() int {
    return len(c.objList)
}

func (c *{{cache_struct_name}}) Get({{primary_key_params}}) (*{{struct_name}}, bool){
    {% if primary_key_length == 1 -%}
    key := {{primary_key[0]}}
    {% else -%}
    var key string
    {% for k in primary_key -%}
    key += fmt.Sprintf("%v-",{{k}})
    {% endfor -%}
    {% endif -%}
    v,ok :=  c.objMap[key]
    return v,ok
}

// 仅限运营后台实时刷新服务器数据用
func (c *{{cache_struct_name}}) Update(v *{{struct_name}}){
    {% if primary_key_length == 1 -%}
    key := v.{{primary_field[0]}}
    {% else -%}
    var key string
    {% for k in primary_field -%}
        key += fmt.Sprintf("%v-",v.{{k}})
    {% endfor -%}
    {% endif -%}
    c.objMap[key] = v
}

{% endif -%}
