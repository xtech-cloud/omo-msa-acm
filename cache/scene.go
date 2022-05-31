package cache

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.acm/proxy/nosql"
	"time"
)

type SceneInfo struct {
	Status uint8
	BaseInfo
	Type uint32
	Scene  string
	Modules []string
}

func AllScenes() []*SceneInfo {
	return cacheCtx.scenes
}

func GetScene(uid string) *SceneInfo {
	for i := 0;i < len(cacheCtx.scenes);i += 1 {
		if cacheCtx.scenes[i].Scene == uid || cacheCtx.scenes[i].UID == uid {
			return cacheCtx.scenes[i]
		}
	}
	db,err := nosql.GetScene(uid)
	if err == nil {
		info := new(SceneInfo)
		info.initInfo(db)
		cacheCtx.scenes = append(cacheCtx.scenes, info)
		return info
	}
	return getSceneByLink(uid)
}

func getSceneByLink(scene string) *SceneInfo {
	db,err := nosql.GetSceneLink(scene)
	if err == nil {
		info := new(SceneInfo)
		info.initInfo(db)
		cacheCtx.scenes = append(cacheCtx.scenes, info)
		return info
	}
	return nil
}

func (mine *SceneInfo)initInfo(db *nosql.SceneLink)  {
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
	mine.CreateTime = db.CreatedTime
	mine.UpdateTime = db.UpdatedTime
	mine.Operator = db.Operator
	mine.Creator = db.Creator
	mine.Scene = db.Scene
	mine.Status = db.Status
	mine.Type = uint32(db.Type)
	mine.Modules = db.Modules
}

func (mine *SceneInfo)Create(tp uint32, links []string) error {
	db := new(nosql.SceneLink)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetSceneNextID()
	db.CreatedTime = time.Now()
	db.UpdatedTime = time.Now()
	db.Scene = mine.Scene
	db.Operator = mine.Operator
	db.Modules = links
	db.Status = 0
	db.Type = uint8(tp)
	err := nosql.CreateScene(db)
	if err == nil {
		mine.initInfo(db)
		cacheCtx.scenes = append(cacheCtx.scenes, mine)
	}
	return err
}

func (mine *SceneInfo)Remove(operator string) error {
	err := nosql.RemoveScene(mine.UID)
	if err == nil {
		for i := 0;i < len(cacheCtx.scenes);i += 1 {
			if cacheCtx.scenes[i].UID == mine.UID {
				if i == len(cacheCtx.scenes) - 1 {
					cacheCtx.scenes = append(cacheCtx.scenes[:i])
				}else{
					cacheCtx.scenes = append(cacheCtx.scenes[:i], cacheCtx.scenes[i+1:]...)
				}
				break
			}
		}
	}
	return err
}

func (mine *SceneInfo)UpdateModules(list []string, operator string) error {
	if list == nil {
		return errors.New("the links is nil")
	}
	err := nosql.UpdateSceneModules(mine.UID, operator, list)
	if err == nil {
		mine.Modules = list
	}
	return err
}

func (mine *SceneInfo)UpdateStatus(st uint8, operator string) error {
	err := nosql.UpdateSceneStatus(mine.UID, operator, st)
	if err == nil {
		mine.Status = st
	}
	return err
}
