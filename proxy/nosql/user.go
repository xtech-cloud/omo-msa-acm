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

	Type uint8 						`json:"type" bson:"type"`
	Status uint8 					`json:"status" bson:"status"`
	User   string                	`json:"user" bson:"user"`
	Owner string 					`json:"owner" bson:"owner"`
	Roles  []string                `json:"roles" bson:"roles"`
	Links  []string 				`json:"links" bson:"links"`
}

func CreateUser(info *UserLink) error {
	_, err := insertOne(TableUsers, info)
	if err != nil {
		return err
	}
	return nil
}

func GetUserNextID() uint64 {
	num, _ := getSequenceNext(TableUsers)
	return num
}

func GetUser(uid string) (*UserLink, error) {
	result, err := findOne(TableUsers, uid)
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
	cursor, err1 := findAll(TableUsers, 0)
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
	result, err := findOneBy(TableUsers, msg)
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

func GetUsersByOwner(owner string) (*UserLink, error) {
	msg := bson.M{"owner":owner}
	result, err := findOneBy(TableUsers, msg)
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

func RemoveUser(uid string) error {
	_, err := deleteOne(TableUsers, uid)
	return err
}

func RemoveUserPermissions(uid, operator string) error {
	msg := bson.M{"roles": make([]string, 0, 1), "links": make([]string, 0, 1), "operator":operator,  "updatedAt": time.Now()}
	_, err := updateOne(TableUsers, uid, msg)
	return err
}

func UpdateUserRoles(uid, operator string, list []string) error {
	msg := bson.M{"roles": list, "operator":operator,  "updatedAt": time.Now()}
	_, err := updateOne(TableUsers, uid, msg)
	return err
}

func UpdateUserLinks(uid, operator string, list []string) error {
	msg := bson.M{"links": list, "operator":operator,  "updatedAt": time.Now()}
	_, err := updateOne(TableUsers, uid, msg)
	return err
}

func UpdateUserStatus(uid, operator string, st uint8) error {
	msg := bson.M{"status": st, "operator":operator,  "updatedAt": time.Now()}
	_, err := updateOne(TableUsers, uid, msg)
	return err
}

func AppendUserRole(uid string, role string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"roles": role}
	_, err := appendElement(TableUsers, uid, msg)
	return err
}

func SubtractUserRole(uid string, role string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"roles": role}
	_, err := removeElement(TableUsers, uid, msg)
	return err
}



