package svc

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"

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

	// login后才可获取的信息
	userId int64
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

func ipToVarBinary(conn net.Conn) ([]byte, error) {
	// 1. 获取 RemoteAddr 并解析 IP
	remoteAddr := conn.RemoteAddr()
	host, _, err := net.SplitHostPort(remoteAddr.String())
	if err != nil {
		return nil, fmt.Errorf("failed to split host:port: %v", err)
	}

	ip := net.ParseIP(host)
	if ip == nil {
		return nil, fmt.Errorf("invalid IP address: %s", host)
	}

	// 2. 转换为 4 字节 (IPv4) 或 16 字节 (IPv6)
	var ipBytes []byte
	if ipv4 := ip.To4(); ipv4 != nil {
		ipBytes = ipv4 // IPv4: 4 字节
	} else {
		ipBytes = ip.To16() // IPv6: 16 字节
	}

	return ipBytes, nil
}

func (s *TcpConn) listen() {
	defer s.conn.Close()
  defer s.server.connections.Delete(s.connId)
	defer s.logout()

	if err := s.login(); err != nil {
		logx.Errorf("Login phase failed: %v", err)
		return
	}

	for {
		// 此为心跳包使用。 TODO: 心跳包
		// conn.SetReadDeadline(time.Now().Add(60 * time.Second))

		reader := bufio.NewReader(s.conn)
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
		s.handlePacket(line)
	}
}

func (s *TcpConn) login() error {
	// TODO 登录 网关在这里要做的事主要就是ip黑名单
	// TODO 以及在登录流程中，客户端只能发1个包，通过之后才能进入以下的无限循环
	ctx, cancel := context.WithTimeout(context.Background(), 30 * time.Second)
	defer cancel()

	ipbytes, err := ipToVarBinary(s.conn)
	if err != nil {
		return err
	}
	urpc := s.server.svcCtx.UserRpc
	connIdMsg := &user.ConnIdMsg{
		ConnId: s.connId,
		ConnIp: ipbytes,
		UserId: 0,
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

		ConnId: s.connId,
		ConnIp: ipbytes,
	}

	loginReply, err := urpc.Login(ctx, loginPacket)
	if err != nil {
		return fmt.Errorf("Login fail: %v", err)
	}

	s.userId = loginReply.UserId
	// TODO 处理aesKey

	return nil
}

func (s *TcpConn) logout() {
	ctx, cancel := context.WithTimeout(context.Background(), 30 * time.Second)
	defer cancel()

	ipbytes, _ := ipToVarBinary(s.conn)
	connIdMsg := &user.ConnIdMsg{
		ConnId: s.connId,
		ConnIp: ipbytes,
		UserId: s.userId,
	}

	urpc := s.server.svcCtx.UserRpc
	urpc.Logout(ctx, connIdMsg)
}

func (s *TcpConn) handlePacket(line []byte) error {
	var rawpacket []any
	err := json.Unmarshal(line, &rawpacket)
	if err != nil {
		return fmt.Errorf("JSON.parse error: %v", err)
	}

	// TODO 记得进行类型不匹配测试
	tpRaw, ok := rawpacket[1].(float64)
	if !ok {
		return fmt.Errorf("JSON error: data[1] should be number: %v", string(line))
	}

	tp := int(tpRaw)

	if tp & t_NOTIFICATION != 0 {
		// TODO: 等其他rpc施工...

		// command := rawpacket[2].(string)
		// jsondata := rawpacket[3].(string)
		// packet := router.Packet {
		// 	Command: command,
		// 	Data: jsonData,
		// 	ConnectionId: s.connId,
		// }
	} else if tp & t_REPLY != 0 {
		reqId := rawpacket[0].(int)
		if reqId != s.expectedReplyId {
			return fmt.Errorf("requestId != expectedReplyId: ignored.")
		}

		s.expectedReplyId = -1
		// TODO
	}

	return nil
}

func (s *TcpConn) Notify(command string, data string) error {
	tmpPacket := []any{
		-2,
		t_NOTIFICATION,
		command,
		data,
	}
	rawLine, err := json.Marshal(tmpPacket)
	if err != nil {
		return fmt.Errorf("JSON.stringify failed: %v", err)
	}
	return s.send(rawLine)
}

func (s *TcpConn) Request(command string, data string, timeout int64, timestamp int64) error {
	s.requestId++
	s.expectedReplyId = s.requestId
	s.replyTimeout = int(timeout)
	s.requestStartTime = timestamp
	if (timestamp < 0) {
		s.requestStartTime = time.Now().UnixMilli()
	}
	s.replyContent = "__notready"

	tmpPacket := []any{
		s.requestId,
		t_REQUEST,
		command,
		data,
		s.replyTimeout,
		s.requestStartTime,
	}
	rawLine, err := json.Marshal(tmpPacket)
	if err != nil {
		return fmt.Errorf("JSON.stringify failed: %v", err)
	}
	return s.send(rawLine)
}

func (s *TcpConn) send(msg []byte) error {
	msg = append(msg, '\n')
	// TODO 压缩&加密传输
	_, err := s.conn.Write(msg)
	return err
}
