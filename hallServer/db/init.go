package db

import (
	"sync"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/lovelly/leaf/log"
	"github.com/jmoiron/sqlx"
)

const driverName = "mysql"

var once = &sync.Once{}
var BaseDB *sqlx.DB
var DB *sqlx.DB
var StatsDB *sqlx.DB

type IDBCnf interface {
	GetBaseDSN()string
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
		DB = initSqlxDB(cnf.GetUserDSN(), "[USER DB] -> ", cnf.GetUserDBMaxOpen(), cnf.GetUserDBMaxIdle())
		StatsDB = initSqlxDB(cnf.GetStatsDSN(), "[STATS_DB] ->", cnf.GetStatsDBMaxOpen(), cnf.GetStatsDBMaxIdle())
		UpdateDB()
		log.Debug("Init DB success.")
	})
}

func initSqlxDB(dbConfig, logHeader string, maxOpen, maxIdle int) *sqlx.DB {
	fmt.Println(dbConfig, logHeader, maxOpen, maxIdle)
	db := sqlx.MustConnect(driverName, dbConfig)
	db.SetMaxOpenConns(maxOpen)
	db.SetMaxIdleConns(maxIdle)
	return db
}


