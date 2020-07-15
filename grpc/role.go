package grpc

import (
	"context"
	"errors"
	pb "github.com/xtech-cloud/omo-msp-acm/proto/acm"
	"omo.msa.acm/cache"
)

type RoleService struct {}

func switchRole(info *cache.RoleInfo) *pb.RoleInfo {
	tmp := new(pb.RoleInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Created = info.CreateTime.Unix()
	tmp.Updated = info.UpdateTime.Unix()
	tmp.Name = info.Name
	tmp.Remark = info.Remark
	tmp.Operator = info.Operator
	tmp.Creator = info.Creator
	tmp.Menus = make([]string, 0, 10)
	for _, info := range info.AllMenus() {
		tmp.Menus = append(tmp.Menus, info.UID)
	}
	return tmp
}

func (mine *RoleService)AddOne(ctx context.Context, in *pb.ReqRoleAdd, out *pb.ReplyRoleInfo) error {
	info := new(cache.RoleInfo)
	info.Name = in.Name
	info.Remark = in.Remark
	info.Creator = in.Operator
	err := info.Create(in.Menus)
	if err == nil {
		out.Info = switchRole(info)
	}else{
		out.Status = pb.ResultStatus_DBException
	}

	return err
}

func (mine *RoleService)GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyRoleInfo) error {
	if len(in.Uid) < 1 {
		out.Status = pb.ResultStatus_Empty
		return errors.New("the role uid is empty")
	}
	info := cache.GetRole(in.Uid)
	if info == nil {
		out.Status = pb.ResultStatus_NotExisted
		return errors.New("the role not found")
	}
	out.Info = switchRole(info)
	return nil
}

func (mine *RoleService)RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	if len(in.Uid) < 1 {
		out.Status = pb.ResultStatus_Empty
		return errors.New("the user uid is empty")
	}
	info := cache.GetRole(in.Uid)
	if info == nil {
		out.Status = pb.ResultStatus_NotExisted
		return errors.New("the role not found")
	}
	err := info.Remove(in.Operator)
	if err != nil {
		out.Status = pb.ResultStatus_DBException
	}
	return err
}

func (mine *RoleService)GetAll(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyRoleList) error {
	out.List = make([]*pb.RoleInfo, 0, 5)
	for _, value := range cache.AllRoles() {
		out.List = append(out.List, switchRole(value))
	}
	return nil
}

func (mine *RoleService)UpdateBase(ctx context.Context, in *pb.ReqRoleUpdate, out *pb.ReplyRoleInfo) error {
	if len(in.Uid) < 1 {
		out.Status = pb.ResultStatus_Empty
		return errors.New("the user uid is empty")
	}
	info := cache.GetRole(in.Uid)
	if info == nil {
		out.Status = pb.ResultStatus_NotExisted
		return errors.New("the role not found")
	}
	err := info.Update(in.Name, in.Remark, in.Operator)
	if err != nil {
		out.Status = pb.ResultStatus_DBException
	}else{
		out.Info = switchRole(info)
	}
	return err
}

func (mine *RoleService)AppendMenu(ctx context.Context, in *pb.ReqRoleMenu, out *pb.ReplyRoleMenu) error {
	if len(in.Role) < 1 {
		out.Status = pb.ResultStatus_Empty
		return errors.New("the user uid is empty")
	}
	info := cache.GetRole(in.Role)
	if info == nil {
		out.Status = pb.ResultStatus_NotExisted
		return errors.New("the role not found")
	}
	err := info.AppendMenu(cache.GetMenu(in.Menu))
	if err != nil {
		out.Status = pb.ResultStatus_DBException
	}
	out.Role = in.Role
	out.Menus = info.Menus()
	return err
}

func (mine *RoleService)SubtractMenu(ctx context.Context, in *pb.ReqRoleMenu, out *pb.ReplyRoleMenu) error {
	if len(in.Role) < 1 {
		out.Status = pb.ResultStatus_Empty
		return errors.New("the user uid is empty")
	}
	info := cache.GetRole(in.Role)
	if info == nil {
		out.Status = pb.ResultStatus_NotExisted
		return errors.New("the role not found")
	}
	err := info.SubtractMenu(in.Menu)
	if err != nil {
		out.Status = pb.ResultStatus_DBException
	}
	out.Role = in.Role
	out.Menus = info.Menus()
	return err
}

