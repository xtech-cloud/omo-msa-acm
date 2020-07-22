package nosql

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type UserLink struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Creator string `json:"creator" bson:"creator"`
	Operator string `json:"operator" bson:"operator"`

	Type uint8 `json:"type" bson:"type"`
	User   string                	`json:"user" bson:"user"`
	Roles  []string                `json:"roles" bson:"roles"`
}

func CreateUser(info *UserLink) error {
	_, err := insertOne(TableUserRoles, info)
	if err != nil {
		return err
	}
	return nil
}

func GetUserNextID() uint64 {
	num, _ := getSequenceNext(TableUserRoles)
	return num
}

func GetUser(uid string) (*UserLink, error) {
	result, err := findOne(TableUserRoles, uid)
	if err != nil {
		return nil, err
	}
	model := new(UserLink)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetAllUsers() ([]*UserLink, error) {
	var items = make([]*UserLink, 0, 100)
	cursor, err1 := findAll(TableUserRoles, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var node = new(UserLink)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetUserLink(user string) (*UserLink, error) {
	msg := bson.M{"user":user}
	result, err := findOneBy(TableUserRoles, msg)
	if err != nil {
		return nil, err
	}
	model := new(UserLink)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func RemoveUser(uid, operator string) error {
	_, err := removeOne(TableUserRoles, uid, operator)
	return err
}

func AppendUserRole(uid string, role string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"roles": role}
	_, err := appendElement(TableUserRoles, uid, msg)
	return err
}

func SubtractUserRole(uid string, role string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"roles": role}
	_, err := removeElement(TableUserRoles, uid, msg)
	return err
}

