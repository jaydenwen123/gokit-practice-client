package main

import (
	"context"
	"github.com/astaxie/beego/logs"
	http2 "github.com/go-kit/kit/transport/http"
	"github.com/jaydenwen123/gokit-practice-client/services"
	"net/url"
)

func main() {

	//todo 改造成从consul获取连接信息
	tgt, err := url.Parse("http://localhost:2345")
	if err != nil {
		logs.Error("the url parse is error:", err.Error())
		panic(err)
	}
	client := http2.NewClient("GET", tgt, services.EncodeRequest, services.DecodeResponse)
	rsp, err := client.Endpoint()(context.Background(), &services.UserRequest{
		Uid:    111,
		Method: "",
	})
	realRsp := rsp.(*services.UserResponse)
	if err != nil {
		logs.Error("client send request error:", err.Error())
	} else {
		logs.Debug("the resp info:", realRsp.Result)
	}

}
