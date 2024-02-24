package app

import (
	"bk_reptile/app/web/coolpc"
	"bk_reptile/app/web/efish"
	"bk_reptile/config"
	"errors"
	"fmt"
	"time"

	"github.com/yangioc/bk_pack/log"
	"github.com/yangioc/bk_pack/proto/dtomsg"
	"github.com/yangioc/bk_pack/util"
)

func (self *Handle) handleReqMessage(msg *dtomsg.Dto_Msg) ([]byte, error) {
	return nil, errors.New("req message request not found.")
}

func (self *Handle) handleNoticMessage(msg *dtomsg.Dto_Msg) error {
	return errors.New("notic message request not found.")
}

func (self *Handle) GetCoolpc() {
	uuid := util.GenStrUUID(config.EnvInfo.NodeNum)
	datas, err := coolpc.GetWeb()
	if err != nil {
		panic(err)
	}

	for _, data := range datas {
		payload, err := util.Marshal(data)
		if err != nil {
			panic(err)
		}

		if err = self.dba.CreateCoolpcData(uuid, payload); err != nil {
			panic(err)
		}
	}
}

func (self *Handle) GetEfish() {
	uuid := util.GenStrUUID(config.EnvInfo.NodeNum)

	date := util.ServerTimeNow()

	locations := []string{"F109", "F200", "F241", "F261", "F270", "F300", "F330", "F360", "F400", "F500", "F513", "F545", "F600", "F630", "F708", "F709", "F722", "F730", "F800", "F820", "F826", "F880", "F916", "F936", "F950"}

	dataGroup := map[string]interface{}{
		"_key": date.Format("20060102"),
	}

	for _, location := range locations {
		data, err := efish.GetDayFishByMarket(location, date)
		if err != nil {
			panic(err)
		}

		dataGroup[location] = data
	}

	payload, _ := util.Marshal(dataGroup)
	if err := self.dba.CreateEfish(uuid, payload); err != nil {
		panic(err)
	}
}

func (self *Handle) GetOldEfish(fishId string, date_start, date_end time.Time) {
	if date_start.After(date_end) {
		log.Error("GetOldEfish time error")
		return
	}
	// date_start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.Local)
	// date_end := time.Date(2024, 1, 30, 0, 0, 0, 0, time.Local)

	datas, err := efish.GetHistory(fishId, date_start, date_end)
	if err != nil {
		panic(err)
	}

	fmt.Println(datas)
}
