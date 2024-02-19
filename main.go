package main

import (
	"bk_reptile/config"
	"bk_reptile/model/msg_nats"
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/yangioc/bk_pack/log"
)

var configPath = flag.String("config", "./env.yaml", "specific config to processing")

func main() {
	if err := config.Init(*configPath); err != nil {
		panic(err)
	}

	// handle_app := app.New(*config.EnvInfo)
	// go func() {
	// 	if err := handle_app.Launch(); err != nil {
	// 		panic(err)
	// 	}
	// }()

	msg_nats.New(context.TODO(), nil)

	// handle_crontab := crontab.New()
	// if err := handle_crontab.AddTask("test1", "*/5 * * * * ?", func() { fmt.Println(util.ServerTimeNow()) }); err != nil {
	// 	panic(err)
	// }

	// handle_crontab.Run()

	log.Info("Service Up.")
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)
	<-c

	log.Info("Service Down.")
}
