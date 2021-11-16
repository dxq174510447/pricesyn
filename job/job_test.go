package job

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"pricesyn/util"
	"testing"
	"time"
)
// 确保两件事情
// task panic会不会影响cron 会影响到主 //要全局包裹 panice错误
// task 执行时间过长会不会影响下一次task运行 不会
// 时间设置一样是否都可以运行 可以

type User struct{
	Name string
	Timer string
}

func TestCronNew2(t *testing.T){
	var as []*User
	as = append(as,&User{
		Name: "a1",
		Timer: "10 * * * * ?",
	})
	as = append(as,&User{
		Name: "a2",
		Timer: "20 * * * * ?",
	})

	c := cron.New(cron.WithSeconds())
	for _,m := range as {
		c.AddFunc(m.Timer, func(u1 *User) func(){
			return func() {
				fmt.Println(util.DateUtil.FormatNow(),"begin",u1.Name)
			}
		}(m))
	}

	c.Start()

	time.Sleep(time.Duration(100)*time.Minute)
}

func TestCronNew(t *testing.T) {
	c := cron.New(cron.WithSeconds())

	flag := 0
	_,err := c.AddFunc("10 * * * * ?", func() {
		fmt.Println(flag,util.DateUtil.FormatNow(),"begin")
		flag1 := flag
		flag = flag +1
		time.Sleep(time.Duration(80)*time.Second)
		fmt.Println(flag1,util.DateUtil.FormatNow(),"end")
	})

	_,err = c.AddFunc("10 * * * * ?", func() {
		fmt.Println("stop")
	})
	if err != nil {
		panic(err)
	}
	c.Start()

	for _,e := range c.Entries() {
		fmt.Println(util.DateUtil.FormatByType(&e.Prev,util.DatePattern1),util.DateUtil.FormatByType(&e.Next,util.DatePattern1))
	}

	time.Sleep(time.Duration(100)*time.Minute)
}


func TestCronNew1(t *testing.T) {
	//t1 := time.NewTicker(time.Second * 10)
	//for m := range t1.C {
	//	fmt.Println(util.DateUtil.FormatByType(&m,util.DatePattern1))
	//}
	t1 := time.Time{}
	fmt.Println(util.DateUtil.FormatByType(&t1,util.DatePattern1))
}