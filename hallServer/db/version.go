package db

import (
	"database/sql"
	"fmt"
	"os/exec"

	"mj/hallServer/conf"

	"github.com/jmoiron/sqlx"
	"github.com/lovelly/leaf/log"
)

const (
	LOCK_ID = 1
)

func UpdateDB() {
	err, up := updateDB()
	if up && conf.Test {
		log.Debug("重新生成配置中，请稍后。。。")
		RanderDB("../db/tools/")
		RanderDB("../../gameServer/db/tools/")
		log.Fatal("更新数据成功，请重启。。。")
	}
	if err != nil {
		log.Fatal("InitDB: %s", err.Error())
	}
}

func RanderDB(path string) {
	cmd := exec.Command("python", "generate_model.py")
	cmd.Dir = path
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Error("RanderDB error :%s", err.Error())
	} else {
		log.Debug("RanderDB out: %s", string(out))
	}
}

// 数据库增量更新
func updateDB() (err error, up bool) {
	log.Debug("Start update db.")
	defer func() {
		_, err := DB.Exec("DELETE  FROM version_locker WHERE  id = ?", LOCK_ID)
		if err != nil {
			log.Debug("%s", err.Error())
		}

		_, err = StatsDB.Exec("DELETE  FROM version_locker WHERE  id = ?", LOCK_ID)
		if err != nil {
			log.Debug("%s", err.Error())
			return
		}
	}()

	//var err error
	// user db
	DB.Exec(`CREATE TABLE if not exists  version_locker (id int(11) NOT NULL,PRIMARY KEY (id)) ENGINE=InnoDB DEFAULT CHARSET=utf8;`)
	DB.Exec(`CREATE TABLE if not exists version (ver int(11) NOT NULL,id int(11) NOT NULL,PRIMARY KEY (id)) ENGINE=InnoDB DEFAULT CHARSET=utf8`)

	var r sql.Result
	r, err = DB.Exec("INSERT  INTO version_locker(id) VALUES(?)", LOCK_ID)
	if err != nil {
		log.Debug("%s", err.Error())
		return err, false
	}
	row, err := r.RowsAffected()
	if err != nil {
		log.Debug("%s", err.Error())
		return err, false
	}
	if row <= 0 {
		log.Debug("%s", err.Error())
		return err, false
	}
	log.Debug("get userdb lock sucess")

	err, up = UpdateSingle(DB, userUpdateSql)
	if err != nil {
		return err, false
	}

	log.Debug("release userdb lock sucess")

	// stats db
	StatsDB.Exec(`CREATE TABLE if not exists  version_locker (id int(11) NOT NULL,PRIMARY KEY (id)) ENGINE=InnoDB DEFAULT CHARSET=utf8;`)
	StatsDB.Exec(`CREATE TABLE if not exists version (ver int(11) NOT NULL,id int(11) NOT NULL,PRIMARY KEY (id)) ENGINE=InnoDB DEFAULT CHARSET=utf8`)
	r, err = StatsDB.Exec("INSERT  INTO version_locker(id) VALUES(?)", LOCK_ID)
	if err != nil {
		log.Debug("%s", err.Error())
		return err, up
	}
	row, err = r.RowsAffected()
	if err != nil {
		log.Debug("%s", err.Error())
		return err, up
	}
	if row <= 0 {
		log.Debug("%s", err.Error())
		return err, up
	}
	log.Debug("get statsdb lock sucess")

	err, up = UpdateSingle(StatsDB, statsUpdateSql)
	if err != nil {
		return err, up
	}

	log.Debug("release statsdb lock sucess")

	return nil, up
}

func UpdateSingle(inst *sqlx.DB, sqls [][]string) (error, bool) {
	// id may have other uses?
	log.Debug("enter updateSingle ,len = %d", len(sqls))

	var ret []int
	err := inst.Select(&ret, "select ver from version where id = 1;")
	if err != nil {
		/*	r := inst.QueryRowx("SHOW TABLES LIKE 'version';")
			have := ""
			r.Scan(&have)
			if have == "version" {
				log.Error("query ver encounter a error.Error: %s", err.Error())*/
		return err, false
		//}
	}
	var ver int
	if len(ret) > 0 {
		ver = ret[0]
	}

	log.Debug("sql version :%d", ver)

	if len(sqls) < ver {
		log.Debug("sql lend %d", len(sqls))
		return nil, false
	}

	// 需要更新的部分
	updateSqls := sqls[ver:]
	if err != nil {
		log.Error("Begin tx encounter a error.Error:%s", err.Error())
		return err, false
	}
	for newIndex, updateSql := range updateSqls {
		tx, err := inst.Begin()
		for _, updateSql_ := range updateSql {
			log.Debug("Exec sql.Sql: %s", updateSql_)
			halder, err := tx.Prepare(updateSql_)
			if err != nil {
				log.Error("Exec tx encounter a error.Error: %s Sql:%s", err.Error(), updateSql_)
				err1 := tx.Rollback()
				if err1 != nil {
					log.Error("Rollback encounter a error.Error: %s", err.Error())
				}
				return err, false
			}
			halder.Exec()
		}

		err = tx.Commit()

		// 刷新version表
		newv := ver + newIndex + 1
		_, err = inst.Exec(fmt.Sprintf("INSERT INTO version (id, ver) VALUES(1, %d)  ON DUPLICATE KEY UPDATE ver=%d ;", newv, newv))
		if err != nil {
			return err, false
		}

		if err != nil {
			log.Error("Commit encounter a error.Error: %s", err.Error())
			return err, false
		}
	}

	return nil, true
}
