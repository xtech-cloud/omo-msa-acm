package grpc

import (
	"context"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-acm/proto/acm"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	"omo.msa.acm/cache"
)

type CatalogService struct{}

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
	tmp.Type = uint32(info.Type)
	tmp.Remark = info.Remark
	tmp.Concepts = info.Concepts
	return tmp
}

func (mine *CatalogService) AddOne(ctx context.Context, in *pb.ReqCatalogAdd, out *pb.ReplyCatalogInfo) error {
	path := "catalog.addOne"
	inLog(path, in)
	if len(in.Name) < 1 || len(in.Key) < 1 {
		out.Status = outError(path, "the catalog name or key is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	if cache.HadCatalogByKey(cache.CatalogType(in.Type), in.Name) {
		out.Status = outError(path, "the catalog name is existed", pbstatus.ResultStatus_Repeated)
		return nil
	}
	info := new(cache.CatalogInfo)
	info.Name = in.Name
	info.Remark = in.Remark
	info.Key = in.Key
	info.Type = cache.CatalogType(in.Type)
	info.Creator = in.Operator
	err := info.Create()
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchCatalog(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *CatalogService) GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyCatalogInfo) error {
	path := "catalog.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the catalog uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.GetCatalog(in.Uid)
	if info == nil {
		out.Status = outError(path, "the catalog not found", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchCatalog(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *CatalogService) RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "catalog.removeOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the catalog uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.GetCatalog(in.Uid)
	if info == nil {
		out.Status = outError(path, "the catalog not found", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	if info.Creator == "system" {
		out.Status = outError(path, "the system catalog not allow to delete", pbstatus.ResultStatus_DBException)
		return nil
	}
	err := info.Remove(in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.Status = outLog(path, out)
	return err
}

func (mine *CatalogService) GetAll(ctx context.Context, in *pb.RequestPage, out *pb.ReplyCatalogList) error {
	path := "catalog.getAll"
	inLog(path, in)
	out.List = make([]*pb.CatalogInfo, 0, 10)
	list := cache.AllCatalogsByType(cache.CatalogType(in.Type))
	for _, value := range list {
		out.List = append(out.List, switchCatalog(value))
	}

	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *CatalogService) UpdateBase(ctx context.Context, in *pb.ReqCatalogUpdate, out *pb.ReplyCatalogInfo) error {
	path := "catalog.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the catalog uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.GetCatalog(in.Uid)
	if info == nil {
		out.Status = outError(path, "the catalog not found", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	if info.Key != in.Key && cache.HadCatalogByKey(info.Type, in.Key) {
		out.Status = outError(path, "the catalog key had existed", pbstatus.ResultStatus_Repeated)
		return nil
	}
	err := info.Update(in.Name, in.Key, in.Remark, in.Operator, in.Concepts)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchCatalog(info)
	out.Status = outLog(path, out)
	return nil
}
