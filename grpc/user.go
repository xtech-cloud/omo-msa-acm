package grpc

import (
	"context"
	"errors"
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
	inLog("user.add", in)
	if len(in.User) < 1 {
		out.Status = pb.ResultStatus_Empty
		return errors.New("the user is empty")
	}
	tmp := cache.GetUser(in.User)
	if tmp != nil {
		out.Info = switchUser(tmp)
		return nil
	}
	info := new(cache.UserInfo)
	info.User = in.User
	info.Operator = in.Operator
	err := info.Create(cache.UserType(in.Type), in.Roles)
	if err == nil {
		out.Info = switchUser(info)
	}else{
		out.Status = pb.ResultStatus_DBException
	}

	return err
}

func (mine *UserService)GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyUserLink) error {
	inLog("user.get", in)
	if len(in.Uid) < 1 {
		out.Status = pb.ResultStatus_Empty
		return errors.New("the user uid is empty")
	}
	info := cache.GetUser(in.Uid)
	if info == nil {
		out.Status = pb.ResultStatus_NotExisted
		return errors.New("the user not found")
	}
	out.Info = switchUser(info)
	return nil
}

func (mine *UserService)RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	inLog("user.remove", in)
	if len(in.Uid) < 1 {
		out.Status = pb.ResultStatus_Empty
		return errors.New("the user uid is empty")
	}
	info := cache.GetUser(in.Uid)
	if info == nil {
		out.Status = pb.ResultStatus_NotExisted
		return errors.New("the user not found")
	}
	err := info.Remove(in.Operator)
	if err != nil {
		out.Status = pb.ResultStatus_DBException
	}
	return err
}

func (mine *UserService)GetList(ctx context.Context, in *pb.ReqUserList, out *pb.ReplyUserList) error {
	inLog("user.list", in)
	out.Users = make([]*pb.UserLink, 0, in.Number)
	for i, value := range cache.AllUsers() {
		t := int32(i) / in.Number + 1
		if t == in.Page {
			out.Users = append(out.Users, switchUser(value))
		}
	}
	return nil
}

func (mine *UserService) IsPermission (ctx context.Context, in *pb.ReqUserPermission, out *pb.ReplyUserPermission) error {
	inLog("user.permission", in)
	if len(in.User) < 1 {
		out.Status = pb.ResultStatus_Empty
		return errors.New("the user uid is empty")
	}
	info := cache.GetUser(in.User)
	if info == nil {
		out.Status = pb.ResultStatus_NotExisted
		return errors.New("the user not found")
	}
	out.User = in.User
	out.Permission = info.IsPermission(in.Path, in.Action)
	return nil
}

func (mine *UserService) AppendRole (ctx context.Context, in *pb.ReqUserAdd, out *pb.ReplyLinkRole) error {
	inLog("user.append", in)
	if len(in.User) < 1 || len(in.Roles) < 1 {
		out.Status = pb.ResultStatus_Empty
		return errors.New("the user uid is empty")
	}
	var user *cache.UserInfo
	user = cache.GetUser(in.User)
	if user == nil {
		info := new(cache.UserInfo)
		info.User = in.User
		err := info.Create(cache.UserType(in.Type), in.Roles)
		if err != nil {
			out.Status = pb.ResultStatus_NotExisted
			return errors.New(err.Error())
		}
		user = info
	}
	var err error
	for _, item := range in.Roles {
		err = user.AppendRole(cache.GetRole(item))
	}

	out.User = in.User
	out.Roles = user.Roles()
	return err
}

func (mine *UserService) SubtractRole (ctx context.Context, in *pb.ReqLinkRole, out *pb.ReplyLinkRole) error {
	inLog("user.subtract", in)
	if len(in.User) < 1 {
		out.Status = pb.ResultStatus_Empty
		return errors.New("the user uid is empty")
	}
	info := cache.GetUser(in.User)
	if info == nil {
		out.Status = pb.ResultStatus_NotExisted
		return errors.New("the user not found")
	}
	err := info.SubtractRole(in.Role)
	out.User = in.User
	out.Roles = info.Roles()
	if err != nil {
		out.Status = pb.ResultStatus_DBException
	}
	return err
}

