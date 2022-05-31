package nosql

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Module struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Creator     string             `json:"creator" bson:"creator"`
	Operator    string             `json:"operator" bson:"operator"`
	Name        string             `json:"name" bson:"name"`
	Remark      string             `json:"remark" bson:"remark"`
	Type        uint8              `json:"type" bson:"type"`
	Menus       []string           `json:"menus" bson:"menus"`
}

func CreateModule(info *Module) error {
	_, err := insertOne(TableModules, info)
	if err != nil {
		return err
	}
	return nil
}

func GetModuleNextID() uint64 {
	num, _ := getSequenceNext(TableModules)
	return num
}

func GetModule(uid string) (*Module, error) {
	result, err := findOne(TableModules, uid)
	if err != nil {
		return nil, err
	}
	model := new(Module)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetAllModules() ([]*Module, error) {
	var items = make([]*Module, 0, 100)
	cursor, err1 := findAll(TableModules, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var node = new(Module)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func RemoveModule(uid, operator string) error {
	_, err := removeOne(TableModules, uid, operator)
	return err
}

func UpdateModuleMenus(uid, operator string, menus []string) error {
	msg := bson.M{"menus": menus, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableModules, uid, msg)
	return err
}

func UpdateModuleBase(uid, name, remark, operator string, menus []string) error {
	msg := bson.M{"name": name, "remark":remark, "operator": operator,"menus": menus, "updatedAt": time.Now()}
	_, err := updateOne(TableModules, uid, msg)
	return err
}
