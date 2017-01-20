package domain

import "time"

type Client struct {
	Id          string `xorm:"id            text pk"`
	Secret      string `xorm:"secret        text notnull"`
	Extra       string `xorm:"extra         text notnull"`
	RedirectUri string `xorm:"redirect_uri  text notnull"`
}

func (c *Client) TableName() string {
	return "client"
}

type Authorize struct {
	Code        string    `xorm:"code           text pk"`
	Client      string    `xorm:"client         text notnull"`
	ExpiresIn   int       `xorm:"expires_in     int  notnull"`
	Scope       string    `xorm:"scope          text notnull"`
	RedirectUri string    `xorm:"redirect_uri   text notnull"`
	State       string    `xorm:"state          text notnull"`
	Extra       string    `xorm:"extra          text notnull"`
	CreatedAt   time.Time `xorm:"created_at     timestampz created"`
}

func (c *Authorize) TableName() string {
	return "authorize"
}

type Access struct {
	AccessToken  string    `xorm:"access_token  text pk"`
	Client       string    `xorm:"client        text notnull"`
	Authorize    string    `xorm:"authorize     text notnull"`
	Previous     string    `xorm:"previous      text notnull"`
	RefreshToken string    `xorm:"refresh_token text notnull"`
	ExpiresIn    int       `xorm:"expires_in    int  notnull"`
	Scope        string    `xorm:"scope         text notnull"`
	RedirectUri  string    `xorm:"redirect_uri  text notnull"`
	Extra        string    `xorm:"extra         text notnull"`
	CreatedAt    time.Time `xorm:"created_at    timestampz created"`
}

func (c *Access) TableName() string {
	return "access"
}

type Refresh struct {
	Token  string `xorm:"token  text pk"`
	Access string `xorm:"access text notnull"`
}

func (c *Refresh) TableName() string {
	return "refresh"
}
