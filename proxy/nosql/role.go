package nosql

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Role struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Creator string `json:"creator" bson:"creator"`
	Operator string `json:"operator" bson:"operator"`

	Name   string                `json:"name" bson:"name"`
	Remark  string                `json:"remark" bson:"remark"`
	Menus  []string `json:"menus" bson:"menus"`
}

func CreateRole(info *Role) error {
	_, err := insertOne(TableRole, info)
	if err != nil {
		return err
	}
	return nil
}

func GetRoleNextID() uint64 {
	num, _ := getSequenceNext(TableRole)
	return num
}

func GetRole(uid string) (*Role, error) {
	result, err := findOne(TableRole, uid)
	if err != nil {
		return nil, err
	}
	model := new(Role)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetAllRoles() ([]*Role, error) {
	var items = make([]*Role, 0, 10)
	cursor, err1 := findAll(TableMenu, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var node = new(Role)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdateRoleBase(uid, name, remark, operator string) error {
	msg := bson.M{"name": name, "remark": remark, "operator":operator, "updatedAt": time.Now()}
	_, err := updateOne(TableRole, uid, msg)
	return err
}

func RemoveRole(uid, operator string) error {
	_, err := removeOne(TableRole, uid, operator)
	return err
}

func AppendRoleMenu(uid string, menu string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"menus": menu}
	_, err := appendElement(TableRole, uid, msg)
	return err
}

func SubtractRoleMenu(uid string, menu string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"menus": menu}
	_, err := removeElement(TableRole, uid, msg)
	return err
}

