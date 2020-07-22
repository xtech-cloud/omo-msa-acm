package grpc

import (
	"context"
	"errors"
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
	inLog("menu.add", in)
	if len(in.Name) < 1 {
		out.Status = pb.ResultStatus_Empty
		return errors.New("the menu name is empty")
	}
	if cache.HadMenuByName(in.Name) {
		out.Status = pb.ResultStatus_Repeated
		return errors.New("the menu name is existed")
	}
	info := new(cache.MenuInfo)
	info.Name = in.Name
	info.Type = in.Type
	info.Path = in.Path
	info.Method = in.Method
	info.Creator = in.Operator
	err := info.Create()
	if err == nil {
		out.Info = switchMenu(info)
	}else{
		out.Status = pb.ResultStatus_DBException
	}
	return err
}

func (mine *MenuService)GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyMenuInfo) error {
	inLog("menu.get", in)
	if len(in.Uid) < 1 {
		out.Status = pb.ResultStatus_Empty
		return errors.New("the menu uid is empty")
	}
	info := cache.GetMenu(in.Uid)
	if info == nil {
		out.Status = pb.ResultStatus_NotExisted
		return errors.New("the menu not found")
	}
	out.Info = switchMenu(info)
	return nil
}

func (mine *MenuService)RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	inLog("menu.remove", in)
	if len(in.Uid) < 1 {
		out.Status = pb.ResultStatus_Empty
		return errors.New("the menu uid is empty")
	}
	info := cache.GetMenu(in.Uid)
	if info == nil {
		out.Status = pb.ResultStatus_NotExisted
		return errors.New("the menu not found")
	}
	err := info.Remove(in.Operator)
	if err != nil {
		out.Status = pb.ResultStatus_DBException
	}
	out.Uid = in.Uid
	return err
}

func (mine *MenuService)GetAll(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyMenuList) error {
	inLog("menu.all", in)
	out.List = make([]*pb.MenuInfo, 0, 10)
	for _, value := range cache.AllMenus() {
		out.List = append(out.List, switchMenu(value))
	}
	return nil
}

func (mine *MenuService)UpdateBase(ctx context.Context, in *pb.ReqMenuUpdate, out *pb.ReplyMenuInfo) error {
	inLog("menu.update", in)
	if len(in.Uid) < 1 {
		out.Status = pb.ResultStatus_Empty
		return errors.New("the menu uid is empty")
	}
	info := cache.GetMenu(in.Uid)
	if info == nil {
		out.Status = pb.ResultStatus_NotExisted
		return errors.New("the menu not found")
	}
	err := info.Update(in.Name, in.Type, in.Path, in.Method, in.Operator)
	if err != nil {
		out.Status = pb.ResultStatus_DBException
	}
	out.Info = switchMenu(info)
	return err
}
