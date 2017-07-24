/// leader election
package consul

import (
	"fmt"
	"time"

	"github.com/lovelly/leaf/log"
	consulapi "github.com/hashicorp/consul/api"
)

const LockKey = "ElectionLeader"

var DefaultLeaderElection = NewLeaderElection()

func NewLeaderElection() *LeaderElection {
	v := &LeaderElection{}
	v.cleanupChannel = make(chan struct{})
	v.stopChannel = make(chan struct{})
	return v
}

type LeaderElection struct {
	lock           *consulapi.Lock
	cleanupChannel chan struct{}
	stopChannel    chan struct{}
	leader         bool
}

func (l *LeaderElection) start() {
	fmt.Println("begin start LeaderElection")
	clean := false
	for !clean {
		select {
		case <-l.cleanupChannel:
			clean = true
		default:
			log.Debug("Running for leader election...")
			intChan, _ := l.lock.Lock(l.stopChannel)
			if intChan != nil {
				log.Debug("Now acting as leader.")
				l.leader = true
				<-intChan
				l.leader = false
				log.Debug("Lost leadership.")
				l.lock.Unlock()
				l.lock.Destroy()
			} else {
				log.Debug("start .... ")
				time.Sleep(10000 * time.Millisecond)
			}
		}
	}
}

func (l *LeaderElection) stop() {
	log.Debug("cleaning up")
	l.cleanupChannel <- struct{}{}
	l.stopChannel <- struct{}{}
	l.lock.Unlock()
	l.lock.Destroy()
	l.leader = false
	log.Debug("cleanup done")
}

func (l *LeaderElection) hasLeader() bool {
	kvpair, _, err := Cli.KV().Get(LockKey, nil)
	return kvpair != nil && err == nil
}

func (l *LeaderElection) startLeaderElection() {
	config := consulapi.DefaultConfig()
	config.Address = Config.GetConsulAddr()
	config.Datacenter = Config.GetConsulDc()
	config.Token = Config.GetConsulToken()
	var err error
	l.lock, err = Cli.LockKey(LockKey)
	if err != nil {
		log.Debug("at startLeaderElection LockKey error", err.Error())
		return
	}
	go l.start()
}
