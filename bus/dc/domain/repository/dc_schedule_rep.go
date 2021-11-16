package repository

import (
	"context"
	"fmt"
	"net/http"
	"pricesyn/bus/dc/domain/model"
	"pricesyn/util"
	"sync"
)

type DcScheduleRep struct {
	client *http.Client
	initLock sync.Once
	url string
}

func (d *DcScheduleRep) init(){
	d.initLock.Do(func() {
		d.client = util.HttpUtil.GetHttpClient(10,10,"")
		//d.url = "https://swgql.gentingcruises.com:3000/graphql"
		d.url = "http://52.77.183.57:3000/graphql"
	})
}

func (d *DcScheduleRep) GetToken(ctx context.Context) (string,error){
	d.init()
	req := &LoginRequest{
		Pwd: "wywK46iu",
	}
	reqMap := make(map[string]interface{})
	reqMap["query"] = loginTpl
	reqMap["variables"] = req

	url := d.url
	response := &LoginResponse{}
	err := util.HttpUtil.PostBody(ctx,d.client,url,reqMap,response,nil)
	if err != nil {
		return "",err
	}
	return response.Data.Login.Token,err
}

func (d *DcScheduleRep) GetSchedule(ctx context.Context,shipCode string,activityDate string,duration int) (*model.DcScheduleVo,error){
	d.init()

	token,err := d.GetToken(ctx)
	if err != nil {
		return nil,err
	}
	header := make(map[string]string)
	header["Cookie"] =  fmt.Sprintf("__token=%s", token)

	req := &ScheduleReqest{
		Datefrom: activityDate,
		Dateto: activityDate,
		Ports: []string{},
		Ships: []string{shipCode},
		MinDur: duration,
		MaxDur:  duration,
	}
	reqMap := make(map[string]interface{})
	reqMap["query"] = scheduleTpl
	reqMap["variables"] = req

	url := d.url
	response := &ScheduleResponse{}
	err = util.HttpUtil.PostBody(ctx,d.client,url,reqMap,response,header)
	if err != nil {
		return nil,err
	}

	if len(response.Data.AvailableVoyages) == 0 {
		if len(response.Errors) > 0 {
			msg := util.StringUtil.ArrayJoin(response.Errors, func(index int) string {
				return response.Errors[index].Message
			})
			return nil,fmt.Errorf(msg)
		}
		return nil,nil
	}
	schedule := &model.DcScheduleVo{}
	schedule.FillByRest(&response.Data.AvailableVoyages[0])
	return schedule,nil
}

var DcScheduleRepImpl DcScheduleRep = DcScheduleRep{}