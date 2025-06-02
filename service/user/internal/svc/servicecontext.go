package svc

import (
	"github.com/Qsgs-Fans/freekill-server/common/cryptoutil"
	"github.com/Qsgs-Fans/freekill-server/service/router/sender"
	"github.com/Qsgs-Fans/freekill-server/service/user/internal/config"
)

type ServiceContext struct {
	Config config.Config

	Sender *sender.Sender
	RsaKeyPair *cryptoutil.RSAKeyPair
}

func NewServiceContext(c config.Config) *ServiceContext {
	kpair, err := cryptoutil.LoadOrGenerateRSAKeys("rsakey/rsa", "rsakey/rsa_pub")
	if err != nil {
		panic(err)
	}

	return &ServiceContext{
		Config: c,
		Sender: sender.NewSender(c.RouterRpc),
		RsaKeyPair: kpair,
	}
}
