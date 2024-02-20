package app

import (
	"bk_reptile/app/web/coolpc"
	"bk_reptile/app/web/efish"
	"bk_reptile/config"
	"fmt"

	"github.com/yangioc/bk_pack/util"
)

func (self *Handle) GetCoolpc() {
	datas, err := coolpc.GetWeb()
	if err != nil {
		panic(err)
	}

	for _, data := range datas {
		uuid := util.GenStrUUID(config.EnvInfo.NodeNum)
		payload, _ := util.Marshal(data)

		if err := self.dba.CreateCoolpcData(uuid, payload); err != nil {
			panic(err)
		}
	}
}

func (self *Handle) GetEfish() {
	// date_start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.Local)
	// date_end := time.Date(2024, 1, 30, 0, 0, 0, 0, time.Local)

	// datas, err := efish.GetHistory("1012", date_start, date_end)
	// if err != nil {
	// 	panic(err)
	// }

	// F109,F200,F241,F261,F270,F300,F330,F360,F400,F500,F513,F545,F600,F630,F708,F709,F722,F730,F800,F820,F826,F880,F916,F936,F950
	todayDatas, err := efish.GetToday("F109")
	if err != nil {
		panic(err)
	}

	// fmt.Println(datas, err)
	fmt.Println(todayDatas)
}
