package controllers

import (
    "github.com/gin-gonic/gin"
    "prosnav.com/web/controllers/oauth"
    "prosnav.com/web/midwares"
)

func RegisterHandlers(verGroup *gin.RouterGroup) {
    verGroup.GET("/authorize", oauth.Authorize)
    verGroup.POST("/authorize", oauth.Authorize)
    verGroup.POST("/token", oauth.RequestAccessToken)
    verGroup.GET("/token", oauth.RequestAccessToken)
    verGroup.GET("/userInfo", midwares.Auth(), oauth.RequestUserInfo)
    verGroup.POST("/users", midwares.Auth(), oauth.QueryUsers)
    verGroup.GET("/workgroup/users", midwares.Auth(), oauth.QueryUsersByLeaderId)
}
