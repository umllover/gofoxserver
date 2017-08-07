package db

import (
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/lovelly/leaf/log"
	"gopkg.in/redis.v4"
)

const driverName = "mysql"

var (
	once      = &sync.Once{}
	BaseDB    *sqlx.DB
	DB        *sqlx.DB
	StatsDB   *sqlx.DB
	AccountDB *sqlx.DB
	RdsDB     *redis.Client // go-reids
)

type IDBCnf interface {
	GetAccoutDSN() string
	GetBaseDSN() string
	GetUserDSN() string
	GetStatsDSN() string
	GetBaseDBMaxOpen() int
	GetBaseDBMaxIdle() int
	GetUserDBMaxOpen() int
	GetUserDBMaxIdle() int
	GetStatsDBMaxOpen() int
	GetStatsDBMaxIdle() int
	GetAccountDBMaxOpen() int
	GetAccountDBMaxIdle() int
	GetRedisAddr() string
}

func InitDB(cnf IDBCnf) {
	once.Do(func() {
		BaseDB = initSqlxDB(cnf.GetBaseDSN(), "[BASE_DB] -> ", cnf.GetBaseDBMaxOpen(), cnf.GetBaseDBMaxIdle())
		DB = initSqlxDB(cnf.GetUserDSN(), "[USER DB] -> ", cnf.GetUserDBMaxOpen(), cnf.GetUserDBMaxIdle())
		StatsDB = initSqlxDB(cnf.GetStatsDSN(), "[STATS_DB] ->", cnf.GetStatsDBMaxOpen(), cnf.GetStatsDBMaxIdle())
		AccountDB = initSqlxDB(cnf.GetAccoutDSN(), "[STATS_DB] ->", cnf.GetAccountDBMaxOpen(), cnf.GetAccountDBMaxIdle())
		log.Debug("Init DB success.")

		UpdateDB()

		//RdsDB = redis.NewClient(&redis.Options{
		//	Addr:     cnf.GetRedisAddr(),
		//	Password: "",
		//	DB:       0,
		//})
	})
}

func initSqlxDB(dbConfig, logHeader string, maxOpen, maxIdle int) *sqlx.DB {
	log.Debug("dbConfig: %s, logHeader: %s, maxOpen: %d, maxIdle: %d", dbConfig, logHeader, maxOpen, maxIdle)
	db := sqlx.MustConnect(driverName, dbConfig)
	db.SetMaxOpenConns(maxOpen)
	db.SetMaxIdleConns(0)
	return db
}
