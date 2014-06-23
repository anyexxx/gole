package server

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"me.qqtu.game/logger"
	"net"
	"net/http"
	"strings"
)

//同质泛滥，为新不破
//1.无服务器区分，全游戏大世界架构
//2.游戏服务器做了负载均衡，增强服务器稳定性以及易于横向扩展
//3.玩家数据按账号区分在不同的redis实例上，对redis的读写与单组服务器相同，保证效率
//4.采用了MongoDB集群，对全游戏活动排名等活动有高效存储保证
//5.目前支持pb和json数据格式以及socket,websocket连接，切换方便

//创建websocket服务
//pattern 服务路径
//port 端口
//dataformat 传输数据格式，为常量DATAFORMAT_JSON等
func CreateWSServer(pattern string, port uint16, dataFormat int) *WSServer {
	ret := &WSServer{Pattern: pattern, Port: port, DataFormat: dataFormat}
	logger.GetLogger().Info("Create WebSocketServer:", logger.Extras{"pattern": pattern, "port": port, "dataformat": dataFormat})
	return ret
}

func CreateSocketServer(host string, port string, dataFormat int) *SocketServer {
	ret := &SocketServer{Host: host, Port: port, DataFormat: dataFormat}
	logger.GetLogger().Info("Create SocketServer:", logger.Extras{"host": host, "port": port, "dataformat": dataFormat})
	return ret
}

func CreateWSClient(url string, origin string, dataFormat int) *WSClient {
	ret := &WSClient{Url: url, Origin: origin, DataFormat: dataFormat}
	logger.GetLogger().Info("Create WebSocketClient:", logger.Extras{"url": url, "origin": origin, "dataformat": dataFormat})
	return ret
}

func CreateSocketClient(host string, port string, dataFormat int) *SocketClient {
	ret := &SocketClient{Host: host, Port: port, DataFormat: dataFormat}
	logger.GetLogger().Info("Create SocketClient:", logger.Extras{"host": host, "port": port, "dataformat": dataFormat})
	return ret
}

//------------ --------------client-------------------------
type WSClient struct {
	Url        string
	Origin     string
	Conn       *WSConn
	DataFormat int
}

func (client *WSClient) Run() {
	ws, err := websocket.Dial(client.Url, "", client.Origin)
	if err == nil {
		logger.GetLogger().Info("WebSocketClient Run!", nil)
		client.Conn = NewWSConnection(ws, client.DataFormat)
		client.Conn.Start()
	} else {
		logger.GetLogger().Error(err.Error(), nil)
	}
}

type SocketClient struct {
	Host       string
	Port       string
	Conn       *SocketConn
	DataFormat int
}

func (client *SocketClient) Run() {
	remote := client.Host + ":" + client.Port
	con, err := net.Dial("tcp", remote)
	if err == nil {
		logger.GetLogger().Info("SocketClient Run!", nil)
		client.Conn = NewSocketConnection(con, client.DataFormat)
		go client.Conn.Start()
	} else {
		logger.GetLogger().Error(err.Error(), nil)
	}
}

//--------------------------websocket--------------------------------
type WSServer struct {
	Pattern    string
	Port       uint16
	DataFormat int
}

func (server *WSServer) Run() {
	defer logger.GetLogger().HandlePanic()

	if server.Pattern == "" {
		server.Pattern = "/socket"
	}
	if strings.Index(server.Pattern, "/") != 0 {
		server.Pattern = "/" + server.Pattern
	}
	http.Handle(server.Pattern, websocket.Handler(server.handle))
	if err := http.ListenAndServe(fmt.Sprint(":", server.Port), nil); err != nil {
		panic(err)
	} else {
		logger.GetLogger().Info("WebSocketServer Run!", logger.Extras{"pattern": server.Pattern})
	}
}

func (server *WSServer) handle(ws *websocket.Conn) {
	wsConn := NewWSConnection(ws, server.DataFormat)
	wsConn.Start()
}

//-------------------------socket--------------------------------------
type SocketServer struct {
	Host       string
	Port       string
	DataFormat int
}

func (server *SocketServer) Run() {
	remote := server.Host + ":" + server.Port
	listen, err := net.Listen("tcp", remote)

	//发生错误关闭连接，输出错误日志
	defer func() {
		listen.Close()
		logger.GetLogger().HandlePanic()
	}()

	if err != nil {
		panic(err)
	}

	logger.GetLogger().Info("SocketServer Run!", logger.Extras{"host": server.Host, "port": server.Port})

	// 等待客户端连接
	for {
		conn, err := listen.Accept()
		if err != nil {
			logger.GetLogger().Error("SocketServer accept error:", logger.Extras{"error": err})
			continue
		}
		socketConn := NewSocketConnection(conn, server.DataFormat)
		go socketConn.Start()
	}
}
