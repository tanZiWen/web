package db

import (
    "gopkg.in/mgo.v2"
    "prosnav.com/common/conf"
    "time"
)

var (
    Mongo *mgo.Session
)


func init() {
    registerInitFun(func() {
        var err error
        Mongo, err = mgo.Dial(conf.String("database.mongodb", "DSN"))
        if err != nil {
            panic(err)
        }
        Mongo.SetPoolLimit(conf.Int("database.mongodb", "MAX_CONNECTION", 50))
        Mongo.SetSocketTimeout(30 * time.Second)
    })
}