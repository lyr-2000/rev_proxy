package config

import (
	"github.com/spf13/viper"
	"net/http/httputil"
	"sync"
)

//设置权重
type BalanceWeightProxy struct {
	RevProxy *httputil.ReverseProxy
	Weight   int
}

var LoadBalanceRule struct {
	Mutex sync.Mutex
	//用于 http请求 负载均衡 ，随机选择一个 转发
	//每个端口一个协程来处理
	Port2Proxy  map[string][]*BalanceWeightProxy
	AllowIpList []string
	DenyIpList  []string
}

//init variables
//func init() {
//	LoadBalanceRule.Url2Proxy = make(map[string][]*httputil.ReverseProxy, 0)
//
//}
func InitHttpConfig() {
	viper.Get("")
}
