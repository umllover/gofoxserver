package db

import (
	"github.com/lovelly/leaf/log"
	"github.com/jmoiron/sqlx"
)

const (
	LOCK_ID = 1
)

// 数据库增量更新
func UpdateDB() error {
	log.Debug("Start update db.")
	defer func() {
		log.Debug("Update db end.")
	}()

	//var err error
	// user db
	r, err := DB.Exec("INSERT  INTO version_locker(id) VALUES(?)", LOCK_ID)
	if err != nil {
		log.Debug("%s", err.Error())
		return  err
	}
	row, err  := r.RowsAffected()
	if err != nil {
		log.Debug("%s", err.Error())
		return  err
	}
	if row <= 0 {
		log.Debug("%s", err.Error())
		return err
	}
	log.Debug("get userdb lock sucess")

	err = UpdateSingle(DB, userUpdateSql)
	if err != nil {
		return err
	}

	r, err = DB.Exec("DELETE  FROM version_locker WHERE  id = ?", LOCK_ID)
	if err != nil {
		log.Debug("%s", err.Error())
		return  err
	}
	log.Debug("release userdb lock sucess")


	// stats db
	r, err = StatsDB.Exec("INSERT  INTO version_locker(id) VALUES(?)", LOCK_ID)
	if err != nil {
		log.Debug("%s", err.Error())
		return  err
	}
	row, err  = r.RowsAffected()
	if err != nil {
		log.Debug("%s", err.Error())
		return  err
	}
	if row <= 0 {
		log.Debug("%s", err.Error())
		return err
	}
	log.Debug("get statsdb lock sucess")

	err = UpdateSingle(StatsDB, statsUpdateSql)
	if err != nil {
		return err
	}

	r, err = StatsDB.Exec("DELETE  FROM version_locker WHERE  id = ?", LOCK_ID)
	if err != nil {
		log.Debug("%s", err.Error())
		return  err
	}
	log.Debug("release statsdb lock sucess")

	return nil
}

func UpdateSingle(inst *sqlx.DB, sqls [][]string) error {
	// id may have other uses?
	log.Debug("enter updateSingle ,len = %d", len(sqls))

	row := inst.QueryRowx("select ver from version where id = 1;")
	ver := 0
	err := row.Scan(&ver)
	if err != nil {
		r := inst.QueryRowx("SHOW TABLES LIKE 'version';")
		have := ""
		r.Scan(&have)
		if have == "version" {
			log.Error("query ver encounter a error.Error: %s", err.Error())
			return err
		}
	}

	log.Debug("%d", ver)

	if len(sqls)-1 <= ver {
		return nil
	}

	// 需要更新的部分
	updateSqls := sqls[ver+1:]
	for i,_ := range updateSqls {
		for j,_ := range updateSqls[i] {
			log.Debug(updateSqls[i][j])
		}
	}

	tx, err := inst.Begin()
	if err != nil {
		log.Error("Begin tx encounter a error.Error:%s", err.Error())
		return err
	}
	for _, updateSql := range updateSqls {
		for _, updateSql_ := range updateSql {
			log.Debug("Exec sql.Sql: %s", updateSql_)
			_, err = tx.Exec(updateSql_)
			if err != nil {
				log.Error("Exec tx encounter a error.Error: %s Sql:%s", err.Error(),  updateSql_)
				err1 := tx.Rollback()
				if err1 != nil {
					log.Error("Rollback encounter a error.Error: %s", err.Error())
				}
				return err
			}
		}
	}
	// 刷新version表
	_, err = tx.Exec("INSERT INTO version (id, ver) VALUES(1, ?)  ON DUPLICATE KEY UPDATE ver=?;", len(sqls)-1, len(sqls)-1)
	if err != nil {
		log.Error("Update version field[ver] encounter a error.Error: %s  ver:%v", err.Error(), len(sqls)-1)
		err1 := tx.Rollback()
		if err1 != nil {
			log.Error("Rollback encounter a error.Error: %s", err.Error())
		}
		return err
	}
	err = tx.Commit()
	if err != nil {
		log.Error("Commit encounter a error.Error: %s", err.Error())
		return err
	}

	return nil
}
