package util

import (
	"fmt"
	"math"
)

type mathUtil struct {
}

func (c *mathUtil) Twrailmon(source float64, percent float64) error {

	//var source float64 =  750
	//var percent float64 = 25

	fmt.Println(math.Floor(5.9))
	fmt.Println(math.Floor(5.6))
	fmt.Println(math.Floor(5.1))
	fmt.Println(math.Floor(5.0))
	fmt.Println(math.Floor(5.01))

	lastPrice := source * (100 - percent) / 100
	lastPrice = math.Floor(lastPrice / 5)
	lastPrice = lastPrice * 5
	fmt.Println(lastPrice)

	return nil
}

var MathUtil mathUtil = mathUtil{}
