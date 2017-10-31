package db

import (
	"../config"
	mgo "gopkg.in/mgo.v2"
)

func GetMongo() *mgo.Database {
	//mongo
	session, err := mgo.Dial(config.Conf.MongoDB.Host)
	if err != nil {
		panic(err)
	}

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	mongo := session.DB(config.Conf.MongoDB.DBName)
	return mongo
}
