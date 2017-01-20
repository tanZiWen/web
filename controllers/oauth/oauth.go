package oauth

import (
    "github.com/RangelReale/osin"
    "github.com/dgrijalva/jwt-go"
    "github.com/gin-gonic/gin"
    "net/http"
    "prosnav.com/common/conf"
    "prosnav.com/common/log"
    "prosnav.com/common/oauth"
    "prosnav.com/web/db"
    "prosnav.com/web/services/userservice"
    "html/template"
    "errors"
    "strings"
    "strconv"
    "prosnav.com/web/domain"
)

var (
    server *osin.Server
    l = log.NewLogger()
)

type AccessTokenGenJWT struct {
    PrivateKey []byte
}

func (c *AccessTokenGenJWT) GenerateAccessToken(data *osin.AccessData, generaterefresh bool) (accesstoken string, refreshtoken string, err error) {
    // generate JWT access token
    token := jwt.New(jwt.GetSigningMethod(oauth.SignedMethod()))
    token.Claims["cid"] = data.Client.GetId()
    token.Claims["aud"] = data.UserData
    token.Claims["exp"] = data.ExpireAt().Unix()

    accesstoken, err = token.SignedString(c.PrivateKey)
    if err != nil {
        l.Error("%v", err)
        return
    }

    if !generaterefresh {
        return
    }

    // generate JWT access token
    token = jwt.New(jwt.GetSigningMethod(oauth.SignedMethod()))
    token.Claims["cid"] = data.Client.GetId()
    token.Claims["aud"] = data.UserData
    token.Claims["at"] = accesstoken
    //token.Claims["exp"] = data.ExpireAt().Unix()

    refreshtoken, err = token.SignedString(c.PrivateKey)
    if err != nil {
        return
    }
    return
}

func Authorize(c *gin.Context) {
    if server == nil {
        serverInstance()
    }
    resp := server.NewResponse()
    defer resp.Close()

    if ar := server.HandleAuthorizeRequest(resp, c.Request); ar != nil {
        if ok := validate(c); !ok {
            return
        }
        if err := loginFilter(c); err != nil {
            c.JSON(http.StatusUnauthorized, err)
            return
        }
        ar.Authorized = true
        tk := c.MustGet("token").(*jwt.Token)
        ar.UserData = tk.Claims["aud"]
        server.FinishAuthorizeRequest(resp, c.Request, ar)
    }

    handleOutput(resp, c)
}

func validate(c *gin.Context) bool {
    tk, err := jwt.ParseFromRequest(c.Request, func(*jwt.Token) (interface{}, error) {
        return oauth.LookupPrivateKey(), nil
    })
    if err != nil {
        forwardLogin(c)
        return false
    }
    if !tk.Valid {
        forwardLogin(c)
        return false
    }
    c.Set("token", tk)
    if tk.Claims["aud"] != nil && !isAuthorized(tk, c.Request.FormValue("client_id")) {
        forwardAuth(c)
        return false
    }
    return true
}

func handleOutput(resp *osin.Response, c *gin.Context) {
    // Add headers
    for k, _ := range resp.Headers {
        c.Writer.Header().Add(k, resp.Headers.Get(k))
    }

    if resp.IsError {
        l.Error("%v", resp.Output)
        c.JSON(500, resp.Output["error_description"])
        return
    }
    if resp.Type == osin.REDIRECT {
        redirectURI, err := resp.GetRedirectUrl()
        if err != nil {
            l.Error("%v", err)
            c.JSON(500, err)
            return
        }
        c.Redirect(http.StatusFound, redirectURI)
        return
    }

    c.JSON(200, resp.Output)
}

func serverInstance() {
    cfg := osin.NewServerConfig()
    cfg.AccessExpiration = int32(conf.Int("sso", "ACCESS_EXPIRATION", 3600))
    cfg.AuthorizationExpiration = int32(conf.Int("sso", "AUTHORIZATION_EXPIRATION", 100))
    cfg.AllowClientSecretInParams = true
    cfg.AllowGetAccessRequest = true
    cfg.AllowedAuthorizeTypes = osin.AllowedAuthorizeType{osin.CODE, osin.TOKEN}
    cfg.AllowedAccessTypes = convertAccessType(conf.Strings("sso", "ALLOWED_ACCESSTYPES", ","))
    server = osin.NewServer(cfg, db.GetStorage())
    server.AccessTokenGen = &AccessTokenGenJWT{PrivateKey:oauth.LookupPrivateKey()}
}

func convertAccessType(typeStrs []string) osin.AllowedAccessType {
    var accessRequestTypes []osin.AccessRequestType
    for _, s := range typeStrs {
        accessRequestTypes = append(accessRequestTypes, osin.AccessRequestType(s))
    }
    return osin.AllowedAccessType(accessRequestTypes)
}

type loginForm struct {
    UserName string `form:"username" 	json:"username"`
    Password string `form:"password" 	json:"password"`
}

