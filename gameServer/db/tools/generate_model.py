#!/usr/bin/env python
#-* coding:UTF-8 -*


import os
import re
import pymysql
from pymysql.cursors import DictCursor
from jinja2 import Template
import json
import platform
import string
import sys
reload(sys)
sys.setdefaultencoding('utf8')

HERE = os.path.dirname(os.path.abspath(__file__))
ROOT_DIR = os.path.dirname(HERE)
BASE_MODEL_DIR = os.path.join(ROOT_DIR, 'model', 'base')
USER_MODEL_DIR = os.path.join(ROOT_DIR, 'model')
STATS_MODEL_DIR = os.path.join(ROOT_DIR, 'model', 'stats')
MODEL_DIRS = [BASE_MODEL_DIR, USER_MODEL_DIR, STATS_MODEL_DIR]
for d in MODEL_DIRS :
    if not os.path.exists(d):
        os.makedirs(d)
MODEL_TEMPLATE_FILE = os.path.join(HERE, "model.tpl")
MODEL_MANAGER_TEMPLATE_FILE = os.path.join(HERE, "model_manager.tpl")
FFJSON_TEMPLATE_FILE = os.path.join(HERE, "model_ffjson.tpl")
TABLE_TEMPLATE_FILE = os.path.join(HERE, "table_map.tpl")
TABLE_LOADER_TEMPLATE_FILE = os.path.join(HERE, "table_loader.tpl")
CONFIG_FILE = os.path.join(ROOT_DIR, "../gameApp/gameServer.json")


def to_camel_case(snake_str):
    components = snake_str.split('_')
    first_char = components[0][0]
    if first_char in ('0', '1', '2', '3', '4', '5', '6', '7', '8', '9'):
        first_char = 'A' + first_char
    name = first_char.upper() + components[0][1:] + \
        "".join(x.title() for x in components[1:])
    return name


go_type_dict = {
    'smallint': 'int8',
    'tinyint': 'int8',
    'varchar': 'string',
    'int': 'int',
    'decimal': 'float64',
    'timestamp': '*time.Time',
    'bigint': 'int64',
    'char': 'string',
    'float': 'float64',
    'text': 'string',
    'longtext': 'string',
    'date': 'string',
    'double': 'float64'
}


def get_go_type(db_name, db_type, data_type):
    ret = re.search("unsigned",data_type)
    if (db_type == "int" and db_name.endswith('time')) or db_name == "v":
        return 'int64'
    t  = go_type_dict.get(db_type)
    if t == 'int64' and ret != None:
        return 'uint64'
    return t


def generate_gen(model_dir):
    print 'generate_gen for model_dir %s' % model_dir
    # github.com/clipperhouse/gen
    os.chdir(model_dir)
    os.system("gen")
    old_dir = os.getcwd()
    os.chdir(old_dir)
    
def rename_go1Togo(path):
    all_files = os.listdir(path)
    for filename in all_files:
        fullname = os.path.join(path, filename)
        if os.path.isfile(fullname) and filename.split('.')[-1] == "go1":
            newfilename = fullname.rsplit('.', 1)[0]
            print fullname, newfilename
            os.rename(fullname, newfilename + ".go")

