package cluster

import (
	"sync"

	"github.com/lovelly/leaf/network/gob"

	"errors"

	"github.com/lovelly/leaf/log"
	"github.com/nsqio/go-nsq"
)

var (
	producer    *nsq.Producer
	consumers   []*nsq.Consumer
	proclose    bool
	prolock     sync.Mutex
	publishChan = make(chan *NsqRequest, 10000)
	encoder     *gob.Encoder
	decoder     *gob.Decoder
	SelfName    string
)

type NsqRequest struct {
	Topc string
	Data []byte
}

type Cluster_config struct {
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
	SelfName           string
}

func Start(cfg *Cluster_config) {
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

	if cfg.SelfName == "" {
		log.Fatal("at nsq start selfname is nil")
	}

	var err error
	SelfName = cfg.SelfName
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
		consumers = append(consumers, consumer)
		if err = consumer.ConnectToNSQDs(cfg.CsmNsqdAddrs); err != nil {
			return
		}
		if err = consumer.ConnectToNSQLookupds(cfg.CsmNsqLookupdAddrs); err != nil {
			return
		}

	}

	go publishLoop()
}

func Publish(topc string, msg interface{}) error {
	prolock.Lock()
	if proclose {
		prolock.Unlock()
		return errors.New("server is close at Publish")
	}
	prolock.Unlock()

	data, err := Processor.Marshal(encoder, msg)
	if err != nil {
		return err
	}

	if len(data) < 1 {
		return errors.New("error at Publish data is nil")
	}
	publishChan <- &NsqRequest{Topc: topc, Data: data[0]}
	return nil
}

func Stop() {
	prolock.Lock()
	proclose = true
	prolock.Unlock()
	for _, v := range consumers {
		v.Stop()
	}
	producer.Stop()
	close(publishChan)
}

func isClose() bool {
	prolock.Lock()
	defer prolock.Unlock()
	return proclose
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

type nsqHandler struct {
}

// HandleMessage - Handles an NSQ message.
func (h *nsqHandler) HandleMessage(message *nsq.Message) error {
	defer func() {
		message.Requeue(-1)
		message.Finish()
	}()
	message.DisableAutoResponse()
	data, err := Processor.Unmarshal(decoder, message.Body)
	if err != nil {
		log.Error("handler msg error:%s, data:%s", err.Error(), string(message.Body))
		return nil
	}
	msg := data.(*S2S_NsqMsg)
	if msg.CallType == callBroadcast && msg.ServerName == SelfName {
		return nil
	}
	switch msg.ReqType {
	case NsqMsgTypeReq:
		handleRequestMsg(msg)
	case NsqMsgTypeRsp:
		handleResponseMsg(msg)
	}
	Processor.Route(data, nil)
	return nil
}

func safePulishg(msg *NsqRequest) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("Publish msg recover error : %s, topc :%s ", string(msg.Data), msg.Topc)
		}
	}()
	err := producer.Publish(msg.Topc, msg.Data)
	if err != nil {
		log.Error("Publish msg error : %s, topc :%s ", string(msg.Data), msg.Topc)
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
