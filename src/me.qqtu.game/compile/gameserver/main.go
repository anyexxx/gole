package main

import (
	"code.google.com/p/goconf/conf"
	"me.qqtu.game/handler"
	"me.qqtu.game/logger"
	"me.qqtu.game/message"
	"me.qqtu.game/mio"
	"me.qqtu.game/notifier"
	"me.qqtu.game/router"
	"me.qqtu.game/server"
	"os"
	"runtime"
)

var config = `
[server]
host=127.0.0.1
port=2000
[mongoDB]
seed=127.0.0.1:27017
user=""
pwd="" 
[notifier]
host=127.0.0.1
port=2001
`

func main() {
	//设置cpu数量
	runtime.GOMAXPROCS(runtime.NumCPU())

	//1.读取服务器配置
	cfg, _ := conf.ReadConfigBytes([]byte(config))
	shost, _ := cfg.GetString("server", "host")
	sport, _ := cfg.GetString("server", "port")
	mseed, _ := cfg.GetString("mongoDB", "seed")
	muser, _ := cfg.GetString("mongoDB", "user")
	mpwd, _ := cfg.GetString("mongoDB", "pwd")
	nhost, _ := cfg.GetString("notifier", "host")
	nport, _ := cfg.GetString("notifier", "port")
	logger.GetLogger().Info("StartServer:1.Load config.", logger.Extras{
		"host":        shost,
		"port":        sport,
		"mongodbURL":  mseed,
		"mongodbUser": muser,
		"mongodbPwd":  mpwd})

	//2. init mongoDB connection
	logger.GetLogger().Info("StartServer:2.Init MongoDB connection.", nil)
	mio.InitMongoDB(mseed, muser, mpwd)

	//3.init router config
	logger.GetLogger().Info("StartServer:3.Init RouterConfig.", nil)
	router.InitRouter(shost, sport)

	//4.init handlers
	logger.GetLogger().Info("StartServer:4.Init Handlers.", nil)
	handler.InitHandlers()

	//5.init notifiers and console
	logger.GetLogger().Info("StartServer:5.Init Console and Notifier connection.", nil)
	notifier.InitNotifierServer(nhost, nport)
	notifier.ServerLogin()

	//6.run server
	logger.GetLogger().Info("StartServer:6.Init Server.", logger.Extras{"pid": os.Getpid()})
	sc := server.CreateSocketServer(router.CurrHost, router.CurrPort, message.DATAFORMAT_JSON)
	sc.Run()
}