def render(conn, db_name, db_map, model_dir, package_name, is_base_db):
    cur = conn.cursor()
    sql = """ SELECT TABLE_TYPE, TABLE_NAME,TABLE_COMMENT FROM TABLES
    WHERE table_schema = %s
    """

    cur.execute(sql, [db_name])
    result = cur.fetchall()

    if result:
        if not os.path.exists(model_dir):
            os.mkdir(model_dir)
        os.system("rm %s/*" % model_dir)
    all_render_dict = {'table_list': []}

    with open(MODEL_TEMPLATE_FILE, "r") as f:
        tpl = f.read()
    with open(FFJSON_TEMPLATE_FILE, "r") as f:
        ffjson_tpl = f.read()
    with open(TABLE_TEMPLATE_FILE, "r") as f:
        table_map_tpl = f.read()
    #with open(MODEL_MANAGER_TEMPLATE_FILE, "r") as f:
    #    model_manager_tpl = f.read()
    with open(TABLE_LOADER_TEMPLATE_FILE, "r") as f:
        table_loader_tpl = f.read()

    # for table_mao
    table_list = []

    for row in result:
        table_name = row['TABLE_NAME']
        is_view = row['TABLE_TYPE'] == 'VIEW'
        table_comment = row['TABLE_COMMENT']
        if "[logic server no use]" in table_comment:
            continue

        sql = """
        SELECT
        COLUMN_NAME,DATA_TYPE, COLUMN_COMMENT,
        COLUMN_DEFAULT,COLUMN_KEY,COLUMN_TYPE,EXTRA
        FROM COLUMNS
        WHERE TABLE_NAME = %s  and TABLE_SCHEMA = %s
        """

        cur.execute(sql, [table_name, db_name])

        columns = cur.fetchall()

        struct_name = to_camel_case(table_name)
        op_struct_name = struct_name[0].lower() + struct_name[1:] + 'Op'
        cache_struct_name = struct_name[0].lower() + struct_name[1:] + 'Cache'
        op_name = op_struct_name[0].upper() + op_struct_name[1:]
        cache_name = cache_struct_name[0].upper() + cache_struct_name[1:]
        has_player_id = False
        primary_key = []
        primary_field = []
        primary_field_type = []
        column_list = []
        primary_key_param_list = []
        primary_keys = []
        for c in columns:
            column_name = c['COLUMN_NAME']
            field_name = to_camel_case(column_name)
            column_type = c['DATA_TYPE']
            data_type = c['COLUMN_TYPE']
            go_type = get_go_type(column_name, column_type, data_type)
            column_key = c['COLUMN_KEY']
            column_comment = c['COLUMN_COMMENT']
            auto_incr = (c['EXTRA'] == 'auto_increment')
            if column_name == 'player_id':
                has_player_id = True
            if column_key == 'PRI':
                primary_key.append(column_name)
                primary_field.append(field_name)
                primary_field_type.append(go_type)
                primary_key_param_list.append((column_name, go_type))
                primary_keys.append(column_name +' ' + go_type)		
            column_list.append({
                'field_name': field_name,
                'type': go_type,
                'name': column_name,
                'comment': column_comment.strip(),
                'auto_incr' : auto_incr,
                'is_pk' : column_key == 'PRI',
            })

        primary_key_params = ','.join(
            ['%s %s' % (k[0], k[1]) for k in primary_key_param_list])
        primary_key_param_names = ','.join([k[0] for k in primary_key_param_list])

        get_by_pk_sql = 'select * from %s where %s ' % (
            table_name, ' and '.join(['%s =:%s' % (k, k) for k in primary_key]))
        get_by_pk_sql2 = 'select * from %s where %s ' % (
            table_name, ' and '.join(['%s=?' % k for k in primary_key]))
        get_by_pk_result = struct_name if len(
            primary_key) == 1 else '[]%s' % struct_name

        insert_sql = 'insert into %s(%s) values(%s)' % \
            (table_name, \
            ','.join([c['COLUMN_NAME']  for c in columns if c['EXTRA'] != 'auto_increment' ]),
            ','.join(['?' for c in columns if c['EXTRA'] != 'auto_increment']))

        insert_update_sql = 'insert into %s(%s) values(%s) ON DUPLICATE KEY UPDATE ' % \
            (table_name, \
            ','.join([c['COLUMN_NAME']  for c in columns if c['EXTRA'] != 'auto_increment' ]),
            ','.join(['?' for c in columns if c['EXTRA'] != 'auto_increment']))			

        update_sql = 'update %s set %s where %s' % \
                (table_name, \
                ','.join([c['COLUMN_NAME']+'=?'  for c in columns if c['COLUMN_KEY'] != 'PRI' ]), \
                ' and '.join(['%s=?' % k for k in primary_key]))

        db_sel = "DB"
        if package_name == "stats":
            db_sel = "StatsDB"
        elif package_name == "base":
            db_sel = "BaseDB"

        render_dict = {
            'is_view' : is_view,
            'insert_sql' : insert_sql,
            'update_sql' : update_sql,
			'insert_update_sql':insert_update_sql,
            'table_name': table_name,
            'struct_name': struct_name,
            'op_struct_name': op_struct_name,
            'op_name': op_name,
            'cache_struct_name' : cache_struct_name,
            'cache_name': cache_name,
            'db_map': db_map,
            'package_name': package_name,
            'table_comment': table_comment.strip(),
            'has_player_id': has_player_id,
            'primary_key': primary_key,
            'primary_field': primary_field,
            'primary_field_type': primary_field_type,
            'primary_key_length': len(primary_key),
            'primary_key_params': primary_key_params,
            'primary_key_param_names': primary_key_param_names,
            'primary_key_param_list': primary_key_param_list,
			'primary_keys': primary_keys,
            'get_by_pk_sql': get_by_pk_sql,
            'get_by_pk_sql2': get_by_pk_sql2,
            'get_by_pk_result': get_by_pk_result,
            'column_list': column_list,
            'is_base_db': is_base_db,
            #'db_name' : db_name,
			"db_sel" : db_sel,
        }

        template = Template(tpl)
        model_file_path = os.path.join(model_dir, '%s.go1' % table_name)
        with open(model_file_path, 'w') as f:
            f.write(template.render(**render_dict))
        os.system("goimports -w %s" % model_file_path)
        all_render_dict['table_list'].append(render_dict)
        all_render_dict['package_name'] = package_name
        print 'write %s success' % model_file_path

        table_list.append({
            'table_name': table_name,
            'db_map': db_map,
            'struct_name': struct_name,
            'is_auto_incr': 'auto_increment' in [c['EXTRA'] for c in columns],
            'primary_key': ','.join(
                ['"%s"' % to_camel_case(c1) for c1 in
                 [c['COLUMN_NAME']
                  for c in columns if c['COLUMN_KEY'] == 'PRI']
                 ]),
        })

    old_dir = os.getcwd()
    os.chdir(model_dir)
    all_model_file_path = os.path.join(model_dir, "model.go")

    #template = Template(ffjson_tpl)
    #with open(all_model_file_path, 'w') as f:
    #    f.write(template.render(**all_render_dict))

    #windows had't /dev/null
    if platform.system().upper() == "windows".upper() :
        nullFile = os.path.join(os.environ["TEMP"], "null.ffjson")
        os.system("ffjson %s > %s" % (all_model_file_path, nullFile))
        os.system("rm %s " % nullFile)
    else:
        os.system("ffjson %s > /dev/null" % all_model_file_path)

    os.system("rm %s " % all_model_file_path)
    rename_go1Togo(model_dir)
    #os.system("rename 's/\.go1/\.go/' *")

    # generate table_map
    table_map_file_path = os.path.join(model_dir, "table_map.go")
    template = Template(table_map_tpl)
    with open(table_map_file_path, 'w') as f:
        f.write(
            template.render(table_list=table_list,
                            package_name=package_name,
                            )
        )
    os.system("goimports -w %s" % table_map_file_path)

    os.system("go install ./")
    os.chdir(old_dir)

    # base db data manager
    if is_base_db:
        # for render_dict in all_render_dict['table_list']:
        #     model_file_path = os.path.join(
        #         model_dir, '%s_manager.go' % render_dict['table_name'])
        #     template = Template(model_manager_tpl)
        #     with open(model_file_path, 'w') as f:
        #         f.write(template.render(**render_dict))
        #     os.system("goimports -w %s" % model_file_path)
        # for table loader
        table_loader_path = os.path.join(model_dir, 'table_loader.go')
        template = Template(table_loader_tpl)
        with open(table_loader_path, 'w') as f:
            f.write(
                template.render(table_list=table_list, package_name=package_name))
        os.system("goimports -w %s" % table_loader_path)

    cur.close()
    conn.close()

    #generate_gen(model_dir)


