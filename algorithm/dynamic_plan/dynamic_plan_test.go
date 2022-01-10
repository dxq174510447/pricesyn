package dynamic_plan

import (
	"context"
	"testing"
)

func TestPackage_Run(t *testing.T) {
	dynaPlan := NewDynamicPlan()
	dynaPlan.AddThing(&Thing{
		Name:   "吉他",
		Size:   1,
		Weight: 1500,
	})
	dynaPlan.AddThing(&Thing{
		Name:   "音响",
		Size:   4,
		Weight: 3000,
	})
	dynaPlan.AddThing(&Thing{
		Name:   "笔记本电脑",
		Size:   3,
		Weight: 2000,
	})
	dynaPlan.AddThing(&Thing{
		Name:   "iphone",
		Size:   1,
		Weight: 2000,
	})
	dynaPlan.AddThing(&Thing{
		Name:   "mp3",
		Size:   1,
		Weight: 1000,
	})

	dynaPlan.AddColumn(&ColumnDef{
		Desc:     "1磅",
		Capacity: 1,
	})
	dynaPlan.AddColumn(&ColumnDef{
		Desc:     "2磅",
		Capacity: 2,
	})
	dynaPlan.AddColumn(&ColumnDef{
		Desc:     "3磅",
		Capacity: 3,
	})
	dynaPlan.AddColumn(&ColumnDef{
		Desc:     "4磅",
		Capacity: 4,
	})
	dynaPlan.Run(context.Background())
	dynaPlan.Print()
}

func TestPackage_Travel(t *testing.T) {
	dynaPlan := NewDynamicPlan()
	dynaPlan.AddThing(&Thing{
		Name:   "威斯密斯特教堂",
		Size:   1,
		Weight: 7,
	})
	dynaPlan.AddThing(&Thing{
		Name:   "环球剧场",
		Size:   1,
		Weight: 6,
	})
	dynaPlan.AddThing(&Thing{
		Name:   "英国国家美术馆",
		Size:   2,
		Weight: 9,
	})
	dynaPlan.AddThing(&Thing{
		Name:   "大英博物馆",
		Size:   4,
		Weight: 9,
	})
	dynaPlan.AddThing(&Thing{
		Name:   "圣保罗大教堂",
		Size:   1,
		Weight: 8,
	})

	dynaPlan.AddColumn(&ColumnDef{
		Desc:     "1个半天",
		Capacity: 1,
	})
	dynaPlan.AddColumn(&ColumnDef{
		Desc:     "2个半天",
		Capacity: 2,
	})
	dynaPlan.AddColumn(&ColumnDef{
		Desc:     "3个半天",
		Capacity: 3,
	})
	dynaPlan.AddColumn(&ColumnDef{
		Desc:     "4个半天",
		Capacity: 4,
	})
	dynaPlan.Run(context.Background())
	dynaPlan.Print()
}
