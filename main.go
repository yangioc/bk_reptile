package main

import (
	"bk_reptile/app"
	"bk_reptile/config"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/yangioc/bk_pack/crontab"
	"github.com/yangioc/bk_pack/log"
	"github.com/yangioc/bk_pack/util"
)

var configPath = flag.String("config", "./env.yaml", "specific config to processing")

func main() {
	if err := config.Init(*configPath); err != nil {
		panic(err)
	}

	handle_app := app.New(*config.EnvInfo)
	go func() {
		if err := handle_app.Launch(); err != nil {
			panic(err)
		}
	}()

	handle_crontab := crontab.Init()
	if err := handle_crontab.NewIntervalTask("test1", "*/5 * * * * ?", func() { fmt.Println(util.ServerTimeNow()) }); err != nil {
		panic(err)
	}

	handle_crontab.Run()

	log.Info("Service Up.")
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)
	<-c

	log.Info("Service Down.")
}
