package db

import (
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/lovelly/leaf/log"
)

const driverName = "mysql"

var once = &sync.Once{}
var BaseDB *sqlx.DB
var DB *sqlx.DB
var StatsDB *sqlx.DB

type IDBCnf interface {
	GetBaseDSN() string
	GetUserDSN() string
	GetStatsDSN() string
	GetBaseDBMaxOpen() int
	GetBaseDBMaxIdle() int
	GetUserDBMaxOpen() int
	GetUserDBMaxIdle() int
	GetStatsDBMaxOpen() int
	GetStatsDBMaxIdle() int
}

func InitDB(cnf IDBCnf) {
	once.Do(func() {
		BaseDB = initSqlxDB(cnf.GetBaseDSN(), "[BASE_DB] -> ", cnf.GetBaseDBMaxOpen(), cnf.GetBaseDBMaxIdle())
		DB = initSqlxDB(cnf.GetUserDSN(), "[USER_DB] -> ", cnf.GetUserDBMaxOpen(), cnf.GetUserDBMaxIdle())
		StatsDB = initSqlxDB(cnf.GetStatsDSN(), "[STATS_DB] ->", cnf.GetStatsDBMaxOpen(), cnf.GetStatsDBMaxIdle())
		log.Debug("Init DB success.")
	})
}

func initSqlxDB(dbConfig, logHeader string, maxOpen, maxIdle int) *sqlx.DB {
	log.Debug("dbConfig: %s, logHeader: %s, maxOpen: %d, maxIdle: %d", dbConfig, logHeader, maxOpen, maxIdle)
	db := sqlx.MustConnect(driverName, dbConfig)
	db.SetMaxOpenConns(maxOpen)
	db.SetMaxIdleConns(maxIdle)
	return db
}
