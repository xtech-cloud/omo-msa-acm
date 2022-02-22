package grpc

import (
	"context"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-acm/proto/acm"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	"omo.msa.acm/cache"
	"strconv"
)

type CatalogService struct {}

func switchCatalog(info *cache.CatalogInfo) *pb.CatalogInfo {
	tmp := new(pb.CatalogInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Created = info.CreateTime.Unix()
	tmp.Updated = info.UpdateTime.Unix()
	tmp.Operator = info.Operator
	tmp.Creator = info.Creator
	tmp.Name = info.Name
	tmp.Key = info.Key
	tmp.Remark = info.Remark
	return tmp
}

func (mine *CatalogService)AddOne(ctx context.Context, in *pb.ReqCatalogAdd, out *pb.ReplyCatalogInfo) error {
	path := "catalog.addOne"
	inLog(path, in)
	if len(in.Name) < 1 || len(in.Key) < 1 {
		out.Status = outError(path,"the catalog name or key is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	if cache.HadCatalogByKey(in.Name) {
		out.Status = outError(path,"the catalog name is existed", pbstatus.ResultStatus_Repeated)
		return nil
	}
	info := new(cache.CatalogInfo)
	info.Name = in.Name
	info.Remark = in.Remark
	info.Key = in.Key
	info.Type = uint8(in.Type)
	info.Creator = in.Operator
	err := info.Create()
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchCatalog(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *CatalogService)GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyCatalogInfo) error {
	path := "catalog.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the catalog uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.GetCatalog(in.Uid)
	if info == nil {
		out.Status = outError(path,"the catalog not found", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchCatalog(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *CatalogService)RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "catalog.removeOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the catalog uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.GetCatalog(in.Uid)
	if info == nil {
		out.Status = outError(path,"the catalog not found", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	if info.Creator == "system" {
		out.Status = outError(path,"the system catalog not allow to delete", pbstatus.ResultStatus_DBException)
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

func (mine *CatalogService)GetAll(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyCatalogList) error {
	path := "catalog.getAll"
	inLog(path, in)
	if in.Uid == "" {
		out.List = make([]*pb.CatalogInfo, 0, 10)
		list := cache.AllCatalogsByType(0)
		for _, value := range list {
			out.List = append(out.List, switchCatalog(value))
		}
	}else{
		tp, err := strconv.ParseUint(in.Uid, 10, 32)
		if err != nil {
			out.Status = outError(path, err.Error(), pbstatus.ResultStatus_FormatError)
			return nil
		}
		out.List = make([]*pb.CatalogInfo, 0, 10)
		list := cache.AllCatalogsByType(uint8(tp))
		for _, value := range list {
			out.List = append(out.List, switchCatalog(value))
		}
	}

	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *CatalogService)UpdateBase(ctx context.Context, in *pb.ReqCatalogUpdate, out *pb.ReplyCatalogInfo) error {
	path := "catalog.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the catalog uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.GetCatalog(in.Uid)
	if info == nil {
		out.Status = outError(path,"the catalog not found", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	if info.Key != in.Key && cache.HadCatalogByKey(in.Key) {
		out.Status = outError(path,"the catalog key had existed", pbstatus.ResultStatus_Repeated)
		return nil
	}
	err := info.Update(in.Name, in.Key, in.Remark, in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchCatalog(info)
	out.Status = outLog(path, out)
	return nil
}
