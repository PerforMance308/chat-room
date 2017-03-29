package handlers

import (
	"github.com/golang/protobuf/proto"
	"server/proto_struct"
	"server/tcp_server"
	"server/logger"
	"server/db"
)

type data struct {
	User string
}
func HandleRoleLoginC2S(c *tcp_server.Client, rqParam proto_struct.RoleLoginC2S) {
	data := &data{User:*rqParam.Account}
	if err := db.MDB("player").Insert(data); err!=nil{
		logger.Logger().Error(err)
	}
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