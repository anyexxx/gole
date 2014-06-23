package main

import (
	"code.google.com/p/go.net/websocket"
	"me.qqtu.game/message"
	// "code.google.com/p/goprotobuf/proto"
	"encoding/json"
	"fmt"
	"io"
	// "me.qqtu.game/pb"
	"me.qqtu.game/rio"
	"net"
	"os"
	"time"
)

var writeStr, readStr = make([]byte, 1024), make([]byte, 1024)

func main() {
	// runWSServer()
	runSocketServer()
	// runReidsServer()
}

func runReidsServer() {
	dbReids := rio.GetRedisByUserId(0)
	dbReids.Conn()
	dbReids.Incr("runReidsServer")
}

func runSocketServer() {
	var (
		host   = "127.0.0.1"
		port   = "2000"
		remote = host + ":" + port
	)
	con, err := net.Dial("tcp", remote)

	if err != nil {
		fmt.Println("Server not found.")
		os.Exit(-1)
	}
	fmt.Println("Connection OK.")

	go read(con)

	msg := new(message.CGLoginData)
	msg.Account = "client1"
	msg.Pwd = "pwd1"
	str, _ := json.Marshal(msg)

	kk := make([]byte, 2)
	kk[0] = 0
	kk[1] = 1
	kk = append(kk, []byte(str)...)
	con.Write(kk)
}

func read(conn net.Conn) {
	for {
		length, err := conn.Read(readStr)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Server closed. Exiting...")
			} else {
				fmt.Printf("Error when read from server. Error:%s\n", err)
			}
			os.Exit(0)
		}
		fmt.Println(string(readStr[2:length]))
		// data := readStr[2:length]
		// pbData := new(pb.BaseMessage)
		// proto.Unmarshal(data, pbData)
		// fmt.Println("recev:", pbData.GetHeight())
	}
}

func Handle(ws *websocket.Conn) {
	var msg []byte = make([]byte, 1024*256)
	for {
		if err := websocket.Message.Receive(ws, &msg); err != nil {
			continue
		}
		fmt.Println("received client:", string(msg[2:]))
	}
}

func runWSServer() {
	origin := "http://localhost/"
	url := "ws://localhost:2000/socket"
	if ws, err := websocket.Dial(url, "", origin); err == nil {
		go Handle(ws)
		ticker := time.NewTicker(time.Second)
		for t := range ticker.C {
			dm := make(map[string]int)
			dm["data"] = t.Second()
			str, _ := json.Marshal(&dm)
			data := make([]byte, 2)
			data[0] = 0
			data[1] = 1
			data = append(data, []byte(str)...)
			websocket.Message.Send(ws, data)
		}
	}
}
