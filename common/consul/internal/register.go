package internal

import (
	"errors"
	"fmt"
	. "mj/common/cost"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/lovelly/leaf/log"
)

// 注册到consul
func register(c *api.Client, service *api.AgentServiceRegistration) (dereg chan bool) {
	var serviceID string

	registered := func() bool {
		if serviceID == "" {
			return false
		}
		services, err := c.Agent().Services()
		if err != nil {
			log.Error("consul: Cannot get service list. error:%s", err.Error())
			return false
		}
		return services[serviceID] != nil
	}

	register := func() {
		if err := c.Agent().ServiceRegister(service); err != nil {
			log.Error("consul: Cannot register in consul. error:%s", err.Error())
			return
		}

		log.Debug(" consul: Registered  with id service:%s", service.ID)
		log.Debug(" consul: Registered  with address:%s", service.Address)
		log.Debug(" consul: Registered  with tags %s", strings.Join(service.Tags, ","))
		if len(service.Checks) > 1 {
			log.Debug(" consul: Registered  with health check to %s", *service.Checks[0])
		}
		if len(service.Checks) > 2 {
			log.Debug("consul: Registered  with health check to  %s", *service.Checks[1])
		}

		serviceID = service.ID
	}

	deregister := func() {
		log.Debug("consul: Deregistering serverId :%s", serviceID)
		c.Agent().ServiceDeregister(serviceID)
	}

	dereg = make(chan bool)
	go func() {
		register()
		for {
			select {
			case <-dereg:
				deregister()
				dereg <- true
				return
			case <-time.After(time.Second * 2):
				if !registered() {
					register()
				}
			}
		}
	}()
	return dereg
}

//构建一个配置用于注册到consul
func buildRoomSvrConfig(Addr string, checkAddr, svrName string, svrID int) (*api.AgentServiceRegistration, error) {
	tcpPort := -1
	list := strings.Split(Addr, ":")
	if len(list) > 1 {
		var err error
		tcpPort, err = strconv.Atoi(list[1])
		if err != nil {
			log.Error("at buildRoomSvrConfig get tcp port error:", err.Error())
			panic("bug")
		}
	}

	consulSvrId := fmt.Sprintf(svrName+"_%v", svrID)
	q := &api.QueryOptions{RequireConsistent: true}
	svcs, _, err := Cli.Catalog().Service(svrName, "", q)
	if err != nil {
		log.Fatal("check regist faild at buildRoomSvrConfig %v", consulSvrId)
		return nil, errors.New("check regist faild at buildRoomSvrConfig")
	}

	for _, v := range svcs {
		if v.ServiceID == consulSvrId {
			log.Fatal("check regist faild at buildRoomSvrConfig 11 %v", consulSvrId)
			return nil, errors.New("check regist faild at buildRoomSvrConfig")
		}

		if v.Address == list[0] && v.ServicePort == tcpPort {
			log.Fatal("check regist faild at buildRoomSvrConfig 22 %v", consulSvrId)
			return nil, errors.New("check regist faild at buildRoomSvrConfig")
		}
	}

	// if ip.To16() != nil {
	// 	checkURL = fmt.Sprintf("http://[%s]:%d/health", ip, port)
	// }

	tag := make([]string, 0)
	strPort := strconv.Itoa(tcpPort)
	tag = append(tag, strPort)
	tag = append(tag, "50000")
	if GamePrefix == svrName {
		tag = append(tag, "this is game server")
	} else {
		tag = append(tag, "this is Hall server")
	}

	log.Debug("check addr == %s", checkAddr)

	chs := make([]*api.AgentServiceCheck, 0)
	chs = append(chs, &api.AgentServiceCheck{ // http port check
		TCP:      checkAddr,
		Interval: "2s",
		Timeout:  "5s",
		DeregisterCriticalServiceAfter: "5s",
	})

	service := &api.AgentServiceRegistration{
		ID:      consulSvrId,
		Name:    svrName,
		Address: list[0],
		Port:    tcpPort,
		Tags:    tag,
		Checks:  chs,
	}

	return service, nil
}
