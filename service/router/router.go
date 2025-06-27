package main

import (
	"flag"

	"github.com/Qsgs-Fans/freekill-server/service/router/internal/config"
	"github.com/Qsgs-Fans/freekill-server/service/router/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
)

var configFile = flag.String("f", "etc/router.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	go ctx.TcpServer.StartListenOnMQ()
	// TODO (go) udpServer
	ctx.TcpServer.Start()
}
