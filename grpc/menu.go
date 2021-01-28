package grpc

import (
	"context"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-acm/proto/acm"
	"omo.msa.acm/cache"
)

type MenuService struct {}

func switchMenu(info *cache.MenuInfo) *pb.MenuInfo {
	tmp := new(pb.MenuInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Created = info.CreateTime.Unix()
	tmp.Updated = info.UpdateTime.Unix()
	tmp.Operator = info.Operator
	tmp.Creator = info.Creator
	tmp.Name = info.Name
	tmp.Type = info.Type
	tmp.Path = info.Path
	tmp.Method = info.Method
	return tmp
}

func (mine *MenuService)AddOne(ctx context.Context, in *pb.ReqMenuAdd, out *pb.ReplyMenuInfo) error {
	path := "menu.addOne"
	inLog(path, in)
	if len(in.Name) < 1 {
		out.Status = outError(path,"the menu uid is empty", pb.ResultCode_Empty)
		return nil
	}
	if cache.HadMenuByName(in.Name) {
		out.Status = outError(path,"the menu name is existed", pb.ResultCode_Repeated)
		return nil
	}
	info := new(cache.MenuInfo)
	info.Name = in.Name
	info.Type = in.Type
	info.Path = in.Path
	info.Method = in.Method
	info.Creator = in.Operator
	err := info.Create()
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultCode_DBException)
		return nil
	}
	out.Info = switchMenu(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *MenuService)GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyMenuInfo) error {
	path := "menu.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the menu uid is empty", pb.ResultCode_Empty)
		return nil
	}
	info := cache.GetMenu(in.Uid)
	if info == nil {
		out.Status = outError(path,"the menu not found", pb.ResultCode_NotExisted)
		return nil
	}
	out.Info = switchMenu(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *MenuService)RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "menu.removeOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the menu uid is empty", pb.ResultCode_Empty)
		return nil
	}
	info := cache.GetMenu(in.Uid)
	if info == nil {
		out.Status = outError(path,"the menu not found", pb.ResultCode_NotExisted)
		return nil
	}
	if info.Creator == "system" {
		out.Status = outError(path,"the system menu not allow to delete", pb.ResultCode_DBException)
		return nil
	}
	err := info.Remove(in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultCode_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.Status = outLog(path, out)
	return err
}

func (mine *MenuService)GetAll(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyMenuList) error {
	inLog("menu.getAll", in)
	out.List = make([]*pb.MenuInfo, 0, 10)
	for _, value := range cache.AllMenus() {
		out.List = append(out.List, switchMenu(value))
	}
	out.Status = outLog("menu.getAll", fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *MenuService)UpdateBase(ctx context.Context, in *pb.ReqMenuUpdate, out *pb.ReplyMenuInfo) error {
	path := "menu.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the menu uid is empty", pb.ResultCode_Empty)
		return nil
	}
	info := cache.GetMenu(in.Uid)
	if info == nil {
		out.Status = outError(path,"the menu not found", pb.ResultCode_NotExisted)
		return nil
	}
	err := info.Update(in.Name, in.Type, in.Path, in.Method, in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultCode_DBException)
		return nil
	}
	out.Info = switchMenu(info)
	out.Status = outLog(path, out)
	return nil
}
