package cluster

import (
	"sync"

	"github.com/lovelly/leaf/network/gob"

	"github.com/lovelly/leaf/log"
	"github.com/nsqio/go-nsq"
)

var (
	producer    *nsq.Producer
	proclose    bool
	prolock     sync.Locker
	encMutex    sync.Locker
	publishChan = make(chan *NsqRequest, 10000)
	encoder     *gob.Encoder
	decoder     *gob.Decoder
)

type NsqRequest struct {
	Topc string
	data string
}

type cluster_config struct {
	LogLv              string
	Channel            string   //唯一标识
	Csmtopics          []string //需要订阅的主题
	CsmUserAgent       string   //消费者的UserAgent
	CsmNsqdAddrs       []string
	CsmNsqLookupdAddrs []string
	CsmMaxInFlight     int
	PdrNsqdAddr        string //生产者需要连接的nsqd地址
	PdrUserAgent       string //生产者的UserAgent
	PdrMaxInFlight     int
}

func Start(cfg *cluster_config) {
	if cfg.PdrMaxInFlight == 0 {
		cfg.PdrMaxInFlight = 100
	}

	if cfg.CsmMaxInFlight == 0 {
		cfg.CsmMaxInFlight = 100
	}

	if cfg.PdrUserAgent == "" {
		cfg.PdrUserAgent = "mqjx_producer"
	}

	if cfg.CsmUserAgent == "" {
		cfg.CsmUserAgent = "mqjx_consumer"
	}

	var err error
	nsqcfg := nsq.NewConfig()
	nsqcfg.UserAgent = cfg.PdrUserAgent
	nsqcfg.MaxInFlight = cfg.PdrMaxInFlight
	if producer, err = nsq.NewProducer(cfg.PdrNsqdAddr, nsqcfg); err != nil {
		log.Fatal("start nsq client error:%s", err.Error())
	}

	loglv := getLogLovel(cfg.LogLv)
	producer.SetLogger(log.GetBaseLogger(), loglv)
	err = producer.Ping()
	if err != nil {
		log.Fatal("ping nsq client error:%s", err.Error())
	}

	for _, tpc := range cfg.Csmtopics {
		nsqcfg := nsq.NewConfig()
		nsqcfg.UserAgent = cfg.CsmUserAgent
		nsqcfg.MaxInFlight = cfg.CsmMaxInFlight
		consumer, err := nsq.NewConsumer(tpc, cfg.Channel, nsqcfg)
		if err != nil {
			log.Fatal(" nsq NewConsumer error:%s", err.Error())
		}

		consumer.SetLogger(log.GetBaseLogger(), loglv)
		consumer.AddHandler(&nsqHandler{})

		if err = consumer.ConnectToNSQDs(cfg.CsmNsqdAddrs); err != nil {
			return
		}
		if err = consumer.ConnectToNSQLookupds(cfg.CsmNsqLookupdAddrs); err != nil {
			return
		}
	}

	go publishLoop()
}

type nsqHandler struct {
}

// HandleMessage - Handles an NSQ message.
func (h *nsqHandler) HandleMessage(message *nsq.Message) error {
	defer func() {
		message.Requeue(-1)
		message.Finish()
	}()
	message.DisableAutoResponse()
	encMutex.Lock()
	data, err := Processor.Marshal(encoder, message.Body)
	encMutex.Unlock()
	if err != nil {
		log.Error("handler msg error:%s, data:%s", err.Error(), string(message.Body))
		return nil
	}

	Processor.Route(data)
	return nil
}

func Publish(msg *NsqRequest) {
	prolock.Lock()
	defer prolock.Unlock()
	if proclose {
		return
	}
	publishChan <- msg
}

func Close() {
	prolock.Lock()
	proclose = true
	prolock.Unlock()
	producer.Stop()
	close(publishChan)
}

func publishLoop() {
	defer func() {
		prolock.Lock()
		proclose = true
		prolock.Unlock()
		producer.Stop()
		close(publishChan)
	}()

	var msg *NsqRequest
	var open bool
	for {
		select {
		case msg, open = <-publishChan:
			if !open {
				return
			}
			if msg == nil {
				continue
			}
			safePulishg(msg)
			return
		}
	}
}

func safePulishg(msg *NsqRequest) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("Publish msg recover error : %s, topc :%s ", msg.data, msg.Topc)
		}
	}()
	err := producer.Publish(msg.Topc, msg.data)
	if err != nil {
		log.Error("Publish msg error : %s, topc :%s ", msg.data, msg.Topc)
	}
}

func getLogLovel(loglv string) nsq.LogLevel {
	switch loglv {
	case "Debug":
		return nsq.LogLevelDebug
	case "Release":
		return nsq.LogLevelInfo
	case "Warn":
		return nsq.LogLevelWarning
	case "Error":
		return nsq.LogLevelError
	}

	return nsq.LogLevelError
}
