package main

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

var (
	userList []string
	mu       sync.Mutex
)

type lotteryController struct {
	Ctx iris.Context
}

func newApp() *iris.Application {
	app := iris.New()
	mvc.New(app.Party("/")).Handle(&lotteryController{})
	return app
}

func main() {
	app := newApp()

	userList = []string{}

	app.Run(iris.Addr(":8080"))
}

// 获取用户数
func (c *lotteryController) Get() string {
	count := len(userList)
	return fmt.Sprintf("当前总共参与抽奖用户数：%d\n", count)
}

// 导入用户
// POST http://localhost:8080/import
// params: users
// curl --data "users=小刚,小明" http://localhost:8080/import
func (c *lotteryController) PostImport() string {
	mu.Lock()
	defer mu.Unlock()
	strUsers := c.Ctx.FormValue("users")
	users := strings.Split(strUsers, ",")
	count1 := len(userList)

	for _, u := range users {
		u = strings.TrimSpace(u)
		if len(u) > 0 {
			userList = append(userList, u)
		}
	}
	count2 := len(userList)

	return fmt.Sprintf("当前总共参与抽奖用户数：%d，成功导入的用户数为:%d\n", count1, count2-count1)
}

// 抽奖
// http://localhost:8080/lucky
func (c *lotteryController) GetLucky() string {
	mu.Lock()
	defer mu.Unlock()
	count := len(userList)

	if count > 1 {
		seed := time.Now().UnixNano()
		index := rand.New(rand.NewSource(seed)).Int31n(int32(count))

		user := userList[index]
		userList = append(userList[0:index], userList[index+1:]...)
		return fmt.Sprintf("恭喜: %s 中奖！,剩余用户数: %d\n", user, count-1)
	} else if count == 1 {
		user := userList[0]
		userList = []string{}
		return fmt.Sprintf("恭喜: %s 中奖！,剩余用户数: %d\n", user, count-1)
	}

	return fmt.Sprintf("没有用户，请导入新用户\n")
}
