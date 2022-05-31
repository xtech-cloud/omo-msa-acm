package cache

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.acm/proxy/nosql"
	"time"
)

type ModuleType uint8

type ModuleInfo struct {
	Type ModuleType
	BaseInfo
	Remark string
	Menus []string
}

func AllModulesByType(tp ModuleType) []*ModuleInfo {
	list := make([]*ModuleInfo, 0, 10)
	for _, item := range cacheCtx.modules {
		if item.Type == tp {
			list = append(list, item)
		}
	}
	return list
}

func GetModule(uid string) *ModuleInfo {
	for i := 0;i < len(cacheCtx.modules);i += 1 {
		if cacheCtx.modules[i].UID == uid {
			return cacheCtx.modules[i]
		}
	}
	db,err := nosql.GetModule(uid)
	if err == nil {
		info := new(ModuleInfo)
		info.initInfo(db)
		cacheCtx.modules = append(cacheCtx.modules, info)
		return info
	}
	return nil
}

func HadModuleByName(tp ModuleType, name string) bool {
	for i := 0;i < len(cacheCtx.modules);i += 1{
		if cacheCtx.modules[i].Type == tp && cacheCtx.modules[i].Name == name {
			return true
		}
	}
	return false
}

func (mine *ModuleInfo)initInfo(db *nosql.Module)  {
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
	mine.CreateTime = db.CreatedTime
	mine.UpdateTime = db.UpdatedTime
	mine.Operator = db.Operator
	mine.Creator = db.Creator
	mine.Name = db.Name
	mine.Remark = db.Remark
	mine.Type = ModuleType(db.Type)
	mine.Menus = db.Menus
}

func (mine *ModuleInfo)Create() error {
	db := new(nosql.Module)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetModuleNextID()
	db.CreatedTime = time.Now()
	db.UpdatedTime = time.Now()
	db.Name = mine.Name
	db.Remark = mine.Remark
	db.Creator = mine.Creator
	db.Menus = mine.Menus
	db.Type = uint8(mine.Type)
	err := nosql.CreateModule(db)
	if err == nil {
		mine.initInfo(db)
		cacheCtx.modules = append(cacheCtx.modules, mine)
	}
	return err
}

func (mine *ModuleInfo)UpdateBase(name, remark, operator string, menus []string) error {
	err := nosql.UpdateModuleBase(mine.UID, name, remark, operator,menus)
	if err == nil {
		mine.Name = name
		mine.Remark = remark
		mine.Operator = operator
	}
	return err
}

func (mine *ModuleInfo)UpdateMens(list []string, operator string) error {
	if list == nil {
		return errors.New("the links is nil")
	}
	err := nosql.UpdateModuleMenus(mine.UID, operator, list)
	if err == nil {
		mine.Menus = list
		mine.Operator = operator
	}
	return err
}

func (mine *ModuleInfo)Remove(operator string) error {
	err := nosql.RemoveModule(mine.UID, operator)
	if err == nil {
		for i := 0;i < len(cacheCtx.modules);i += 1 {
			if cacheCtx.modules[i].UID == mine.UID {
				if i == len(cacheCtx.scenes) - 1 {
					cacheCtx.modules = append(cacheCtx.modules[:i])
				}else{
					cacheCtx.modules = append(cacheCtx.modules[:i], cacheCtx.modules[i+1:]...)
				}
				break
			}
		}
	}
	return err
}