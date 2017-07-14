package consul

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/lovelly/leaf/log"
)

//定时获取consul上所有检点的健康信息
func watchServices(client *api.Client, serverName string, status []string) {
	var lastIndex uint64
	for {
		q := &api.QueryOptions{RequireConsistent: true, WaitIndex: lastIndex}
		checks, meta, err := client.Health().Checks(serverName, q)
		if err != nil {
			log.Error("[WARN] consul: Error fetching health state.%s", err.Error())
			time.Sleep(time.Second * UpdateConfigTicke)
			continue
		}

		if meta.LastIndex == lastIndex {
			continue
		}

		lastIndex = meta.LastIndex

		//log.Debug("[INFO] consul: Health changed  LastIndex:%d,, new inx: %d, server name %s, data:%v", lastIndex, meta.LastIndex, serverName, checks)
		newSvrs := servicesConfig(client, passingServices(checks, status))
		log.Debug("%v, %v", len(newSvrs) > 0, len(checks) < 1)
		if len(newSvrs) > 0 || len(checks) < 1 {
			ChanRPC.Go("AddServerInfo", newSvrs)
		}
	}
}

//获取健康信息
func servicesConfig(client *api.Client, checks []*api.HealthCheck) map[string]*CacheInfo {
	configs := make(map[string]*CacheInfo)
	m := map[string]map[string]bool{}

	for _, check := range checks {
		name, id := check.ServiceName, check.ServiceID

		if _, ok := m[name]; !ok {
			m[name] = map[string]bool{}
		}
		m[name][id] = true
	}

	for name, passing := range m { //去除不健康的服务
		if name == "" || len(passing) == 0 {
			continue
		}

		q := &api.QueryOptions{RequireConsistent: true}
		svcs, _, err := client.Catalog().Service(name, "", q)
		if err != nil {
			log.Error("[WARN] consul: Error getting catalog service %s name = %s  error = %s", name, err.Error())
			continue
		}

		filterMap := make(map[string]string)
		hasFunc := func(host string) string {
			cid, ok := filterMap[host]
			if ok {
				return cid
			}
			return ""
		}

		for _, svc := range svcs {
			// 去除不健康的服务
			if _, ok := passing[svc.ServiceID]; !ok {
				continue
			}
			config := &CacheInfo{}
			if len(svc.ServiceTags) >= 2 {
				max, err := strconv.Atoi(svc.ServiceTags[1])
				if err == nil {
					config.MaxCount = max
				}
			}

			config.Host = fmt.Sprintf("%v:%v", svc.ServiceAddress, svc.ServicePort)
			config.Csid = svc.ServiceID
			config.tags = svc.ServiceTags
			oldid := hasFunc(config.Host)
			if oldid != "" {
				for i := 0; i < 5; i++ {
					log.Error("config error same addr , please check you config file host ==%s, serverId 1 = %d, serverId 2  = %d",
						config.Host, config.Csid, oldid)
				}
				continue
			}
			filterMap[config.Host] = config.Csid
			configs[config.Csid] = config
		}
	}
	return configs
}

//监控不健康的服务
func WatchAllFaildServices(client *api.Client, ServiceName string) {
	var lastIndex uint64
	for {
		q := &api.QueryOptions{RequireConsistent: true, WaitIndex: lastIndex}
		checks, meta, err := client.Health().State("critical", q)
		if err != nil {
			log.Error("[WARN] consul: Error WatchFaildSvices health state.  error =%s", err.Error())
			time.Sleep(time.Second * UpdateConfigTicke)
			continue
		}
		if lastIndex == meta.LastIndex {
			continue
		}

		lastIndex = meta.LastIndex
		if len(checks) < 1 {
			continue
		}
		log.Debug(" at WatchAllFaildServices .... lastIndex :%d, new indx:%d, serverName:%s data:%v", lastIndex, meta.LastIndex, ServiceName, checks)
		status := make([]string, 1)
		status = append(status, "critical")
		checkIds := passingServices(checks, status)
		newFaildSvr := make(map[string]string)
		for _, v := range checkIds {
			if v.ServiceName == ServiceName {
				newFaildSvr[v.ServiceID] = v.ServiceName
			}
		}

		if len(newFaildSvr) > 0 {
			ChanRPC.Go("NotifySvrFaild", newFaildSvr)
		}
	}
}

// 报告健康状态
func submitCheck(client *api.Client, serverid string) {
	client.Agent().PassTTL(serverid, "pass ok le")
}
