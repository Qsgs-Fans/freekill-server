package main

import (
	"flag"
	"fmt"

	"github.com/Qsgs-Fans/freekill-server/service/router/internal/config"
	"github.com/Qsgs-Fans/freekill-server/service/router/internal/server"
	"github.com/Qsgs-Fans/freekill-server/service/router/internal/svc"
	"github.com/Qsgs-Fans/freekill-server/service/router/router"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/router.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	go ctx.TcpServer.Start()
	// TODO udpServer

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		router.RegisterRouterServer(grpcServer, server.NewRouterServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
