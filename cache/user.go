package cache

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.acm/proxy/nosql"
	"time"
)

type UserType uint8

const (
	UserIdle = 0
	UserForbidden = 1
)

type UserInfo struct {
	Status uint8
	BaseInfo
	Type UserType
	User  string
	Owner string
	Links []string
	roles []*RoleInfo
}

func AllUsers() []*UserInfo {
	return cacheCtx.users
}

func GetUser(uid string) *UserInfo {
	for i := 0;i < len(cacheCtx.users);i += 1 {
		if cacheCtx.users[i].UID == uid {
			return cacheCtx.users[i]
		}
	}
	db,err := nosql.GetUser(uid)
	if err == nil {
		info := new(UserInfo)
		info.initInfo(db)
		cacheCtx.users = append(cacheCtx.users, info)
		return info
	}
	return nil
}

func GetUserByOwner(owner, user string) *UserInfo {
	for i := 0;i < len(cacheCtx.users);i += 1 {
		if cacheCtx.users[i].Owner == owner && cacheCtx.users[i].User == user {
			return cacheCtx.users[i]
		}
	}
	db,err := nosql.GetUserLink(owner, user)
	if err == nil {
		info := new(UserInfo)
		info.initInfo(db)
		cacheCtx.users = append(cacheCtx.users, info)
		return info
	}
	return nil
}

func (mine *UserInfo)initInfo(db *nosql.UserLink)  {
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
	mine.CreateTime = db.CreatedTime
	mine.UpdateTime = db.UpdatedTime
	mine.Operator = db.Operator
	mine.Creator = db.Creator
	mine.Owner = db.Owner
	mine.User = db.User
	mine.Status = db.Status
	mine.Type = UserType(db.Type)
	mine.roles = make([]*RoleInfo, 0, len(db.Roles))
	for _, role := range db.Roles {
		info := GetRole(role)
		if info != nil {
			mine.roles = append(mine.roles, info)
		}
	}
	if mine.Owner == "" {
		mine.Owner = "system"
		_ = nosql.UpdateUserOwner(mine.UID, "system")
	}
}

func (mine *UserInfo)Create(tp UserType, owner string, roles, links []string) error {
	db := new(nosql.UserLink)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetUserNextID()
	db.CreatedTime = time.Now()
	db.UpdatedTime = time.Now()
	db.User = mine.User
	db.Operator = mine.Operator
	db.Owner = owner
	db.Roles = roles
	db.Status = UserIdle
	db.Links = links
	if db.Links == nil {
		db.Links = make([]string, 0, 1)
	}
	db.Type = uint8(tp)
	if db.Roles == nil {
		db.Roles = make([]string, 0, 1)
	}

	err := nosql.CreateUser(db)
	if err == nil {
		mine.initInfo(db)
		cacheCtx.users = append(cacheCtx.users, mine)
	}
	return err
}

func (mine *UserInfo)Remove(operator string) error {
	err := nosql.RemoveUser(mine.UID)
	if err == nil {
		for i := 0;i < len(cacheCtx.users);i += 1 {
			if cacheCtx.users[i].UID == mine.UID {
				if i == len(cacheCtx.users) - 1 {
					cacheCtx.users = append(cacheCtx.users[:i])
				}else{
					cacheCtx.users = append(cacheCtx.users[:i], cacheCtx.users[i+1:]...)
				}
				break
			}
		}
	}
	return err
}

func (mine *UserInfo)IsPermission(path string) bool {
	if mine.Status == UserForbidden {
		return false
	}
	//for _, role := range mine.roles {
	//	if role.hadMenu(path) {
	//		return true
	//	}
	//}
	return true
}

func (mine *UserInfo)HadRole(uid string) bool {
	for i := 0;i < len(mine.roles);i += 1{
		if mine.roles[i].UID == uid {
			return true
		}
	}
	return false
}

func (mine *UserInfo)AllRoles() []*RoleInfo {
	return mine.roles
}

func (mine *UserInfo)Roles() []string {
	list := make([]string, 0, len(mine.roles))
	for _, role := range mine.roles {
		list = append(list, role.UID)
	}
	return list
}

func (mine *UserInfo)UpdateLinks(list []string, operator string) error {
	if list == nil {
		return errors.New("the links is nil")
	}
	err := nosql.UpdateUserLinks(mine.UID, operator, list)
	if err == nil {
		mine.Links = list
	}
	return err
}

func (mine *UserInfo)UpdateStatus(st uint8, operator string) error {
	if st > UserForbidden || st < UserIdle {
		return errors.New("the user status error")
	}
	err := nosql.UpdateUserStatus(mine.UID, operator, st)
	if err == nil {
		mine.Status = st
	}
	return err
}

func (mine *UserInfo)UpdateRoles(list []string, operator string) error {
	if list == nil {
		return errors.New("the roles is nil")
	}
	array := make([]string, 0, len(list))
	roles := make([]*RoleInfo, 0, len(list))
	for i := 0;i < len(list);i +=1 {
		role := GetRole(list[i])
		if role != nil {
			roles = append(roles, role)
			array = append(array, list[i])
		}
	}
	err := nosql.UpdateUserRoles(mine.UID, operator, array)
	if err == nil {
		mine.roles = roles
	}
	return err
}

func (mine *UserInfo)AppendRole(info *RoleInfo) error {
	if info == nil {
		return errors.New("the role is nil")
	}
	if mine.HadRole(info.UID) {
		return nil
	}
	err := nosql.AppendUserRole(mine.UID, info.UID)
	if err == nil {
		mine.roles = append(mine.roles, info)
	}
	return err
}

func (mine *UserInfo)SubtractRole(role string) error {
	if len(role) < 1 {
		return errors.New("the role uid is empty")
	}
	if !mine.HadRole(role) {
		return nil
	}
	err := nosql.SubtractUserRole(mine.UID, role)
	if err == nil {
		for i := 0;i < len(mine.roles);i += 1 {
			if mine.roles[i].UID == role {
				mine.roles = append(mine.roles[:i], mine.roles[i+1:]...)
				break
			}
		}
	}
	return err
}

