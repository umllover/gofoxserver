package {{package_name}}

//This file is generate by scripts,don't edit it
//

import (
	"mj/hallServer/db"
    "github.com/lovelly/leaf/log"
)

func LoadBaseData() {
	var start = time.Now().UnixNano()
    {% for table in table_list -%}
    {{table.struct_name}}Cache.LoadAll()
    {% endfor -%}

    {% for table in table_list -%}
    db.BaseDataCaches["{{table.struct_name}}"] = {{table.struct_name}}Cache
    {% endfor -%}

	log.Debug("loadBaseData %v  %v %v", {{table_list | length }},  time.Now().UnixNano()-start, "ns")
}
