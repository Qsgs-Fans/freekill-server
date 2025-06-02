package sender

import (
	"context"
	"sync"

	"github.com/Qsgs-Fans/freekill-server/service/router/router"
	"github.com/Qsgs-Fans/freekill-server/service/router/routerclient"
	"github.com/zeromicro/go-zero/zrpc"
)

type Sender struct {
	RouterRpc routerclient.Router

	config zrpc.RpcClientConf
	once sync.Once
}

func NewSender(config zrpc.RpcClientConf) *Sender {
	return &Sender{
		config: config,
	}
}

func (s *Sender) checkConnected() {
	s.once.Do(func() {
		s.RouterRpc = routerclient.NewRouter(zrpc.MustNewClient(s.config))
	})
}

func (s *Sender) Notify(ctx context.Context, packet *router.Packet) error {
	s.checkConnected()
	_, err := s.RouterRpc.NotifyClient(ctx, packet)
	return err
}

func (s *Sender) Request(ctx context.Context, packet *router.RequestPacket) error {
	s.checkConnected()
	_, err := s.RouterRpc.RequestClient(ctx, packet)
	return err
}
