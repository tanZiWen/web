package midwares

import (
    "github.com/gin-gonic/gin"
    "strings"
    "net/http"
    "github.com/dgrijalva/jwt-go"
    "prosnav.com/common/oauth"
)

func Auth() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := checkBearerAuth(c)
        if token == "" {
            l.Error("JWT not found")
            c.JSON(http.StatusUnauthorized, "JWT not found")
            c.Abort()
            return
        }
        tk, err := jwt.Parse(token, func(*jwt.Token) (interface{}, error) {
            return oauth.LookupPrivateKey(), nil
        })
        if err != nil {
            l.Error("%v\n", err)
            c.JSON(http.StatusUnauthorized, "Unauthorized user.")
            c.Abort()
            return
        }
        if !tk.Valid {
            c.JSON(http.StatusUnauthorized, "Invalid access token.")
            c.Abort()
            return
        }
        c.Set("token", tk)
        c.Next()
    }
}

func checkBearerAuth(c *gin.Context) string {
    authHeader := c.Request.Header.Get("Authorization")
    authForm := c.Query("code")
    if authHeader == "" && authForm == "" {
        return ""
    }
    token := authForm
    if authHeader != "" {
        s := strings.SplitN(authHeader, " ", 2)
        if (len(s) != 2 || s[0] != "Bearer") && token == "" {
            return ""
        }
        token = s[1]
    }
    return token
}