package svc

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/Qsgs-Fans/freekill-server/service/router/router"
	"github.com/Qsgs-Fans/freekill-server/service/user/user"
	"github.com/zeromicro/go-zero/core/logx"
)

// enum PacketType
const (
	t_REQUEST int = 0x100
  t_REPLY = 0x200
	t_NOTIFICATION = 0x400
	// 原版的SRC_CLIENT之类的没用，无视
	// 而且C++的Router::handlePacket里面也只检查type了
)

type TcpConn struct {
	conn net.Conn
	// 此为一个UUID字符串。在此服务中对应一个TcpConn，在其他服务则对应玩家
	connId string
	server *TcpServer

	requestId int
	expectedReplyId int
	replyTimeout int
	requestStartTime int64
	replyContent string

	// TODO Aes密钥
}

func NewTcpConn(conn net.Conn, connId string, server *TcpServer) *TcpConn {
	return &TcpConn{
		conn: conn,
		connId: connId, 
		server: server,

		expectedReplyId: -1,
	}
}

func (self *TcpConn) listen() {
	defer self.conn.Close()
  defer self.server.connections.Delete(self.connId)

	if err := self.login(); err != nil {
		logx.Errorf("Login phase failed: %v", err)
		return
	}

	for {
		// 此为心跳包使用。 TODO: 心跳包
		// conn.SetReadDeadline(time.Now().Add(60 * time.Second))

		reader := bufio.NewReader(self.conn)
		line, err := reader.ReadBytes('\n')

		if err != nil {
			if err == io.EOF {
				// TODO: 告诉别的rpc连接断开了
				// logx.Infof("disconnected")
			} else {
				logx.Errorf("Read error: %v", err)
			}
			return
		}

		// TODO: read rate limiter

		// error怎么办呢？
		self.handlePacket(line)
	}
}

func (s *TcpConn) login() error {
	// TODO 登录 网关在这里要做的事主要就是ip黑名单
	// TODO 以及在登录流程中，客户端只能发1个包，通过之后才能进入以下的无限循环
	ctx, cancel := context.WithTimeout(context.Background(), 30 * time.Second)
	defer cancel()

	urpc := s.server.svcCtx.UserRpc
	connIdMsg := &user.ConnIdMsg{
		ConnId: s.connId,
	}
	if _, err := urpc.NewConn(ctx, connIdMsg); err != nil {
		return err
	}

	reader := bufio.NewReader(s.conn)
	rawPacket, err := reader.ReadBytes('\n')
	if err != nil {
		return fmt.Errorf("failed to read login packet: %v", err)
	}

	var packetData []any
	err = json.Unmarshal(rawPacket, &packetData)
	if err != nil {
		return fmt.Errorf("JSON.stringify failed: %v", err)
	}

	loginBytes := packetData[3].(string)
	var loginData []any
	err = json.Unmarshal([]byte(loginBytes), &loginData)
	if err != nil {
		return fmt.Errorf("JSON.stringify failed: %v", err)
	}
	loginPacket := &user.LoginRequest{
		Username: loginData[0].(string),
		Password: loginData[1].(string),
		Md5: loginData[2].(string),
		Version: loginData[3].(string),
		Deviceid: loginData[4].(string),
	}
	// TODO 处理aesKey
	_, err = urpc.Login(ctx, loginPacket)
	if err != nil {
		return fmt.Errorf("Login fail: %v", err)
	}

	return nil
}

func (self *TcpConn) handlePacket(line []byte) error {
	var rawpacket []any
	err := json.Unmarshal(line, &rawpacket)
	if err != nil {
		return fmt.Errorf("JSON.parse error: %v", err)
	}

	// TODO 记得进行类型不匹配测试
	tp := rawpacket[1].(int)

	if tp & t_NOTIFICATION != 0 {
		// TODO: 等其他rpc施工...

		// command := rawpacket[2].(string)
		// jsondata := rawpacket[3].(string)
		// packet := router.Packet {
		// 	Command: command,
		// 	Data: jsonData,
		// 	ConnectionId: self.connId,
		// }
	} else if tp & t_REPLY != 0 {
		reqId := rawpacket[0].(int)
		if reqId != self.expectedReplyId {
			return fmt.Errorf("requestId != expectedReplyId: ignored.")
		}

		self.expectedReplyId = -1
		// TODO
	}

	return nil
}

func (self *TcpConn) Notify(packet *router.Packet) error {
	tmpPacket := []any{
		-2,
		t_NOTIFICATION,
		packet.Command,
		packet.Data,
	}
	rawLine, err := json.Marshal(tmpPacket)
	if err != nil {
		return fmt.Errorf("JSON.stringify failed: %v", err)
	}
	return self.send(rawLine)
}

func (self *TcpConn) Request(packet *router.RequestPacket) error {
	self.requestId++
	self.expectedReplyId = self.requestId
	self.replyTimeout = int(packet.Timeout)
	self.requestStartTime = packet.Timestamp
	if (packet.Timestamp < 0) {
		self.requestStartTime = time.Now().UnixMilli()
	}
	self.replyContent = "__notready"

	tmpPacket := []any{
		self.requestId,
		t_REQUEST,
		packet.Command,
		packet.Data,
		self.replyTimeout,
		self.requestStartTime,
	}
	rawLine, err := json.Marshal(tmpPacket)
	if err != nil {
		return fmt.Errorf("JSON.stringify failed: %v", err)
	}
	return self.send(rawLine)
}

func (self *TcpConn) send(msg []byte) error {
	msg = append(msg, '\n')
	// TODO 压缩&加密传输
	_, err := self.conn.Write(msg)
	return err
}
