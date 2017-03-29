package db

import (
	"fmt"
	"server/logger"
	"gopkg.in/mgo.v2"
)

type MongoServerConf struct {
	Hosts                  string
	User                   string
	Password               string
	AuthenticationDatabase string
}

type MongoDbConf struct {
	DBName string
	MongoServerConf
	session *mgo.Session
}

func (mc *MongoDbConf) Connect() {
	var connUrl string

	if mc.AuthenticationDatabase != "" {
		if mc.User == "" {
			connUrl = fmt.Sprintf("%s/%s", mc.Hosts, mc.AuthenticationDatabase)
		} else {
			connUrl = fmt.Sprintf("%s:%s@%s/%s", mc.User, mc.Password, mc.Hosts, mc.AuthenticationDatabase)
		}
	} else {
		if mc.User == "" {
			connUrl = fmt.Sprintf("%s/%s", mc.Hosts, mc.DBName)
		} else {
			connUrl = fmt.Sprintf("%s:%s@%s/%s", mc.User, mc.Password, mc.Hosts, mc.DBName)
		}

	}

	logger.Logger().Notice("Connecting Mongodb", connUrl)
	session, err := mgo.Dial(connUrl)
	if err != nil {
		confStr := fmt.Sprintf("Failed to open Mongodb:", mc.DBName, ", ", mc.MongoServerConf)
		logger.Logger().Error(confStr, err)
	}
	session.DB(mc.DBName)
	session.SetMode(mgo.Monotonic, true)
	session.SetSafe(&mgo.Safe{WMode: "majority"})
	mc.session = session
}
