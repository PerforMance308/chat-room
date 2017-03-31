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
	msgId, msg := EncodeRoleLoginS2C()
	c.Write(msgId, msg)
}

func EncodeRoleLoginS2C() (uint16, []byte) {
	msg := &proto_struct.RoleLoginS2C{Res:proto.Bool(true)}
	data, err := proto.Marshal(msg)
	if err != nil {
		logger.Logger().Error("encode EncodeRoleLoginS2C error")
	}
	return uint16(proto_struct.RequestId_role_login_s2c), data
}