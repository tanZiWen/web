package midwares

import (
    "github.com/gin-gonic/gin"
    "prosnav.com/common/conf"
    "time"
)

func LoggerHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        c.Next()
        end := time.Now()
        latency := end.Sub(start)
        clientIp := c.ClientIP()
        method := c.Request.Method
        statusCode := c.Writer.Status()
        appName := conf.AppName
        l.Info("[%s] %d  %v | %s %s %s", appName, statusCode, latency,
            clientIp, method, c.Request.URL.Path)
    }
}
