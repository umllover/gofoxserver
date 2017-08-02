package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	. "mj/common/cost"

	"github.com/lovelly/leaf/log"
)

var Server struct {
	LogLevel    string
	LogPath     string
	WSAddr      string
	CertFile    string
	KeyFile     string
	TCPAddr     string
	ListenAddr  string
	HttpAddr    string
	MaxConnNum  int
	ConsolePort int
	ProfilePath string

	AccountDbHost     string
	AccountDbPort     int
	AccountDbName     string
	AccountDbUsername string
	AccountDbPassword string
	BaseDbHost        string
	BaseDbPort        int
	BaseDbName        string
	BaseDbUsername    string
	BaseDbPassword    string
	UserDbHost        string
	UserDbPort        int
	UserDbName        string
	UserDbUsername    string
	UserDbPassword    string
	StatsDbHost       string
	StatsDbPort       int
	StatsDbName       string
	StatsDbUsername   string
	StatsDbPassword   string
	ConsulAddr        string

	NsqdAddrs       []string
	NsqLookupdAddrs []string
	PdrNsqdAddr     string

	ConnAddrs       map[string]string
	AuthServerUrl   string
	PendingWriteNum int
	WatchAddr       string
	NodeId          int
}

func ServerName() string {
	return fmt.Sprintf(HallPrefixFmt, Server.NodeId)
}

func ServerNsqCahnnel() string {
	return fmt.Sprintf(HallCahnnelFmt, Server.NodeId)
}

func Init() {
	data, err := ioutil.ReadFile("./hallServer.json")
	if err != nil {
		log.Fatal("%v", err)
	}
	err = json.Unmarshal(data, &Server)
	if err != nil {
		log.Fatal("%v", err)
	}

}

const (
	default_db_max_open      = 32
	default_db_max_idle      = 2
	default_stat_log_workers = 64
)

type DBConfig struct{}

func (c *DBConfig) GetAccoutDSN() string {
	s := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		Server.AccountDbUsername, Server.AccountDbPassword, Server.AccountDbHost, Server.AccountDbPort, Server.AccountDbName, "parseTime=true&interpolateParams=true")
	return s
}

func (c *DBConfig) GetBaseDSN() string {
	s := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		Server.BaseDbUsername, Server.BaseDbPassword, Server.BaseDbHost, Server.BaseDbPort, Server.BaseDbName, "parseTime=true&interpolateParams=true")
	return s
}

func (c *DBConfig) GetUserDSN() string {
	s := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		Server.UserDbUsername, Server.UserDbPassword, Server.UserDbHost, Server.UserDbPort, Server.UserDbName, "parseTime=true&charset=utf8mb4&interpolateParams=true")
	return s
}

func (c *DBConfig) GetStatsDSN() string {
	s := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		Server.StatsDbUsername, Server.StatsDbPassword, Server.StatsDbHost, Server.StatsDbPort, Server.StatsDbName, "parseTime=true&interpolateParams=true")
	return s
}

func (c *DBConfig) GetBaseDBMaxOpen() int {
	return default_db_max_open
}

func (c *DBConfig) GetBaseDBMaxIdle() int {
	return default_db_max_idle
}

func (c *DBConfig) GetUserDBMaxOpen() int {
	return default_db_max_open
}

func (c *DBConfig) GetUserDBMaxIdle() int {
	return default_db_max_idle
}

func (c *DBConfig) GetStatsDBMaxOpen() int {
	return default_db_max_open
}

func (c *DBConfig) GetStatsDBMaxIdle() int {
	return default_db_max_idle
}

func (c *DBConfig) GetStatsDBWorkers() int {
	return default_stat_log_workers
}

func (c *DBConfig) GetAccountDBMaxIdle() int {
	return default_db_max_idle
}

func (c *DBConfig) GetAccountDBMaxOpen() int {
	return default_db_max_open
}

//consul config
type ConsulConfig struct{}

func (c *ConsulConfig) GetConsulAddr() string {
	return Server.ConsulAddr
}
func (c *ConsulConfig) GetConsulToken() string {
	return ""
}
func (c *ConsulConfig) GetConsulDc() string {
	return "dc1"
}
func (c *ConsulConfig) GetAddress() string {
	return Server.ListenAddr
}
func (c *ConsulConfig) GetNodeID() int {
	return Server.NodeId
}

func (c *ConsulConfig) GetSvrName() string {
	return ServerName()
}
func (c *ConsulConfig) GetWatchSvrName() string {
	return GamePrefix
}
func (c *ConsulConfig) GetWatchFaildSvrName() string {
	return GamePrefix
}
func (c *ConsulConfig) GetRegistSelf() bool {
	return true
}

func (c *ConsulConfig) GetCheckAddress() string {
	return Server.WSAddr
}
