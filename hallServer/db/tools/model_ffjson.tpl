package {{package_name}}

{% for table in table_list -%}
//{{table.table_comment}}
type {{table.struct_name}} struct {
    {% for column in table.column_list -%}
    // {{column.comment}}
    {% if column.name == "update_time" -%}
    {{column.field_name}} {{column.type}} `db:"{{column.name}}" json:"-"` 
    {% else -%}
    {{column.field_name}} {{column.type}} `db:"{{column.name}}" json:"{{column.name}}"` 
    {% endif -%}
    {% endfor -%}
}
{% endfor -%}
