package app

import (
	"bk_reptile/app/web/coolpc"
	"bk_reptile/config"
	"bk_reptile/service/dba"

	"github.com/yangioc/bk_pack/util"
)

var _instans *Handle

type Handle struct {
	dba dba.IHandle
}

func New(setting config.Env) *Handle {
	_instans = &Handle{
		dba: dba.New(setting),
	}
	return _instans
}

func (self *Handle) Launch() error {
	return self.dba.Launch()
}

func (self *Handle) Getcoolpc() {

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
