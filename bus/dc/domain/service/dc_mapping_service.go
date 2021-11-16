package service

import (
	"context"
	"pricesyn/bus/dc/domain/model"
	"pricesyn/bus/dc/domain/repository"
	"sort"
)

type DcMappingService struct {
	dcMappingRepImp *repository.DcMappingRep
}

func (d *DcMappingService) FindTop(ctx context.Context,top int)([]*model.DcMappingEntity,error) {
	result,err := d.dcMappingRepImp.All(ctx)

	if err != nil {
		return nil,err
	}
	if len(result) == 0 {
		return result,nil
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].GetTimes() > result[j].GetTimes()
	})
	if len(result) <= top {
		return result,nil
	}
	return result[0:top],nil
}


var DcMappingServiceImp DcMappingService = DcMappingService{
	dcMappingRepImp: &repository.DcMappingRepImp,
}