package grpc

import (
	"context"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-acm/proto/acm"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	"omo.msa.acm/cache"
)

type SceneService struct {}

func switchScene(info *cache.SceneInfo) *pb.SceneLink {
	tmp := new(pb.SceneLink)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Created = info.CreateTime.Unix()
	tmp.Updated = info.UpdateTime.Unix()
	tmp.Operator = info.Operator
	tmp.Creator = info.Creator
	tmp.Scene = info.Scene
	tmp.Status = uint32(info.Status)
	tmp.Links = info.Modules
	return tmp
}

func (mine *SceneService)AddOne(ctx context.Context, in *pb.ReqSceneAdd, out *pb.ReplySceneLink) error {
	path := "scene.addOne"
	inLog(path, in)
	if len(in.Scene) < 1 {
		out.Status = outError(path,"the user uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	t := cache.GetScene(in.Scene)
	if t != nil {
		out.Status = outError(path,"the user had repeated", pbstatus.ResultStatus_Repeated)
		return nil
	}
	info := new(cache.SceneInfo)
	info.Scene = in.Scene
	info.Operator = in.Operator
	err := info.Create(in.Type, in.Links)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchScene(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *SceneService)GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplySceneLink) error {
	path := "scene.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the user uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.GetScene(in.Uid)

	if info == nil {
		out.Status = outError(path,"the user not found", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchScene(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *SceneService)RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "scene.removeOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the scene uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.GetScene(in.Uid)
	if info == nil {
		//out.Status = outError(path,"the user not found", pbstatus.ResultStatus_NotExisted)
		out.Status = outLog(path, out)
		return nil
	}
	err := info.Remove(in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *SceneService)GetList(ctx context.Context, in *pb.RequestPage, out *pb.ReplySceneList) error {
	path := "scene.getList"
	inLog(path, in)
	out.List = make([]*pb.SceneLink, 0, in.Number)
	all := cache.AllScenes()
	for _, value := range all {
		out.List = append(out.List, switchScene(value))
	}

	outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *SceneService) UpdateStatus (ctx context.Context, in *pb.ReqSceneStatus, out *pb.ReplyInfo) error {
	path := "scene.updateStatus"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the scene or uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	var user *cache.SceneInfo
	user = cache.GetScene(in.Uid)
	if user == nil {
		out.Status = outError(path,"not found the scene", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := user.UpdateStatus(uint8(in.Status), in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *SceneService) UpdateModules (ctx context.Context, in *pb.RequestList, out *pb.ReplyList) error {
	path := "scene.updateLinks"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the scene or uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	var info *cache.SceneInfo
	info = cache.GetScene(in.Uid)
	if info == nil {
		out.Status = outError(path,"the scene not found", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.UpdateModules(in.List, in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Uid = info.UID
	out.List = info.Modules
	out.Status = outLog(path, out)
	return nil
}

