package cache

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.acm/proxy/nosql"
	"time"
)

type RoleInfo struct {
	BaseInfo
	Remark string
	menus  []*MenuInfo
}

func AllRoles() []*RoleInfo {
	return cacheCtx.roles
}

func GetRole(uid string) *RoleInfo {
	for i := 0;i < len(cacheCtx.roles);i += 1 {
		if cacheCtx.roles[i].UID == uid {
			return cacheCtx.roles[i]
		}
	}
	db,err := nosql.GetRole(uid)
	if err == nil {
		Role := new(RoleInfo)
		Role.initInfo(db)
		cacheCtx.roles = append(cacheCtx.roles, Role)
		return Role
	}
	return nil
}

func HadRoleByName(name string) bool {
	for i := 0;i < len(cacheCtx.roles);i += 1{
		if cacheCtx.roles[i].Name == name {
			return true
		}
	}
	return false
}

func (mine *RoleInfo)initInfo(db *nosql.Role)  {
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
	mine.CreateTime = db.CreatedTime
	mine.UpdateTime = db.UpdatedTime
	mine.Name = db.Name
	mine.Remark = db.Remark
	mine.Operator = db.Operator
	mine.Creator = db.Creator
	mine.updateMenus(db.Menus)
}

func (mine *RoleInfo)Create(menus []string) error {
	db := new(nosql.Role)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetRoleNextID()
	db.CreatedTime = time.Now()
	db.UpdatedTime = time.Now()
	db.Name = mine.Name
	db.Remark = mine.Remark
	db.Menus = menus
	db.Creator = mine.Creator
	if db.Menus == nil {
		db.Menus = make([]string, 0, 1)
	}

	err := nosql.CreateRole(db)
	if err == nil {
		mine.initInfo(db)
		cacheCtx.roles = append(cacheCtx.roles, mine)
	}
	return err
}

func (mine *RoleInfo)updateMenus(menus []string)  {
	if menus == nil {
		mine.menus = make([]*MenuInfo, 0, 1)
		return
	}
	if len(mine.menus) > 0 {
		mine.menus = mine.menus[:0]
	}
	for _, menu := range menus {
		info := GetMenu(menu)
		if info != nil {
			mine.menus = append(mine.menus, info)
		}
	}
}

func (mine *RoleInfo)Update(name, remark, operator string, menus []string) error {
	err := nosql.UpdateRoleBase(mine.UID, name, remark, operator)
	if err == nil {
		mine.Name = name
		mine.Remark = remark
		mine.Operator = operator
		mine.updateMenus(menus)
	}
	return err
}

func (mine *RoleInfo)Remove(operator string) error {
	err := nosql.RemoveRole(mine.UID, operator)
	if err == nil {
		for i := 0;i < len(cacheCtx.roles);i += 1 {
			if cacheCtx.roles[i].UID == mine.UID {
				cacheCtx.roles = append(cacheCtx.roles[:i], cacheCtx.roles[i+1:]...)
				break
			}
		}
	}
	return err
}

func (mine *RoleInfo)hadMenu(path string, act string) bool {
	for i := 0;i < len(mine.menus);i += 1{
		if len(act) < 1 {
			if mine.menus[i].Path == path {
				return true
			}
		}else{
			if mine.menus[i].Path == path && mine.menus[i].Method == act {
				return true
			}
		}
	}
	return false
}

func (mine *RoleInfo)HadMenu(uid string) bool {
	for i := 0;i < len(mine.menus);i += 1{
		if mine.menus[i].UID == uid {
			return true
		}
	}
	return false
}

func (mine *RoleInfo)AllMenus() []*MenuInfo {
	return mine.menus
}

func (mine *RoleInfo)Menus() []string {
	list := make([]string, 0, len(mine.menus))
	for _, role := range mine.menus {
		list = append(list, role.UID)
	}
	return list
}

func (mine *RoleInfo)AppendMenu(menu *MenuInfo) error {
	if menu == nil {
		return errors.New("the menu is nil")
	}
	if mine.HadMenu(menu.UID) {
		return nil
	}
	err := nosql.AppendRoleMenu(mine.UID, menu.UID)
	if err == nil {
		mine.menus = append(mine.menus, menu)
	}
	return err
}

func (mine *RoleInfo)SubtractMenu(menu string) error {
	if len(menu) < 1 {
		return errors.New("the menu uid is empty")
	}
	if !mine.HadMenu(menu) {
		return nil
	}
	err := nosql.SubtractRoleMenu(mine.UID, menu)
	if err == nil {
		for i := 0;i < len(mine.menus);i += 1 {
			if mine.menus[i].UID == menu {
				mine.menus = append(mine.menus[:i], mine.menus[i+1:]...)
				break
			}
		}
	}
	return err
}
