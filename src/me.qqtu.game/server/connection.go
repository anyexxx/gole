package server

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"io"
	"me.qqtu.game/logger"
	"me.qqtu.game/message"
	"net"
)

//统一各种协议的连接
type Conn interface {
	read()
	write()
	Close()
	SendMessage(msg message.Message)
	User() *User //与连接绑定的客户端
	SetUser(u *User)
	Start()
}

type ConnBase struct {
	dataFormat int
	send       chan message.Message
	user       *User
}

//向连接发送消息
func (c *ConnBase) SendMessage(msg message.Message) {
	if msg != nil {
		c.send <- msg
	}
}

func (c *ConnBase) Close() {
	defer func() {
		if c.user != nil {
			c.user.CloseHandler(c.user)
		}
		if err := recover(); err != nil {
			logger.GetLogger().Info("Makesure Connection Closed!", nil)
		}
	}()
	close(c.send)
}

func (c *ConnBase) User() *User {
	return c.user
}

func (c *ConnBase) SetUser(u *User) {
	c.user = u
}

type SocketConn struct {
	ConnBase
	conn net.Conn
}

func (c *SocketConn) read() {
	recv := make([]byte, 1024*200)

	//make sure close connection
	defer c.Close()

	for {
		length, err := c.conn.Read(recv)

		if err != nil {
			if err == io.EOF {
				logger.GetLogger().Info("Socket connection close!", logger.Extras{"addr": c.conn.RemoteAddr().String()})
			} else {
				logger.GetLogger().Error("Socket connection read error!", logger.Extras{"error": err.Error(), "conn": c.conn.RemoteAddr().String()})
			}
			break
		}

		//按格式读取消息
		msg := message.ReadMessage(recv[:length], c.dataFormat)
		msgId := msg.MessageId()
		fmt.Println("recev messageId:", msgId)
		if msgId >= 0 && HasMsgHandler(msgId) {
			handler := GetMsgHandler(msgId)
			go handler(c, msg)
		}
	}
}

func (c *SocketConn) write() {
	//make sure conn close
	defer c.Close()

	for msg := range c.send {
		if _, err := c.conn.Write(msg.ProtocolBytes()); err != nil {
			logger.GetLogger().Error("Socket connection write error!", logger.Extras{"error": err.Error(), "conn": c.conn.RemoteAddr().String()})
			break
		}
	}
}

func (c *SocketConn) Close() {

	//close send channel
	c.ConnBase.Close()

	//close connection
	c.conn.Close()
}

func (c *SocketConn) Start() {
	go c.write()
	c.read()
}

type WSConn struct {
	ConnBase
	ws *websocket.Conn
}

func (c *WSConn) read() {
	defer c.Close()

	recv := make([]byte, 1024)
	for {
		if err := websocket.Message.Receive(c.ws, &recv); err != nil {
			logger.GetLogger().Error("WebSocket connection read error!", logger.Extras{"error": err, "conn": c.ws.RemoteAddr().String()})
			break
		}

		//按格式读取消息
		msg := message.ReadMessage(recv, c.dataFormat)
		msgId := msg.MessageId()
		if msgId >= 0 && HasMsgHandler(msgId) {
			handler := GetMsgHandler(msgId)
			go handler(c, msg)
		}
	}
}

func (c *WSConn) write() {
	defer c.Close()

	for msg := range c.send {
		if err := websocket.Message.Send(c.ws, msg.ProtocolBytes()); err != nil {
			logger.GetLogger().Error("WebSocket connection write error!", logger.Extras{"error": err, "conn": c.ws.RemoteAddr().String()})
			break
		}
	}
}

func (c *WSConn) Close() {
	//
	c.ConnBase.Close()
	//
	c.ws.Close()
}

func (c *WSConn) Start() {
	go c.write()
	c.read()
}

//处理websocket新连接
func NewWSConnection(ws *websocket.Conn, dataFormat int) *WSConn {
	logger.GetLogger().Info("New Websocket connection!", logger.Extras{"conn": ws.RemoteAddr().String()})
	wsConn := &WSConn{ConnBase: ConnBase{dataFormat: dataFormat, send: make(chan message.Message, 256)}, ws: ws}
	return wsConn
}

//处理socket新连接
func NewSocketConnection(conn net.Conn, dataFormat int) *SocketConn {
	logger.GetLogger().Info("New socket connection!", logger.Extras{"conn": conn.RemoteAddr().String()})
	socketConn := &SocketConn{ConnBase: ConnBase{dataFormat: dataFormat, send: make(chan message.Message, 256)}, conn: conn}
	return socketConn
}
