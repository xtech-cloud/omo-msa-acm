package cache

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.acm/proxy/nosql"
	"time"
)

type PermissionAction uint8

type MenuInfo struct {
	BaseInfo
	Type string
	Path string
	Method string
}

func AllMenus() []*MenuInfo {
	return cacheCtx.menus
}

func GetMenu(uid string) *MenuInfo {
	for i := 0;i < len(cacheCtx.menus);i += 1 {
		if cacheCtx.menus[i].UID == uid {
			return cacheCtx.menus[i]
		}
	}
	db,err := nosql.GetMenu(uid)
	if err == nil {
		info := new(MenuInfo)
		info.initInfo(db)
		cacheCtx.menus = append(cacheCtx.menus, info)
		return info
	}
	return nil
}

func HadMenuByName(name string) bool {
	for i := 0;i < len(cacheCtx.menus);i += 1{
		if cacheCtx.menus[i].Name == name {
			return true
		}
	}
	return false
}

func (mine *MenuInfo)initInfo(db *nosql.Menu)  {
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

func (mine *MenuInfo)Create() error {
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
		cacheCtx.menus = append(cacheCtx.menus, mine)
	}
	return err
}

func (mine *MenuInfo)Update(name, kind, path, act, operator string) error {
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

func (mine *MenuInfo)Remove(operator string) error {
	err := nosql.RemoveMenu(mine.UID, operator)
	if err == nil {
		for i := 0;i < len(cacheCtx.menus);i += 1 {
			if cacheCtx.menus[i].UID == mine.UID {
				cacheCtx.menus = append(cacheCtx.menus[:i], cacheCtx.menus[i+1:]...)
				break
			}
		}
	}
	return err
}