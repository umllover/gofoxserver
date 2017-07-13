package cluster

import (
	"github.com/lovelly/leaf/log"
	nsq "github.com/nsqio/go-nsq"

	"sync"
)

var (
	clientsMutex sync.Mutex
	clients      = make(map[string]*NsqClient)
)

type NsqClient struct {
	Addr       string
	ServerName string
}

func AddClient(c *NsqClient) {
	log.Debug("at cluster AddClient %s, %s", c.ServerName, c.Addr)
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	clients[c.ServerName] = c
}

func RemoveClient(serverName string) {
	_, ok := clients[serverName]
	if ok {
		log.Debug("at cluster _removeClient %s", serverName)
		delete(clients, serverName)
	}
}

