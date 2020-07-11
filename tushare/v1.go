package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
)

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
	resp := gothsgt()
	return fmt.Sprintf("code:%d,resp:%s", 200, resp)
}

type tushareReq struct {
	ApiName string            `json:"api_name"`
	Token   string            `json:"token"`
	Params  map[string]string `json:"params"`
	Fields  string            `json:"fields"`
}

// 沪深股通
func gothsgt() string {
	req := &tushareReq{
		ApiName: "moneyflow_hsgt",
		Token:   "f992f5edc3e674a770ce37e5a6c1a15c047d0da547ef02b783a7abf0",
		Params:  map[string]string{"trade_date": "20200710"},
		Fields:  "",
	}
	return Post("http://api.tushare.pro", req, "")
}

// 发送POST请求
// url：         请求地址
// data：        POST请求提交的数据
// contentType： 请求体格式，如：application/json
// content：     请求放回的内容
func Post(url string, data interface{}, contentType string) string {

	// 超时时间：5秒
	client := &http.Client{Timeout: 5 * time.Second}
	jsonStr, _ := json.Marshal(data)
	resp, err := client.Post(url, contentType, bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)
	return string(result)
}
