package grpc

import (
	"context"
	"fmt"
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
	path := "role.addOne"
	inLog(path, in)
	if len(in.Name) < 1 {
		out.Status = outError(path,"the role name is empty", pb.ResultCode_Empty)
		return nil
	}
	if cache.HadRoleByName(in.Name) {
		out.Status = outError(path,"the role name is existed", pb.ResultCode_Repeated)
		return nil
	}
	info := new(cache.RoleInfo)
	info.Name = in.Name
	info.Remark = in.Remark
	info.Creator = in.Operator
	err := info.Create(in.Menus)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultCode_DBException)
		return nil
	}
	out.Info = switchRole(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *RoleService)GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyRoleInfo) error {
	path := "role.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the role uid is empty", pb.ResultCode_Empty)
		return nil
	}
	info := cache.GetRole(in.Uid)
	if info == nil {
		out.Status = outError(path,"the role not found", pb.ResultCode_NotExisted)
		return nil
	}
	out.Info = switchRole(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *RoleService)RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "role.removeOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the role uid is empty", pb.ResultCode_Empty)
		return nil
	}
	info := cache.GetRole(in.Uid)
	if info == nil {
		out.Status = outError(path,"the role not found", pb.ResultCode_NotExisted)
		return nil
	}
	if info.Creator == "system" {
		out.Status = outError(path,"the system role not allow to delete", pb.ResultCode_DBException)
		return nil
	}
	err := info.Remove(in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultCode_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.Status = outLog(path, out)
	return nil
}

func (mine *RoleService)GetAll(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyRoleList) error {
	path := "role.getAll"
	inLog(path, in)
	out.List = make([]*pb.RoleInfo, 0, 5)
	for _, value := range cache.AllRoles() {
		out.List = append(out.List, switchRole(value))
	}
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *RoleService)UpdateBase(ctx context.Context, in *pb.ReqRoleUpdate, out *pb.ReplyRoleInfo) error {
	path := "role.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the role uid is empty", pb.ResultCode_Empty)
		return nil
	}
	info := cache.GetRole(in.Uid)
	if info == nil {
		out.Status = outError(path,"the role not found", pb.ResultCode_NotExisted)
		return nil
	}
	err := info.Update(in.Name, in.Remark, in.Operator, in.Menus)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultCode_DBException)
		return nil
	}
	out.Info = switchRole(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *RoleService)AppendMenu(ctx context.Context, in *pb.ReqRoleMenus, out *pb.ReplyRoleMenu) error {
	path := "role.appendMenu"
	inLog(path, in)
	if len(in.Role) < 1 {
		out.Status = outError(path,"the role uid is empty", pb.ResultCode_Empty)
		return nil
	}
	info := cache.GetRole(in.Role)
	if info == nil {
		out.Status = outError(path,"the role not found", pb.ResultCode_NotExisted)
		return nil
	}
	var err error
	for _, menu := range in.Menus {
		err = info.AppendMenu(cache.GetMenu(menu))
		if err != nil {
			break
		}
	}
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultCode_DBException)
		return nil
	}
	out.Role = in.Role
	out.Menus = info.Menus()
	out.Status = outLog(path, out)
	return nil
}

func (mine *RoleService)SubtractMenu(ctx context.Context, in *pb.ReqRoleMenus, out *pb.ReplyRoleMenu) error {
	path := "role.subtractMenu"
	inLog(path, in)
	if len(in.Role) < 1 {
		out.Status = outError(path,"the role uid is empty", pb.ResultCode_Empty)
		return nil
	}
	info := cache.GetRole(in.Role)
	if info == nil {
		out.Status = outError(path,"the role not found", pb.ResultCode_NotExisted)
		return nil
	}
	var err error
	for _, menu := range in.Menus {
		err = info.SubtractMenu(menu)
		if err != nil {
			break
		}
	}
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultCode_DBException)
		return nil
	}
	out.Role = in.Role
	out.Menus = info.Menus()
	out.Status = outLog(path, out)
	return nil
}

