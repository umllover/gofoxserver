package internal

import (
	"errors"
	"github.com/lovelly/leaf/log"
	"github.com/hashicorp/consul/api"
)

const (
	UpdateConfigTicke = 1 //定时获取最新负载配置的时间 单位秒
)

type Rgconfig interface {
	GetConsulAddr() string
	GetConsulToken() string
	GetConsulDc() string
	GetAddress() string
	GetServerID() int
	GetSvrName() string
	GetWatchSvrName() string
	GetWatchFaildSvrName() string
	GetRegistSelf()bool
	GetCheckAddress() string
}

type CacheInfo struct {
	Host     string   //ip和端口
	Csid     string   //consul里面的id
	weight   int      //权重
	MaxCount int      //最大承载数量
	tags     []string //标签信息
}

var (
	Config Rgconfig
	Cli *api.Client
	Dc     string
	dereg  chan bool
)

//初始化consul
func InitConsul(scheme string) {
	if Config == nil {
		log.Error("at consul Init your moust set config")
	}

	var err error
	Cli, err = api.NewClient(&api.Config{Address: Config.GetConsulAddr(), Scheme: scheme, Token: Config.GetConsulToken(), Datacenter: Config.GetConsulDc()})
	if err != nil {
		log.Error("at NewConsuled NewClient error :%s",  err.Error())
		return
	}

	//ping the agent
	Dc, err = datacenter(Cli)
	if err != nil {
		log.Error("at NewConsuled datacenter error:$s",  err.Error())
		return
	}

	log.Debug("consul: Connecting to %s in datacenter %s", Config.GetConsulAddr(), Config.GetConsulDc())
}

//注册服务
func  Register() error {
	service, err := buildRoomSvrConfig(Config.GetAddress(),Config.GetCheckAddress(), Config.GetSvrName(), Config.GetServerID())
	if err != nil {
		return err
	}

	dereg = register(Cli, service)
	return nil
}

//取消服务
func deregDeregister() error {
	if dereg != nil {
		dereg <- true // trigger deregistration
		<-dereg       // wait for completion
	}
	return nil
}


//开启一个获取房间健康配置信息的协程
func WatchServices( serverName string) {
	status := make([]string, 0)
	status = append(status, "passing")
	go watchServices(Cli, serverName, status)
}

//关注失败的服务
func WatchAllFaild(serverName string){
	go WatchAllFaildServices(Cli, serverName)
}

//开启一个定时获取房间负载信息的协程
func  WatchManual(kvpath string) {
	log.Debug("[ consul: Watching KV path %s",  kvpath)
	go watchKV(Cli, kvpath)
}


//报告一个检查通过
func deregPassTTL(serverId string) {
	submitCheck(Cli, serverId)
}

//删除kv
func deregDelKv(kvPath string) {
	DelKV(Cli, kvPath)
}

// datacenter returns the datacenter of the local agent
func datacenter(c *api.Client) (string, error) {
	self, err := c.Agent().Self()
	if err != nil {
		return "", err
	}

	Config, ok := self["Config"]
	if !ok {
		return "", errors.New("consul: self.Config not found")
	}
	dc, ok := Config["Datacenter"].(string)
	if !ok {
		return "", errors.New("consul: self.Datacenter not found")
	}
	return dc, nil
}
