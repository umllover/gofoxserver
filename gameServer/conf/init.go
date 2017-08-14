package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	. "mj/common/cost"
	"strconv"
	"strings"

	"github.com/lovelly/leaf/log"
)

const (
	default_db_max_open      = 10
	default_db_max_idle      = 10
	default_stat_log_workers = 64
)

var Server struct {
	LogLevel        string
	LogPath         string
	WSAddr          string
	CertFile        string
	KeyFile         string
	TCPAddr         string
	HttpAddr        string
	MaxConnNum      int
	ConsolePort     int
	ProfilePath     string
	RoomModuleCount int

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

	NsqdAddrs       []string
	NsqLookupdAddrs []string
	PdrNsqdAddr     string

	ConsulAddr      string
	RedisAddr       string
	RedisPwd        string
	ListenAddr      string
	ConnAddrs       map[string]string
	PendingWriteNum int
	ValidKind       string
	WatchAddr       string
	NodeId          int
}

var ValidKind = map[int]bool{}

func ServerName() string {
	return fmt.Sprintf(GamePrefixFmt, Server.NodeId)
}

func ServerNsqCahnnel() string {
	return fmt.Sprintf(GameChannelFmt, Server.NodeId)
}

func Init(filePaths ...string) {
	var filePaht = "./gameServer.json"
	if len(filePaths) > 0 {
		filePaht = filePaths[0]
	}
	data, err := ioutil.ReadFile(filePaht)
	if err != nil {
		log.Fatal("%v", err)
	}
	err = json.Unmarshal(data, &Server)
	if err != nil {
		log.Fatal("%v", err)
	}

	list := strings.Split(Server.ValidKind, ",")
	for _, v := range list {
		ikind, err := strconv.Atoi(v)
		if err != nil {
			continue
		}
		ValidKind[ikind] = true
	}

}

type DBConfig struct{}

func (c *DBConfig) GetAccoutDSN() string {
	s := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		Server.AccountDbUsername, Server.AccountDbPassword, Server.AccountDbHost, Server.AccountDbPort, Server.AccountDbName, "parseTime=true&interpolateParams=true&charset=utf8mb4")
	return s
}

func (c *DBConfig) GetRedisAddr() string {
	return Server.RedisAddr
}

func (c *DBConfig) GetRedisPwd() string {
	return Server.RedisPwd
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
	return HallPrefix
}
func (c *ConsulConfig) GetWatchFaildSvrName() string {
	return HallPrefix
}
func (c *ConsulConfig) GetRegistSelf() bool {
	return true
}

func (c *ConsulConfig) GetCheckAddress() string {
	return Server.WSAddr
}

///////////////////////
func GetServerAddrAndPort() (string, int) {
	l := strings.Split(Server.WSAddr, ":")
	if len(l) < 1 {
		l = strings.Split(Server.TCPAddr, ":")
	}

	if len(l) < 2 {
		log.Debug("not foud sver addr at GetServerAddrAndPort")
		return "", 0
	}

	port, err := strconv.Atoi(l[1])
	if err != nil {
		log.Debug("not foud sver port at GetServerAddrAndPort")
		return "", 0
	}
	return l[0], port
}
