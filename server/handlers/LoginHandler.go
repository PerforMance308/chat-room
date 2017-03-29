package handlers

import (
	"github.com/golang/protobuf/proto"
	"server/proto_struct"
	"server/tcp_server"
	"server/logger"
)

func HandleRoleLoginC2S(c *tcp_server.Client, rqParam proto_struct.RoleLoginC2S) {
	c.Write(EncodeRoleLoginS2C())
}

func EncodeRoleLoginS2C() *[]byte {
	msg := &proto_struct.RoleLoginS2C{}
	data, err := proto.Marshal(msg)
	if err != nil {
		logger.Logger().Error("encode EncodeRoleLoginS2C error")
	}
	return &data
}