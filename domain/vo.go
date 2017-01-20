package domain

import (
    "strings"
    "fmt"
)

type UserForm struct {
    PageIndex int   `json:"pageIndex" form:"pageIndex"`
    PageCount int   `json:"pageCount" form:"pageCount"`
    UserInfo  interface{} `json:"userInfo" form:"userInfo"`
    Area      string  `json:"area" form:"area"`
    RealName  string `json:"realName" form:"realName"`
}

type ArrayString []string
func (s *ArrayString) FromDB(data []byte) error {
    str := fmt.Sprintf("%s", string(data))
    str = str[1:len(str) - 1]
    if len(strings.TrimSpace(str)) == 0 {
        return nil
    }
    var out []string
    for _, el := range strings.Split(str, ",") {
        out = append(out, strings.TrimSpace(el))
    }
    *s = ArrayString(out)
    return nil
}

func (s *ArrayString)ToDB() ([]byte, error) {
    out := strings.Join([]string(*s), `,`)
    out = fmt.Sprintf(`{%s}`, out)
    return []byte(out), nil
}


type UpmUser struct {
    UserId             int64 `xorm:"user_id" json:"_id"`
    WorkGroupId        int64 `xorm:"workgroup_id" json:"workgroupId,omitempty"`
    WorkGroupRoleCodes ArrayString `xorm:"workgroup_role_codes" json:"workgroupRoleCodes,omitempty"`
    Email              string `xorm:"email" json:"email,omitempty"`
    RoleCodes          ArrayString `xorm:"role_codes" json:"roleCodes"`
    Status             string `xorm:"status" json:"status,omitempty"`
    UserName           string `xorm:"user_name" json:"username,omitempty"`
    Password           string `xorm:"password" json:"-"`
    RealName           string `xorm:"real_name" json:"realName,omitempty"`
    Sex                string `xorm:"sex" json:"sex,omitempty"`
    Position           string `xorm:"position" json:"position,omitempty"`
    Employeeid         int64 `xorm:"employee_id" json:"employeeid,omitempty"`
    WorkNo             string `xorm:"work_no" json:"workno,omitempty"`
    EmployeeCode       string `xorm:"employee_code" json:"employeeCode,omitempty"`
    ExtNo              string `xorm:"ext_no" json:"extno,omitempty"`
    OrgCode            string `xorm:"org_code" json:"orgCode,omitempty"`
    Online             bool `xorm:"online" json:"online,omitempty"`
    Area               string `xorm:"area" json:"area,omitempty"`
    RoleMap            map[string]interface{} `xorm:"-" json:"roleMap,omitempty"`
    RoleList           []map[string]interface{} `xorm:"-" json:"roleList,omitempty"`
    FunctionMap        map[string]interface{} `xorm:"-" json:"functionMap,omitempty"`
    FunctionList       []map[string]interface{} `xorm:"-" json:"functionList,omitempty"`
    Org                map[string]interface{} `xorm:"-" xorm:"org" json:"org,omitempty"`
}

func (_ *UpmUser) TableName() string {
    return "upm_user"
}