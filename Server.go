package main

import (
	"logger"
	"tcp_server"
	"sync"
	"handlers"
)

func main(){
	var wg sync.WaitGroup

	logger.InitLogger()

	server := tcp_server.NewServer(&wg)

	for id, fun := range handlers.HandlerList() {
		server.AddRequestHandler(uint16(id), fun)
	}

	go server.Start()

	wg.Wait()
}
