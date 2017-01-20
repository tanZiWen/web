package userservice

import (
    "github.com/go-errors/errors"
    "prosnav.com/common/log"
    "prosnav.com/common/utils"
    "prosnav.com/web/db"
    "gopkg.in/mgo.v2/bson"
    "prosnav.com/web/domain"
)

var l = log.NewLogger()

func Login(username, password string) error {
    if username == "" || password == "" {
        return errors.New("Username or password is empty.")
    }
    encryptPassword := utils.EncryptPassword(password)
    user := new(domain.UpmUser)
    user.UserName = username
    user.Password = encryptPassword
    num, err := db.OMSEngine.Count(user)
    if err != nil {
        l.Error("Login failed.\n%v", err)
        return err
    }
    if num != 1 {
        l.Error("Count of user: %s is %d", username, num)
        return errors.New("User doesn't exist.")
    }
    return nil
}

func QueryUser(username string) (*domain.UpmUser, error) {
    sql := `select user_id, workgroup_id, workgroup_role_codes, email,  role_codes from upm_user`
    var rows []domain.UpmUser
    err := db.Engine.Sql(sql).Find(&rows)
    if err != nil {
        l.Debug("query user error %v", err)
        return nil, err
    }
    if len(rows) == 0 {
        return nil, nil
    }
    return &rows[0], nil
}

func AddUser(user *domain.UpmUser) (int64, error) {
    return db.OMSEngine.InsertOne(user)
}

func QueryUserInfo(username string, appCodes []string) (*domain.UpmUser, error) {
    user := new(domain.UpmUser)
    user.UserName = username
    if exists, err := db.OMSEngine.Get(user); err != nil || !exists{
        if err !=  nil {
            l.Error("Find userinfo failed.\n%v", err)
            return nil, err
        }
        l.Error("The user with username: %s does not exist.\n", username)
        return nil, err
    }
    var org map[string]interface{}
    if err := db.Mongo.DB("").C("upm_organization").Find(bson.M{"code":user.OrgCode}).One(&org); err != nil {
        l.Error("Find userinfo failed.\n%v", err)
        return nil, err
    }
    var roleList []map[string]interface{}
    roleCodes := []string(user.RoleCodes)
    l.Debug("roleCodes:", roleCodes, "\n", "appCodes:", appCodes)
    if err := db.Mongo.DB("").C("upm_role").Find(bson.M{"code":bson.M{"$in":roleCodes}, "appCode":bson.M{"$in":appCodes}}).All(&roleList); err != nil {
        l.Error("Query roles failed.\n%v", err)
        return nil, err
    }
    var funcCodes []interface{}
    for _, role := range roleList {
        funcCodes = append(funcCodes, role["fnCodes"].([]interface{})...)
    }
    var funcList []map[string]interface{}
    if err := db.Mongo.DB("").C("upm_function").Find(bson.M{"code":bson.M{"$in":funcCodes}, "appCode":bson.M{"$in":appCodes}}).All(&funcList); err != nil {
        l.Error("Query functions failed.\n%v", err)
        return nil, err
    }
    user.Org = org
    if roleList != nil && len(roleList) > 0{
        user.RoleList = roleList
        user.RoleMap = parseMap("code", roleList)
    }
    if funcList != nil && len(funcList) > 0 {
        user.FunctionList = funcList
        user.FunctionMap = parseMap("code", funcList)
    }

    db.Mongo.Refresh()
    return user, nil
}

func parseMap(keyName string, src []map[string]interface{}) map[string]interface{} {
    var dest = map[string]interface{}{}
    for _, inst := range src {
        dest[inst[keyName].(string)] = inst
    }
    return dest
}
/**
    {
        pageIndex: 1,
        pageCount: 10,
        userInfo: "张",
        realName: "",
        area: "SH"
    }

*/

func paramFilter(params domain.UserForm) bson.M {
    conds := []bson.M{}
    areaConds := []bson.M{}
    if params.UserInfo != nil {
        switch params.UserInfo.(type) {
        case string: conds = append(conds, bson.M{"realName": bson.M{"$regex": bson.RegEx{params.UserInfo.(string), "i"} }})
        default: conds = append(conds, bson.M{"employeeid": params.UserInfo})
        }
    }
    if len(params.Area) > 0 {
        areaConds = append(areaConds, bson.M{"area": params.Area})
    }
    if len(params.RealName) > 0 {
        areaConds = append(areaConds, bson.M{"realName": params.RealName})
    }
    if len(conds) == 0 {
        if len(areaConds) == 0 {
            return bson.M{}
        }
        return bson.M{"$and": areaConds}
    }
    if len(areaConds) > 0 {
        conds = append(conds, bson.M{"$and": areaConds})
    }
    return bson.M{"$or": conds}
}

var user_column = bson.M{"_id":1, "realName": 2, "position":3, "area": 4}

func QueryUsers(params domain.UserForm)(map[string]interface{}, error) {
    var (
        users []map[string]interface{}
    )
    user := paramFilter(params)
    from := (params.PageIndex - 1) * params.PageCount
    if err:= db.Mongo.DB("").C("upm_user").Find(user).Select(user_column).Skip(from).Limit(params.PageCount).All(&users); err != nil {
        l.Error("Query users failed.\n%v", err)
        return nil, err
    }
    count, err := db.Mongo.DB("").C("upm_user").Find(user).Count()
    if err != nil {
        l.Error("Query users failed.\n%v", err)
        return nil, err
    }
    result := map[string]interface{} {
        "count": count,
        "users": users,
    }
    return result, nil
}

const LEADER = "leader"
type user struct {
    UserId int `bson:"userid"`
    WorkerRoles []string `bson:"WorkerRoles"`
}
type workGroup struct {
    Id int `bson:"_id"`
    Workers []user `bson:"workers"`
}
func QueryUsersByLeaderId(userid int)([]map[string]interface{}, error) {
    var (
        users []map[string]interface{}
        wg workGroup
    )
    cond := bson.M{
        "workers": bson.M{
            "$elemMatch": bson.M{
                "workerRoles":"leader",
                "userid": userid,
            },
        },
    }
    if err := db.Mongo.DB("").C("upm_workgroup").Find(cond).Select(bson.M{"workers": 1}).One(&wg);err !=  nil {
        l.Error("Query workgroup leader's userid : %d failed.\n%v", userid, err)
        return nil, err
    }
    var useridList []int
    for _, u := range wg.Workers {
        if utils.In(LEADER, u.WorkerRoles) {
            continue
        }
        useridList = append(useridList, u.UserId)
    }

    if err := db.Mongo.DB("").C("upm_user").Find(bson.M{"_id": bson.M{"$in": useridList}}).Select(user_column).All(&users); err != nil {
        l.Error("Query users with userid list: %v failed.\n", useridList, err)
        return nil, err
    }
    return users, nil
}

func QueryUserById(userId int)(map[string]interface{}, error) {
    var user map[string]interface{}
    if err := db.Mongo.DB("").C("upm_user").FindId(userId).One(&user); err != nil {
        l.Error("Query user with _id: %d failed.\n%v", userId, err)
        return nil, err
    }
    return user, nil
}