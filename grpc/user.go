package grpc

import (
	"context"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-acm/proto/acm"
	"omo.msa.acm/cache"
)

type UserService struct {}

func switchUser(info *cache.UserInfo) *pb.UserLink {
	tmp := new(pb.UserLink)
	tmp.Uid = info.UID
	tmp.User = info.User
	tmp.Roles = info.Roles()
	return tmp
}

func (mine *UserService)AddOne(ctx context.Context, in *pb.ReqUserAdd, out *pb.ReplyUserLink) error {
	path := "user.addOne"
	inLog(path, in)
	if len(in.User) < 1 {
		out.Status = outError(path,"the user uid is empty", pb.ResultCode_Empty)
		return nil
	}
	tmp := cache.GetUser(in.User)
	if tmp != nil {
		out.Info = switchUser(tmp)
		out.Status = outLog(path, out)
		return nil
	}
	info := new(cache.UserInfo)
	info.User = in.User
	info.Operator = in.Operator
	err := info.Create(cache.UserType(in.Type), in.Roles)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultCode_DBException)
		return nil
	}
	out.Info = switchUser(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *UserService)GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyUserLink) error {
	path := "user.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the user uid is empty", pb.ResultCode_Empty)
		return nil
	}
	info := cache.GetUser(in.Uid)
	if info == nil {
		out.Status = outError(path,"the user not found", pb.ResultCode_NotExisted)
		return nil
	}
	out.Info = switchUser(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *UserService)RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "user.removeOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the user uid is empty", pb.ResultCode_Empty)
		return nil
	}
	info := cache.GetUser(in.Uid)
	if info == nil {
		out.Status = outError(path,"the user not found", pb.ResultCode_NotExisted)
		return nil
	}
	err := info.Remove(in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultCode_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *UserService)GetList(ctx context.Context, in *pb.ReqUserList, out *pb.ReplyUserList) error {
	path := "user.getList"
	inLog(path, in)
	out.Users = make([]*pb.UserLink, 0, in.Number)
	for i, value := range cache.AllUsers() {
		t := int32(i) / in.Number + 1
		if t == in.Page {
			out.Users = append(out.Users, switchUser(value))
		}
	}
	outLog(path, fmt.Sprintf("the length = %d", len(out.Users)))
	return nil
}

func (mine *UserService) IsPermission (ctx context.Context, in *pb.ReqUserPermission, out *pb.ReplyUserPermission) error {
	path := "user.isPermission"
	inLog(path, in)
	if len(in.User) < 1 {
		out.Status = outError(path,"the user uid is empty", pb.ResultCode_Empty)
		return nil
	}
	info := cache.GetUser(in.User)
	if info == nil {
		out.Status = outError(path,"the user not found", pb.ResultCode_NotExisted)
		return nil
	}
	out.User = in.User
	out.Permission = info.IsPermission(in.Path, in.Action)
	out.Status = outLog(path, out)
	return nil
}

func (mine *UserService) AppendRole (ctx context.Context, in *pb.ReqUserAdd, out *pb.ReplyLinkRole) error {
	path := "user.appendRole"
	inLog(path, in)
	if len(in.User) < 1 || len(in.Roles) < 1 {
		out.Status = outError(path,"the user uid is empty", pb.ResultCode_Empty)
		return nil
	}
	var user *cache.UserInfo
	user = cache.GetUser(in.User)
	if user == nil {
		info := new(cache.UserInfo)
		info.User = in.User
		err := info.Create(cache.UserType(in.Type), in.Roles)
		if err != nil {
			out.Status = outError(path,err.Error(), pb.ResultCode_DBException)
			return nil
		}
		user = info
	}
	var err error
	for _, item := range in.Roles {
		err = user.AppendRole(cache.GetRole(item))
	}
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultCode_DBException)
		return nil
	}

	out.User = in.User
	out.Roles = user.Roles()
	out.Status = outLog(path, out)
	return nil
}

func (mine *UserService) SubtractRole (ctx context.Context, in *pb.ReqLinkRole, out *pb.ReplyLinkRole) error {
	path := "user.appendRole"
	inLog(path, in)
	if len(in.User) < 1 {
		out.Status = outError(path,"the user uid is empty", pb.ResultCode_Empty)
		return nil
	}
	info := cache.GetUser(in.User)
	if info == nil {
		out.Status = outError(path,"the user not found", pb.ResultCode_NotExisted)
		return nil
	}
	err := info.SubtractRole(in.Role)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultCode_DBException)
		return nil
	}
	out.User = in.User
	out.Roles = info.Roles()
	out.Status = outLog(path, out)
	return nil
}

