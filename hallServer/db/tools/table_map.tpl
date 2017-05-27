package {{package_name}}

//This file is generate by scripts,don't edit it
//

import (

)

func InitTableMap() {
	/*
	    {% for table in table_list -%}
	    {% if table.is_auto_incr -%}
		db.{{table.db_map}}.AddTableWithName({{table.struct_name}}{}, "{{table.table_name}}").SetKeys(true, {{table.primary_key}})
	    {% else -%}
		db.{{table.db_map}}.AddTableWithName({{table.struct_name}}{}, "{{table.table_name}}").SetKeys(false,{{table.primary_key}})
	    {% endif -%}
	    {% endfor -%}
	*/
}
