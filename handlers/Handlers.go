package handlers

import "proto_struct"

func HandlerList() map[proto_struct.RequestId]interface{} {
	list := make(map[proto_struct.RequestId]interface{})
	list[proto_struct.RequestId_role_login_c2s] = HandleRoleLoginC2S
	return list
}
