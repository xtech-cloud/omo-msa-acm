package nosql

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Menu struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Creator string `json:"creator" bson:"creator"`
	Operator string `json:"operator" bson:"operator"`

	Name   string `json:"name" bson:"name"`
	Type   string `json:"type" bson:"type"`
	Path   string `json:"path" bson:"path"`
	Method string `json:"method" bson:"method"`
}

func CreateMenu(info *Menu) error {
	_, err := insertOne(TableMenu, info)
	if err != nil {
		return err
	}
	return nil
}

func GetMenuNextID() uint64 {
	num, _ := getSequenceNext(TableMenu)
	return num
}

func GetMenu(uid string) (*Menu, error) {
	result, err := findOne(TableMenu, uid)
	if err != nil {
		return nil, err
	}
	model := new(Menu)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetAllMenus() ([]*Menu, error) {
	var items = make([]*Menu, 0, 50)
	cursor, err1 := findAll(TableMenu, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var node = new(Menu)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdateMenuBase(uid, name, kind, path, act, operator string) error {
	msg := bson.M{"name": name, "type": kind, "path": path, "method": act, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableMenu, uid, msg)
	return err
}

func RemoveMenu(uid, operator string) error {
	_, err := removeOne(TableMenu, uid, operator)
	return err
}
