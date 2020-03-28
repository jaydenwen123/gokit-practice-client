package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"net/http"
)

//EncodeRequest 客户端对数据进行编码，然后请求服务器
func EncodeRequest(ctx context.Context, req *http.Request, r interface{}) error {
	request, ok := r.(*UserRequest)
	if !ok {
		return fmt.Errorf("the request transport failed")
	}
	uid := request.Uid
	req.URL.Path += fmt.Sprintf("/user/%d", uid)
	return nil
}

//DecodeResponse 客户端拿到服务器返回的数据进行解码
func DecodeResponse(ctx context.Context, response *http.Response) (rsp interface{}, err error) {
	if response.StatusCode > 400 {
		logs.Error("there is not data")
		return nil, fmt.Errorf("there is not data")
	}
	var resp UserResponse
	e := json.NewDecoder(response.Body).Decode(&resp)
	if e != nil {
		logs.Error("the DecodeResponse decode data error")
		return nil, e
	}
	return &resp, nil
}
