package options

type TCPConf struct {
	Host string
	Port int
}

type DBConf struct {
	Hosts                  string
	User                   string
	Password               string
	AuthenticationDatabase string
	Name                   string
}

type ServerConf struct {
	TcpConf TCPConf
	DbConf  DBConf
}
