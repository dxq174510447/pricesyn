package repository

import (
	"context"
	"pricesyn/bus/dc/domain/model"
	"pricesyn/file"
)

type DcMappingRep struct {

}

func (d *DcMappingRep) All(ctx context.Context) ([]*model.DcMappingEntity,error){
	path := "/Users/klook/code/company/pricesyn/bus/dc/domain/repository/config_pkg.csv"
	r := &file.FileRowRead{}
	r.Parse(path)
	var result []*model.DcMappingEntity
	err := r.ScanCsvRow(func(rowIndex int, columns []string) {
		if rowIndex == 0 {
			return
		}
		r := &model.DcMappingEntity{}
		r.FillByRow(columns)
		result = append(result,r)
	})
	return result,err
}

var DcMappingRepImp DcMappingRep = DcMappingRep{}