func loginFilter(c *gin.Context) error {
    var form loginForm
    c.Bind(&form)
    tk := c.MustGet("token").(*jwt.Token)
    if len(tk.Claims) == 0 {
        err := userservice.Login(form.UserName, form.Password)
        if err != nil {
            l.Error("%v\n", err)
            return err
        }
        tk.Claims["aud"] = form.UserName
    }
    return nil
}

func isAuthorized(tk *jwt.Token, clientId string) bool {
    cid := tk.Claims["cid"]
    if cid == nil {
        return false
    }
    return cid.(string) == clientId
}

func forwardLogin(c *gin.Context) {
    tpl, err := template.ParseFiles("template/login.html")
    if err != nil {
        c.JSON(500, err)
    }
    token := jwt.New(jwt.GetSigningMethod(oauth.SignedMethod()))
    signStr, _ := token.SignedString(oauth.LookupPrivateKey())
    data := map[string]interface{}{
        "target": c.Request.RequestURI,
        "token":  signStr,
    }
    tpl.Execute(c.Writer, data)
    c.Abort()
}

func forwardAuth(c *gin.Context) {
    tpl, err := template.ParseFiles("template/login.html")
    if err != nil {
        c.JSON(500, err)
    }
    token := c.MustGet("token").(*jwt.Token)
    signStr, _ := token.SignedString(oauth.LookupPrivateKey())
    data := map[string]interface{}{
        "target":   c.Request.RequestURI,
        "username": token.Claims["aud"],
        "token":    signStr,
    }
    tpl.Execute(c.Writer, data)
    c.Abort()
}

const ASSERTION_TYPE = "urn:ietf:params:oauth:grant-type:jwt-bearer"

func RequestAccessToken(c *gin.Context) {
    if server == nil {
        serverInstance()
    }
    resp := server.NewResponse()
    defer resp.Close()
    c.Request.ParseForm()
    if ar := server.HandleAccessRequest(resp, c.Request); ar != nil {
        switch ar.Type {
        case osin.AUTHORIZATION_CODE:
        case osin.REFRESH_TOKEN:
        case osin.CLIENT_CREDENTIALS:
        case osin.PASSWORD:
            l.Debug("Start login...")
            if err := userservice.Login(ar.Username, ar.Password); err != nil {
                l.Error("%v\n", err)
                c.JSON(500, "Username or password incorrect.")
                return
            }
            ar.UserData = ar.Username
        case osin.ASSERTION:
            if err := processAssertion(ar); err != nil {
                c.JSON(500, err)
                return
            }
        }
        ar.Authorized = true
        server.FinishAccessRequest(resp, c.Request, ar)
    }
    handleOutput(resp, c)
}

func processAssertion(ar *osin.AccessRequest) error {
    if !(ar.AssertionType == ASSERTION_TYPE) {
        return errors.New("Unsupported assertion type.")
    }
    tk, err := jwt.Parse(ar.Assertion, func(*jwt.Token) (interface{}, error) {
        return oauth.LookupPrivateKey(), nil
    })
    if err != nil {
        l.Error("%v\n", err)
        return err
    }
    if !tk.Valid {
        return errors.New("Invalid access token.")
    }
    if tk.Claims["aud"] == nil {
        return errors.New("Unauthorized user")
    }
    ar.UserData = tk.Claims["aud"]
    return nil
}

var appMap = map[string][]string{
    "PUP": []string{"CRM", "UPM"},
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

func RequestUserInfo(c *gin.Context) {
    if server == nil {
        serverInstance()
    }
    tk := c.MustGet("token").(*jwt.Token)
    cid := tk.Claims["cid"].(string)
    username := tk.Claims["aud"].(string)
    appCodes := appMap[strings.ToUpper(cid)]
    if appCodes == nil {
        appCodes = []string{strings.ToUpper(cid)}
    }
    user, err := userservice.QueryUserInfo(username, appCodes)
    if err != nil {
        l.Error("%v", err)
        c.JSON(500, "Query user failed")
        return
    }
    l.Debug("userinfo: %+v", user)
    c.JSON(200, user)
}

func QueryUsers(c *gin.Context) {
    var data domain.UserForm
    c.Bind(&data)
    users, err := userservice.QueryUsers(data)
    if err != nil {
        l.Error("%v", err)
        c.JSON(500, "Query user failed")
        return
    }
    c.JSON(200, users)
}

func QueryUsersByLeaderId(c *gin.Context) {
    userid, err := strconv.Atoi(c.Query("leaderId"))
    if err != nil {
        l.Error("%v", err)
        c.JSON(500, "Query user failed")
        return
    }
    users, err := userservice.QueryUsersByLeaderId(userid)
    if err != nil {
        l.Error("%v", err)
        c.JSON(500, "Query user failed")
        return
    }
    c.JSON(200, users)
}

func QueryUserById(c *gin.Context) {
    userid, err := strconv.Atoi(c.Query("userid"))
    if err != nil {
        l.Error("%v", err)
        c.JSON(500, "Query user failed")
        return
    }
    user, err := userservice.QueryUserById(userid)
    if err != nil {
        l.Error("%v", err)
        c.JSON(500, "Query user failed")
        return
    }
    c.JSON(200, user)
}