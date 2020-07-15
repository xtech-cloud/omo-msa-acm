package cache

import (
	"github.com/casbin/casbin/v2"
	mongodbadapter "github.com/casbin/mongodb-adapter/v2"
	"omo.msa.acm/config"
	"omo.msa.acm/proxy/nosql"
	"time"
)

type BaseInfo struct {
	ID         uint64 `json:"-"`
	UID        string `json:"uid"`
	Name       string `json:"name"`
	Creator string
	Operator string
	CreateTime time.Time
	UpdateTime time.Time
}

type cacheContext struct {
	users    []*UserInfo
	roles    []*RoleInfo
	menus    []*MenuInfo
	enforcer *casbin.Enforcer
}

var cacheCtx *cacheContext

func InitData() error {
	cacheCtx = &cacheContext{}
	cacheCtx.roles = make([]*RoleInfo, 0, 10)
	cacheCtx.users = make([]*UserInfo, 0, 100)
	cacheCtx.menus = make([]*MenuInfo, 0, 100)

	err := nosql.InitDB(config.Schema.Database.IP, config.Schema.Database.Port, config.Schema.Database.Name, config.Schema.Database.Type)
	if nil != err {
		return err
	}

	url := config.Schema.Database.IP+":"+config.Schema.Database.Port
	a,err1 := mongodbadapter.NewAdapter(url)
	if err1 != nil {
		return err1
	}
	e, err2 := casbin.NewEnforcer("conf/acm.conf", a)
	if err2 != nil {
		return err2
	}
	cacheCtx.enforcer = e

	users,err2 := nosql.GetAllUsers()
	if err2 == nil {
		for _, user := range users {
			t := new(UserInfo)
			t.initInfo(user)
			cacheCtx.users = append(cacheCtx.users, t)
		}
	}
	return cacheCtx.enforcer.LoadPolicy()
}

