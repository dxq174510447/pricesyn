package influx

import (
	"context"
	"fmt"
	"math"
	"pricesyn/file"
	"pricesyn/util"
	"testing"
	"time"
)

func TestInfluxClient_WriteMsg(t *testing.T) {
	ctx := context.Background()

	client := InfluxClient{
		Token: "KJkWWuCI8frlZm-7juZeBBUfnYfObPiKhXd66cTIRvzKQP_T-R_vMx-P1CtGRGq6chsNKVeqnXsGywHDuYxeHw==",
		Url:   "http://10.2.10.12:8086",
	}

	path := "/Users/klook/Downloads/db_export_2084_result.000000000.csv"
	r := &file.FileRowRead{}
	r.Parse(path)
	r.ScanCsvRow(func(rowIndex int, columns []string) {
		if rowIndex == 0 {
			return
		}
		//id,create_time,kl_order_no,selling_currency,ticket_status,activity_date,supplier_type,guest_quantity
		//1,2021/9/23 7:16,3511277215,HKD,1,2021/12/19,2,3
		//id := columns[0]
		createTime := columns[1]
		t1, _ := util.DateUtil.Cover2Time(createTime, "2006/1/2 15:04")
		t2 := t1.Add(time.Hour * 8)
		t1 = &t2

		//orderNo := columns[2]
		sellingCurrency := columns[3]
		ticketStatus := columns[4]
		activityDate := columns[5]
		t3, _ := util.DateUtil.Cover2Time(activityDate, "2006/1/2")
		supplierType := columns[6]
		guestQuantity := columns[7]

		duration := t3.Sub(*t1).Hours() / 24
		tags := make(map[string]string)
		tags["currency"] = sellingCurrency
		tags["status"] = ticketStatus
		tags["supplier"] = supplierType
		tags["quantity"] = guestQuantity
		tags["period"] = fmt.Sprintf("%d", int(math.Ceil(duration)))
		tags["activity_date"] = util.DateUtil.FormatByType(t3, util.DatePattern1)
		tags["begin"] = util.DateUtil.FormatByType(t1, util.DatePattern1)

		client.WriteMsg(ctx, "klook", "ktest", "cruise", tags, t1)
	})

	//time.Sleep(time.Second * 100)
}
