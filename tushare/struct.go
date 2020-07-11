package main

type TushareReq struct {
	ApiName string            `json:"api_name"`
	Token   string            `json:"token"`
	Params  map[string]string `json:"params"`
	Fields  string            `json:"fields"`
}

type TushareResp struct {
	RequestId string                 `json:"request_id"`
	Code      int                    `json:"code"`
	Data      map[string]interface{} `json:"data"`
	HasMore   bool                   `json:"has_more"`
}
