package app

import (
	"bk_reptile/config"
	"bk_reptile/model/msg_nats"
	"bk_reptile/service/dba"
	"context"
	"errors"

	"github.com/yangioc/bk_pack/log"
)

var ctxDoneError = errors.New("ctx Done")

var _instans *Handle

type Handle struct {
	ctx      context.Context
	dba      dba.IHandle
	messageQ *msg_nats.Handler
}

func New(setting config.Env, messageQ *msg_nats.Handler) *Handle {
	_instans = &Handle{
		ctx:      context.TODO(),
		dba:      dba.New(setting),
		messageQ: messageQ,
	}
	return _instans
}

func (self *Handle) Launch() {
	go func() {
		for {
			if err := self.dba.Launch(); err != nil {
				log.Error(err)
			}
		}
	}()

	// go func() {
	// 	for {
	// 		err := self.messageSub()
	// 		if err == ctxDoneError {
	// 			return
	// 		} else {
	// 			log.Errorf("app error: %v", err)
	// 		}
	// 	}
	// }()

}
