package messageq

import (
	"bk_reptile/config"
	"fmt"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/yangioc/bk_pack/log"
)

type Handle struct {
	conn   *nats.Conn
	subMap sync.Map
}

func New() *Handle {
	return &Handle{}
}

func (self *Handle) Launch(config config.Env) error {
	var err error
	if self.conn, err = nats.Connect(config.Nats.Addr, nats.UserInfo(config.Nats.UserName, config.Nats.Password)); err != nil {
		return err
	}
	log.Infof("[MessageQ][Launch] Connect success, address: %s", config.Nats.Addr)
	return nil
}

func (self *Handle) Publish(topic string, data []byte) error {
	log.Debugf("[MessageQ][Publish] topic: %s", topic)
	return self.conn.Publish(topic, data)
}

func (self *Handle) RequestAsync(topic, reply string, data []byte) error {
	log.Debugf("[MessageQ][RequestAsync] topic: %s reply: %s", topic, reply)
	return self.conn.PublishRequest(topic, reply, data)
}

func (self *Handle) Request(topic string, data []byte, timeout time.Duration) (*nats.Msg, error) {
	log.Debugf("[MessageQ][Request] topic: %s timeout: %v", topic, timeout)
	return self.conn.Request(topic, data, timeout)
}

func (self *Handle) Subscribe(topic string, fun func(msg *nats.Msg)) error {
	log.Debugf("[MessageQ][Subscribe] %s", topic)
	sub, err := self.conn.Subscribe(topic, fun)
	if err != nil {
		return err
	}

	if _, isload := self.subMap.LoadOrStore(topic, sub); isload {
		return fmt.Errorf("topic %s exist", topic)
	}

	return nil
}

func (self *Handle) Unsubscribe(topic string) error {
	log.Debugf("[MessageQ][Unsubscribe] %s", topic)
	Isub, isLoad := self.subMap.LoadAndDelete(topic)
	if !isLoad {
		return nil
	}

	sub := Isub.(*nats.Subscription)
	return sub.Unsubscribe()
}
