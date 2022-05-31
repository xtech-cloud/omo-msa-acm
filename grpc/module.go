package grpc

import (
	"context"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-acm/proto/acm"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	"omo.msa.acm/cache"
)

type ModuleService struct {}

func switchModule(info *cache.ModuleInfo) *pb.ModuleInfo {
	tmp := new(pb.ModuleInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Created = info.CreateTime.Unix()
	tmp.Updated = info.UpdateTime.Unix()
	tmp.Operator = info.Operator
	tmp.Creator = info.Creator
	tmp.Name = info.Name
	tmp.Type = uint32(info.Type)
	tmp.Remark = info.Remark
	tmp.Menus = info.Menus
	return tmp
}

func (mine *ModuleService)AddOne(ctx context.Context, in *pb.ReqModuleAdd, out *pb.ReplyModuleInfo) error {
	path := "module.addOne"
	inLog(path, in)
	if len(in.Name) < 1 {
		out.Status = outError(path,"the module uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	if cache.HadModuleByName(cache.ModuleType(in.Type), in.Name) {
		out.Status = outError(path,"the module name is existed", pbstatus.ResultStatus_Repeated)
		return nil
	}
	info := new(cache.ModuleInfo)
	info.Name = in.Name
	info.Type = cache.ModuleType(in.Type)
	info.Remark = in.Remark
	info.Creator = in.Operator
	info.Menus = in.Menus
	err := info.Create()
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchModule(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *ModuleService)GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyModuleInfo) error {
	path := "module.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the module uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.GetModule(in.Uid)
	if info == nil {
		out.Status = outError(path,"the module not found", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchModule(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *ModuleService)RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "module.removeOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the module uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.GetModule(in.Uid)
	if info == nil {
		out.Status = outError(path,"the module not found", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	if info.Creator == "system" {
		out.Status = outError(path,"the system module not allow to delete", pbstatus.ResultStatus_DBException)
		return nil
	}
	err := info.Remove(in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.Status = outLog(path, out)
	return err
}

func (mine *ModuleService)GetAll(ctx context.Context, in *pb.RequestPage, out *pb.ReplyModuleList) error {
	inLog("module.getAll", in)
	out.List = make([]*pb.ModuleInfo, 0, 10)
	for _, value := range cache.AllModulesByType(cache.ModuleType(0)) {
		out.List = append(out.List, switchModule(value))
	}
	out.Status = outLog("module.getAll", fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *ModuleService)UpdateBase(ctx context.Context, in *pb.ReqModuleUpdate, out *pb.ReplyInfo) error {
	path := "module.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the module uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.GetModule(in.Uid)
	if info == nil {
		out.Status = outError(path,"the module not found", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	if in.Name != info.Name && cache.HadModuleByName(info.Type, in.Name) {
		out.Status = outError(path,"the module name is existed", pbstatus.ResultStatus_Repeated)
		return nil
	}
	err := info.UpdateBase(in.Name, in.Remark, in.Operator, in.Menus)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *ModuleService)UpdateMenus(ctx context.Context, in *pb.RequestList, out *pb.ReplyList) error {
	path := "module.updateMenus"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the module uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.GetModule(in.Uid)
	if info == nil {
		out.Status = outError(path,"the module not found", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	err := info.UpdateMens(in.List, in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = info.Menus
	out.Status = outLog(path, out)
	return nil
}
