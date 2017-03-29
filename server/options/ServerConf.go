package options

import "server/db"

type TCPConf struct {
	Host string
	Port int
}

type ServerConf struct {
	TcpConf TCPConf
	DbConf  db.DBConf
}
