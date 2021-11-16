package model

import "pricesyn/bus/dc/domain/repository"

type DcScheduleCategory struct {
	Description string
	Id          string
}

type DcScheduleVo struct {
	Reference   string
	Code        string
	Name        string
	Id          string
	Key         string
	Categorys  []DcScheduleCategory
}

func (v *DcScheduleVo) FillByRest(voyages *repository.AvailableVoyageRow) *DcScheduleVo {
	v.Reference  = voyages.Reference
	v.Code       = voyages.Pkg.Code
	v.Name       = voyages.Pkg.Name
	v.Id         = voyages.Pkg.Id
	v.Key        = voyages.Pkg.Key
	for _,category := range voyages.AvailableCategories {
		v.Categorys = append(v.Categorys,DcScheduleCategory{
			Description: category.CabinCategory.Description,
			Id: category.CabinCategory.Id,
		})
	}
	return v
}
