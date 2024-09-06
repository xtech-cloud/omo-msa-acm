package cache

import (
	"github.com/casbin/casbin/v2"
	"omo.msa.acm/config"
	"omo.msa.acm/proxy/nosql"
	"time"
)

type BaseInfo struct {
	ID         uint64 `json:"-"`
	UID        string `json:"uid"`
	Name       string `json:"name"`
	Creator    string
	Operator   string
	CreateTime time.Time
	UpdateTime time.Time
}

type cacheContext struct {
	roles    []*RoleInfo
	apiMenus []*APIMenuInfo
	catalogs []*CatalogInfo
	modules  []*ModuleInfo
	scenes   []*SceneInfo
	enforcer *casbin.Enforcer
}

var cacheCtx *cacheContext

func InitData() error {
	cacheCtx = &cacheContext{}
	cacheCtx.roles = make([]*RoleInfo, 0, 10)
	cacheCtx.apiMenus = make([]*APIMenuInfo, 0, 100)
	cacheCtx.catalogs = make([]*CatalogInfo, 0, 100)
	cacheCtx.modules = make([]*ModuleInfo, 0, 100)
	cacheCtx.scenes = make([]*SceneInfo, 0, 50)
	err := nosql.InitDB(config.Schema.Database.IP, config.Schema.Database.Port, config.Schema.Database.Name, config.Schema.Database.Type)
	if nil != err {
		return err
	}

	//url := config.Schema.Database.IP+":"+config.Schema.Database.Port
	//a,err1 := mongodbadapter.NewAdapter(url)
	//if err1 != nil {
	//	return err1
	//}
	//e, err2 := casbin.NewEnforcer("conf/acm.conf", a)
	//if err2 != nil {
	//	return err2
	//}
	//cacheCtx.enforcer = e

	roles, err := nosql.GetAllRoles()
	if err == nil {
		for _, role := range roles {
			t := new(RoleInfo)
			t.initInfo(role)
			cacheCtx.roles = append(cacheCtx.roles, t)
		}
	}

	modules, err1 := nosql.GetAllModules()
	if err1 == nil {
		for _, menu := range modules {
			t := new(ModuleInfo)
			t.initInfo(menu)
			cacheCtx.modules = append(cacheCtx.modules, t)
		}
	}

	catalogs, err3 := nosql.GetAllCatalogs()
	if err3 == nil {
		for _, menu := range catalogs {
			t := new(CatalogInfo)
			t.initInfo(menu)
			cacheCtx.catalogs = append(cacheCtx.catalogs, t)
		}
	}

	scenes, err4 := nosql.GetAllScenes()
	if err4 == nil {
		for _, temp := range scenes {
			t := new(SceneInfo)
			t.initInfo(temp)
			cacheCtx.scenes = append(cacheCtx.scenes, t)
		}
	}
	//return cacheCtx.enforcer.LoadPolicy()
	return nil
}

