package conf

import (
	"encoding/json"
	"github.com/lovelly/leaf/log"
	"io/ioutil"
	. "mj/common/cost"
	"fmt"
	"strings"
	"strconv"
)

const (
	default_db_max_open      = 32
	default_db_max_idle      = 2
	default_stat_log_workers = 64
)

var Server struct {
	LogLevel        string
	LogPath         string
	WSAddr          string
	CertFile        string
	KeyFile         string
	TCPAddr         string
	MaxConnNum      int
	ConsolePort     int
	ProfilePath     string
	RoomModuleCount int

	BaseDbHost string
	BaseDbPort int
	BaseDbName string
	BaseDbUsername string
	BaseDbPassword string
	UserDbHost string
	UserDbPort int
	UserDbName string
	UserDbUsername string
	UserDbPassword string
	StatsDbHost string
	StatsDbPort int
	StatsDbName string
	StatsDbUsername string
	StatsDbPassword string

	ConsulAddr 		string
	ListenAddr      string
	ConnAddrs       map[string]string
	PendingWriteNum int
	ValidKind string
	ServerId int
}

var ValidKind = map[int]bool{}

func ServerName()string{
	return fmt.Sprintf(GamePrefix + "_%d", Server.ServerId)
}

func Init() {
	data, err := ioutil.ReadFile("./gameServer.json")
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


type DBConfig struct {}

func (c *DBConfig) GetBaseDSN()string {
	s := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		Server.BaseDbUsername, Server.BaseDbPassword, Server.BaseDbHost, Server.BaseDbPort, Server.BaseDbName, "parseTime=true&charset=utf8mb4")
	return s
}

func (c *DBConfig) GetUserDSN() string {
	s := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		Server.UserDbUsername, Server.UserDbPassword, Server.UserDbHost, Server.UserDbPort, Server.UserDbName, "parseTime=true&charset=utf8mb4")
	return s
}

func (c *DBConfig) GetStatsDSN() string {
	s := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		Server.StatsDbUsername, Server.StatsDbPassword, Server.StatsDbHost, Server.StatsDbPort, Server.StatsDbName, "parseTime=true&charset=utf8mb4")
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


//consul config
type  ConsulConfig struct {}


func (c *ConsulConfig)GetConsulAddr() string{
	return Server.ConsulAddr
}
func (c *ConsulConfig)GetConsulToken() string{
	return ""
}
func (c *ConsulConfig)GetConsulDc() string{
	return "dc1"
}
func (c *ConsulConfig)GetAddress() string{
	return Server.ListenAddr
}
func (c *ConsulConfig)GetServerID() int{
	return Server.ServerId
}

func (c *ConsulConfig)GetSvrName() string{
	return GamePrefix
}
func (c *ConsulConfig)GetWatchSvrName() string{
	return HallPrefix
}
func (c *ConsulConfig)GetWatchFaildSvrName() string{
	return HallPrefix
}
func (c *ConsulConfig)GetRegistSelf()bool{
	return true
}

func (c *ConsulConfig) GetCheckAddress() string{
	return Server.WSAddr
}


///////////////////////
func GetServerAddrAndPort()(string, int){
	l := strings.Split(Server.WSAddr, ":")
	if len(l) < 1 {
		l =  strings.Split(Server.TCPAddr, ":")
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
	return l[0],port
}