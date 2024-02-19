package msg_nats

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

type Handler struct {
	ctx       context.Context
	writeFlow chan []byte
	conn      *nats.Conn

	subMap map[string]*nats.Subscription
}

// connect 嘗試連線，成功會開啟 publish() 生命週期
func (self *Handler) connect(urls string) error {
	conn, err := nats.Connect(urls)
	if err != nil {
		return err
	}

	self.conn = conn
	return err
}

func (self *Handler) Pub(subj string, data []byte) error {
	if err := self.conn.Publish(subj, data); err != nil {
		return err
	}

	return nil
}

///////////////////////////////////////////////////////
//	使用 func 訂閱底層將使用 goroutine 執行
//
//	Sub: 訂閱
//	Sync: 同步式訂閱, 只在前一筆訊號處理完以後才會處理下一筆
//	Chan: Chan 回傳的訂閱方式
//	Group: 訂閱群組, 只會由群組內任一定樂者接收
///////////////////////////////////////////////////////

func (self *Handler) Sub(subj string, cb func(msg *nats.Msg)) error {
	sub, err := self.conn.Subscribe(subj, cb)
	if err != nil {
		return err
	}

	self.subMap[subj] = sub
	return nil
}

func (self *Handler) SubChan(subj string, ch chan *nats.Msg) error {
	sub, err := self.conn.ChanSubscribe(subj, ch)
	if err != nil {
		return err
	}

	self.subMap[subj] = sub
	return nil
}

func (self *Handler) SubGroup(subj, group string, cb func(msg *nats.Msg)) error {
	sub, err := self.conn.QueueSubscribe(subj, group, cb)
	if err != nil {
		return err
	}

	self.subMap[fmt.Sprintf("%s%s", subj, group)] = sub
	return nil
}

func (self *Handler) SubGroupChan(subj, group string, ch chan *nats.Msg) error {
	sub, err := self.conn.ChanQueueSubscribe(subj, group, ch)
	if err != nil {
		return err
	}

	self.subMap[fmt.Sprintf("%s%s", subj, group)] = sub
	return nil
}

func (self *Handler) SubSync(subj string) error {
	sub, err := self.conn.SubscribeSync(subj)
	if err != nil {
		return err
	}

	for {
		// Wait for a message
		msg, err := sub.NextMsg(time.Minute)

		if errors.Is(nats.ErrTimeout, err) {
			continue
		} else if err != nil {
			return err
		}

		// Use the response
		fmt.Printf("SubSync Reply: %s\n", msg.Data)
	}
}

func (self *Handler) SubGroupSync(subj, group string, cb func(msg *nats.Msg)) error {
	sub, err := self.conn.QueueSubscribeSync(subj, group)
	if err != nil {
		return err
	}

	for {
		// Wait for a message
		msg, err := sub.NextMsg(time.Minute)

		if errors.Is(nats.ErrTimeout, err) {
			continue
		} else if err != nil {
			return err
		}

		// Use the response
		fmt.Printf("SubGroupSync Reply: %s\n", msg.Data)
	}
}

func (self *Handler) SubGroupChanSync(subj, group string, ch chan *nats.Msg) error {
	sub, err := self.conn.QueueSubscribeSyncWithChan(subj, group, ch)
	if err != nil {
		return err
	}

	for {
		// Wait for a message
		msg, err := sub.NextMsg(time.Minute)

		if errors.Is(nats.ErrTimeout, err) {
			continue
		} else if err != nil {
			return err
		}

		// Use the response
		fmt.Printf("SubGroupChanSync Reply: %s\n", msg.Data)
	}
}
