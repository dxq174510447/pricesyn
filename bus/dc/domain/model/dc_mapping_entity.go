package model

import (
	"context"
	"strconv"
)

type DcMappingEntity struct {
	supplierType            string
	activityDate            string
	times                   int
	duration                string
	guestQuantity           string
	categoryLoc             string
	pricedCategoryCode      string
	agentId                 string
	regionCode              string
	configField             string
	shipCode                string
	schedule *DcScheduleVo
}

func (d *DcMappingEntity) GetTimes() int {
	return d.times
}

func (d *DcMappingEntity) FillByRow(columns []string) *DcMappingEntity{
	d.supplierType      = columns[0]
	d.activityDate      = columns[1]
	if columns[2] == ""{
		d.times = 0
	}else{
		d.times,_           = strconv.Atoi(columns[2])
	}
	d.duration          = columns[3]
	d.guestQuantity     = columns[4]
	d.categoryLoc       = columns[5]
	d.pricedCategoryCode= columns[6]
	d.agentId           = columns[7]
	d.regionCode        = columns[8]
	d.configField       = columns[9]
	d.shipCode          = columns[10]
	return d
}

func (d *DcMappingEntity) GetSchedule(ctx context.Context) (*DcScheduleVo,error){
	if d.schedule == nil {
		
	}
	return d.schedule,nil
}