package middleware

import (
	"errors"
	"fmt"
	"log"
	"myproxyHttp/httpproxy/loadbalance"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"
)

//负载均衡
const (
	Random           byte = 0 //随机
	RoundRobinSimple byte = 1 //简单轮询
)

var (
//proxyGuard = 0
//mu         sync.Mutex
)

type LbProxyHandler struct {
	From        *url.URL
	ToList      []*url.URL
	Handler     []*RewriteUrlMiddleWare
	AliveHost   int                  //存活节点个数
	consisthash *ConsistHashHandler  //一致性hash
	lbrule      loadbalance.LbObject //轮询算法
	LastModify  time.Time
	sync.RWMutex
}

func (u *LbProxyHandler) checkObject() {
	if u.lbrule == nil {
		u.Lock()
		if u.lbrule == nil {
			//普通的轮询算法
			u.lbrule = loadbalance.NewRoundRobinSimple()
		}

		u.Unlock()
	}
	//u must not nil

}
func (u *LbProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	n := len(u.Handler)
	if n <= 0 {
		_, _ = fmt.Fprintln(w, "gateway error , no host is alive")
		return
	}
	if n == 1 {
		//直接请求
		u.Handler[0].ServeHTTP(w, r)
		return
	}
	//一致性hash算法
	if u.consisthash != nil {
		ok := u.consisthash.Get(w, r)
		if ok == nil {
			_, _ = fmt.Fprintln(w, "gateway error , no host is alive")
			return
		}
		ok.ServeHTTP(w, r)
		// 使用一致性hash来处理 http请求
		return
	}

	//轮询或者随机算法
	//make sure loadbalance not null
	u.checkObject()
	u.RLock()
	defer u.RUnlock()
	index := u.lbrule.NextIndex(w, r, u.AliveHost)

	u.Handler[index].ServeHTTP(w, r)

}

func NewLbProxyHandler(fromU string, toUrlList []string) (*LbProxyHandler, error) {
	if len(toUrlList) == 0 {
		return nil, errors.New("no proxy target List")
	}
	from, err := url.Parse(fromU)
	if err != nil {
		return nil, err
	}
	standardUrl(from)
	var toList = make([]*url.URL, len(toUrlList), len(toUrlList))
	var handler = make([]*RewriteUrlMiddleWare, len(toUrlList), len(toUrlList))
	var res = &LbProxyHandler{
		From:    from,
		ToList:  toList,
		Handler: handler,
	}
	res.LastModify = time.Now()
	//var proxyList = make([]*httputil.ReverseProxy, len(toUrlList), len(toUrlList))
	for i, v := range toUrlList {
		toList[i], err = url.Parse(v)
		if err != nil {
			return nil, err
		}
		//处理 后面的 /
		standardUrl(toList[i])
		handler[i] = NewUrlRewriteMiddleWare(from, toList[i])
		// 主机数量大于1
		if len(toUrlList) > 1 {
			handler[i].Proxy.ModifyResponse = res.WrapModifyResponse(handler[i].Proxy.ModifyResponse)
			handler[i].Proxy.ErrorHandler = res.WrapErrorHandler(handler[i].Proxy.ErrorHandler)
		}

	}
	res.AliveHost = len(handler)
	//default
	res.lbrule = loadbalance.NewSimpleRandomLb()
	return res, nil

}

//清理一致性hash节点
func (u *LbProxyHandler) _clearConsistHashRing() {
	//var downd []*RewriteUrlMiddleWare
	u.AliveHost = len(u.Handler)
	for _, v := range u.Handler {
		v.LiveCheck()
		if v.Down {

			if u.consisthash.Has(v.To.Host) {
				u.consisthash.RemoveProxy(v)
			}
			u.AliveHost--
		} else {
			// is up
			if !u.consisthash.Has(v.To.Host) {
				u.consisthash.Ring.AddNode(v.To.Host, v)
			}
		}
	}

}

//清理无用的 主机
func (u *LbProxyHandler) clearDiedProxy() {
	u.Lock()
	defer u.Unlock()

	//配置的只有一个主机，负载均衡不生效
	if len(u.Handler) <= 1 {
		return
	}

	//修改主机存活列表，记录修改时间
	defer func() {
		u.LastModify = time.Now()
	}()
	//使用一致性hash算法负载
	if u.consisthash != nil {
		//调整一致性hash的环结构
		u._clearConsistHashRing()
		return
	}
	//负载均衡
	if u.lbrule == nil {
		return
	}

	n := len(u.Handler)
	//if u.AliveHost > 0 {
	//for i, _ := range u.Handler {
	i := 0
	for i < n {
		_ = u.Handler[i].LiveCheck()
		if u.Handler[i].Down {
			//交换节点
			u.Handler[n-1], u.Handler[i] = u.Handler[i], u.Handler[n-1]
			n--
		} else {
			// up
			log.Printf("up host is %v\n", u.Handler[i].To)
			i++
			//to next i
		}
	}
	u.AliveHost = n
	//log.Printf("%#v\n", u.Handler)
	//}

}

func (u *LbProxyHandler) UseConsistHash() {
	u.consisthash = NewConsistHashHandler(u.Handler)
}
func (u *LbProxyHandler) UseLoadBalanceRule(f loadbalance.LbObject) {
	u.lbrule = f
}

//异常处理
func (u *LbProxyHandler) WrapErrorHandler(in func(w http.ResponseWriter, r *http.Request, err error)) func(writer http.ResponseWriter, request *http.Request, err error) {

	return func(w http.ResponseWriter, r *http.Request, err error) {

		if err != nil {

			switch err.(type) {
			case *net.OpError:
				op := err.(*net.OpError)
				if op.Op == "dial" {
					select {
					case ClearThread.Channel <- u: //通知清理协程清理掉down了的主机
					default:
					}
				}
			}

		}

		in(w, r, err)
	}
}

//
func (u *LbProxyHandler) WrapModifyResponse(in func(w *http.Response) error) func(w *http.Response) error {
	return func(w *http.Response) error {
		if u.AliveHost < len(u.Handler) && time.Now().Sub(u.LastModify).Minutes() > 4 {
			select {
			case ClearThread.Channel <- u: //通知清理协程清理或者回复down了的主机
			default:
			}
		}
		//return httpcallback.ModifyResponse(response)
		return in(w)
	}
}
