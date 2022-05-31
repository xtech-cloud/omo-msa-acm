package nosql

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type SceneLink struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Creator string `json:"creator" bson:"creator"`
	Operator string `json:"operator" bson:"operator"`

	Type uint8 						`json:"type" bson:"type"`
	Status uint8 					`json:"status" bson:"status"`
	Scene   string                	`json:"scene" bson:"scene"`
	Modules  []string                `json:"modules" bson:"modules"`
}

func CreateScene(info *SceneLink) error {
	_, err := insertOne(TableScenes, info)
	if err != nil {
		return err
	}
	return nil
}

func GetSceneNextID() uint64 {
	num, _ := getSequenceNext(TableScenes)
	return num
}

func GetScene(uid string) (*SceneLink, error) {
	result, err := findOne(TableScenes, uid)
	if err != nil {
		return nil, err
	}
	model := new(SceneLink)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetAllScenes() ([]*SceneLink, error) {
	var items = make([]*SceneLink, 0, 100)
	cursor, err1 := findAll(TableScenes, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var node = new(SceneLink)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetSceneLink(scene string) (*SceneLink, error) {
	msg := bson.M{"scene":scene}
	result, err := findOneBy(TableScenes, msg)
	if err != nil {
		return nil, err
	}
	model := new(SceneLink)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func RemoveScene(uid string) error {
	_, err := deleteOne(TableScenes, uid)
	return err
}

func RemoveScenePermissions(uid, operator string) error {
	msg := bson.M{"modules": make([]string, 0, 1), "operator":operator,  "updatedAt": time.Now()}
	_, err := updateOne(TableScenes, uid, msg)
	return err
}

func UpdateSceneModules(uid, operator string, menus []string) error {
	msg := bson.M{"modules": menus, "operator":operator,  "updatedAt": time.Now()}
	_, err := updateOne(TableScenes, uid, msg)
	return err
}

func UpdateSceneStatus(uid, operator string, st uint8) error {
	msg := bson.M{"status": st, "operator":operator,  "updatedAt": time.Now()}
	_, err := updateOne(TableScenes, uid, msg)
	return err
}

func AppendSceneModule(uid, module string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"modules": module}
	_, err := appendElement(TableScenes, uid, msg)
	return err
}

func SubtractSceneModule(uid, module string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"modules": module}
	_, err := removeElement(TableScenes, uid, msg)
	return err
}



