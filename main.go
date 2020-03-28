package main

import (
	"context"
	"github.com/go-kit/kit/sd/lb"
	"io"
	"net/url"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	http2 "github.com/go-kit/kit/transport/http"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/jaydenwen123/gokit-practice-client/services"
)

func main() {

	//client := consul.NewClient()

	//直连模式
	//endPoint := ConnectDirect()

	//todo 改造成从consul获取连接信息
	eps, err := GetConnectEndpointerWithConsul()
	if err != nil {
		logs.Error("GetConnectEndpointerWithConsul from consul occurs error:", err.Error())
		return
	}
	logs.Debug("========the round robin way....")
	balancer := lb.NewRoundRobin(eps)
	for {
		endPoint, err := balancer.Endpoint()
		if err != nil {
			logs.Error("the ConnectWithConsul get endpoint error:", err.Error())
			return
		}
		rsp, err := endPoint(context.Background(), &services.UserRequest{
			Uid:    112,
			Method: "",
		})
		realRsp := rsp.(*services.UserResponse)
		if err != nil {
			logs.Error("client send request error:", err.Error())
		} else {
			logs.Debug("the resp info:", realRsp.Result)
		}
		time.Sleep(time.Second * 1)
	}

}

//ConnectDirect 直连模式
func ConnectDirect() endpoint.Endpoint {
	tgt, err := url.Parse("http://localhost:2345")
	if err != nil {
		logs.Error("the url parse is error:", err.Error())
		panic(err)
	}
	client := http2.NewClient("GET", tgt, services.EncodeRequest, services.DecodeResponse)
	endPoint := client.Endpoint()
	return endPoint
}

//ConnectWithConsul 连接consul获取服务器连接信息
func ConnectWithConsul() ([]endpoint.Endpoint, error) {
	ep, err := GetConnectEndpointerWithConsul()
	if err != nil {
		logs.Error("the GetConnectEndpointerWithConsul error:%s", err.Error())
		return nil, err
	}

	eps, err := ep.Endpoints()
	if err != nil {
		logs.Error("the get  Endpoints occurs error:", err.Error())
		return nil, err
	}
	logs.Debug("the endpoints len:", len(eps))
	return eps, err
}

//ConnectWithConsul 连接consul获取服务器连接信息
func GetConnectEndpointerWithConsul() (*sd.DefaultEndpointer, error) {

	logger := log.NewNopLogger()
	//client = newTestClient(consulState)
	config := consulapi.DefaultConfig()
	config.Address = "127.0.0.1:8500"
	c, err := consulapi.NewClient(config)
	if err != nil {
		logs.Error("consulapi.NewClient error")
		return nil, err
	}
	client := consul.NewClient(c)

	s := consul.NewInstancer(client, logger, "userService", []string{"primary"}, true)
	//创建factory
	f := func(serviceUrl string) (endpoint.Endpoint, io.Closer, error) {
		tgt, e := url.Parse("http://" + serviceUrl)
		if e != nil {
			logs.Error("the url parse is error:", e.Error())
			return nil, nil, e
		}
		cc := http2.NewClient("GET", tgt, services.EncodeRequest, services.DecodeResponse)
		return cc.Endpoint(), nil, nil
	}
	ep := sd.NewEndpointer(s, f, logger)
	return ep, nil
}
