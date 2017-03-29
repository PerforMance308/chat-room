package main

import (
	"server/logger"
	"server/tcp_server"
	"sync"
	"server/handlers"
	"server/util"
	"server/options"
)

func main() {
	var wg sync.WaitGroup

	opt := "options/server.json"
	conf := util.GetJSONConfig(&opt, &options.ServerConf{}).(*options.ServerConf)

	logger.InitLogger()

	server := tcp_server.NewServer(conf.TcpConf, &wg)

	for id, fun := range handlers.HandlerList() {
		server.AddRequestHandler(uint16(id), fun)
	}

	conf.DbConf.InitDB()

	go server.Start()

	wg.Wait()
}
