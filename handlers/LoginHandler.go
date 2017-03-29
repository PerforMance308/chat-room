package handlers

import (
	"tcp_server"
	"proto_struct"
	"github.com/golang/protobuf/proto"
	"logger"
)

func HandleRoleLoginC2S(c *tcp_server.Client, rqParam proto_struct.RoleLoginC2S) {
	c.Write(EncodeRoleLoginS2C())
}

func EncodeRoleLoginS2C() *[]byte {
	msg := &proto_struct.RoleLoginS2C{}
	data, err := proto.Marshal(msg)
	if err != nil {
		logger.Error("encode EncodeRoleLoginS2C error")
	}
	return &data
}