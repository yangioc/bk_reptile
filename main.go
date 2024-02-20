package main

import (
	"bk_reptile/app"
	"bk_reptile/config"
	"bk_reptile/messageq"
	"bk_reptile/model/msg_nats"
	"bk_reptile/tmpproto/dtoschedule"
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/yangioc/bk_pack/log"
	"github.com/yangioc/bk_pack/proto/dtomsg"
	"github.com/yangioc/bk_pack/util"
	"google.golang.org/protobuf/proto"
)

var configPath = flag.String("config", "./env.yaml", "specific config to processing")

func main() {
	if err := config.Init(*configPath); err != nil {
		panic(err)
	}

	// 設定 log
	log.Level = log.Level_Info // 預設
	if logLevel, ok := log.LevelToStringMap[config.EnvInfo.Log.Level]; ok {
		log.Level = logLevel
	}

	// messageq 介面
	// type1
	handle_messageq := msg_nats.New(context.TODO(), *config.EnvInfo, nil)

	//type2
	// handle_messageq := messageq.New()
	// if err := handle_messageq.Launch(*config.EnvInfo); err != nil {
	// 	panic(err)
	// }
	// addtask(handle_messageq)

	// 核心服務
	handle_app := app.New(*config.EnvInfo, handle_messageq)
	handle_app.Launch()

	// handle_crontab := crontab.New()
	// if err := handle_crontab.AddTask("test1", "*/5 * * * * ?", func() { fmt.Println(util.ServerTimeNow()) }); err != nil {
	// 	panic(err)
	// }

	// handle_crontab.Run()

	// test
	// handle_app.GetCoolpc()
	// handle_app.GetEfish()
	///////

	log.Info("Service Up.")
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)
	<-c

	log.Info("Service Down.")
}

func addtask(handle_messageq *messageq.Handle) {
	callback := "reptile.coolpc.run"
	req := dtoschedule.NewTaskReq{
		Name:     "coolpc",
		Spec:     "*/5 * * * * ?",
		CallBack: callback,
	}

	reqdata, err := proto.Marshal(&req)
	if err != nil {
		panic(err)
	}

	msg := dtomsg.Dto_Msg{
		Type:    "req",
		Request: "addtask",
		Data:    reqdata,
	}
	msgData, err := util.Marshal(msg)
	if err != nil {
		panic(err)
	}

	payload, err := util.MsgEncode(&dtomsg.Dto_Base{
		UUID: util.GenStrUUID(4),
		// From:           "",
		// Router:         "",
		StartTime:      util.ServerTimeNow().UnixMicro(),
		ExpirationTime: util.ServerTimeNow().Add(5 * time.Hour).UnixMicro(),
		Payload:        msgData,
	})
	if err != nil {
		panic(err)
	}

	handle_messageq.Subscribe(callback, func(msg *nats.Msg) {
		fmt.Println("on coolpc run")
	})
	handle_messageq.Publish("schdule.rep.a1", payload)
}
