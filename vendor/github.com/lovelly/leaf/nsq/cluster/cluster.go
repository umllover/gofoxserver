package cluster

import (
	"bytes"
	"log"

	nsq "github.com/nsqio/go-nsq"
	"github.com/segmentio/go-queue"
)

type cluster_conf interface {
	getNsqAddress() string
	getConCurrency() int
	getLookupAddress() string
	getTopic() string
	getChannel() string
}

func Start(cfg cluster_conf) {
	done := make(chan bool)
	b := new(bytes.Buffer)
	l := log.New(b, "", 0)

	c := queue.NewConsumer("events", "ingestion")
	c.SetLogger(l, nsq.LogLevelDebug)

	c.Set("nsqd", ":5001")
	c.Set("nsqds", []interface{}{":5001"})
	c.Set("concurrency", 5)
	c.Set("max_attempts", 10)
	c.Set("max_in_flight", 150)
	c.Set("default_requeue_delay", "15s")

	err := c.Start(nsq.HandlerFunc(func(msg *nsq.Message) error {
		done <- true
		return nil
	}))

	//assert.Equal(t, nil, err)

	go func() {
		p, err := nsq.NewProducer(":5001", nsq.NewConfig())
		p.Publish("events", []byte("hello"))
	}()

	<-done
	//assert.Equal(t, nil, c.Stop())
}
