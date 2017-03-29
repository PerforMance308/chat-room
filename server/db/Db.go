package db

import (
	"gopkg.in/mgo.v2"
	"fmt"
	"server/options"
	"server/logger"
)

type dbStruct struct {
	session *mgo.Session
	name    string
}

var Db *dbStruct

func InitDB(db options.DBConf) {
	var connUrl string

	if db.AuthenticationDatabase != "" {
		if db.User == "" {
			connUrl = fmt.Sprintf("%s/%s", db.Hosts, db.AuthenticationDatabase)
		} else {
			connUrl = fmt.Sprintf("%s:%s@%s/%s", db.User, db.Password, db.Hosts, db.AuthenticationDatabase)
		}
	} else {
		if db.User == "" {
			connUrl = fmt.Sprintf("%s/%s", db.Hosts, db.Name)
		} else {
			connUrl = fmt.Sprintf("%s:%s@%s/%s", db.User, db.Password, db.Hosts, db.Name)
		}

	}

	logger.Logger().Notice("Connecting Mongodb", connUrl)
	sess, err := mgo.Dial(connUrl)
	if err != nil {
		confStr := fmt.Sprintf("Failed to open Mongodb:", db.Name)
		logger.Logger().Error(confStr, err)
	}
	sess.SetMode(mgo.Monotonic, true)

	dbstruct := &dbStruct{
		session: sess,
		name:    db.Name,
	}

	Db = dbstruct
}

func MDB(table string) *mgo.Collection {
	return Db.session.DB(Db.name).C(table)
}
