package db

import (
	"time"
	"gopkg.in/redis.v2"
	"server/util"
	"server/logger"
)

type DBConf struct {
	MongoDbConf map[string]*MongoDbConf
	RedisDbConf map[string]*RedisConf
}

func (db *DBConf) InitDB() {
	for key, conf := range db.MongoDbConf {
		logger.Logger().Notice("mongo connect to", key)
		conf.Connect()
	}

	var dOptions = redis.Options{
		Network:      "tcp",
		Addr:         "localhost:6379",
		Password:     "",
		DB:           0,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		PoolSize:     20,
		IdleTimeout:  60 * time.Second}

	for key, conf := range db.RedisDbConf {
		dst := &redis.Options{}
		util.MergeSimpleStruct(dst, dOptions)
		util.MergeSimpleStruct(dst, conf.RedisDBConf.Options)
		conf.RedisDBConf.Options = dst
		c := redis.NewClient(conf.RedisDBConf.Options)
		logger.Logger().Notice("redis connect to db", key, conf.RedisDBConf.Options.Addr, conf.RedisDBConf.Options.DB)
		conf.RedisDBConf.client = c
	}
}
