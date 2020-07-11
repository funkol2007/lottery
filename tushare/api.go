package main

import (
	"encoding/json"
	"fmt"
	"strconv"

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
	app.Run(iris.Addr(":8888"))
}

// 获取沪股通、深股通、港股通每日 资金流向数据 http://localhost:8888/hsgt
func (c *lotterController) GetHsgt() string {
	var (
		res  = make(map[string]string)
		err  = &Error{}
		resp *TushareResp
	)
	defer func() {
		if err != nil {
			fmt.Printf("GetHsgt error code :%d, mssage:%s", err.Code(), err.String())
		}
	}()
	resp, err = moneyflowHsgt()
	if err != nil {
		return "empty data"
	}
	fields := resp.Data["fields"].([]interface{})
	items := resp.Data["items"].([]interface{})
	if len(items) == 0 {
		return "empty data"
	}
	first := items[0].([]interface{})
	if len(fields) != len(first) {
		err.code = ERROR_NOTMATCH
		err.message = "data length not match"
		return "empty data"
	}
	for i, d := range fields {
		k := d.(string)
		v := first[i]
		switch k {
		case trade_date:
			res["日期"] = v.(string)
		case ggt_ss:
			val := strconv.FormatFloat(v.(float64), 'f', 1, 64)
			res["港股通（上海）"] = val
		case ggt_sz:
			val := strconv.FormatFloat(v.(float64), 'f', 1, 64)
			res["港股通（深圳）"] = val
		case hgt:
			val := strconv.FormatFloat(v.(float64), 'f', 1, 64)
			res["沪股通（百万元）"] = val
		case sgt:
			val := strconv.FormatFloat(v.(float64), 'f', 1, 64)
			res["深股通（百万元）"] = val
		case north_money:
			val := strconv.FormatFloat(v.(float64), 'f', 1, 64)
			res["北上资金（百万元）"] = val
		case south_money:
			val := strconv.FormatFloat(v.(float64), 'f', 1, 64)
			res["南下资金（百万元）"] = val
		}
	}
	r, e := json.Marshal(res)
	if e != nil {
		err = &Error{
			code:    ERROR_MARSHAL,
			message: e.Error(),
		}
		return ""
	}

	return string(r)
}

// 沪深股通
func moneyflowHsgt() (tushareResp *TushareResp, err *Error) {
	err = &Error{
		code: OK,
	}
	tushareResp = &TushareResp{}
	req := &TushareReq{
		ApiName: MoneyflowHsgt,
		Token:   TOKEN,
		Params:  map[string]string{trade_date: PretradeDate()},
		Fields:  "",
	}
	resp := Post("http://api.tushare.pro", req, "")
	e := json.Unmarshal([]byte(resp), &tushareResp)
	if e != nil {
		err.code = ERROR_UNMARSHAL
		err.message = e.Error()
		return nil, err
	}
	return tushareResp, nil
}
