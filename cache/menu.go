package cache

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.acm/proxy/nosql"
	"time"
)

type PermissionAction uint8

type APIMenuInfo struct {
	BaseInfo
	Type string
	Path string
	Method string
}

func AllMenus() []*APIMenuInfo {
	return cacheCtx.apiMenus
}

func GetMenu(uid string) *APIMenuInfo {
	for i := 0;i < len(cacheCtx.apiMenus);i += 1 {
		if cacheCtx.apiMenus[i].UID == uid {
			return cacheCtx.apiMenus[i]
		}
	}
	db,err := nosql.GetMenu(uid)
	if err == nil {
		info := new(APIMenuInfo)
		info.initInfo(db)
		cacheCtx.apiMenus = append(cacheCtx.apiMenus, info)
		return info
	}
	return nil
}

func HadMenuByName(name string) bool {
	for i := 0;i < len(cacheCtx.apiMenus);i += 1{
		if cacheCtx.apiMenus[i].Name == name {
			return true
		}
	}
	return false
}

func (mine *APIMenuInfo)initInfo(db *nosql.Menu)  {
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
	mine.CreateTime = db.CreatedTime
	mine.UpdateTime = db.UpdatedTime
	mine.Operator = db.Operator
	mine.Creator = db.Creator
	mine.Name = db.Name
	mine.Type = db.Type
	mine.Path = db.Path
	mine.Method = db.Method
}

func (mine *APIMenuInfo)Create() error {
	db := new(nosql.Menu)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetMenuNextID()
	db.CreatedTime = time.Now()
	db.UpdatedTime = time.Now()
	db.Name = mine.Name
	db.Method = mine.Method
	db.Path = mine.Path
	db.Type = mine.Type
	db.Creator = mine.Creator
	err := nosql.CreateMenu(db)
	if err == nil {
		mine.initInfo(db)
		cacheCtx.apiMenus = append(cacheCtx.apiMenus, mine)
	}
	return err
}

func (mine *APIMenuInfo)Update(name, kind, path, act, operator string) error {
	err := nosql.UpdateMenuBase(mine.UID, name, kind, path, act, operator)
	if err == nil {
		mine.Name = name
		mine.Type = kind
		mine.Path = path
		mine.Method = act
		mine.Operator = operator
	}
	return err
}

func (mine *APIMenuInfo)Remove(operator string) error {
	err := nosql.RemoveMenu(mine.UID, operator)
	if err == nil {
		for i := 0;i < len(cacheCtx.apiMenus);i += 1 {
			if cacheCtx.apiMenus[i].UID == mine.UID {
				cacheCtx.apiMenus = append(cacheCtx.apiMenus[:i], cacheCtx.apiMenus[i+1:]...)
				break
			}
		}
	}
	return err
}