def run():
    with open(CONFIG_FILE, "r") as f:
        config = json.loads(f.read())
    print "config:", config

    base_db_conn = pymysql.connect(host=config['BaseDbHost'],
                                   port=config['BaseDbPort'],
                                   charset='utf8',
                                   user=config['BaseDbUsername'],
                                   passwd=config['BaseDbPassword'],
                                   db='information_schema', cursorclass=DictCursor)
    render(base_db_conn, config['BaseDbName'],
           'BaseDBMap', BASE_MODEL_DIR, 'base', True)

    user_db_conn = pymysql.connect(host=config['UserDbHost'],
                                   port=config['UserDbPort'],
                                   charset='utf8',
                                   user=config['UserDbUsername'],
                                   passwd=config['UserDbPassword'],
                                   db='information_schema', cursorclass=DictCursor)
    render(user_db_conn, config['UserDbName'],
           'DBMap', USER_MODEL_DIR, 'model', False)

    stats_db_conn = pymysql.connect(host=config['StatsDbHost'],
                                   port=config['StatsDbPort'],
                                   charset='utf8',
                                   user=config['StatsDbUsername'],
                                   passwd=config['StatsDbPassword'],
                                   db='information_schema', cursorclass=DictCursor)
    render(stats_db_conn, config['StatsDbName'],
           'StatsDBMap', STATS_MODEL_DIR, 'stats', False)


if __name__ == "__main__":
    #print to_camel_case('2_day_liveness')
    #print to_camel_case('30_day_liveness')
    run()
