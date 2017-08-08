package db

import (
	"sync"

	redis "gopkg.in/redis.v4"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/lovelly/leaf/log"
)

const driverName = "mysql"

var (
	once    = &sync.Once{}
	BaseDB  *sqlx.DB
	DB      *sqlx.DB
	StatsDB *sqlx.DB
	RdsDB   *redis.Client // go-reids
)

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
	GetRedisAddr() string
	GetRedisPwd() string
}

func InitDB(cnf IDBCnf) {
	once.Do(func() {
		BaseDB = initSqlxDB(cnf.GetBaseDSN(), "[BASE_DB] -> ", cnf.GetBaseDBMaxOpen(), cnf.GetBaseDBMaxIdle())
		DB = initSqlxDB(cnf.GetUserDSN(), "[USER_DB] -> ", cnf.GetUserDBMaxOpen(), cnf.GetUserDBMaxIdle())
		StatsDB = initSqlxDB(cnf.GetStatsDSN(), "[STATS_DB] ->", cnf.GetStatsDBMaxOpen(), cnf.GetStatsDBMaxIdle())
		log.Debug("Init DB success.")
	})

	//RdsDB = redis.NewClient(&redis.Options{
	//	Addr:     cnf.GetRedisAddr(),
	//	Password: cnf.GetRedisPwd(),
	//	DB:       0,
	//})
	//
	//ret := RdsDB.Ping()
	//if ret.Err() != nil {
	//	log.Fatal("connect redis error ")
	//}

}

func initSqlxDB(dbConfig, logHeader string, maxOpen, maxIdle int) *sqlx.DB {
	log.Debug("dbConfig: %s, logHeader: %s, maxOpen: %d, maxIdle: %d", dbConfig, logHeader, maxOpen, maxIdle)
	db := sqlx.MustConnect(driverName, dbConfig)
	db.SetMaxOpenConns(maxOpen)
	db.SetMaxIdleConns(0)
	return db
}
