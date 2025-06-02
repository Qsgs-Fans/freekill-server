package other

import (
	"bufio"
	"encoding/json"
	"io"
	"net"
	"time"

	"github.com/Qsgs-Fans/freekill-server/service/router/router"
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

	// TODO 登录

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

		self.handlePacket(line)
	}
}

func (self *TcpConn) handlePacket(line []byte) {
	var rawpacket []any
	err := json.Unmarshal(line, &rawpacket)
	if err != nil {
		// 好像不用发log，发不发呢
		return
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
			// 发log吗?
			return
		}

		self.expectedReplyId = -1
		// TODO
	}
}

func (self *TcpConn) notify(packet *router.Packet) {
	tmpPacket := []any{
		-2,
		t_NOTIFICATION,
		packet.Command,
		packet.Data,
	}
	rawLine, err := json.Marshal(tmpPacket)
	if err != nil {
		// 发log?
		return
	}
	self.send(rawLine)
}

func (self *TcpConn) request(packet *router.Packet, timeout int, timestamp int64) {
	self.requestId++
	self.expectedReplyId = self.requestId
	self.replyTimeout = timeout
	self.requestStartTime = timestamp
	if (timestamp < 0) {
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
		// 发log?
		return
	}
	// 错误处理？
	self.send(rawLine)
}

func (self *TcpConn) send(msg []byte) error {
	msg = append(msg, '\n')
	// TODO 压缩&加密传输
	_, err := self.conn.Write(msg)
	return err
}
