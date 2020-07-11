package main

import (
	"encoding/json"
	"time"
)

func Today() string {
	return time.Now().Format(FormatDay)
}

// 上一个交易日
func PretradeDate() string {
	tushareResp := &TushareResp{}
	req := &TushareReq{
		ApiName: TradeCal,
		Token:   TOKEN,
		Params:  map[string]string{"exchange": "SSE", start_date: Today(), end_date: Today()},
		Fields:  pretrade_date,
	}
	resp := Post("http://api.tushare.pro", req, "")
	e := json.Unmarshal([]byte(resp), &tushareResp)
	if e != nil {
		return ""
	}
	items := tushareResp.Data["items"].([]interface{})
	if len(items) == 0 {
		return ""
	}
	first := items[0].([]interface{})
	if len(first) == 0 {
		return ""
	}
	return first[0].(string)
}
