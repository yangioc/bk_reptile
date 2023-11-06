package dba

import (
	"context"
	"fmt"

	"github.com/yangioc/bk_pack/dto"
	"github.com/yangioc/bk_pack/util"
)

func (self *Handle) CreateCoolpcData(uuid string, payload []byte) error {
	msg, err := util.PackDBAReq(uuid, "createcoolpcdata", payload)
	if err != nil {
		panic(err)
	}

	if err := self.websocket.Send(context.TODO(), msg); err != nil {
		panic(err)
	}

	return nil
}

func (self *Handle) resChanNew(uuid string) (chan *dto.Dto_DBA_Res, error) {
	resChan := make(chan *dto.Dto_DBA_Res)
	if _, isLoad := self.requestTracing.LoadOrStore(uuid, resChan); isLoad {
		return nil, fmt.Errorf("[Error][getResChan] %v.", uuid)
	}

	return resChan, nil
}

func (self *Handle) resChanLoadAndDelete(uuid string) (chan *dto.Dto_DBA_Res, bool) {
	resChan, ok := self.requestTracing.LoadAndDelete(uuid)
	if !ok {
		return nil, false
	}
	return resChan.(chan *dto.Dto_DBA_Res), true
}
