package db

import "gopkg.in/redis.v2"

type RedisConf struct {
	RedisDBConf
	Slaves []*RedisDBConf
}

type RedisDBConf struct {
	Key     string
	client *redis.Client
	*redis.Options
}

func (rc RedisDBConf) Client() *redis.Client {
	return rc.client
}