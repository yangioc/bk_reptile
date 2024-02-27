package dba

import (
	"fmt"

	"github.com/yangioc/bk_pack/proto/dtomsg"
)

func (self *Handle) CreateCoolpcData(uuid string, payload []byte) error {
	return self.CommonCreate(uuid, "notice", "create.coolpcdata", payload)
}

func (self *Handle) resChanNew(uuid string) (chan *dtomsg.Dto_Msg_Res, error) {
	resChan := make(chan *dtomsg.Dto_Msg_Res)
	if _, isLoad := self.requestTracing.LoadOrStore(uuid, resChan); isLoad {
		return nil, fmt.Errorf("[Error][getResChan] %v.", uuid)
	}

	return resChan, nil
}

func (self *Handle) resChanLoadAndDelete(uuid string) (chan *dtomsg.Dto_Msg_Res, bool) {
	resChan, ok := self.requestTracing.LoadAndDelete(uuid)
	if !ok {
		return nil, false
	}
	return resChan.(chan *dtomsg.Dto_Msg_Res), true
}

func (self *Handle) CreateEfish(uuid string, payload []byte) error {
	return self.CommonCreate(uuid, "notice", "create.efish", payload)
}

func (self *Handle) CreateStockAnalysis(uuid string, payload []byte) error {
	return self.CommonCreate(uuid, "notice", "create.stockanalysis", payload)
}

func (self *Handle) CreateStockIndex(uuid string, payload []byte) error {
	return self.CommonCreate(uuid, "notice", "create.stockindex", payload)
}

func (self *Handle) CreateStockMarket(uuid string, payload []byte) error {
	return self.CommonCreate(uuid, "notice", "create.stockmarket", payload)
}

func (self *Handle) CreateStockClosePrice(uuid string, payload []byte) error {
	return self.CommonCreate(uuid, "notice", "create.stockclosePrice", payload)
}

func (self *Handle) CreateStockThreefoundationTotal(uuid string, payload []byte) error {
	return self.CommonCreate(uuid, "notice", "create.stockthreefoundationtotal", payload)
}

func (self *Handle) CreateThreefoundationStockDay(uuid string, payload []byte) error {
	return self.CommonCreate(uuid, "notice", "create.stockthreefoundationstockday", payload)
}
