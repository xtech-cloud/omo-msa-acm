package nosql

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Catalog struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Creator     string             `json:"creator" bson:"creator"`
	Operator    string             `json:"operator" bson:"operator"`

	Name   string `json:"name" bson:"name"`
	Remark string `json:"remark" bson:"remark"`
	Key    string `json:"key" bson:"key"`
	Type   uint8  `json:"type" bson:"type"`
}

func CreateCatalog(info *Catalog) error {
	_, err := insertOne(TableCatalog, info)
	if err != nil {
		return err
	}
	return nil
}

func GetCatalogNextID() uint64 {
	num, _ := getSequenceNext(TableCatalog)
	return num
}

func GetCatalog(uid string) (*Catalog, error) {
	result, err := findOne(TableCatalog, uid)
	if err != nil {
		return nil, err
	}
	model := new(Catalog)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetAllCatalogs() ([]*Catalog, error) {
	var items = make([]*Catalog, 0, 50)
	cursor, err1 := findAll(TableCatalog, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var node = new(Catalog)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdateCatalogBase(uid, name, key, remark, operator string) error {
	msg := bson.M{"name": name, "key": key, "remark": remark, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableCatalog, uid, msg)
	return err
}

func RemoveCatalog(uid, operator string) error {
	_, err := removeOne(TableCatalog, uid, operator)
	return err
}
