package svc

import (
	"encoding/json"
	"net"
	"sync"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/Qsgs-Fans/freekill-server/service/router/sender"

	amqp "github.com/rabbitmq/amqp091-go"
)

type TcpServer struct {
	svcCtx *ServiceContext
	lister net.Listener
	connections sync.Map // k: connId, v: *TcpConn
	// TODO: rate limit

	mqconn *amqp.Connection
}

func NewTcpServer(svcCtx *ServiceContext) *TcpServer {
	return &TcpServer {
		svcCtx: svcCtx,
	}
}

func (s *TcpServer) StartListenOnMQ() {
	var err error
	s.mqconn, err = amqp.Dial(s.svcCtx.Config.Amqp)
	if err != nil {
		logx.Errorf("Could not connect to RabbitMQ: %v", err)
		return
	}
	defer s.mqconn.Close()

	ch, err := s.mqconn.Channel()
	if err != nil {
		logx.Errorf("RabbitMQ channel creation failed: %v", err)
		return
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"freekill-router", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		logx.Errorf("RabbitMQ failed to declare a queue: %v", err)
		return
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		logx.Errorf("RabbitMQ failed to register a consumer: %v", err)
		return
	}

	for d := range msgs {
		var packet sender.Packet
		err := json.Unmarshal(d.Body, &packet)
		if err != nil {
			logx.Errorf("RabbitMQ: message unmarshal failed: %v", err)
			continue
		}

		connId := packet.ConnectionId
		connRaw, ok := s.connections.Load(connId)
		if !ok {
			logx.Errorf("RabbitMQ: cannot find conn by connId %v", connId)
			continue
		}
		conn, ok := connRaw.(*TcpConn)
		if !ok {
			continue
		}
		if packet.Type == "notify" {
			conn.Notify(packet.Command, packet.Data)
		} else if packet.Type == "request" {
			conn.Request(packet.Command, packet.Data, int64(packet.Timeout), packet.Timestamp)
		}
	}
}

func (s *TcpServer) Start() {
	var err error

	s.lister, err = net.Listen("tcp", s.svcCtx.Config.TcpListenOn)
	if err != nil {
		logx.Errorf("Tcp listen error: %v", err)
		return
	}
	defer s.lister.Close()
	logx.Infof("Tcp listening at %v", s.svcCtx.Config.TcpListenOn)

	for {
		conn, err := s.lister.Accept()
		if err != nil {
			logx.Errorf("Tcp Accept error: %v", err)
			continue
		}

		// TODO: SYN rate limiter

		connId := uuid.New().String()
		tcpConn := NewTcpConn(conn, connId, s)
		s.connections.Store(connId, tcpConn)

		go tcpConn.listen()
	}
}

func (s *TcpServer) GetConn(connId string) (*TcpConn) {
	val, ok := s.connections.Load(connId)
	if !ok {
		return nil
	}
	return val.(*TcpConn)
}

func (s *TcpServer) notifyLoginRequest(connId string) {
	// TODO: 等userRpc施工...
}
