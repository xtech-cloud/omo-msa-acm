package cache

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.acm/proxy/nosql"
	"time"
)

type UserInfo struct {
	BaseInfo
	User  string
	roles []*RoleInfo
}

func AllUsers() []*UserInfo {
	return cacheCtx.users
}

func getUser(uid string) *UserInfo {
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

func GetUser(user string) *UserInfo {
	for i := 0;i < len(cacheCtx.users);i += 1 {
		if cacheCtx.users[i].User == user {
			return cacheCtx.users[i]
		}
	}
	db,err := nosql.GetUserLink(user)
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
	mine.User = db.User
	mine.roles = make([]*RoleInfo, 0, len(db.Roles))
	for _, role := range db.Roles {
		info := GetRole(role)
		if info != nil {
			mine.roles = append(mine.roles, info)
		}
	}
}

func (mine *UserInfo)Create(roles []string) error {
	db := new(nosql.UserLink)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetUserNextID()
	db.CreatedTime = time.Now()
	db.UpdatedTime = time.Now()
	db.User = mine.User
	db.Roles = roles
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
	err := nosql.RemoveUser(mine.UID, operator)
	if err == nil {
		for i := 0;i < len(cacheCtx.users);i += 1 {
			if cacheCtx.users[i].UID == mine.UID {
				cacheCtx.users = append(cacheCtx.users[:i], cacheCtx.users[i+1:]...)
				break
			}
		}
	}
	return err
}

func (mine *UserInfo)IsPermission(path string, action string) bool {
	for _, role := range mine.roles {
		if role.hadMenu(path, action) {
			return true
		}
	}
	return false
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

func (mine *UserInfo)AppendRole(info *RoleInfo) error {
	if info == nil {
		return errors.New("the role is nil")
	}
	if mine.HadRole(info.UID) {
		return errors.New("the role had exist for user")
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
		return errors.New("the role not existed for user")
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

