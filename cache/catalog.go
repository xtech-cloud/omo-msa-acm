package cache

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.acm/proxy/nosql"
	"time"
)

const (
	CatalogTypeUser  CatalogType = 0
	CatalogTypeScene CatalogType = 1
)

type CatalogType uint8

type CatalogInfo struct {
	Type CatalogType
	BaseInfo
	Remark   string
	Key      string
	Concepts []string
}

func AllCatalogsByType(tp CatalogType) []*CatalogInfo {
	list := make([]*CatalogInfo, 0, 10)
	for _, catalog := range cacheCtx.catalogs {
		if catalog.Type == tp {
			list = append(list, catalog)
		}
	}
	return list
}

func GetCatalog(uid string) *CatalogInfo {
	for i := 0; i < len(cacheCtx.catalogs); i += 1 {
		if cacheCtx.catalogs[i].UID == uid {
			return cacheCtx.catalogs[i]
		}
	}
	db, err := nosql.GetCatalog(uid)
	if err == nil {
		info := new(CatalogInfo)
		info.initInfo(db)
		cacheCtx.catalogs = append(cacheCtx.catalogs, info)
		return info
	}
	return nil
}

func HadCatalogByKey(tp CatalogType, key string) bool {
	for i := 0; i < len(cacheCtx.catalogs); i += 1 {
		if cacheCtx.catalogs[i].Type == tp && cacheCtx.catalogs[i].Key == key {
			return true
		}
	}
	return false
}

func (mine *CatalogInfo) initInfo(db *nosql.Catalog) {
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
	mine.CreateTime = db.CreatedTime
	mine.UpdateTime = db.UpdatedTime
	mine.Operator = db.Operator
	mine.Creator = db.Creator
	mine.Name = db.Name
	mine.Remark = db.Remark
	mine.Key = db.Key
	mine.Type = CatalogType(db.Type)
	mine.Concepts = db.Concepts
}

func (mine *CatalogInfo) Create() error {
	db := new(nosql.Catalog)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetCatalogNextID()
	db.CreatedTime = time.Now()
	db.UpdatedTime = time.Now()
	db.Name = mine.Name
	db.Remark = mine.Remark
	db.Creator = mine.Creator
	db.Key = mine.Key
	db.Type = uint8(mine.Type)
	db.Concepts = make([]string, 0, 1)
	err := nosql.CreateCatalog(db)
	if err == nil {
		mine.initInfo(db)
		cacheCtx.catalogs = append(cacheCtx.catalogs, mine)
	}
	return err
}

func (mine *CatalogInfo) Update(name, key, remark, operator string, concepts []string) error {
	err := nosql.UpdateCatalogBase(mine.UID, name, key, remark, operator, concepts)
	if err == nil {
		mine.Name = name
		mine.Remark = remark
		mine.Key = key
		mine.Concepts = concepts
		mine.Operator = operator
	}
	return err
}

func (mine *CatalogInfo) Remove(operator string) error {
	err := nosql.RemoveCatalog(mine.UID, operator)
	if err == nil {
		for i := 0; i < len(cacheCtx.catalogs); i += 1 {
			if cacheCtx.catalogs[i].UID == mine.UID {
				cacheCtx.catalogs = append(cacheCtx.catalogs[:i], cacheCtx.catalogs[i+1:]...)
				break
			}
		}
	}
	return err
}
