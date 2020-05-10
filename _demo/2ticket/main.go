package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

/*
1. 即开即得型
2. 双色球自选型
*/

type lotterController struct {
	Ctx iris.Context
}

func newApp() *iris.Application {
	app := iris.New()
	mvc.New(app.Party("/")).Handle(&lotterController{})
	return app
}

func main() {
	app := newApp()
	app.Run(iris.Addr(":8080"))
}

// 即开即得型 http://localhost:8080
func (c *lotterController) Get() string {
	var price string
	seed := time.Now().UnixNano()
	code := rand.New(rand.NewSource(seed)).Int31n(10)
	switch {
	case code == 1:
		price = "一等奖"
	case code >= 2 && code <= 3:
		price = "二等奖"
	case code >= 3 && code <= 6:
		price = "三等奖"
	default:
		return fmt.Sprintf("尾号为1获得一等奖<br/>尾号为2或者3获得二等奖<br/>尾号为4、5、6获得三等奖<br/>code=%d,很遗憾没有中奖<br/>", code)
	}
	return fmt.Sprintf("尾号为1获得一等奖<br/>尾号为2或者3获得二等奖<br/>尾号为4、5、6获得三等奖<br/>code=%d,恭喜你获得%s<br/>", code, price)
}

// 双色球自选型 http://localhost:8080/price
func (c *lotterController) GetPrice() string {
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))

	var price [7]int

	// 前六个红色球
	for i := 0; i != 6; i++ {
		price[i] = r.Intn(33) + 1
	}

	// 最后一个为蓝色球
	price[6] = r.Intn(16) + 1

	return fmt.Sprintf("今日开奖号码是：%d", price)
}
