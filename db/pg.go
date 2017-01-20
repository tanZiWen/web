package db

import (
    "database/sql"
    "github.com/go-xorm/core"
    "github.com/go-xorm/xorm"
    _ "github.com/lib/pq"
    "github.com/ory-am/osin-storage/storage/postgres"
    "prosnav.com/common/conf"
    "fmt"
)

var (
    Engine  *xorm.Engine
    OMSEngine *xorm.Engine
    storage *postgres.Storage
)

type logAdapter struct{}

func (da *logAdapter) Debug(v ...interface{}) (err error) {
    if len(v) > 1 {
        l.Debug(v[0].(string), v[1:]...)
        return
    }
    l.Debug(v[0].(string))
    return
}

func (da *logAdapter) Debugf(format string, v ...interface{}) (err error) {
    l.Debug(format, v...)
    return
}

func (da *logAdapter) Err(v ...interface{}) (err error) {
    if len(v) > 1 {
        l.Error(v[0].(string), v[1:]...)
        return
    }
    l.Error(v[0].(string))
    return
}
func (da *logAdapter) Errf(format string, v ...interface{}) (err error) {
    l.Error(format, v...)
    return
}
func (da *logAdapter) Info(v ...interface{}) (err error) {
    if len(v) > 1 {
        l.Info("%v", v)
        return
    }
    l.Info(v[0].(string), v[1:]...)
    return
}
func (da *logAdapter) Infof(format string, v ...interface{}) (err error) {
    l.Info(format, v...)
    return
}
func (da *logAdapter) Warning(v ...interface{}) (err error) {
    if len(v) > 1 {
        l.Warn(v[0].(string), v[1:]...)
        return
    }
    l.Warn(v[0].(string), v[1:]...)
    return
}
func (da *logAdapter) Warningf(format string, v ...interface{}) (err error) {
    l.Warn(format, v...)
    return
}

func (da *logAdapter) Level() core.LogLevel {
    return core.LOG_DEBUG
}
func (da *logAdapter) SetLevel(l core.LogLevel) (err error) {
    return nil
}

func (da *logAdapter) ShowSQL(show ...bool) (err error) {
    return nil
}

func (da *logAdapter) IsShowSQL() bool {
    return true
}

func GetStorage() *postgres.Storage {
    if storage == nil {
        storage = postgres.New(Engine.DB().DB)
    }
    return storage
}

func init() {
    registerInitFun(func() {
        var err error
        fmt.Println(conf.String("database.postgresql", "DSN"));
        Engine, err = xorm.NewEngine(sql.Drivers()[0], conf.String("database.postgresql", "DSN"))
        if err != nil {
            panic(err)
        }
        if conf.ENV != "release" {
            Engine.ShowSQL()
        }
        Engine.SetMaxOpenConns(conf.Int("database.postgresql", "MAX_CONNECTION", 10))
        Engine.SetMaxIdleConns(conf.Int("database.postgresql", "MAX_IDLE_CONNECTION", 50))

        OMSEngine, err = xorm.NewEngine(sql.Drivers()[0], conf.String("database.oms", "DSN"))
        if err != nil {
            panic(err)
        }
        if conf.ENV != "release" {
            OMSEngine.ShowSQL()
        }
        OMSEngine.SetMaxOpenConns(conf.Int("database.oms", "MAX_CONNECTION", 10))
        OMSEngine.SetMaxIdleConns(conf.Int("database.oms", "MAX_IDLE_CONNECTION", 50))
    })
}
