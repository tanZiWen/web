package main

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "prosnav.com/common/conf"
    "prosnav.com/common/log"
    "prosnav.com/web/controllers"
    "prosnav.com/web/db"
    "prosnav.com/web/midwares"
    "prosnav.com/common/oauth"
    "github.com/fvbock/endless"
)

func start() {
    app := gin.New()
    gin.SetMode(conf.ENV)
    appContext := fmt.Sprintf("/%s", conf.AppName)
    rootGroup := app.Group(appContext)
    version := fmt.Sprintf("/%s", conf.AppVer)
    verGroup := rootGroup.Group(version)
    verGroup.Use(
        midwares.LoggerHandler(),
        midwares.ErrorHandler(),
    )

    controllers.RegisterHandlers(verGroup)
    endless.ListenAndServe(conf.ListeningPort, app)
}

func init() {
    conf.Init("app.ini")
    log.Init()
    db.Init()
    oauth.Init()
    midwares.Init()
}

func main() {
    start()
}
